package main

type Event struct {
	Email      string `json:"email"`
	Event      string `json:"event"`
	Timestamp  int64  `json:"timestamp"`
	ResourceID string `json:"resource_id"`
	Subject    string `json:"subject"`
	Owner      string `json:"owner"`
	//
	Data struct {
		Authorization struct {
			Action string `json:"action"`
		} `jsons:"authorization"`
		Claims struct {
			Name  string `json:"name"`
			Email string `json:"http://schemas.xmlsoap.org/ws/2005/05/identity/claims/emailaddress"`
		} `jsons:"claims"`
		HttpRequest struct {
			Method string `json:"method"`
		} `json:"httpRequest"`
		ValidationCode string `json:"validationCode"`
		Status         string `json:"status"`
	} `jsons:"data"`
	EventType string `json:"eventType"`
}
