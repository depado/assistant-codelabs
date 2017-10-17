package dialogflow

// Entities is a slice of Entity
type Entities []Entity

// Entity is the representation of an entity
type Entity struct {
	Value    string   `json:"value"`
	Synonyms []string `json:"synonyms"`
}
