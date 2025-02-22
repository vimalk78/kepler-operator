#!/usr/bin/env bash
set -e -u -o pipefail

trap cleanup INT

PROJECT_ROOT="$(git rev-parse --show-toplevel)"
declare -r PROJECT_ROOT

source "$PROJECT_ROOT/hack/utils.bash"

declare -r OPERATOR="kepler-operator"
declare -r OLM_CATALOG="kepler-operator-catalog"
declare -r VERSION="0.0.0-e2e"
declare -r OPERATOR_DEPLOY_YAML="config/manager/manager.yaml"
declare -r OPERATOR_DEPLOY_NAME="kepler-operator-controller-manager"

declare IMG_BASE="${IMG_BASE:-localhost:5001/$OPERATOR}"
declare OPERATOR_IMG="$IMG_BASE/$OPERATOR:$VERSION"
declare BUNDLE_IMG="$IMG_BASE/$OPERATOR-bundle:$VERSION"

declare CI_MODE=false
declare NO_DEPLOY=false
declare NO_BUILDS=false
declare SHOW_USAGE=false
declare LOGS_DIR="tmp/e2e"
declare OPERATORS_NS="operators"
declare TEST_TIMEOUT="15m"
# declare -a PORT_FORWARDED_PIDS=()

cleanup() {
	info "Cleaning up ..."
	# shell check  ignore word splitting when using jobs -p
	# shellcheck disable=SC2046
	[[ -z "$(jobs -p)" ]] || kill $(jobs -p) || true

	return 0
}

delete_olm_subscription() {
	header "Delete Old Deployments"

	$CI_MODE && {
		ok "skipping deletion of old deployment in CI mode"
		return 0
	}

	kubectl delete -n "$OPERATORS_NS" csv --all || true
	kubectl delete -n "$OPERATORS_NS" installplan,subscriptions,catalogsource \
		-l operators.coreos.com/$OPERATOR.operators= || true
	kubectl delete -n "$OPERATORS_NS" installplan,subscriptions,catalogsource \
		-l operators.coreos.com/$OPERATOR.openshift-operators= || true
	kubectl delete -n "$OPERATORS_NS" catalogsource $OLM_CATALOG || true
}

build_bundle() {
	header "Build Operator Bundle"

	$NO_BUILDS && {
		info "skipping building of images"
		return 0
	}

	run make operator-build bundle bundle-build \
		OPERATOR_IMG="$OPERATOR_IMG" \
		BUNDLE_IMG="$BUNDLE_IMG" \
		VERSION="$VERSION"
}

push_bundle() {
	header "Push Operator Bundle Images"
	$NO_BUILDS && {
		info "skipping pushing images"
		return 0
	}

	run make operator-push bundle-push \
		OPERATOR_IMG="$OPERATOR_IMG" \
		BUNDLE_IMG="$BUNDLE_IMG" \
		VERSION="$VERSION"

}

run_bundle() {
	header "Running Bundle"

	run ./tmp/bin/operator-sdk run bundle "$BUNDLE_IMG" \
		--install-mode AllNamespaces --namespace "$OPERATORS_NS" --skip-tls
}

log_events() {
	local ns="$1"
	shift
	kubectl get events -w \
		-o custom-columns=FirstSeen:.firstTimestamp,LastSeen:.lastTimestamp,Count:.count,From:.source.component,Type:.type,Reason:.reason,Message:.message \
		-n "$ns" | tee "$LOGS_DIR/$ns-events.log"
}

watch_operator_errors() {
	local err_log="$1"
	shift

	kubectl logs -f -n "$OPERATORS_NS" "deploy/$OPERATOR_DEPLOY_NAME" |
		grep -i error | tee "$err_log"
}

run_e2e() {
	header "Running e2e tests"

	local error_log="$LOGS_DIR/operator-errors.log"

	log_events "$OPERATORS_NS" &
	log_events "e2e-tests" &
	watch_operator_errors "$error_log" &

	local ret=0
	go test -v -failfast -timeout $TEST_TIMEOUT \
		./tests/e2e/... | tee "$LOGS_DIR/e2e.log" || ret=1

	# terminate both log_events
	{ jobs -p | xargs -I {} -- pkill -TERM -P {}; } || true
	wait
	sleep 1 #  wait for the Termiated logs to be written

	if [[ "$ret" -ne 0 ]]; then
		# logging of errors may not be immediate, so it is better to read logs again
		# than dumping the $error_log file
		sleep 2
		info "Operator Error Logs"
		line 50
		kubectl logs -n "$OPERATORS_NS" "deploy/$OPERATOR_DEPLOY_NAME" |
			grep -i error | tee "$error_log"
		line 50
	fi

	return $ret
}

parse_args() {
	### while there are args parse them
	while [[ -n "${1+xxx}" ]]; do
		case $1 in
		-h | --help)
			SHOW_USAGE=true
			break
			;; # exit the loop
		--no-deploy)
			NO_DEPLOY=true
			shift
			;;
		--no-builds)
			NO_BUILDS=true
			shift
			;;
		--ci)
			CI_MODE=true
			shift
			;;
		--image-base)
			shift
			IMG_BASE="$1"
			shift
			;;
		--ns)
			shift
			OPERATORS_NS=$1
			shift
			;;
		*) return 1 ;; # show usage on everything else
		esac
	done

	# NOTE: if you modify this update the declare at the top of the file
	OPERATOR_IMG="$IMG_BASE/$OPERATOR:$VERSION"
	BUNDLE_IMG="$IMG_BASE/$OPERATOR-bundle:$VERSION"
	return 0
}

print_usage() {
	local scr
	scr="$(basename "$0")"

	read -r -d '' help <<-EOF_HELP || true
		Usage:
		  $scr
		  $scr  --no-deploy
		  $scr  -h|--help


		Options:
		  -h|--help        show this help
		  --ci             run in CI mode
		  --no-deploy      do not build and deploy Operator; useful for rerunning tests
		  --no-builds      skip building operator images; useful when operator image is already
		                   built and pushed
		  --ns NAMESPACE   namespace to deploy operators (default: $OPERATORS_NS)
		                   E.g. running against openshift use --ns openshift-operators


	EOF_HELP

	echo -e "$help"
	return 0
}

init_logs_dir() {
	rm -rf "$LOGS_DIR-prev"
	mv "$LOGS_DIR" "$LOGS_DIR-prev" || true
	mkdir -p "$LOGS_DIR"
}

restart_operator() {
	header "Restart Operator deployment"

	ensure_deploy_img_is_always_pulled || return 1
	local deployment="deployment/$OPERATOR_DEPLOY_NAME"

	info "scale down Operator"
	run kubectl scale -n "$OPERATORS_NS" --replicas=0 "$deployment"
	run kubectl wait -n "$OPERATORS_NS" --for=delete \
		pods -l app.kubernetes.io/component=operator --timeout=60s

	info "scale up"
	kubectl scale -n "$OPERATORS_NS" --replicas=1 "$deployment"
	wait_for_operator "$OPERATORS_NS"

	ok "Operator deployment restarted"

}

update_cluster_mon_crds() {
	# try replacing any installed crds; failure is often because the
	# CRDs are absent and in that case, try creating and fail if that fails

	run kubectl apply --server-side --force-conflicts \
		-k config/crd

	run kubectl wait --for=condition=Established crds --all --timeout=120s

	return 0
}

deploy_operator() {
	header "Build and Deploy Operator"

	delete_olm_subscription || true
	ensure_imgpullpolicy_always_in_yaml
	update_cluster_mon_crds
	build_bundle
	push_bundle
	run_bundle
	wait_for_operator "$OPERATORS_NS"
}

ensure_imgpullpolicy_always_in_yaml() {
	$CI_MODE && {
		ok "skipping check of imagePullPolicy in deployment yaml"
		return 0
	}

	local pull_policy
	pull_policy=$(grep '\s\+imagePullPolicy:' "$OPERATOR_DEPLOY_YAML" | tr -d ' ' | cut -f2 -d:)

	[[ "$pull_policy" != "Always" ]] && {
		info "Modify $OPERATOR_DEPLOY_YAML imagePullPolicy -> Always"
		info "  ❯ sed -e 's|imagePullPolicy: .*|imagePullPolicy: Always|g' -i $OPERATOR_DEPLOY_YAML"
		warn "Deployment's imagePullPolicy must be Always instead of $pull_policy"
		return 1
	}

	ok "Operator deployment yaml imagePullPolicy is Always"
}

ensure_deploy_img_is_always_pulled() {
	$CI_MODE && {
		ok "skipping check of imagePullPolicy of Operator deployment"
		return 0
	}

	local pull_policy
	pull_policy=$(kubectl get deploy/$OPERATOR_DEPLOY_NAME \
		-n "$OPERATORS_NS" \
		-ojsonpath='{.spec.template.spec.containers[1].imagePullPolicy}')

	if [[ "$pull_policy" != "Always" ]]; then
		info "Edit $OPERATOR_DEPLOY_YAML imagePullPolicy and redeploy"
		info "  ❯ sed -e 's|imagePullPolicy: .*|imagePullPolicy: Always|g' -i $OPERATOR_DEPLOY_YAML"
		warn "Deployment's imagePullPolicy must be Always instead of $pull_policy"
		return 1
	fi
	ok "Operator deployment imagePullPolicy is Always"
}

reset_env() {
	kubectl delete --wait ns e2e-tests || true
}

# wait_for_operator requires the namespace where the operator is installed
wait_for_operator() {
	header "Waiting  for Kepler Operator ($OPERATORS_NS) to be Ready"
	local deployment="deploy/$OPERATOR_DEPLOY_NAME"

	wait_until 30 10 "operator to run" \
		kubectl -n "$OPERATORS_NS" rollout status "$deployment"

	run kubectl wait -n "$OPERATORS_NS" --for=condition=Available \
		--timeout=300s "$deployment"

	ok "Operator up and running"
}

print_config() {
	header "Test Configuration"
	cat <<-EOF
		  image base:      $IMG_BASE
		  operator image:  $OPERATOR_IMG
		  bundle:          $BUNDLE_IMG
		  CI Mode:         $CI_MODE
		  Skip Builds:     $NO_BUILDS
		  Skip Deploy:     $NO_DEPLOY
		  Operator namespace: $OPERATORS_NS
		  Logs directory: $LOGS_DIR

	EOF
	line 50
}

main() {
	parse_args "$@" || die "parse args failed"
	$SHOW_USAGE && {
		print_usage
		exit 0
	}

	cd "$PROJECT_ROOT"

	# delete the e2e-tests but continue deploying the operator
	reset_env & # note must wait before running tests
	init_logs_dir
	print_config

	if ! $NO_DEPLOY; then
		deploy_operator
	else
		restart_operator || die "restarting operator failed 🤕"
	fi

	# wait for the deletion to complete before running tests
	wait

	local -i ret=0
	run_e2e || ret=$?
	info "e2e test - exit code: $ret"
	line 50 heavy

	return $ret
}

main "$@"
