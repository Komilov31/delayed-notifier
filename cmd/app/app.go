package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/Komilov31/delayed-notifier/docs"
	"github.com/Komilov31/delayed-notifier/internal/cache/redis"
	"github.com/Komilov31/delayed-notifier/internal/config"
	"github.com/Komilov31/delayed-notifier/internal/handler"
	"github.com/Komilov31/delayed-notifier/internal/rabbitmq"
	"github.com/Komilov31/delayed-notifier/internal/repository"
	"github.com/Komilov31/delayed-notifier/internal/sender"
	"github.com/Komilov31/delayed-notifier/internal/service"

	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

func Run() error {
	zlog.Init()

	dbString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Cfg.Postgres.Host,
		config.Cfg.Postgres.Port,
		config.Cfg.Postgres.User,
		config.Cfg.Postgres.Password,
		config.Cfg.Postgres.Name,
	)
	opts := &dbpg.Options{MaxOpenConns: 10, MaxIdleConns: 5}
	db, err := dbpg.New(dbString, []string{}, opts)
	if err != nil {
		log.Fatal("could not init db: " + err.Error())
	}
	repository := repository.New(db)
	cache := redis.New()
	queue := rabbitmq.New()
	sender := sender.New()
	service := service.New(repository, cache, queue, sender)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		zlog.Logger.Info().Msgf("recieved shutting signal %v. Shuting down", sig)
		cancel()
	}()

	go func() {
		if err := service.PublishReadyNotifications(ctx); err != nil {
			log.Fatal("error while publishing notifications: ", err)
		}
	}()

	go func() {
		if err := service.ConsumeMessages(ctx); err != nil {
			log.Fatal("could not start consumer: ", err)
		}
	}()

	handler := handler.New(service)

	router := ginext.New()
	registerRoutes(router, handler)

	zlog.Logger.Info().Msg("succesfully started server on " + config.Cfg.HttpServer.Address)
	return router.Run(config.Cfg.HttpServer.Address)
}

func registerRoutes(engine *ginext.Engine, handler *handler.Handler) {
	// Register static files
	engine.LoadHTMLFiles("/app/static/index.html")
	engine.Static("/static", "/app/static")

	// POST requests
	engine.POST("/notify", handler.CreateNotification)

	// GET requests
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	engine.GET("/", handler.GetMainPage)
	engine.GET("/notify/:id", handler.GetNotificationStatus)
	engine.GET("/notify", handler.GetAllNotifications)

	// DELETE request
	engine.DELETE("notify/:id", handler.UpdateNotificationStatus)

}
