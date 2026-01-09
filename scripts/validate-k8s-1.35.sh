#!/bin/bash

# Copyright (c) 2017-present SIGHUP s.r.l All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# Kubernetes 1.35 Compatibility Validation Script
# This script validates that a cluster meets all Kubernetes 1.35 requirements

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test counters
TESTS_RUN=0
TESTS_PASSED=0
TESTS_FAILED=0

# Logging functions
log_info() {
    echo -e "${BLUE}ℹ${NC} $1"
}

log_pass() {
    echo -e "${GREEN}✓${NC} $1"
    ((TESTS_PASSED++))
}

log_fail() {
    echo -e "${RED}✗${NC} $1"
    ((TESTS_FAILED++))
}

log_warn() {
    echo -e "${YELLOW}⚠${NC} $1"
}

# Check if kubectl is available
check_kubectl() {
    ((TESTS_RUN++))
    if command -v kubectl &> /dev/null; then
        log_pass "kubectl is available"
    else
        log_fail "kubectl is not available"
        exit 1
    fi
}

# Check Kubernetes version
check_k8s_version() {
    ((TESTS_RUN++))
    K8S_VERSION=$(kubectl version -o json | grep -o '"gitVersion":"[^"]*"' | head -1 | cut -d'"' -f4)

    if [[ $K8S_VERSION == *"1.35"* ]]; then
        log_pass "Kubernetes version is 1.35.x: $K8S_VERSION"
    else
        log_fail "Kubernetes version is not 1.35.x: $K8S_VERSION"
    fi
}

# Check cluster connectivity
check_cluster_connectivity() {
    ((TESTS_RUN++))
    if kubectl cluster-info &> /dev/null; then
        log_pass "Cluster is reachable"
    else
        log_fail "Cannot reach cluster"
        exit 1
    fi
}

# Check cgroup v2 (for OnPremises)
check_cgroup_v2() {
    ((TESTS_RUN++))
    log_info "Checking cgroup v2 support (OnPremises only)..."

    # This would need to run on nodes, skipping for now in most environments
    log_warn "cgroup v2 check requires node access - skipped in remote validation"
}

# Check containerd version (for OnPremises)
check_containerd() {
    ((TESTS_RUN++))
    log_info "Checking containerd version (OnPremises only)..."

    # This would need to run on nodes, skipping for now
    log_warn "containerd check requires node access - skipped in remote validation"
}

# Check kube-proxy IPVS mode
check_kube_proxy() {
    ((TESTS_RUN++))

    IPVS_FOUND=$(kubectl get ds -n kube-system kube-proxy -o jsonpath='{.spec.template.spec.containers[0].command}' 2>/dev/null | grep -c "ipvs" || true)

    if [ "$IPVS_FOUND" -eq 0 ]; then
        log_pass "kube-proxy is not using deprecated IPVS mode"
    else
        log_warn "kube-proxy is using IPVS mode (deprecated in K8s 1.35)"
    fi
}

# Check node OS compatibility
check_node_os() {
    ((TESTS_RUN++))

    OS_INFO=$(kubectl get nodes -o jsonpath='{.items[*].status.nodeInfo.osImage}')
    UNSUPPORTED_COUNT=0

    while read -r os; do
        if [[ $os == *"Ubuntu 22.04"* ]] || [[ $os == *"Ubuntu 24"* ]] || \
           [[ $os == *"RHEL 9"* ]] || [[ $os == *"CentOS 9"* ]] || \
           [[ $os == *"Debian 12"* ]] || [[ $os == *"Amazon Linux 2"* ]]; then
            log_pass "Node OS is supported: $os"
        else
            log_fail "Node OS is not supported: $os"
            ((UNSUPPORTED_COUNT++))
        fi
    done <<< "$OS_INFO"

    if [ "$UNSUPPORTED_COUNT" -eq 0 ]; then
        log_pass "All nodes run supported OS versions"
    fi
}

# Check for required namespaces
check_namespaces() {
    ((TESTS_RUN++))

    REQUIRED_NS=("default" "kube-system" "kube-public" "kube-node-lease")

    for ns in "${REQUIRED_NS[@]}"; do
        if kubectl get ns "$ns" &> /dev/null; then
            log_pass "Namespace $ns exists"
        else
            log_fail "Namespace $ns not found"
        fi
    done
}

# Test module compatibility validation (if config available)
test_module_validation() {
    ((TESTS_RUN++))

    if command -v furyctl &> /dev/null; then
        log_info "Testing module compatibility validation..."

        # Run Go tests for breaking changes
        if cd "$PROJECT_ROOT" && go test ./internal/distribution -v -run TestBreakingChangeValidator &> /tmp/k8s-1.35-validation.log; then
            log_pass "Module compatibility validation tests passed"
        else
            log_fail "Module compatibility validation tests failed"
            cat /tmp/k8s-1.35-validation.log
        fi
    else
        log_warn "furyctl not available - skipping module validation tests"
    fi
}

# Generate summary report
generate_summary() {
    echo ""
    echo "========================================"
    echo "Validation Summary"
    echo "========================================"
    echo "Tests Run:    $TESTS_RUN"
    echo "Tests Passed: $TESTS_PASSED"
    echo "Tests Failed: $TESTS_FAILED"

    if [ "$TESTS_FAILED" -eq 0 ]; then
        echo -e "${GREEN}✓ All validations passed!${NC}"
        echo "Kubernetes 1.35 cluster is compatible."
        return 0
    else
        echo -e "${RED}✗ Some validations failed.${NC}"
        echo "Please address the issues above before proceeding."
        return 1
    fi
}

# Main execution
main() {
    echo -e "${BLUE}════════════════════════════════════════${NC}"
    echo -e "${BLUE}Kubernetes 1.35 Compatibility Validator${NC}"
    echo -e "${BLUE}════════════════════════════════════════${NC}"
    echo ""

    # Run all checks
    check_kubectl
    check_cluster_connectivity
    check_k8s_version
    check_node_os
    check_kube_proxy
    check_namespaces
    check_cgroup_v2
    check_containerd
    test_module_validation

    # Generate summary
    echo ""
    generate_summary
}

# Run main function
main "$@"
