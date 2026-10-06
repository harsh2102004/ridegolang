package scheduler

import (
	"context"
	"log"
	"org/api-core/internal/workflow/service"
	"time"
)

type Scheduler struct {
	workflowService *service.WorkflowService
}

func NewScheduler(workflowService *service.WorkflowService) *Scheduler {
	return &Scheduler{
		workflowService: workflowService,
	}
}

func (s *Scheduler) Start(ctx context.Context) {
	log.Println("scheduler started")

	ticker := time.NewTicker(1 * time.Minute)

	go func() {
		for {
			select {
			case <-ticker.C:
				s.run(ctx)

			case <-ctx.Done():
				log.Println("scheduler stopped")
				return
			}
		}
	}()
}

func (s *Scheduler) run(ctx context.Context) {
	workflows, err := s.workflowService.ListDueScheduledWorkflows(ctx)
	if err != nil {
		log.Println("scheduler list error:", err)
		return
	}

	for _, workflow := range workflows {
		wf := workflow

		log.Println("running scheduled workflow:", wf.ID)

		go func() {
			err := s.workflowService.RunScheduledWorkflow(ctx, wf)

			if err != nil {
				log.Println("scheduled workflow failed:", err)
			}
		}()
	}
}
