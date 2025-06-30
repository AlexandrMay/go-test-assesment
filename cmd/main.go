package main

import (
	"context"
	"fmt"
	_ "go-test-assesment/docs"
	"go-test-assesment/internal/cat"
	httpCat "go-test-assesment/internal/cat/delivery/http"
	catRepo "go-test-assesment/internal/cat/repository"
	catUsecase "go-test-assesment/internal/cat/usecase"
	httpMission "go-test-assesment/internal/mission/delivery/http"
	missionRepo "go-test-assesment/internal/mission/repository"
	missionUsecase "go-test-assesment/internal/mission/usecase"

	"go-test-assesment/pkg/logger"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func waitForDatabase(dsn string, maxWait time.Duration) (*pgxpool.Pool, error) {
	start := time.Now()
	for {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		pool, err := pgxpool.New(ctx, dsn)
		if err == nil {
			if pingErr := pool.Ping(ctx); pingErr == nil {
				log.Println("Successfully connected to the database.")
				return pool, nil
			}
			pool.Close()
		}

		if time.Since(start) > maxWait {
			return nil, fmt.Errorf("timed out waiting for database to be ready")
		}

		log.Println("Waiting for the database to be ready...")
		time.Sleep(2 * time.Second)
	}
}

func main() {
	dsn := os.Getenv("DATABASE_URL")
	fmt.Println(dsn)

	pool, err := waitForDatabase(dsn, 30*time.Second)
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}
	defer pool.Close()

	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Use(logger.Logger())

	catRepository := catRepo.NewPostgresCatRepository(pool)
	breedValidator := cat.NewCatAPIValidator()
	catUC := catUsecase.NewCatUsecase(catRepository, breedValidator)
	httpCat.NewCatHandler(r, catUC)

	missionRepository := missionRepo.NewMissionPostgres(pool)
	missionUC := missionUsecase.NewMissionUsecase(missionRepository)
	missionHandler := httpMission.NewHandler(missionUC)
	missionHandler.RegisterRoutes(r)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctxShutdown); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
