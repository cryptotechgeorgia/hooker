package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/cryptotechgeorgia/sdk/notifier"
)

type HookResponse map[string]interface{}

func (r HookResponse) Action() string {
	return r["action"].(string)
}

func (r HookResponse) TaskDescription() string {
	return r["data"].(map[string]interface{})["description"].(string)
}

func (r HookResponse) CreatedBy() string {
	return r["data"].(map[string]interface{})["owner"].(map[string]interface{})["username"].(string)
}

func (r HookResponse) PermaLink() string {
	return r["data"].(map[string]interface{})["permalink"].(string)
}

func (r HookResponse) AssignedToUserName() string {
	return r["data"].(map[string]interface{})["assigned_to"].(map[string]interface{})["username"].(string)
}

type NotificationData struct {
	Action      string
	Owner       string
	Description string
	PermaLink   string
}

func NotifyHandler(client *notifier.Notifier, destination map[string]string) http.HandlerFunc {
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

		// retrieve username to telegram id map value
		telegramId, ok := destination[resp.AssignedToUserName()]
		if !ok {
			log.Printf("username %s  does not match internal config %+v\n", resp.AssignedToUserName(), destination)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(500)
			w.Write([]byte("error"))
			return
		}
		client.WithDestination(telegramId)

		notifBytes, err := json.Marshal(notif)
		if err != nil {
			log.Printf("request marshall error  %s\n", resp.AssignedToUserName())
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(500)
			w.Write([]byte("error"))
			return
		}

		client.Notify(r.Context(), notifBytes, 0)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("OK"))
	}
}
