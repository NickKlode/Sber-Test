package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sber-test"
	"sber-test/internal/handler"
	"sber-test/internal/repository"
	"sber-test/internal/service"
	"syscall"

	_ "github.com/lib/pq"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// @title Sber Test
// @version 1.0
// @description Api Server Todo List

// @host localhost:8080
// @BasePath /

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("error init config. %s", err.Error())
	}
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error reading .env file. %s", err.Error())
	}
	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		log.Fatalf("filed to init db. %s", err.Error())
	}
	repo := repository.NewRepository(db)
	service := service.NewService(repo)
	handler := handler.NewHandler(service)

	srv := new(sber.Server)
	go func() {
		if err := srv.Run(viper.GetString("port"), handler.InitRoutes()); err != nil {
			log.Fatalf("error while running http server. %s", err.Error())
		}
	}()
	log.Println("Server Started")

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Println("Server Shutting Down")

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Printf("error while shutting down. %s", err.Error())
	}

	if err := db.Close(); err != nil {
		log.Printf("error while close db connection. %s", err.Error())
	}

}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
