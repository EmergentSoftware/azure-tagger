package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	msgraphs "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/serviceprincipals"
)

// getServicePrincipalNameFromId retrieves name of service principal from id
func getServicePrincipalNameFromId(cred *azidentity.DefaultAzureCredential, event Event, logger *Logger) (*string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Minute)
	defer cancel()
	var err error
	if event.Data.Authorization.Evidence.PrincipalType != "ServicePrincipal" {
		logger.Error.Printf("Service Principal to App Name function called but type PrincipalType is not ServicePrincipal but: %v", event.Data.Authorization.Evidence.PrincipalType)
		return to.Ptr(""), fmt.Errorf("Service Principal to App Name function called but PrincipalType type is not ServicePrincipal but: %v", event.Data.Authorization.Evidence.PrincipalType)
	}
	// Create Microsoft Graph client
	client, err := msgraphs.NewGraphServiceClientWithCredentials(cred, nil)
	if err != nil {
		logger.Error.Printf("Failed to create graph client: %v", err)
		return to.Ptr(""), fmt.Errorf("Failed to create graph client: %v", err)
	}

	// Filter the Service Principal list by App ID
	req, err := client.ServicePrincipals().Get(ctx, &serviceprincipals.ServicePrincipalsRequestBuilderGetRequestConfiguration{
		QueryParameters: &serviceprincipals.ServicePrincipalsRequestBuilderGetQueryParameters{Filter: to.Ptr(fmt.Sprintf("appId eq '%s'", event.Data.Claims.Appid))},
	})
	if err != nil {
		logger.Error.Printf("Error getting service principal: %v", err)
		return to.Ptr(""), fmt.Errorf("Error getting service principal: %v", err)
	}
	spList := req.GetValue()

	// Check if the result is not empty
	if len(spList) == 0 {
		logger.Warn.Printf("No service principal found for app ID: %s, or incorrect access permission for MI of azure tagger app", event.Data.Claims.Appid)
		return to.Ptr(""), nil
	}

	// Get the first service principal (assuming only one is found for App ID)
	sp := spList[0]
	logger.Info.Printf("App Name: %s\n", *sp.GetDisplayName())
	return sp.GetDisplayName(), err
}
