package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
)

// getTagsToSet is a function to tag the resource with Owner and Date tags
func getTagsToSet(client *armresources.TagsClient, event *Event, newTags map[string]*string, logger *Logger) (map[string]*string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Minute)
	defer cancel()
	var err error
	existingTagsResponse, err := client.GetAtScope(
		ctx,
		event.Subject,
		nil,
	)
	if err != nil {
		logger.Error.Printf("failed to fetch resource tags: %v", err)
		return nil, fmt.Errorf("failed to fetch resource tags: %w", err)
	}
	existingTags := existingTagsResponse.Properties.Tags

	logger.Info.Println("existing tags:", convertMap(existingTagsResponse.Properties.Tags))
	// Filter out new tags that already exist
	tagsToSet := filterExistingTags(existingTags, newTags)
	return tagsToSet, err
}

// filterExistingTags is a helper function to filter out tags that already exist by name
func filterExistingTags(existingTags, newTags map[string]*string) map[string]*string {
	tagsToSet := make(map[string]*string)

	// Add new tags that do not already exist
	for key, value := range newTags {
		if _, exists := existingTags[key]; !exists {
			tagsToSet[key] = value
		}
	}
	if len(tagsToSet) == 0 {
		return tagsToSet
	}

	// Copy existing tags into the merged map (existing tags take precedence)
	for key, value := range existingTags {
		tagsToSet[key] = value
	}

	return tagsToSet
}
