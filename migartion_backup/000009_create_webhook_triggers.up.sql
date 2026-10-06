CREATE TABLE webhook_triggers (
    id UUID PRIMARY KEY,
    workflow_id UUID NOT NULL REFERENCES workflows(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    webhook_url_id TEXT NOT NULL UNIQUE,
    frontend_node_id TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);