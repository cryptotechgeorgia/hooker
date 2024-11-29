package notifier

import (
	"context"
	"encoding/json"
	"log"

	"github.com/cryptotechgeorgia/hooker/sdk/redis"

	"github.com/google/uuid"
)

const (
	Source                  = "hooker"
	DefaultDestinationType  = "telegram"
	DefaultTemplateName     = "taiga-task-assignment"
	DefaultTemplateLanguage = "en"
)

type Config struct {
	Users map[string]string
}

func (c Config) Destination(userName string) string {
	return c.Users[userName]
}

type NotifyRequest struct {
	Template        string `json:"template"`
	Language        string `json:"language"`
	Payload         string `json:"payload"`
	Subject         string `json:"subject"`
	Source          string `json:"source"`
	DestinationType string `json:"destinationType"`
	Destination     string `json:"destination"`
	Expire          struct {
		Key string  `json:"key"`
		Ttl float64 `json:"ttl"`
	} `json:"expire"`
}

type NotifyErr struct {
	Error string
}

type NotifierClient struct {
	redisClient *redis.Client
	opts        Config
}

func NewNotifierClient(rClient *redis.Client, opts Config) *NotifierClient {
	return &NotifierClient{
		redisClient: rClient,
		opts:        opts,
	}
}

func (n *NotifierClient) Notify(ctx context.Context, message interface{}, userName string) {
	msgBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("notifier: Notify marshall payload %s\n", err.Error())
		return
	}

	log.Printf("msg bytes are %s\n", msgBytes)
	log.Printf("destination is %s \n", n.opts.Users[userName])

	req := NotifyRequest{
		Template:        DefaultTemplateName,
		Language:        DefaultTemplateLanguage,
		Payload:         string(msgBytes),
		Source:          Source,
		DestinationType: DefaultDestinationType,
		Destination:     n.opts.Destination(userName),
		Expire: struct {
			Key string  "json:\"key\""
			Ttl float64 "json:\"ttl\""
		}{Key: uuid.NewString(), Ttl: 0},
	}

	b, err := json.Marshal(req)
	if err != nil {
		log.Printf("notifier: Notify marshall request %s\n", err.Error())
		return
	}
	_, err = n.redisClient.Publish(ctx, string(b)).Result()
	if err != nil {
		log.Printf("notifier:  Notify  redis publish %s\n", err.Error())
	}
}
