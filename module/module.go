package module

type SendMsgMode struct {
	Completion struct {
		Prompt   string `json:"prompt"`
		Timezone string `json:"timezone"`
		Model    string `json:"model"`
	} `json:"completion"`
	OrganizationUUID string   `json:"organization_uuid"`
	ConversationUUID string   `json:"conversation_uuid"`
	Text             string   `json:"text"`
	Attachments      []string `json:"attachments"`
}
type SendUuidMode struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

type UuidResponse struct {
	UUID      string `json:"uuid"`
	Name      string `json:"name"`
	Summary   string `json:"summary"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
type MessageLimit struct {
	Type string `json:"type"`
}

type MsgResponse struct {
	Completion   string       `json:"completion"`
	StopReason   *string      `json:"stop_reason"`
	Model        string       `json:"model"`
	Truncated    bool         `json:"truncated"`
	Stop         *string      `json:"stop"`
	LogID        string       `json:"log_id"`
	Exception    interface{}  `json:"exception"`
	MessageLimit MessageLimit `json:"messageLimit"`
}

type Choice struct {
	Index        int               `json:"index"`
	Delta        map[string]string `json:"delta"`
	FinishReason interface{}       `json:"finish_reason"`
}

type JsonData struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created string   `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
}
type UuidBody struct {
	UUID string `json:"uuid"`
}

type OpenAIRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIResponse struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Organization struct {
	Uuid         string            `json:"uuid"`
	Name         string            `json:"name"`
	JoinToken    string            `json:"join_token"`
	CreatedAt    string            `json:"created_at"`
	UpdatedAt    string            `json:"updated_at"`
	Capabilities []string          `json:"capabilities"`
	Settings     map[string]string `json:"settings"`
	ActiveFlags  []interface{}     `json:"active_flags"`
}

type AssistantRequest struct {
	Model             string `json:"model"`
	Prompt            string `json:"prompt"`
	MaxTokensToSample int    `json:"max_tokens_to_sample"`
	Stream            bool   `json:"stream"`
}
type ClaudeRes struct {
	Completion string  `json:"completion"`
	StopReason *string `json:"stop_reason"`
	Model      string  `json:"model"`
}
