package tests

import (
	"testing"

	"github.com/hex/gatekeeper-go/internal/auth"
	"github.com/hex/gatekeeper-go/internal/policy"
)

func TestPolicyEvaluateAllow(t *testing.T) {
	engine := &policy.Engine{
		DefaultMode: "deny",
		Rules: []policy.Rule{
			{
				Name:      "ops-read-incidents",
				Subjects:  []string{"role:ops"},
				Actions:   []string{"read"},
				Resources: []string{"incidents"},
				Constraints: map[string][]string{
					"environment": {"prod"},
				},
				Effect: "allow",
			},
		},
	}

	decision := engine.Evaluate(policy.Input{
		Subject: auth.Subject{ID: "ops-bot", Roles: []string{"ops"}},
		Action:  "read",
		Resource: "incidents",
		Attributes: map[string]string{
			"environment": "prod",
		},
	})

	if decision.Effect != "allow" {
		t.Fatalf("expected allow, got %s", decision.Effect)
	}
}

func TestPolicyEvaluateDefaultDeny(t *testing.T) {
	engine := &policy.Engine{DefaultMode: "deny"}
	decision := engine.Evaluate(policy.Input{
		Subject:  auth.Subject{ID: "guest"},
		Action:   "delete",
		Resource: "incidents",
	})
	if decision.Effect != "deny" {
		t.Fatalf("expected deny, got %s", decision.Effect)
	}
}
