package main

import (
	"context"
	"fmt"
	"item-service/internal/app"
	handler "item-service/internal/transport/http/v1"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

func main() {
	// Syscall
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-c
		cancel()
	}()

	// logger
	logrus.SetOutput(os.Stdout)
	logrus.SetFormatter(&logrus.JSONFormatter{})
	lvl, err := logrus.ParseLevel("info")
	if err != nil {
		panic(err)
	}
	logrus.SetLevel(lvl)

	// postgresql
	connPoolCfg := pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:     "localhost",
			Database: "test",
			Port:     5432,
			User:     "test",
			Password: "test",
		},
		AfterConnect:   nil,
		MaxConnections: 20,
		AcquireTimeout: 30 * time.Second,
	}

	connPool, err := pgx.NewConnPool(connPoolCfg)
	if err != nil {
		panic(err)
	}
	defer connPool.Close()

	nativeDB := stdlib.OpenDBFromPool(connPool, stdlib.OptionPreferSimpleProtocol(true))

	db := sqlx.NewDb(nativeDB, "pgx")
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(5 * time.Minute)
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	// redis
	redisOpt := &redis.Options{
		WriteTimeout:    time.Duration(15) * time.Second,
		ReadTimeout:     time.Duration(15) * time.Second,
		Addr:            "localhost:6379",
		MaxIdleConns:    25,
		ConnMaxIdleTime: 5 * time.Minute,
	}

	redisClient := redis.NewClient(redisOpt)
	err = redisClient.Ping(ctx).Err()
	if err != nil {
		panic(err)
	}

	// application
	application := app.New(db, redisClient)
	httpHandler := handler.New(application)

	// http server
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)

	mux.Post("/items", httpHandler.CreateItem)
	mux.Get("/items/{id}", httpHandler.GetItemByID)

	httpServer := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		WriteTimeout: time.Duration(15) * time.Second,
		ReadTimeout:  time.Duration(15) * time.Second,
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			fmt.Println(err)
		}
	}()

	logrus.Info("started on port :8080")
	<-ctx.Done()
	logrus.Info("stopped")

	err = httpServer.Shutdown(context.Background())
	if err != nil {
		logrus.Error(err)
	}
}
