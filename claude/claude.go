package claude

import (
	"encoding/json"
	"tebot/pkgs/config"
	"tebot/pkgs/initfunc"
	"tebot/pkgs/logtool"
)

var (
	claudeApi       string
	claudeToken     string
	claudechannelID string
)

func init() {
	initfunc.RegisterInitFunc(ClaudeInit)
}

func ClaudeInit() {
	claudeApi = config.Claudeconfig.Api
	claudeToken = config.Claudeconfig.Token
	claudechannelID = config.Claudeconfig.ChannelID
}

func GetMessFromClaude(prompt, conversationID string) chan ClaudeResponseStruct {
	client := InitSse(claudeApi)
	client.Headers = map[string]string{
		"Authorization": "Bearer " + claudechannelID + "@" + claudeToken,
		"Content-Type":  "application/json",
	}

	client.Connect(prompt, conversationID)

	var ClaudeResponse = make(chan ClaudeResponseStruct, 50)
	go func() {

		defer close(ClaudeResponse)
		var r = ClaudeResponseStruct{}
		for event := range client.EventChannel {
			if event.Err != nil {
				r.DecodeErr = event.Err
				ClaudeResponse <- r
				return
			}
			err := json.Unmarshal([]byte(event.Data), &r)
			if err != nil {
				logtool.SugLog.Info("Error while decoding JSON:", err)
				continue
			}
			ClaudeResponse <- r

		}
	}()
	return ClaudeResponse
}
