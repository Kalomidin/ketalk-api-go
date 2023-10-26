package handler

import (
	"context"
	"fmt"
	"ketalk-api/common"
	"ketalk-api/pkg/config"
	auth_handler "ketalk-api/pkg/handler/auth"
	item_handler "ketalk-api/pkg/handler/item"
	user_handler "ketalk-api/pkg/handler/user"
	auth_manager "ketalk-api/pkg/manager/auth"
	auth_repo "ketalk-api/pkg/manager/auth/repository"
	item_manager "ketalk-api/pkg/manager/item"
	item_repo "ketalk-api/pkg/manager/item/repository"
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

func InitHandlers(ctx context.Context, ginEngine *gin.Engine, cfg config.Config) error {
	middleware, err := NewMiddleware(ctx, &cfg.DB)
	if err != nil {
		return err
	}

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

	// run migrations
	if err := runMigrations(db, &cfg.DB, userRepo, authRepo, itemRepo, itemImageRepo, userItemRepo); err != nil {
		return err
	}

	userPort := user_manager.NewUserPort(userRepo)

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

	return nil
}

func runMigrations(db *gorm.DB,
	dbConfig postgres.ConfigPostgres,
	userRepo user_repo.Repository, authRepo auth_repo.Repository,
	itemRepo item_repo.ItemRepository,
	itemImageRepo item_repo.ItemImageRepository,
	userItemRepo item_repo.UserItemRepository,
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
	return nil
}
