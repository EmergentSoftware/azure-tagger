package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

// HTTP handler for SendGrid events
func handleSendGridEvents(w http.ResponseWriter, r *http.Request) {
	var events []Event
	// Read the request body
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		ErrorLogger.Printf("Failed to read request body: %v", err)
		http.Error(w, "Failed to read request", http.StatusInternalServerError)
		return
	}

	// Log the raw request body
	InfoLogger.Printf("Raw request body: %s", string(bodyBytes))

	// Reconstruct the request body for decoding since io.ReadAll consumes the body
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	if err := json.NewDecoder(r.Body).Decode(&events); err != nil {
		ErrorLogger.Printf("Bad request: %v", r.Body)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Check if the request is a validation event
	// This is the verification layer during the installation of EventGrid event
	if isValidationEvent := EventGridValidation(w, &bodyBytes, events); isValidationEvent {
		return
	}

	for _, event := range events {
		if event.EventType == "Microsoft.Resources.ResourceWriteSuccess" &&
			event.Data.Status == "Succeeded" &&
			event.Data.HttpRequest.Method == "PUT" { // Only process creation events
			if contains(EventTaggingExcludedActions, event.Data.Authorization.Action) {
				InfoLogger.Printf("filtered event: %v", event)
				continue
			}
			InfoLogger.Printf("event: %v", event)
			if err := tagResource(event); err != nil {
				ErrorLogger.Printf("failed to tag resource: %v", err)
			}
		}
	}

	w.WriteHeader(http.StatusOK)
}
