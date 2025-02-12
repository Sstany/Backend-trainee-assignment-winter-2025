package app

import (
	"database/sql"

	_ "github.com/lib/pq"
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

	logger.Info("start migration")

	if err := repo.Migrate(db, cfg.Migrations); err != nil {
		panic(err)
	}

	logger.Info("start api", zap.String("apiVersion", "v1"))

	authRepo := repo.NewAuth(db, logger.Named("auth-repo"))
	passHasher := password.NewHasherBcrypt(logger.Named("pass-hasher"))
	userUsecase := usecase.NewUser(authRepo, passHasher)

	handler.NewServer(logger.Named("http"), userUsecase, cfg.Address).Start()
}
