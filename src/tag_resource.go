package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
)

// Function to tag the resource with Owner and Date tags
func tagResource(event Event, logger *Logger) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Minute)
	defer cancel()
	AzureTagerPrefix := os.Getenv("AZURE_TAGGER_PREFIX")
	if AzureTagerPrefix == "" {
		AzureTagerPrefix = "AzTagger"
	}
	// AZURE_CLIENT_ID is used as Managed Identity reference
	// it's picked up from env vars automatically
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		logger.Error.Printf("failed to obtain a credential: %v", err)
		return fmt.Errorf("failed to obtain a credential: %w", err)
	}

	// client, err := armresources.NewTagsClient(os.Getenv("AZURE_SUBSCRIPTION_ID"), cred, nil)
	// use event.Data.SubscriptionId from the request, to allow multi-subscription events
	client, err := armresources.NewTagsClient(event.Data.SubscriptionId, cred, nil)
	if err != nil {
		logger.Error.Printf("failed to create resources client: %v", err)
		return fmt.Errorf("failed to create resources client: %w", err)
	}

	tags := map[string]*string{
		fmt.Sprintf("%sCreationDate", AzureTagerPrefix): to.Ptr(event.EventTime),
	}
	if event.Data.Claims.Email != "" {
		tags[fmt.Sprintf("%sEmail", AzureTagerPrefix)] = to.Ptr(event.Data.Claims.Email)
	} else if event.Data.Claims.ClaimsName != "" {
		tags[fmt.Sprintf("%sEmail", AzureTagerPrefix)] = to.Ptr(event.Data.Claims.ClaimsName)
	}
	if event.Data.Claims.Name != "" {
		tags[fmt.Sprintf("%sUserName", AzureTagerPrefix)] = to.Ptr(event.Data.Claims.Name)
	} else {
		spName, err := getServicePrincipalNameFromId(cred, event, logger)
		if err == nil {
			tags[fmt.Sprintf("%sUserName", AzureTagerPrefix)] = spName
		} else {
			tags[fmt.Sprintf("%sUserName", AzureTagerPrefix)] = to.Ptr(event.Data.Claims.Appid)
		}
	}

	tagsToSet := make(map[string]*string)
	tagsToSet, err = getTagsToSet(client, &event, tags, logger)
	if err != nil {
		return err
	}
	if len(tagsToSet) == 0 {
		logger.Info.Printf("All specified tags already exist. No update needed.")
		return nil
	}
	response, err := client.CreateOrUpdateAtScope(
		ctx,
		event.Subject,
		armresources.TagsResource{
			Properties: &armresources.Tags{
				Tags: tagsToSet,
			},
		},
		nil,
	)
	if err != nil {
		logger.Error.Printf("failed to update resource tags: %v", err)
		return fmt.Errorf("failed to update resource tags: %w", err)
	}
	logger.Info.Println("created tags:", convertMap(response.Properties.Tags))

	return nil
}
