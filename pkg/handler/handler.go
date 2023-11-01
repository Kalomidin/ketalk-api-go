package handler

import (
	"context"
	"fmt"
	"ketalk-api/common"
	"ketalk-api/pkg/config"
	auth_handler "ketalk-api/pkg/handler/auth"
	conversation_handler "ketalk-api/pkg/handler/conversation"
	item_handler "ketalk-api/pkg/handler/item"
	user_handler "ketalk-api/pkg/handler/user"
	auth_manager "ketalk-api/pkg/manager/auth"
	auth_repo "ketalk-api/pkg/manager/auth/repository"
	conversation_manager "ketalk-api/pkg/manager/conversation"
	conn_redis "ketalk-api/pkg/manager/conversation/redis"
	con_repo "ketalk-api/pkg/manager/conversation/repository"
	conversation_repo "ketalk-api/pkg/manager/conversation/repository"
	"ketalk-api/pkg/manager/conversation/ws"
	item_manager "ketalk-api/pkg/manager/item"
	item_repo "ketalk-api/pkg/manager/item/repository"
	"ketalk-api/pkg/manager/port"
	"log"

	"ketalk-api/pkg/manager/middleware"

	user_manager "ketalk-api/pkg/manager/user"
	user_repo "ketalk-api/pkg/manager/user/repository"
	"ketalk-api/storage"

	"ketalk-api/pkg/provider"
	"ketalk-api/pkg/provider/google"
	"ketalk-api/postgres"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewMiddleware(ctx context.Context, dbConfig postgres.ConfigPostgres) (common.Middleware, error) {
	db, err := postgres.InitDB(ctx, dbConfig)
	if err != nil {
		return nil, err
	}
	userRepo := user_repo.NewRepository(ctx, db)
	userPort := user_manager.NewUserPort(userRepo)
	return middleware.NewMiddleware(userPort), nil
}

func InitHandlers(
	ctx context.Context,
	middleware common.Middleware,
	redis conn_redis.RedisClient,
	ginEngine *gin.Engine,
	cfg config.Config) error {

	db, err := postgres.InitDB(ctx, &cfg.DB)
	if err != nil {
		return err
	}
	blobStorage := storage.NewAzureBlobStorage(cfg.AzureBlobStorage)

	userRepo := user_repo.NewRepository(ctx, db)
	authRepo := auth_repo.NewRepository(ctx, db)
	itemRepo := item_repo.NewItemRepository(ctx, db)
	itemImageRepo := item_repo.NewItemImageRepository(ctx, db)
	userItemRepo := item_repo.NewUserItemRepository(db)

	conversationRepo := conversation_repo.NewConversationRepository(db)
	memberRepo := conversation_repo.NewMemberRepository(db)
	messageRepo := conversation_repo.NewMessageRepository(db)

	// run migrations
	if err := runMigrations(db, &cfg.DB,
		userRepo,
		authRepo,
		itemRepo,
		itemImageRepo,
		userItemRepo,
		conversationRepo,
		memberRepo,
		messageRepo,
	); err != nil {
		return err
	}

	userPort := user_manager.NewUserPort(userRepo)
	itemPort := item_manager.NewItemPort(itemRepo, itemImageRepo)

	googleClient := google.NewGoogleClient(cfg.Google)
	providerClient := provider.NewProviderClient(googleClient)

	authManager := auth_manager.NewAuthManager(authRepo, userPort, providerClient, cfg.Auth)
	authHandler := auth_handler.NewHandler(authManager)

	userManager := user_manager.NewUserManager(userRepo, blobStorage)
	userHandler := user_handler.NewHandler(userManager)

	itemManager := item_manager.NewItemManager(itemRepo, itemImageRepo, userItemRepo, userPort, blobStorage)
	itemHandler := item_handler.NewHandler(itemManager)

	authHttpHandler := auth_handler.NewHttpHandler(ctx, authHandler, middleware)
	authHttpHandler.Init(ctx, ginEngine)

	userHttpHandler := user_handler.NewHttpHandler(ctx, userHandler, middleware)
	userHttpHandler.Init(ctx, ginEngine)

	itemHttpHandler := item_handler.NewHttpHandler(ctx, itemHandler, middleware)
	itemHttpHandler.Init(ctx, ginEngine)

	conversationManager := conversation_manager.NewConversationManager(ctx, conversationRepo, memberRepo, messageRepo, itemPort, blobStorage, userPort, redis)
	conversationHandler := conversation_handler.NewHandler(conversationManager)

	conversationHttpHandler := conversation_handler.NewHttpHandler(conversationHandler, middleware)
	conversationHttpHandler.Init(ctx, ginEngine)

	// run web socket server properly instead of inside the handler
	go initWebSocketServer(ctx, userPort, middleware, redis, db, cfg)

	return nil
}

func runMigrations(db *gorm.DB,
	dbConfig postgres.ConfigPostgres,
	userRepo user_repo.Repository, authRepo auth_repo.Repository,
	itemRepo item_repo.ItemRepository,
	itemImageRepo item_repo.ItemImageRepository,
	userItemRepo item_repo.UserItemRepository,
	conversationRepo conversation_repo.ConversationRepository,
	memberRepo conversation_repo.MemberRepository,
	messageRepo conversation_repo.MessageRepository,
) error {
	if resp := db.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", dbConfig.GetSchema())); resp.Error != nil {
		return resp.Error
	}

	err := userRepo.MigrateUser()
	if err != nil {
		return err
	}
	err = authRepo.Migrate()
	if err != nil {
		return err
	}
	err = itemRepo.Migrate()
	if err != nil {
		return err
	}
	err = itemImageRepo.Migrate()
	if err != nil {
		return err
	}

	err = userItemRepo.Migrate()
	if err != nil {
		return err
	}

	err = conversationRepo.Migrate()
	if err != nil {
		return err
	}
	err = memberRepo.Migrate()
	if err != nil {
		return err
	}
	err = messageRepo.Migrate()
	if err != nil {
		return err
	}

	return nil
}

func initWebSocketServer(
	ctx context.Context,
	userPort port.UserPort,
	middleware common.Middleware,
	redis conn_redis.RedisClient,
	db *gorm.DB,
	cfg config.Config,
) error {
	messageRepo := con_repo.NewMessageRepository(db)
	memberRepo := con_repo.NewMemberRepository(db)

	// run web socket server
	wsServer, err := ws.NewWebSocketServer(
		ctx,
		userPort,
		messageRepo,
		memberRepo,
		middleware,
		cfg.Auth,
		redis,
	)
	if err != nil {
		return err
	}

	if err := wsServer.Serve(cfg.WebSocketServer.Port); err != nil {
		log.Printf("failed to start web socket server: %v\n", err)
	}
	return nil
}
