package aggregator

import "encoding/json"

// ComponentResponse stores details of an individual component
type ComponentResponse struct {
	ID        string `json:"id"`
	Status    int    `json:"status"`
	Body      string `json:"body"`
	Summary   string `json:"summary"`
	Mandatory bool   `json:"mandatory"`
}

type result struct {
	Summary    string              `json:"summary"`
	Components []ComponentResponse `json:"components"`
}

func finalSummary(components []ComponentResponse) string {
	for _, c := range components {
		if c.Mandatory == true && c.Summary == "failure" {
			return "failure"
		}
	}

	return "success"
}

// Process function
func Process(cr []ComponentResponse) ([]byte, error) {
	j, err := json.Marshal(result{finalSummary(cr), cr})
	if err != nil {
		return []byte{}, err
	}

	return j, nil
}
