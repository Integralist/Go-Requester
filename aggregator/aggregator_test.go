package aggregator

import (
	"encoding/json"
	"testing"
)

var successfulComponent = ComponentResponse{
	ID:        "foo",
	Status:    200,
	Body:      "<p>foo</p>",
	Summary:   "success",
	Mandatory: true,
}

var failureComponent = ComponentResponse{
	ID:        "bar",
	Status:    404,
	Body:      "<p>bar</p>",
	Summary:   "failure",
	Mandatory: true,
}

func TestSuccessSummary(t *testing.T) {
	collection := []ComponentResponse{
		successfulComponent,
		successfulComponent,
	}

	response, _ := Process(collection)

	var jsonResponse result

	err := json.Unmarshal(response, &jsonResponse)
	if err != nil {
		t.Errorf("There was an unexpected error: %s", err.Error())
	}

	verify("success", jsonResponse.Summary, t)
}

func TestFailureSummary(t *testing.T) {
	collection := []ComponentResponse{
		successfulComponent,
		failureComponent,
	}

	response, _ := Process(collection)

	var jsonResponse result

	err := json.Unmarshal(response, &jsonResponse)
	if err != nil {
		t.Errorf("There was an unexpected error: %s", err.Error())
	}

	verify("failure", jsonResponse.Summary, t)
}

func verify(response, expectation string, t *testing.T) {
	if response != expectation {
		t.Errorf("The response:\n '%s'\ndidn't match the expectation:\n '%s'", response, expectation)
	}
}
