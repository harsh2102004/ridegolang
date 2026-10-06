CREATE TABLE workflow_edges (
    id UUID PRIMARY KEY,
    workflow_id UUID NOT NULL REFERENCES workflows(id) ON DELETE CASCADE,
    source_step_id UUID NOT NULL,
    target_step_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);