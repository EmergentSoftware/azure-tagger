package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	// Initialize thread-safe loggers
	var logger *Logger = &Logger{
		Info:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		Warn:  log.New(os.Stdout, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile),
		Error: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
	// FUNCTIONS_CUSTOMHANDLER_PORT contains dynamic port passed by Azure Functions Runtime for each runninng instance of Function App
	customHandlerPort, exists := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT")
	if !exists {
		customHandlerPort = "8080"
	}
	http.HandleFunc("/api/SendGridEvents", handleSendGridEventsProxy(logger))
	logger.Info.Println("Go server Listening on: ", customHandlerPort)
	logger.Error.Fatal(http.ListenAndServe(":"+customHandlerPort, nil))
}
