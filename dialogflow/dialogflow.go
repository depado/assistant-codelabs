package dialogflow

import "time"

// Request is the main dialogflow request
type Request struct {
	Lang            string          `json:"lang"`
	Status          Status          `json:"status"`
	Timestamp       time.Time       `json:"timestamp"`
	SessionID       string          `json:"sessionId"`
	Result          Result          `json:"result"`
	ID              string          `json:"id"`
	OriginalRequest OriginalRequest `json:"originalRequest"`
}

// Status is the status of the request
type Status struct {
	ErrorType string `json:"errorType"`
	Code      int    `json:"code"`
}

// Result is the result of the request
type Result struct {
	Parameters       interface{}   `json:"parameters"`
	Contexts         []interface{} `json:"contexts"`
	ResolvedQuery    string        `json:"resolvedQuery"`
	Source           string        `json:"source"`
	Score            float64       `json:"score"`
	Speech           string        `json:"speech"`
	Fulfillment      Fulfillment   `json:"fulfillment"`
	ActionIncomplete bool          `json:"actionIncomplete"`
	Action           string        `json:"action"`
	Metadata         Metadata      `json:"metadata"`
}

type Fulfillment struct {
	Messages interface{} `json:"messages"`
	Speech   string      `json:"speech"`
}

type Message struct {
	Speech string `json:"speech"`
	Type   int    `json:"type"`
}

type Metadata struct {
	IntentID                  string `json:"intentId"`
	WebhookForSlotFillingUsed string `json:"webhookForSlotFillingUsed"`
	IntentName                string `json:"intentName"`
	WebhookUsed               string `json:"webhookUsed"`
}

// OriginalRequest is the initial request
type OriginalRequest struct {
	Source string `json:"source"`
	Data   Data   `json:"data"`
}

type Data struct {
	Inputs       Inputs      `json:"inputs"`
	User         User        `json:"user"`
	Conversation interface{} `json:"conversation"`
}

// User simply holds the user ID
type User struct {
	UserID string `json:"user_id"`
}

// Conversation holds information about the conversation
type Conversation struct {
	ConversationID    string `json:"conversation_id"`
	Type              int    `json:"type"`
	ConversationToken string `json:"conversation_token"`
}

type RawInput struct {
	Query     string `json:"query"`
	InputType int    `json:"input_type"`
}

type Argument struct {
	TextValue string `json:"text_value"`
	RawText   string `json:"raw_text"`
	Name      string `json:"name"`
}

// Input are the user inputs
type Input struct {
	RawInputs []RawInput `json:"raw_inputs"`
	Intent    string     `json:"intent"`
	Arguments []Argument `json:"arguments"`
}

// Inputs is a slice of Input
type Inputs []Input
