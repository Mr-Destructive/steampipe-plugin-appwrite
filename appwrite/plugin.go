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
			"appwrite_bucket":     tableAppwriteBucket(ctx),
			"appwrite_collection": tableAppwriteCollection(ctx),
			"appwrite_database":   tableAppwriteDatabase(ctx),
			"appwrite_document":   tableAppwriteDocument(ctx),
			"appwrite_deployment": tableAppwriteDeployment(ctx),
			"appwrite_execution":  tableAppwriteExecution(ctx),
			"appwrite_file":       tableAppwriteFile(ctx),
			"appwrite_function":   tableAppwriteFunction(ctx),
			"appwrite_health":     tableAppwriteHealth(ctx),
			"appwrite_user":       tableAppwriteUser(ctx),
		},
	}
	return p
}
