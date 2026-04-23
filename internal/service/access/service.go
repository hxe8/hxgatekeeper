package access

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/hex/gatekeeper-go/internal/audit"
	"github.com/hex/gatekeeper-go/internal/auth"
	"github.com/hex/gatekeeper-go/internal/policy"
)

type Request struct {
	Subject    auth.Subject
	Action     string
	Resource   string
	Attributes map[string]string
}

type Response struct {
	Decision string `json:"decision"`
	Reason   string `json:"reason"`
	TraceID  string `json:"trace_id"`
}

type Service struct {
	Engine *policy.Engine
	Audit  *audit.Store
}

func (s *Service) Authorize(req Request) Response {
	traceID := newTraceID()
	decision := s.Engine.Evaluate(policy.Input{
		Subject:    req.Subject,
		Action:     req.Action,
		Resource:   req.Resource,
		Attributes: req.Attributes,
	})

	s.Audit.Append(audit.Event{
		TraceID:   traceID,
		Subject:   req.Subject.ID,
		Action:    req.Action,
		Resource:  req.Resource,
		Decision:  decision.Effect,
		Rule:      decision.Rule,
		CreatedAt: time.Now().UTC(),
	})

	return Response{
		Decision: decision.Effect,
		Reason:   decision.Reason,
		TraceID:  traceID,
	}
}

func newTraceID() string {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "trace-unavailable"
	}
	return hex.EncodeToString(buf)
}
