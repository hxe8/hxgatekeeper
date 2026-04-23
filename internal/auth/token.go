package auth

import (
	"errors"
	"net/http"
	"strings"
)

type Subject struct {
	ID    string
	Roles []string
}

func ParseBearerToken(header string) (Subject, error) {
	parts := strings.SplitN(strings.TrimSpace(header), " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return Subject{}, errors.New("missing bearer token")
	}
	return ParseServiceToken(parts[1])
}

func ParseServiceToken(raw string) (Subject, error) {
	raw = strings.TrimSpace(raw)
	if !strings.HasPrefix(raw, "svc:") {
		return Subject{}, errors.New("unsupported token format")
	}

	segments := strings.SplitN(strings.TrimPrefix(raw, "svc:"), "|", 2)
	subjectID := strings.TrimSpace(segments[0])
	if subjectID == "" {
		return Subject{}, errors.New("empty subject")
	}

	subject := Subject{ID: subjectID}
	if len(segments) == 2 {
		for _, tokenPart := range strings.Split(segments[1], ";") {
			kv := strings.SplitN(strings.TrimSpace(tokenPart), "=", 2)
			if len(kv) != 2 {
				continue
			}
			if kv[0] == "roles" {
				for _, role := range strings.Split(kv[1], ",") {
					role = strings.TrimSpace(role)
					if role != "" {
						subject.Roles = append(subject.Roles, role)
					}
				}
			}
		}
	}
	return subject, nil
}

func SubjectFromRequest(r *http.Request) (Subject, error) {
	return ParseBearerToken(r.Header.Get("Authorization"))
}
