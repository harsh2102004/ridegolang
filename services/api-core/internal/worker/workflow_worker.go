package worker

import (
	"context"
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"org/api-core/internal/queue"
	workflowServicePkg "org/api-core/internal/workflow/service"
)

type WorkflowWorker struct {
	queue   *queue.RedisQueue
	service *workflowServicePkg.WorkflowService
}

func NewWorkflowWorker(
	queue *queue.RedisQueue,
	service *workflowServicePkg.WorkflowService,
) *WorkflowWorker {
	return &WorkflowWorker{
		queue:   queue,
		service: service,
	}
}

func (w *WorkflowWorker) Start(ctx context.Context) {
	log.Println("workflow worker started")

	for {
		select {
		case <-ctx.Done():
			log.Println("workflow worker stopped")
			return
		default:
		}

		payload, err := w.queue.Pop(ctx)
		if err != nil {
			if ctx.Err() != nil {
				log.Println("workflow worker stopped")
				return
			}

			if err == redis.Nil {
				continue
			}

			log.Println("failed to pop workflow job:", err)
			continue
		}

		var job queue.WorkflowJob
		if err := json.Unmarshal(payload, &job); err != nil {
			log.Println("failed to unmarshal workflow job:", err)
			continue
		}

		workflowID, err := uuid.Parse(job.WorkflowID)
		if err != nil {
			log.Println("invalid workflow id:", err)
			continue
		}

		userID, err := uuid.Parse(job.UserID)
		if err != nil {
			log.Println("invalid user id:", err)
			continue
		}

		runID, err := uuid.Parse(job.RunID)
		if err != nil {
			log.Println("invalid run id:", err)
			continue
		}

		log.Println("executing queued workflow:", job.WorkflowID)

		_, err = w.service.ExecuteWorkflowRun(
			ctx,
			workflowID,
			userID,
			runID,
			job.Input,
		)

		if err != nil {
			log.Println("workflow execution failed:", err)
		}
	}
}