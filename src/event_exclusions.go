package main

var EventTaggingExcludedActions = []string{
	"Microsoft.Resources/tags/write",
	"Microsoft.Compute/restorePointCollections/restorePoints/write",
	"Microsoft.Authorization/roleAssignments/write",
	// "Microsoft.Sql/servers/firewallRules/write",
	// "Microsoft.Web/sites/config/write",
	// "Microsoft.Logic/workflows/write",
	// "Microsoft.EventGrid/eventSubscriptions/write",
	// "Microsoft.Web/serverFarms/write",
	// "Microsoft.ManagedIdentity/userAssignedIdentities/write",
	// "Microsoft.OperationalInsights/workspaces/write",
	// "microsoft.insights/actionGroups/write",
	// "Microsoft.Insights/components/write",
}
