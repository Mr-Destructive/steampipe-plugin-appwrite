package appwrite

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name: "steampipe-plugin-appwrite",
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
			Schema:      ConfigSchema,
		},
		DefaultTransform: transform.FromGo().NullIfZero(),
		DefaultGetConfig: &plugin.GetConfig{
			ShouldIgnoreError: isNotFoundError,
		},
		TableMap: map[string]*plugin.Table{
			"appwrite_accounts":    tableAccounts(ctx),
			"appwrite_buckets":     tableBuckets(ctx),
			"appwrite_collections": tableCollections(ctx),
			"appwrite_databases":   tableDatabases(ctx),
			"appwrite_documents":   tableDocuments(ctx),
			"appwrite_deployments": tableDeployments(ctx),
			"appwrite_executions":  tableExecutions(ctx),
			"appwrite_files":       tableFiles(ctx),
			"appwrite_functions":   tableFunctions(ctx),
		},
	}
	return p
}
