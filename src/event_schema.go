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
			Action   string `json:"action"`
			Evidence struct {
				PrincipalType string `json:"principalType"` // [ ServicePrincipal, User ]
			}
		} `jsons:"authorization"`
		Claims struct {
			Appid            string `json:"appid"`
			Appidacr         string `json:"appidacr"` // 1 - App Registration; 0 or 2 - Enterprise Applications
			Idtyp            string `json:"idtyp"`    // [ app, user ]
			Name             string `json:"name"`
			ClaimsName       string `json:"http://schemas.xmlsoap.org/ws/2005/05/identity/claims/name"` // email if idtyp is user, empty otherwise
			Email            string `json:"http://schemas.xmlsoap.org/ws/2005/05/identity/claims/emailaddress"`
			Objectidentifier string `json:"http://schemas.xmlsoap.org/ws/2005/05/identity/claims/objectidentifier"`
		} `jsons:"claims"`
		HttpRequest struct {
			Method string `json:"method"`
		} `json:"httpRequest"`
		ValidationCode string `json:"validationCode"`
		Status         string `json:"status"`
		SubscriptionId string `json:"subscriptionId"`
		Tenantid       string `json:"tenantId"`
	} `jsons:"data"`
	EventType string `json:"eventType"`
	EventTime string `json:"eventTime"`
}
