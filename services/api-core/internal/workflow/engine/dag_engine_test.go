package engine

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/google/uuid"
)

type mockExecutor struct {
	execCount int32
	failUntil int32
}

func (m *mockExecutor) Execute(ctx context.Context, node Node, input []byte) ([]byte, error) {
	current := atomic.AddInt32(&m.execCount, 1)
	if current <= m.failUntil {
		return nil, errors.New("simulated transient failure")
	}
	return []byte(`{"status":"ok"}`), nil
}

func TestDAGValidationCycle(t *testing.T) {
	nodeA := uuid.New()
	nodeB := uuid.New()

	nodes := []Node{
		{ID: nodeA, StepType: StepTypeHTTPRequest},
		{ID: nodeB, StepType: StepTypeDelay},
	}

	edges := []Edge{
		{SourceStepID: nodeA, TargetStepID: nodeB},
		{SourceStepID: nodeB, TargetStepID: nodeA}, // cycle!
	}

	res := ValidateDAG(nodes, edges)
	if res.IsValid {
		t.Fatalf("expected cycle detection, but DAG was validated as valid")
	}
}

func TestDAGConcurrentBranches(t *testing.T) {
	nodeA := uuid.New()
	nodeB := uuid.New()
	nodeC := uuid.New()
	nodeD := uuid.New()

	nodes := []Node{
		{ID: nodeA, StepType: StepTypeWebhookTrigger},
		{ID: nodeB, StepType: StepTypeHTTPRequest},
		{ID: nodeC, StepType: StepTypeDelay},
		{ID: nodeD, StepType: StepTypeEmail},
	}

	edges := []Edge{
		{SourceStepID: nodeA, TargetStepID: nodeB},
		{SourceStepID: nodeA, TargetStepID: nodeC},
		{SourceStepID: nodeB, TargetStepID: nodeD},
		{SourceStepID: nodeC, TargetStepID: nodeD},
	}

	val := ValidateDAG(nodes, edges)
	if !val.IsValid {
		t.Fatalf("DAG validation failed unexpectedly: %v", val.Errors)
	}

	mock := &mockExecutor{}
	dagEngine := NewDAGEngine(mock)

	events := make([]ExecutionEvent, 0)
	results, err := dagEngine.ExecuteGraph(context.Background(), uuid.New(), nodes, edges, []byte(`{}`), func(e ExecutionEvent) {
		events = append(events, e)
	})

	if err != nil {
		t.Fatalf("ExecuteGraph failed: %v", err)
	}

	if len(results) != 4 {
		t.Fatalf("expected 4 node results, got %d", len(results))
	}

	for id, res := range results {
		if res.Status != StatusSuccess {
			t.Errorf("node %s did not succeed: %s", id, res.Status)
		}
	}
}

func TestDAGRetryWithBackoff(t *testing.T) {
	nodeA := uuid.New()
	nodes := []Node{
		{
			ID:       nodeA,
			StepType: StepTypeHTTPRequest,
			RetryConfig: RetryConfig{
				Enabled:         true,
				MaxAttempts:     3,
				Delay:           10 * time.Millisecond,
				BackoffExponent: 2.0,
			},
		},
	}

	mock := &mockExecutor{failUntil: 2} // fails 2 times, succeeds on 3rd attempt
	dagEngine := NewDAGEngine(mock)

	results, err := dagEngine.ExecuteGraph(context.Background(), uuid.New(), nodes, nil, nil, func(e ExecutionEvent) {})
	if err != nil {
		t.Fatalf("expected retry to succeed on 3rd attempt, got error: %v", err)
	}

	if results[nodeA].Attempts != 3 {
		t.Fatalf("expected 3 attempts, got %d", results[nodeA].Attempts)
	}
}
