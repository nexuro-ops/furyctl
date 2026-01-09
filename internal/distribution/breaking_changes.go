// Copyright (c) 2017-present SIGHUP s.r.l All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package distribution

import (
	"fmt"
	"strings"

	"github.com/Al-Pragliola/go-version"
	"github.com/sighupio/furyctl/internal/semver"
)

// ModuleBreakingChange represents a breaking change for a specific module
type ModuleBreakingChange struct {
	Module       string // Module name (e.g., "auth", "logging")
	MinVersion   string // Minimum module version required for K8s version
	Description  string // Description of the breaking change
	MigrationURL string // URL with migration instructions (optional)
}

// KubernetesBreakingChanges represents breaking changes for a specific K8s version
type KubernetesBreakingChanges struct {
	KubernetesVersion string                    // e.g., "1.35.0"
	Changes           []ModuleBreakingChange
}

// Get all known breaking changes by Kubernetes version
func getKubernetesBreakingChanges() []KubernetesBreakingChanges {
	return []KubernetesBreakingChanges{
		{
			KubernetesVersion: "1.35.0",
			Changes: []ModuleBreakingChange{
				{
					Module:      "networking",
					MinVersion:  "v3.0.0",
					Description: "Kubernetes 1.35 requires networking module v3.0.0+ for CNI compatibility improvements",
				},
				{
					Module:      "monitoring",
					MinVersion:  "v4.0.1",
					Description: "Kubernetes 1.35 requires monitoring module v4.0.1+ for metric changes",
				},
				{
					Module:      "auth",
					MinVersion:  "v0.6.0",
					Description: "Kubernetes 1.35 requires auth module v0.6.0+ for WebSocket RBAC validation",
				},
				{
					Module:      "logging",
					MinVersion:  "v5.2.0",
					Description: "Kubernetes 1.35 requires logging module v5.2.0+ for new audit log formats",
				},
			},
		},
	}
}

// BreakingChangeValidator validates module versions for breaking changes
type BreakingChangeValidator struct {
	kubernetesVersion string
	modules           map[string]string // module name -> version
}

// NewBreakingChangeValidator creates a new breaking change validator
func NewBreakingChangeValidator(kubernetesVersion string, modules map[string]string) *BreakingChangeValidator {
	return &BreakingChangeValidator{
		kubernetesVersion: kubernetesVersion,
		modules:           modules,
	}
}

// Validate checks if current module versions are compatible with the target Kubernetes version
func (v *BreakingChangeValidator) Validate() ([]string, error) {
	var warnings []string
	var errors []string

	// Parse target Kubernetes version
	targetK8sVersion, err := semver.NewVersion(v.kubernetesVersion)
	if err != nil {
		return nil, fmt.Errorf("invalid kubernetes version format: %w", err)
	}

	// Get all breaking changes
	allChanges := getKubernetesBreakingChanges()

	// Check each breaking change to see if it applies to our target version
	for _, bc := range allChanges {
		bcVersion, err := semver.NewVersion(bc.KubernetesVersion)
		if err != nil {
			continue
		}

		// Only apply breaking changes if target version is >= breaking change version
		if targetK8sVersion.GreaterThanOrEqual(bcVersion) {
			for _, change := range bc.Changes {
				currentModuleVersion, exists := v.modules[change.Module]

				// Skip if module not specified in current config
				if !exists {
					continue
				}

				// Parse minimum required version
				minVersion, err := semver.NewVersion(change.MinVersion)
				if err != nil {
					continue
				}

				// Parse current module version
				currentVersion, err := semver.NewVersion(currentModuleVersion)
				if err != nil {
					errors = append(errors, fmt.Sprintf(
						"Module '%s' has invalid version format '%s'",
						change.Module, currentModuleVersion,
					))
					continue
				}

				// Check if current version meets minimum requirement
				if currentVersion.LessThan(minVersion) {
					msg := fmt.Sprintf(
						"Module '%s' version %s is incompatible with Kubernetes %s. "+
							"Required: %s or later. %s",
						change.Module,
						currentModuleVersion,
						bc.KubernetesVersion,
						change.MinVersion,
						change.Description,
					)
					if change.MigrationURL != "" {
						msg += fmt.Sprintf(" See: %s", change.MigrationURL)
					}
					errors = append(errors, msg)
				}
			}
		}
	}

	if len(errors) > 0 {
		return errors, fmt.Errorf("module compatibility check failed: %d incompatibilities detected", len(errors))
	}

	return warnings, nil
}

// CheckModuleCompatibility is a convenience function to validate module compatibility
func CheckModuleCompatibility(kubernetesVersion string, modules map[string]string) error {
	validator := NewBreakingChangeValidator(kubernetesVersion, modules)
	errors, err := validator.Validate()

	if err != nil {
		// Format all errors into a single message
		if len(errors) > 0 {
			return fmt.Errorf("%w:\n  - %s", err, strings.Join(errors, "\n  - "))
		}
		return err
	}

	return nil
}
