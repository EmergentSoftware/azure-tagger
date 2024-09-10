package main

import (
	"log"
	"net/http"
	"os"
)

// Custom logger with levels
var (
	InfoLogger  = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarnLogger  = log.New(os.Stdout, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
)

func main() {
	// FUNCTIONS_CUSTOMHANDLER_PORT contains dynamic port passed by Azure Functions Runtime for each runninng instance of Function App
	customHandlerPort, exists := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT")
	if !exists {
		customHandlerPort = "8080"
	}
	http.HandleFunc("/api/SendGridEvents", handleSendGridEvents)
	InfoLogger.Println("Go server Listening on: ", customHandlerPort)
	ErrorLogger.Fatal(http.ListenAndServe(":"+customHandlerPort, nil))
}
