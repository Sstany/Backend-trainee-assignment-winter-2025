package app

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq" // Postgres driver.
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"shop/config"
	"shop/internal/adapter/password"
	"shop/internal/adapter/repo"
	"shop/internal/app/usecase"
	handler "shop/internal/controller/http"
)

func Run(cfg *config.Config) {
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.Level.SetLevel(zapcore.Level(cfg.Log))
	loggerConfig.EncoderConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder
	logger, err := loggerConfig.Build(zap.AddStacktrace(zapcore.ErrorLevel), zap.AddCaller())
	if err != nil {
		panic(err)
	}

	logger.Info("start avito shop")

	defer func() {
		if lErr := logger.Sync(); lErr != nil {
			panic(err)
		}
	}()

	logger.Info("open database")

	db, err := sql.Open("postgres", cfg.ConnectionString)
	if err != nil {
		panic(err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	logger.Info("start migration")

	if err = repo.Migrate(db, cfg.Migrations); err != nil {
		panic(err)
	}

	logger.Info("start api", zap.String("apiVersion", "v1"))

	authRepo, err := repo.NewAuth(db, logger.Named("auth-repo"))
	if err != nil {
		panic(err)
	}

	passHasher := password.NewHasherBcrypt(logger.Named("pass-hasher"))
	secretRepo := repo.NewSecret(logger.Named("secret-repo"), cfg.SigningKeyPath, cfg.JWTIssuer)
	shopRepo := repo.NewShop(logger.Named("shop-repo"))
	balanceRepo := repo.NewBalance(db, logger.Named("balance-repo"))
	inventoryRepo := repo.NewInventory(db, logger.Named("inventory-repo"))
	transactionController := repo.NewTransactionSQL(db, logger.Named("transaction-ctrl"))
	userTransactionRepo := repo.NewUserTransaction(db, logger.Named("user-tranaction"))
	userUsecase := usecase.NewUser(
		shopRepo,
		balanceRepo,
		inventoryRepo,
		transactionController,
		userTransactionRepo,
		logger.Named("user"),
	)

	authUsecase, err := usecase.NewAuth(
		authRepo,
		balanceRepo,
		passHasher,
		secretRepo,
		logger.Named("auth"),
	)

	handler.NewServer(logger.Named("http"), userUsecase, authUsecase, cfg.Address).Start()
}
