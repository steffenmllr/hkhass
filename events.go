package hkhass

import (
	"encoding/json"
	"fmt"
	"github.com/r3labs/sse"
)

type hkevent struct {
	EventType string `json:"event_type"`
	Data      struct {
		EntityID string `json:"entity_id"`
		OldState struct {
			EntityID string `json:"entity_id"`
			State    string `json:"state"`
		} `json:"old_state"`
		NewState struct {
			EntityID string `json:"entity_id"`
			State    string `json:"state"`
		} `json:"new_state"`
	} `json:"data"`
}

func ListenToEvents(client *HAClient, host string, port int, password string) error {
	BaseURL := fmt.Sprintf("http://%s:%d/api/stream", host, port)
	sclient := sse.NewClient(BaseURL)
	sclient.Headers = map[string]string{
		"Content-Type": "application/json",
		"x-ha-access":  password,
	}

	sclient.Subscribe("messages", func(msg *sse.Event) {
		// Got some data!
		if msg.Data != nil && string(msg.Data) != "ping" {
			var message hkevent
			if err := json.Unmarshal(msg.Data, &message); err != nil {
				panic(err)
			}

			if (message.EventType == "state_changed") {
				client.emitter.Publish(message.Data.EntityID, message.Data.NewState.State, message.Data.OldState.State)
			}
		}

	})

	return nil
}
