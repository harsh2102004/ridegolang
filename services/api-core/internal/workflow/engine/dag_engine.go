package engine

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

type StepType string

const (
	StepTypeWebhookTrigger StepType = "webhook_trigger"
	StepTypeHTTPRequest    StepType = "http_request"
	StepTypeCondition      StepType = "condition"
	StepTypeDelay          StepType = "delay"
	StepTypeEmail          StepType = "email"
	StepTypeAI             StepType = "ai"
)

type Node struct {
	ID             uuid.UUID       `json:"id"`
	FrontendNodeID string          `json:"frontend_node_id"`
	StepType       StepType        `json:"step_type"`
	StepOrder      int             `json:"step_order"`
	Config         json.RawMessage `json:"config"`
	RetryConfig    RetryConfig     `json:"retry_config"`
}

type Edge struct {
	ID              uuid.UUID `json:"id"`
	SourceStepID    uuid.UUID `json:"source_step_id"`
	TargetStepID    uuid.UUID `json:"target_step_id"`
	ConditionBranch string    `json:"condition_branch,omitempty"` // "true", "false", or ""
}

type RetryConfig struct {
	Enabled         bool          `json:"enabled"`
	MaxAttempts     int           `json:"max_attempts"`
	Delay           time.Duration `json:"delay"`
	BackoffExponent float64       `json:"backoff_exponent"`
}

type NodeStatus string

const (
	StatusPending   NodeStatus = "pending"
	StatusRunning   NodeStatus = "running"
	StatusSuccess   NodeStatus = "success"
	StatusFailed    NodeStatus = "failed"
	StatusSkipped   NodeStatus = "skipped"
	StatusCancelled NodeStatus = "cancelled"
)

type NodeResult struct {
	NodeID       uuid.UUID  `json:"node_id"`
	Status       NodeStatus `json:"status"`
	Input        []byte     `json:"input,omitempty"`
	Output       []byte     `json:"output,omitempty"`
	ErrorMessage string     `json:"error_message,omitempty"`
	Attempts     int        `json:"attempts"`
	DurationMs   int64      `json:"duration_ms"`
}

type NodeExecutor interface {
	Execute(ctx context.Context, node Node, input []byte) ([]byte, error)
}

type ExecutionEvent struct {
	RunID     uuid.UUID  `json:"run_id"`
	NodeID    uuid.UUID  `json:"node_id"`
	StepType  StepType   `json:"step_type"`
	Status    NodeStatus `json:"status"`
	Message   string     `json:"message"`
	Output    []byte     `json:"output,omitempty"`
	Error     string     `json:"error,omitempty"`
	Timestamp time.Time  `json:"timestamp"`
}

type DAGValidationResult struct {
	IsValid  bool     `json:"is_valid"`
	Errors   []string `json:"errors"`
	StartIDs []uuid.UUID
}

// ValidateDAG verifies graph connectivity and detects cycles using Kahn's algorithm.
func ValidateDAG(nodes []Node, edges []Edge) DAGValidationResult {
	res := DAGValidationResult{IsValid: true}

	if len(nodes) == 0 {
		res.IsValid = false
		res.Errors = append(res.Errors, "workflow must contain at least one node")
		return res
	}

	nodeMap := make(map[uuid.UUID]bool, len(nodes))
	inDegree := make(map[uuid.UUID]int, len(nodes))
	adj := make(map[uuid.UUID][]uuid.UUID, len(nodes))

	for _, n := range nodes {
		nodeMap[n.ID] = true
		inDegree[n.ID] = 0
	}

	for _, e := range edges {
		if !nodeMap[e.SourceStepID] {
			res.IsValid = false
			res.Errors = append(res.Errors, fmt.Sprintf("edge source %s does not exist", e.SourceStepID))
		}
		if !nodeMap[e.TargetStepID] {
			res.IsValid = false
			res.Errors = append(res.Errors, fmt.Sprintf("edge target %s does not exist", e.TargetStepID))
		}
		inDegree[e.TargetStepID]++
		adj[e.SourceStepID] = append(adj[e.SourceStepID], e.TargetStepID)
	}

	queue := make([]uuid.UUID, 0)
	for id, deg := range inDegree {
		if deg == 0 {
			queue = append(queue, id)
			res.StartIDs = append(res.StartIDs, id)
		}
	}

	if len(res.StartIDs) == 0 {
		res.IsValid = false
		res.Errors = append(res.Errors, "cycle detected: no entry point / root node found")
		return res
	}

	visitedCount := 0
	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		visitedCount++

		for _, neighbor := range adj[curr] {
			inDegree[neighbor]--
			if inDegree[neighbor] == 0 {
				queue = append(queue, neighbor)
			}
		}
	}

	if visitedCount != len(nodes) {
		res.IsValid = false
		res.Errors = append(res.Errors, "cycle detected in workflow execution graph")
	}

	return res
}

type DAGEngine struct {
	executor NodeExecutor
}

func NewDAGEngine(executor NodeExecutor) *DAGEngine {
	return &DAGEngine{
		executor: executor,
	}
}

// ExecuteGraph coordinates concurrent branch execution with sync.WaitGroup and channel notification.
func (e *DAGEngine) ExecuteGraph(
	ctx context.Context,
	runID uuid.UUID,
	nodes []Node,
	edges []Edge,
	initialInput []byte,
	eventEmitter func(event ExecutionEvent),
) (map[uuid.UUID]*NodeResult, error) {

	validation := ValidateDAG(nodes, edges)
	if !validation.IsValid {
		return nil, fmt.Errorf("DAG validation failed: %v", validation.Errors)
	}

	nodeMap := make(map[uuid.UUID]Node, len(nodes))
	inDegree := make(map[uuid.UUID]int, len(nodes))
	childEdges := make(map[uuid.UUID][]Edge)

	for _, n := range nodes {
		nodeMap[n.ID] = n
		inDegree[n.ID] = 0
	}

	for _, ed := range edges {
		inDegree[ed.TargetStepID]++
		childEdges[ed.SourceStepID] = append(childEdges[ed.SourceStepID], ed)
	}

	results := make(map[uuid.UUID]*NodeResult)
	nodeOutputs := make(map[uuid.UUID][]byte)
	var mu sync.Mutex

	readyChan := make(chan uuid.UUID, len(nodes)+10)
	for _, startID := range validation.StartIDs {
		nodeOutputs[startID] = initialInput
		readyChan <- startID
	}

	var wg sync.WaitGroup
	var execErr error
	var errOnce sync.Once

	nodesDone := 0
	totalNodes := len(nodes)

	dispatch := func(n Node, in []byte) {
		defer wg.Done()

		start := time.Now()
		eventEmitter(ExecutionEvent{
			RunID:     runID,
			NodeID:    n.ID,
			StepType:  n.StepType,
			Status:    StatusRunning,
			Message:   fmt.Sprintf("Step %s started", n.StepType),
			Timestamp: start,
		})

		output, attempts, err := e.executeWithRetry(ctx, n, in)
		duration := time.Since(start).Milliseconds()

		mu.Lock()
		nodesDone++

		if err != nil {
			results[n.ID] = &NodeResult{
				NodeID:       n.ID,
				Status:       StatusFailed,
				Input:        in,
				ErrorMessage: err.Error(),
				Attempts:     attempts,
				DurationMs:   duration,
			}
			mu.Unlock()

			eventEmitter(ExecutionEvent{
				RunID:     runID,
				NodeID:    n.ID,
				StepType:  n.StepType,
				Status:    StatusFailed,
				Message:   fmt.Sprintf("Step %s failed: %v", n.StepType, err),
				Error:     err.Error(),
				Timestamp: time.Now(),
			})

			errOnce.Do(func() {
				execErr = fmt.Errorf("node %s failed: %w", n.ID, err)
			})
			return
		}

		results[n.ID] = &NodeResult{
			NodeID:     n.ID,
			Status:     StatusSuccess,
			Input:      in,
			Output:     output,
			Attempts:   attempts,
			DurationMs: duration,
		}
		nodeOutputs[n.ID] = output

		var conditionResult bool
		if n.StepType == StepTypeCondition {
			var condOut struct {
				Result bool `json:"result"`
			}
			if jsonErr := json.Unmarshal(output, &condOut); jsonErr == nil {
				conditionResult = condOut.Result
			}
		}

		// Unlock state and decrement in-degrees of children
		for _, outEdge := range childEdges[n.ID] {
			targetID := outEdge.TargetStepID

			if n.StepType == StepTypeCondition {
				if conditionResult && outEdge.ConditionBranch == "false" {
					inDegree[targetID]--
					continue
				}
				if !conditionResult && outEdge.ConditionBranch == "true" {
					inDegree[targetID]--
					continue
				}
				nodeOutputs[targetID] = in
			} else {
				nodeOutputs[targetID] = output
			}

			inDegree[targetID]--
			if inDegree[targetID] == 0 {
				readyChan <- targetID
			}
		}
		mu.Unlock()

		eventEmitter(ExecutionEvent{
			RunID:     runID,
			NodeID:    n.ID,
			StepType:  n.StepType,
			Status:    StatusSuccess,
			Output:    output,
			Message:   fmt.Sprintf("Step %s completed", n.StepType),
			Timestamp: time.Now(),
		})
	}

	for {
		mu.Lock()
		done := (nodesDone >= totalNodes) || (execErr != nil)
		mu.Unlock()

		if done {
			break
		}

		select {
		case <-ctx.Done():
			return results, ctx.Err()

		case nodeID := <-readyChan:
			node := nodeMap[nodeID]

			mu.Lock()
			input := nodeOutputs[nodeID]
			if len(input) == 0 {
				input = node.Config
			}
			mu.Unlock()

			wg.Add(1)
			go dispatch(node, input)
		}
	}

	wg.Wait()
	return results, execErr
}

func (e *DAGEngine) executeWithRetry(ctx context.Context, node Node, input []byte) ([]byte, int, error) {
	maxAttempts := node.RetryConfig.MaxAttempts
	if !node.RetryConfig.Enabled || maxAttempts <= 0 {
		maxAttempts = 1
	}

	baseDelay := node.RetryConfig.Delay
	if baseDelay <= 0 {
		baseDelay = 50 * time.Millisecond
	}

	exponent := node.RetryConfig.BackoffExponent
	if exponent <= 1.0 {
		exponent = 2.0
	}

	var lastErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		select {
		case <-ctx.Done():
			return nil, attempt, ctx.Err()
		default:
		}

		output, err := e.executor.Execute(ctx, node, input)
		if err == nil {
			return output, attempt, nil
		}

		lastErr = err
		if attempt < maxAttempts {
			delay := time.Duration(float64(baseDelay) * float64(attempt))
			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return nil, attempt, ctx.Err()
			}
		}
	}

	return nil, maxAttempts, errors.New(lastErr.Error())
}
