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

	"github.com/caarlos0/env/v8"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type Postgres struct {
	User     string `env:"PG_USER,required"`
	Password string `env:"PG_PASSWORD,required"`
	Host     string `env:"PG_HOST,required"`
	Port     int    `env:"PG_PORT,required"`
	Db       string `env:"PG_DB,required"`
}

type Redis struct {
	User     string `env:"REDIS_USER,required"`
	Password string `env:"REDIS_PASSWORD,required"`
	Host     string `env:"REDIS_HOST,required"`
	Port     string `env:"REDIS_PORT,required"`
}

type Options struct {
	Debug    bool   `env:"DEBUG" envDefault:"false"`
	HttPort  string `env:"HTTP_PORT" envDefault:"8080"`
	LogLevel string `env:"LOG_LEVEL" envDefault:"info"`

	Postgres Postgres
	Redis    Redis
}

func main() {

	// Syscall
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-c
		cancel()
	}()

	// config
	var cfg Options
	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}

	// logger
	logrus.SetOutput(os.Stdout)
	logrus.SetFormatter(&logrus.JSONFormatter{})
	lvl, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		panic(err)
	}
	logrus.SetLevel(lvl)

	// postgresql
	connPoolCfg := pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			PreferSimpleProtocol: true,
			Host:                 cfg.Postgres.Host,
			Port:                 uint16(cfg.Postgres.Port),
			Database:             cfg.Postgres.Db,
			User:                 cfg.Redis.User,
			Password:             cfg.Redis.Password,
		},
		AfterConnect:   nil,
		MaxConnections: 40,
		AcquireTimeout: 30 * time.Second,
	}

	connPool, err := pgx.NewConnPool(connPoolCfg)
	if err != nil {
		panic(err)
	}
	defer connPool.Close()

	nativeDB := stdlib.OpenDBFromPool(connPool, stdlib.OptionPreferSimpleProtocol(true))

	db := sqlx.NewDb(nativeDB, "pgx")
	defer db.Close()

	db.SetMaxOpenConns(40)
	db.SetMaxIdleConns(40)
	db.SetConnMaxLifetime(5 * time.Minute)

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	// redis
	redisOpt := &redis.Options{
		WriteTimeout:    time.Duration(15) * time.Second,
		ReadTimeout:     time.Duration(15) * time.Second,
		Addr:            fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		MaxIdleConns:    40,
		ConnMaxIdleTime: 5 * time.Minute,
		Password:        cfg.Redis.Password,
		Username:        cfg.Redis.User,
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
		Addr:         fmt.Sprintf(":%s", cfg.HttPort),
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
