package policy_test

import (
	"testing"

	"github.com/cautem/quarshell-core/policy"
)

func TestInferenceProfiles(t *testing.T) {
	inf := &policy.Inference{
		Profiles: []policy.ProviderProfile{{
			ID: "vllm", Host: "host.cautem.internal", Port: 8000,
			EnvKeys: []string{"OPENAI_API_KEY"}, Refresh: "env",
		}},
	}
	rules, err := policy.ExpandInferenceRules(inf)
	if err != nil {
		t.Fatal(err)
	}
	if len(rules) != 1 || rules[0].Host != "host.cautem.internal" {
		t.Fatalf("rules=%v", rules)
	}
	keys := policy.EnvKeysForInference(inf)
	if len(keys) != 1 || keys[0] != "OPENAI_API_KEY" {
		t.Fatalf("keys=%v", keys)
	}
}
