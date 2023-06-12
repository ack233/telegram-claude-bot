package claude

type Payload struct {
	Action          string    `json:"action"`
	Messages        []Message `json:"messages"`
	ConversationID  string    `json:"conversation_id"`
	ParentMessageID string    `json:"parent_message_id"`
	Model           string    `json:"model"`
}

type Message struct {
	ID      string  `json:"id"`
	Role    string  `json:"role"`
	Author  Author  `json:"author"`
	Content Content `json:"content"`
}

type Author struct {
	Role string `json:"role"`
}

type Content struct {
	ContentType string   `json:"content_type"`
	Parts       []string `json:"parts"`
}

type ClaudeResponseStruct struct {
	Message struct {
		ID     string `json:"id"`
		Role   string `json:"role"`
		Author struct {
			Role string `json:"role"`
		} `json:"author"`
		Content struct {
			ContentType string   `json:"content_type"`
			Parts       []string `json:"parts"`
		} `json:"content"`
	} `json:"message"`
	ConversationID string  `json:"conversation_id"`
	Error          *string `json:"error"`
	DecodeErr      error
}

type Result struct {
	Data string
	Err  error
}

type Client struct {
	URL          string
	EventChannel chan Result
	Headers      map[string]string
}
