package main

import (
	"context"
	"fmt"

	"log"
	"net/http"

	swaggerFiles "github.com/swaggo/files"

	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"

	_ "ketalk-api/docs"
	"ketalk-api/pkg/handler"
	conn_redis "ketalk-api/pkg/manager/conversation/redis"
)

func main() {
	// initialize config
	var cfg Config
	if err := cfg.Load(); err != nil {
		log.Printf("failed to read config\n")
	}

	ctx := context.Background()

	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				log.Printf("main panicked because: %v\n", err)
			} else {
				log.Printf("main panicked because: %v\n", r)
			}
		}
	}()
	err := run(ctx, cfg)
	if err != nil {
		log.Printf("server error: %v\n", err)
	}
}

func run(ctx context.Context, cfg Config) error {
	router := gin.Default()

	redis, err := conn_redis.Init(ctx, cfg.Config.Redis)
	if err != nil {
		return err
	}

	middleware, err := handler.NewMiddleware(ctx, &cfg.Config.DB)
	if err != nil {
		return err
	}
	router.Use(middleware.AuthMiddleware(cfg.Config.Auth))
	router.Use(middleware.HttpMiddleware()) // not sure why i added ??
	if err := handler.InitHandlers(ctx, middleware, redis, router, cfg.Config); err != nil {
		return err
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.NewHandler()))

	// run the server
	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: router,
	}
	fmt.Printf("server is running on port %d\n", cfg.Server.Port)

	if err := srv.ListenAndServe(); err != nil {
		return err
	}
	return nil
}
