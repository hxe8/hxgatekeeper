package tests

import (
	"testing"

	"github.com/hex/gatekeeper-go/internal/auth"
)

func TestParseServiceToken(t *testing.T) {
	subject, err := auth.ParseServiceToken("svc:deploy-agent|roles=ops,deploy")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if subject.ID != "deploy-agent" {
		t.Fatalf("unexpected subject id: %s", subject.ID)
	}
	if len(subject.Roles) != 2 {
		t.Fatalf("expected two roles, got %d", len(subject.Roles))
	}
}
