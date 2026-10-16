// Package resourceadmission validates operator policy for caller supplied
// backend configuration that can attach externally managed resources.
package resourceadmission

import (
	"fmt"
	"strings"
)

const workspacePlaceholder = "${workspace}"

var reservedLabelKeys = map[string]struct{}{
	"openshell.ai/managed-by":                       {},
	"openshell.ai/gateway-id":                       {},
	"openshell.ai/sandbox-workspace":                {},
	"cautem.internal/workspace":                     {},
	"cautem.internal/gateway-id":                    {},
	"cautem.internal/driver-config-used":            {},
	"cautem.internal/resource-admission-identities": {},
}

// Config is the operator-owned label policy for caller attached volumes.
// A supplied RequiredLabels map replaces the defaults.
type Config struct {
	Enabled        bool              `json:"enabled"`
	RequiredLabels map[string]string `json:"required_labels"`
}

// Default returns the pinned v0.1.4 admission policy.
func Default() Config {
	return Config{
		Enabled: true,
		RequiredLabels: map[string]string{
			"openshell.ai/sandbox-attachable":           "true",
			"openshell.ai/sandbox-attachable-workspace": workspacePlaceholder,
		},
	}
}

// Validate checks keys, values, reserved metadata, and policy completeness.
func (c Config) Validate() error {
	if c.Enabled && len(c.RequiredLabels) == 0 {
		return fmt.Errorf("resource_admission.required_labels must not be empty while enabled")
	}
	workspaceOnly := len(c.RequiredLabels) > 0
	for key, value := range c.RequiredLabels {
		if _, reserved := reservedLabelKeys[key]; reserved {
			return fmt.Errorf("resource admission label key is reserved for driver-owned metadata: %s", key)
		}
		if !validLabelKey(key) {
			return fmt.Errorf("invalid resource admission label key: %s", key)
		}
		if value == workspacePlaceholder {
			continue
		}
		workspaceOnly = false
		if !validLabelValue(value) {
			return fmt.Errorf("invalid resource admission label value for %s", key)
		}
	}
	if c.Enabled && workspaceOnly {
		return fmt.Errorf("resource_admission.required_labels must include a shared approval label")
	}
	return nil
}

// Admit permits a named volume only when its engine-owned labels match the
// configured policy for the verified workspace identity.
func (c Config) Admit(workspace string, labels map[string]string) error {
	if err := c.Validate(); err != nil {
		return fmt.Errorf("invalid resource admission policy: %w", err)
	}
	if !c.Enabled {
		return nil
	}
	for key, expected := range c.RequiredLabels {
		if expected == workspacePlaceholder {
			if workspace == "" || !validLabelValue(workspace) {
				return fmt.Errorf("resource not admitted: invalid workspace identity")
			}
			expected = workspace
		}
		if labels[key] != expected {
			return fmt.Errorf("external resource not admitted by required labels")
		}
	}
	return nil
}

// RejectUnlabelable denies an attachment when the backend cannot prove its
// identity through trusted engine labels.
func (c Config) RejectUnlabelable(kind string) error {
	if err := c.Validate(); err != nil {
		return fmt.Errorf("invalid resource admission policy: %w", err)
	}
	if c.Enabled {
		return fmt.Errorf("%s cannot be attached while resource admission is enabled: no trusted label resolver", kind)
	}
	return nil
}

func validLabelValue(value string) bool {
	if value == "" {
		return true
	}
	if len(value) > 63 || !isASCIIAlphaNumeric(value[0]) || !isASCIIAlphaNumeric(value[len(value)-1]) {
		return false
	}
	for _, ch := range []byte(value) {
		if !isASCIIAlphaNumeric(ch) && ch != '-' && ch != '_' && ch != '.' {
			return false
		}
	}
	return true
}

func validLabelKey(key string) bool {
	prefix, name, hasPrefix := strings.Cut(key, "/")
	if !hasPrefix {
		name = key
	}
	if !validLabelValue(name) || name == "" {
		return false
	}
	if !hasPrefix {
		return true
	}
	if prefix == "" || len(prefix) > 253 {
		return false
	}
	for _, part := range strings.Split(prefix, ".") {
		if part == "" || len(part) > 63 || !isASCIILowerAlphaNumeric(part[0]) || !isASCIILowerAlphaNumeric(part[len(part)-1]) {
			return false
		}
		for _, ch := range []byte(part) {
			if !isASCIILowerAlphaNumeric(ch) && ch != '-' {
				return false
			}
		}
	}
	return true
}

func isASCIIAlphaNumeric(ch byte) bool {
	return ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z' || ch >= '0' && ch <= '9'
}

func isASCIILowerAlphaNumeric(ch byte) bool {
	return ch >= 'a' && ch <= 'z' || ch >= '0' && ch <= '9'
}
