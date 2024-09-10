package main

import (
	"encoding/json"
	"net/http"
)

// Define the structure for Event Grid validation request
type EventGridValidationEvent struct {
	ValidationCode string `json:"validationCode"`
	EventType      string `json:"eventType"`
}

type ValidationResponse struct {
	ValidationResponse string `json:"validationResponse"`
}

// EventGridValidation checks if the request is a validation event
// This is the verification layer during the installation of EventGrid event
func EventGridValidation(w http.ResponseWriter, bodyBytes *[]byte, events []Event) bool {
	if len(events) > 0 && events[0].EventType == "Microsoft.EventGrid.SubscriptionValidationEvent" {
		// Pick the validation code from the request
		validationCode := events[0].Data.ValidationCode
		response := ValidationResponse{
			ValidationResponse: validationCode,
		}
		InfoLogger.Printf("Microsoft.EventGrid.SubscriptionValidationEvent received for validation code: %s", string(*bodyBytes))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		// Respond with the validation code
		// There's a potential pitfall when we cannot encode a response, but we must set the status before writing the response
		json.NewEncoder(w).Encode(response)
		InfoLogger.Printf("Validated Event Grid subscription with code: %s", validationCode)
		return true
	}
	return false
}
