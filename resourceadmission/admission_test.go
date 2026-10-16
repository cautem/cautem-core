package resourceadmission

import "testing"

func TestDefaultPolicyBindsResourceToWorkspace(t *testing.T) {
	policy := Default()
	labels := map[string]string{
		"openshell.ai/sandbox-attachable":           "true",
		"openshell.ai/sandbox-attachable-workspace": "team-a",
	}
	if err := policy.Admit("team-a", labels); err != nil {
		t.Fatalf("admit matching workspace volume: %v", err)
	}
	if err := policy.Admit("team-b", labels); err == nil {
		t.Fatal("cross-workspace volume was admitted")
	}
	if err := policy.Admit("team-a", nil); err == nil {
		t.Fatal("unlabeled volume was admitted")
	}
}

func TestPolicyValidationAndDisabledMode(t *testing.T) {
	policy := Config{Enabled: true}
	if err := policy.Validate(); err == nil {
		t.Fatal("enabled empty policy was accepted")
	}
	policy = Config{Enabled: true, RequiredLabels: map[string]string{"cautem.internal/workspace": "team-a", "example.com/ok": "yes"}}
	if err := policy.Validate(); err == nil {
		t.Fatal("driver-owned label was accepted")
	}
	policy = Config{Enabled: false}
	if err := policy.RejectUnlabelable("bind mount"); err != nil {
		t.Fatalf("disabled admission rejected bind mount: %v", err)
	}
}

func TestRejectUnlabelableUnderDefaultPolicy(t *testing.T) {
	if err := Default().RejectUnlabelable("host bind mount"); err == nil {
		t.Fatal("unlabelable host mount was admitted")
	}
}
