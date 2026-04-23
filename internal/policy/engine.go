package policy

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/hex/gatekeeper-go/internal/auth"
)

type Rule struct {
	Name        string              `json:"name"`
	Subjects    []string            `json:"subjects"`
	Actions     []string            `json:"actions"`
	Resources   []string            `json:"resources"`
	Constraints map[string][]string `json:"constraints"`
	Effect      string              `json:"effect"`
}

type Input struct {
	Subject    auth.Subject
	Action     string
	Resource   string
	Attributes map[string]string
}

type Decision struct {
	Effect string
	Rule   string
	Reason string
}

type Engine struct {
	Rules       []Rule
	DefaultMode string
}

func LoadFromFile(path, defaultMode string) (*Engine, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read policy file: %w", err)
	}
	var rules []Rule
	if err := json.Unmarshal(data, &rules); err != nil {
		return nil, fmt.Errorf("decode policy file: %w", err)
	}
	return &Engine{Rules: rules, DefaultMode: strings.ToLower(defaultMode)}, nil
}

func (e *Engine) Evaluate(in Input) Decision {
	for _, rule := range e.Rules {
		if !matchesSubjects(rule.Subjects, in.Subject) {
			continue
		}
		if !matches(rule.Actions, in.Action) {
			continue
		}
		if !matches(rule.Resources, in.Resource) {
			continue
		}
		if !matchesConstraints(rule.Constraints, in.Attributes) {
			continue
		}
		return Decision{
			Effect: normalizeEffect(rule.Effect, e.DefaultMode),
			Rule:   rule.Name,
			Reason: "matched policy rule",
		}
	}
	return Decision{
		Effect: normalizeEffect(e.DefaultMode, "deny"),
		Rule:   "default",
		Reason: "no matching policy rule",
	}
}

func normalizeEffect(effect, fallback string) string {
	effect = strings.ToLower(strings.TrimSpace(effect))
	if effect == "allow" || effect == "deny" {
		return effect
	}
	return strings.ToLower(strings.TrimSpace(fallback))
}

func matches(values []string, target string) bool {
	target = strings.TrimSpace(target)
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "*" || value == target {
			return true
		}
	}
	return false
}

func matchesSubjects(subjectRefs []string, subject auth.Subject) bool {
	for _, ref := range subjectRefs {
		ref = strings.TrimSpace(ref)
		if ref == "*" || ref == subject.ID {
			return true
		}
		if strings.HasPrefix(ref, "role:") {
			needle := strings.TrimPrefix(ref, "role:")
			for _, role := range subject.Roles {
				if role == needle {
					return true
				}
			}
		}
	}
	return false
}

func matchesConstraints(expected map[string][]string, attrs map[string]string) bool {
	for key, allowed := range expected {
		value, ok := attrs[key]
		if !ok {
			return false
		}
		if !matches(allowed, value) {
			return false
		}
	}
	return true
}
