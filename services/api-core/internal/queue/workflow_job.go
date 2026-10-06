package queue

type WorkflowJob struct {
	WorkflowID string `json:"workflow_id"`
	UserID     string `json:"user_id"`
	RunID      string `json:"run_id"`
	Input      []byte `json:"input"`
	Source     string `json:"source"`
}
