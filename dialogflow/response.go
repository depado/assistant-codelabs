package dialogflow

// Response the structure of the expected response for dialogflow
type Response struct {
	Speech      string        `json:"speech"`
	DisplayText string        `json:"displayText"`
	Data        interface{}   `json:"data"`
	ContextOut  []interface{} `json:"contextOut"`
	Source      string        `json:"source"`
}
