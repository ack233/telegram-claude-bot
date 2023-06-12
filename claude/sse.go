package claude

import (
	"bufio"
	"fmt"
	"tebot/pkgs/logtool"
	"tebot/pkgs/requests"

	"github.com/google/uuid"
	"github.com/launchdarkly/eventsource"
)

func InitSse(url string) Client {
	return Client{
		URL:          url,
		EventChannel: make(chan Result, 50),
	}
}

func (c *Client) Connect(prompt, conversationID string) {
	payload := Payload{
		Action: "next",
		Messages: []Message{{
			ID:     uuid.New().String(),
			Role:   "user",
			Author: Author{Role: "user"},
			Content: Content{
				ContentType: "text",
				Parts:       []string{prompt},
			},
		}},
		ConversationID:  conversationID,
		ParentMessageID: uuid.New().String(),
		Model:           "claude-unknown-version",
	}

	resp, err := requests.Client.R().
		Notparse().
		SetBody(payload).
		SetHeaders(c.Headers).
		Post(c.URL)

	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	// EventSource Decoder
	go func() {
		defer resp.RawBody().Close()
		defer close(c.EventChannel)
		reader := bufio.NewReader(resp.RawBody())
		decoder := eventsource.NewDecoder(reader)

		for {
			ev, err := decoder.Decode()
			if err != nil {
				logtool.SugLog.Info(fmt.Sprintf("Failed to decode event: %v", err))
				break
			}
			if ev.Data() == "" {
				continue
			}
			if ev.Data() == "[DONE]" {
				break
			}
			c.EventChannel <- Result{Data: ev.Data(), Err: err}
			logtool.SugLog.Debug("Received http event")

		}
	}()

}
