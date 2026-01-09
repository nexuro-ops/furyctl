// Copyright (c) 2017-present SIGHUP s.r.l All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package create

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

// Kubernetes135ChecksError represents a collection of errors from 1.35 specific checks
type Kubernetes135ChecksError struct {
	Blockers []string
	Warnings []string
}

func (e Kubernetes135ChecksError) Error() string {
	msg := ""
	if len(e.Blockers) > 0 {
		msg += fmt.Sprintf("BLOCKER ISSUES (must fix):\n")
		for _, b := range e.Blockers {
			msg += fmt.Sprintf("  - %s\n", b)
		}
	}
	if len(e.Warnings) > 0 {
		msg += fmt.Sprintf("WARNINGS (recommended to fix):\n")
		for _, w := range e.Warnings {
			msg += fmt.Sprintf("  - %s\n", w)
		}
	}
	return msg
}

func (e Kubernetes135ChecksError) HasBlockers() bool {
	return len(e.Blockers) > 0
}

// CheckKubernetes135Compatibility performs Kubernetes 1.35 specific validation checks for KFDDistribution
// Note: KFDDistribution-specific checks focus on cluster-level requirements only, as it operates on existing clusters
func (p *PreFlight) CheckKubernetes135Compatibility() error {
	logrus.Info("Checking Kubernetes 1.35 specific requirements...")

	checks := Kubernetes135ChecksError{
		Blockers: []string{},
		Warnings: []string{},
	}

	// Check IPVS deprecation (WARNING)
	if err := p.checkKubeProxyIPVS(); err != nil {
		checks.Warnings = append(checks.Warnings, err.Error())
	}

	// Check node OS compatibility (WARNING)
	if err := p.checkNodeOSCompatibility(); err != nil {
		checks.Warnings = append(checks.Warnings, err.Error())
	}

	// If there are blockers, return error
	if checks.HasBlockers() {
		logrus.Error("Kubernetes 1.35 compatibility checks failed:")
		return checks
	}

	// Log warnings if any
	if len(checks.Warnings) > 0 {
		logrus.Warn("Kubernetes 1.35 compatibility warnings:")
		for _, w := range checks.Warnings {
			logrus.Warnf("  - %s", w)
		}
	}

	logrus.Info("Kubernetes 1.35 compatibility checks passed")
	return nil
}

// checkKubeProxyIPVS checks if IPVS mode is being used (deprecated in 1.35)
func (p *PreFlight) checkKubeProxyIPVS() error {
	logrus.Debug("Checking kube-proxy mode...")

	// Only check if we have a valid kubeconfig
	if p.kubeRunner == nil {
		logrus.Debug("kube-proxy: skipping check - no kubeconfig available")
		return nil
	}

	// Try to get kube-proxy daemonset if it exists
	out, err := p.kubeRunner.Exec(
		"get",
		"ds",
		"-n", "kube-system",
		"-l", "component=kube-proxy",
		"-o", "jsonpath={.items[*].spec.template.spec.containers[*].command}",
	)

	if err != nil {
		logrus.Debug("kube-proxy: cannot determine mode - assuming iptables")
		return nil
	}

	output := strings.TrimSpace(string(out))
	if strings.Contains(output, "--proxy-mode=ipvs") || strings.Contains(output, "ipvs") {
		return fmt.Errorf(
			"kube-proxy: IPVS mode detected - IPVS is deprecated in Kubernetes 1.35 and will be removed in future versions. "+
				"Plan migration to nftables for next cluster update",
		)
	}

	logrus.Debug("kube-proxy mode is compatible - OK")
	return nil
}

// checkNodeOSCompatibility checks if nodes run supported OS versions
func (p *PreFlight) checkNodeOSCompatibility() error {
	logrus.Debug("Checking node OS compatibility...")

	// Only check if we have a valid kubeconfig
	if p.kubeRunner == nil {
		logrus.Debug("node OS: skipping check - no kubeconfig available")
		return nil
	}

	// Get node OS info
	out, err := p.kubeRunner.Exec(
		"get",
		"nodes",
		"-o", "jsonpath={.items[*].status.nodeInfo.osImage}",
	)

	if err != nil {
		logrus.Debug("node OS: cannot determine OS - assuming compatible")
		return nil
	}

	osImages := strings.Split(strings.TrimSpace(string(out)), " ")
	unsupportedOS := []string{}

	for _, osImage := range osImages {
		if osImage == "" {
			continue
		}

		osLower := strings.ToLower(osImage)
		isSupported := false

		// Check for supported OS versions
		if strings.Contains(osLower, "ubuntu") && strings.Contains(osLower, "22.04") {
			isSupported = true
		} else if strings.Contains(osLower, "ubuntu") && (
			strings.Contains(osLower, "24.") || strings.Contains(osLower, "25.")) {
			isSupported = true
		} else if strings.Contains(osLower, "rhel") && (
			strings.Contains(osLower, "9.") || strings.Contains(osLower, "10.")) {
			isSupported = true
		} else if strings.Contains(osLower, "centos") && strings.Contains(osLower, "9") {
			isSupported = true
		} else if strings.Contains(osLower, "debian") && (
			strings.Contains(osLower, "12") || strings.Contains(osLower, "13")) {
			isSupported = true
		} else if strings.Contains(osLower, "amazonlinux") || strings.Contains(osLower, "al2") {
			// Amazon Linux 2 is supported on KFD Distribution
			isSupported = true
		}

		if !isSupported {
			unsupportedOS = append(unsupportedOS, osImage)
		}
	}

	if len(unsupportedOS) > 0 {
		return fmt.Errorf(
			"node OS: some nodes run unsupported OS versions: %s. "+
				"Kubernetes 1.35 requires Ubuntu 22.04+, RHEL/CentOS 9+, Debian 12+, or Amazon Linux 2",
			strings.Join(unsupportedOS, ", "),
		)
	}

	logrus.Debug("node OS versions are compatible - OK")
	return nil
}
