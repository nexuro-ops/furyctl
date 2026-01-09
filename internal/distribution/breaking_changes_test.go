// Copyright (c) 2017-present SIGHUP s.r.l All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package distribution

import (
	"testing"
)

func TestBreakingChangeValidator_ValidModules(t *testing.T) {
	tests := []struct {
		name              string
		kubernetesVersion string
		modules           map[string]string
		shouldErr         bool
	}{
		{
			name:              "Kubernetes 1.35.0 with compatible modules",
			kubernetesVersion: "1.35.0",
			modules: map[string]string{
				"auth":       "v0.6.0",
				"logging":    "v5.2.0",
				"networking": "v3.0.0",
				"monitoring": "v4.0.1",
			},
			shouldErr: false,
		},
		{
			name:              "Kubernetes 1.35.5 with compatible modules",
			kubernetesVersion: "1.35.5",
			modules: map[string]string{
				"auth":       "v0.6.0",
				"logging":    "v5.2.0",
				"networking": "v3.0.0",
				"monitoring": "v4.0.1",
			},
			shouldErr: false,
		},
		{
			name:              "Kubernetes 1.35.0 with newer modules",
			kubernetesVersion: "1.35.0",
			modules: map[string]string{
				"auth":       "v0.7.0",
				"logging":    "v5.3.0",
				"networking": "v3.1.0",
				"monitoring": "v4.1.0",
			},
			shouldErr: false,
		},
		{
			name:              "Kubernetes 1.34.0 (before 1.35) with newer modules",
			kubernetesVersion: "1.34.0",
			modules: map[string]string{
				"auth":       "v0.6.0",
				"logging":    "v5.2.0",
			},
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := NewBreakingChangeValidator(tt.kubernetesVersion, tt.modules)
			_, err := validator.Validate()

			if (err != nil) != tt.shouldErr {
				t.Errorf("Validate() error = %v, shouldErr %v", err, tt.shouldErr)
			}
		})
	}
}

func TestBreakingChangeValidator_IncompatibleModules(t *testing.T) {
	tests := []struct {
		name              string
		kubernetesVersion string
		modules           map[string]string
		shouldErr         bool
		expectedErrors    int
	}{
		{
			name:              "Kubernetes 1.35.0 with old auth module",
			kubernetesVersion: "1.35.0",
			modules: map[string]string{
				"auth":       "v0.5.0",
				"logging":    "v5.2.0",
				"networking": "v3.0.0",
				"monitoring": "v4.0.1",
			},
			shouldErr:      true,
			expectedErrors: 1,
		},
		{
			name:              "Kubernetes 1.35.0 with multiple old modules",
			kubernetesVersion: "1.35.0",
			modules: map[string]string{
				"auth":       "v0.5.0",
				"logging":    "v5.1.0",
				"networking": "v2.9.0",
				"monitoring": "v4.0.0",
			},
			shouldErr:      true,
			expectedErrors: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := NewBreakingChangeValidator(tt.kubernetesVersion, tt.modules)
			errors, err := validator.Validate()

			if (err != nil) != tt.shouldErr {
				t.Errorf("Validate() error = %v, shouldErr %v", err, tt.shouldErr)
			}

			if tt.shouldErr && len(errors) != tt.expectedErrors {
				t.Errorf("Expected %d errors, got %d errors: %v", tt.expectedErrors, len(errors), errors)
			}
		})
	}
}

func TestCheckModuleCompatibility(t *testing.T) {
	tests := []struct {
		name              string
		kubernetesVersion string
		modules           map[string]string
		shouldErr         bool
	}{
		{
			name:              "Valid configuration",
			kubernetesVersion: "1.35.0",
			modules: map[string]string{
				"auth":       "v0.6.0",
				"logging":    "v5.2.0",
				"networking": "v3.0.0",
				"monitoring": "v4.0.1",
			},
			shouldErr: false,
		},
		{
			name:              "Invalid Kubernetes version format",
			kubernetesVersion: "invalid",
			modules: map[string]string{
				"auth": "v0.6.0",
			},
			shouldErr: true,
		},
		{
			name:              "Empty modules",
			kubernetesVersion: "1.35.0",
			modules:           map[string]string{},
			shouldErr:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckModuleCompatibility(tt.kubernetesVersion, tt.modules)

			if (err != nil) != tt.shouldErr {
				t.Errorf("CheckModuleCompatibility() error = %v, shouldErr %v", err, tt.shouldErr)
			}
		})
	}
}
