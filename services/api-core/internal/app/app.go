package app

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	tokenstore "org/api-core/internal/auth/tokenstore"
	"org/api-core/internal/db"
	"org/api-core/internal/logger"
	"org/api-core/internal/queue"
	"org/api-core/internal/router"
	"org/api-core/internal/scheduler"
	"org/api-core/internal/worker"
	workflowRepository "org/api-core/internal/workflow/repository"
	workflowServicePkg "org/api-core/internal/workflow/service"

	geminiRepository "org/api-core/internal/gemini/repository"
	geminiServicePkg "org/api-core/internal/gemini/service"
)

func New() *gin.Engine {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	db.Init()
	tokenstore.InitRedis()

	appLogger := logger.New()

	appLogger.Info("api-core starting", "service", "api-core")

	r := gin.Default()

	router.RegisterRoutes(r, appLogger)

	workflowRepo := workflowRepository.NewSQLCWorkflowRepository(
		db.QueriesInstance,
	)

	workflowQueue := queue.NewRedisQueue(
		tokenstore.RDB,
	)

	geminiRepo := geminiRepository.NewSQLCGeminiRepository(
		db.QueriesInstance,
	)

	geminiService := geminiServicePkg.NewGeminiService(
		geminiRepo,
	)

	workflowService := workflowServicePkg.NewWorkflowService(
		workflowRepo,
		workflowQueue,
		geminiService,
		appLogger,
	)

	workflowWorker := worker.NewWorkflowWorker(
		workflowQueue,
		workflowService,
	)

	go workflowWorker.Start(context.Background())

	schedulerWorker := scheduler.NewScheduler(workflowService)
	schedulerWorker.Start(context.Background())

	return r
}
