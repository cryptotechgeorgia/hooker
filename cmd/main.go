package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cryptotechgeorgia/hooker/hooker"
	"github.com/cryptotechgeorgia/hooker/sdk/filerotate"
	"github.com/cryptotechgeorgia/hooker/sdk/notifier"
	"github.com/cryptotechgeorgia/hooker/sdk/redis"
	"github.com/spf13/viper"
)

var (
	logFile *filerotate.File
)

type Config struct {
	Notifier struct {
		Users map[string]string
	}
	Redis struct {
		Address  string
		Password string
		DB       int
		Channel  string
	}
	Logger struct {
		Dir string
	}
}

func main() {
	if err := run(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func run() error {
	conf, err := loadConfiguration("config")
	if err != nil {
		log.Fatalf("while reading conf %s\n", err.Error())
		os.Exit(1)
	}

	filenameSuffix := ".log"
	err = openLogFile(conf.Logger.Dir, filenameSuffix, nil)
	if err != nil {
		return err
	}

	log.SetOutput(logFile)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	notifierConfig := notifier.Config{
		Users: conf.Notifier.Users,
	}
	redisConfig := redis.Config{
		Address:        conf.Redis.Address,
		Password:       conf.Redis.Password,
		DB:             conf.Redis.DB,
		DefaultChannel: conf.Redis.Channel,
	}

	notifierClient := notifier.NewNotifierClient(redis.NewClient(redisConfig), notifierConfig)
	http.HandleFunc("/notify", hooker.NotifyHandler(notifierClient))

	server := http.Server{
		Addr: ":8080",
	}

	go func() {
		log.Println("starting server on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("error starting server: %v", err)
		}
	}()

	<-shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("shutting down server: %v", err)
	}

	log.Println("server shut down successfully")
	logFile.Close()

	return nil
}

func loadConfiguration(confName string) (Config, error) {
	viper.SetConfigName(confName)
	viper.AddConfigPath(".")
	viper.SetConfigType("json")

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, err
	}

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, err
	}

	var cfg Config
	err := viper.Unmarshal(&cfg)
	if err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func openLogFile(dir string, fileNameSuffix string, onClose func(string, bool)) error {
	w, err := filerotate.NewDaily(dir, fileNameSuffix, onClose)
	if err != nil {
		return err
	}
	logFile = w
	return nil
}
