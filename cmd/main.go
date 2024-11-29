package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/cryptotechgeorgia/hooker/sdk/notifier"
	"github.com/cryptotechgeorgia/hooker/sdk/redis"
	"github.com/spf13/viper"
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

func NotifyHandler(client *notifier.NotifierClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var resp HookResponse
		if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
			w.Write([]byte(err.Error()))
			log.Printf("error while marshalling : %s\n", err.Error())
			return
		}

		notif := NotificationData{
			Action:      resp.Action(),
			Owner:       resp.CreatedBy(),
			Description: resp.TaskDescription(),
			PermaLink:   resp.PermaLink(),
		}
		client.Notify(r.Context(), notif, resp.AssignedToUserName())
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("ok"))
	}
}

func main() {
	log.SetOutput(os.Stdout)

	conf, err := loadConfiguration("config")
	if err != nil {
		log.Fatalf("while reading conf %s\n", err.Error())
		os.Exit(1)
	}

	notifierConfig := notifier.Config{
		Users: conf.Notifier.Users,
	}
	redisConfig := redis.Config{
		Address:        conf.Redis.Address,
		Password:       conf.Redis.Password,
		DB:             conf.Redis.DB,
		DefaultChannel: conf.Redis.Channel,
	}

	client := notifier.NewNotifierClient(redis.NewClient(redisConfig), notifierConfig)

	http.HandleFunc("/notify", NotifyHandler(client))
	fmt.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
