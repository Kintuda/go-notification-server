package http

import (
	"github.com/Kintuda/notification-server/pkg/config"
	"github.com/Kintuda/notification-server/pkg/notification"
	"github.com/Kintuda/notification-server/pkg/postgres"
	"github.com/Kintuda/notification-server/pkg/queue"
	"github.com/gin-gonic/gin"
)

type Router struct {
	Engine              *gin.Engine
	Cfg                 *config.AppConfig
	DB                  *postgres.Pool
	RabbitMQ            *queue.RabbitMQProvider
	NotifierService     *notification.NotifierService
	NotificationService *notification.NotificationService
}

func NewRouter(cfg *config.AppConfig, db *postgres.Pool, rabbit queue.RabbitMQProvider) (*Router, error) {
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	repository := postgres.NewPostgresRepository(db)
	notificationRepository := postgres.NewRepository[postgres.NotificationRepositoryPg](repository)
	notificationService := notification.NewNotificationService(notificationRepository)
	notifierService := notification.NewNotifierService(notificationRepository)

	r := Router{
		Engine:              router,
		Cfg:                 cfg,
		DB:                  db,
		RabbitMQ:            &rabbit,
		NotifierService:     notifierService,
		NotificationService: notificationService,
	}

	return &r, nil
}

func RegisterRoutes(r *Router) {
	r.Engine.Use(ErrorHandler())

	v1 := r.Engine.Group("v1")
	v1.POST("/notifier", r.RegisterNotifier)
	v1.POST("/notification", r.SendNotification)
}
