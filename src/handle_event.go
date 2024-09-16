package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

// handleSendGridEventsProxy is a closure proxy function for handleSendGridEvents
func handleSendGridEventsProxy(logger *Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleSendGridEvents(w, r, logger)
	}
}

// handleSendGridEvents is a HTTP handler for SendGrid events
func handleSendGridEvents(w http.ResponseWriter, r *http.Request, logger *Logger) {
	var events []Event
	// Read the request body
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Error.Printf("Failed to read request body: %v", err)
		http.Error(w, "Failed to read request", http.StatusInternalServerError)
		return
	}

	// Log the raw request body
	logger.Info.Printf("Raw request body: %s", string(bodyBytes))

	// Reconstruct the request body for decoding since io.ReadAll consumes the body
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	if err := json.NewDecoder(r.Body).Decode(&events); err != nil {
		logger.Error.Printf("Bad request: %v", r.Body)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Check if the request is a validation event
	// This is the verification layer during the installation of EventGrid event
	if isValidationEvent := EventGridValidation(w, &bodyBytes, events, logger); isValidationEvent {
		return
	}

	for _, event := range events {
		if event.EventType == "Microsoft.Resources.ResourceWriteSuccess" &&
			event.Data.Status == "Succeeded" &&
			event.Data.HttpRequest.Method == "PUT" { // Only process creation events
			if contains(EventTaggingExcludedActions, event.Data.Authorization.Action) {
				logger.Info.Println("filtered event:", event)
				continue
			}
			logger.Info.Println("event:", event)
			if err := tagResource(event, logger); err != nil {
				logger.Error.Printf("failed to tag resource: %v", err)
			}
		}
	}

	w.WriteHeader(http.StatusOK)
}
