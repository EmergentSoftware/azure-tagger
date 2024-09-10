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
func tagResource(event Event) error {
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
		ErrorLogger.Printf("failed to obtain a credential: %v", err)
		return fmt.Errorf("failed to obtain a credential: %w", err)
	}

	client, err := armresources.NewTagsClient(os.Getenv("AZURE_SUBSCRIPTION_ID"), cred, nil)
	if err != nil {
		ErrorLogger.Printf("failed to create resources client: %v", err)
		return fmt.Errorf("failed to create resources client: %w", err)
	}

	dateTag := time.Now().UTC().Format(time.RFC3339)

	tags := map[string]*string{
		fmt.Sprintf("%sEmail", AzureTagerPrefix):        to.Ptr(event.Data.Claims.Email),
		fmt.Sprintf("%sUserName", AzureTagerPrefix):     to.Ptr(event.Data.Claims.Name),
		fmt.Sprintf("%sCreationDate", AzureTagerPrefix): to.Ptr(dateTag),
	}
	tagsToSet := make(map[string]*string)
	tagsToSet, err = getTagsToSet(client, &event, tags)
	if err != nil {
		return err
	}
	if len(tagsToSet) == 0 {
		InfoLogger.Printf("All specified tags already exist. No update needed.")
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
		ErrorLogger.Printf("failed to update resource tags: %v", err)
		return fmt.Errorf("failed to update resource tags: %w", err)
	}
	InfoLogger.Println("created tags:", convertMap(response.Properties.Tags))

	return nil
}
