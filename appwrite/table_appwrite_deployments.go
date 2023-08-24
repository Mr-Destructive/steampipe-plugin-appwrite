package appwrite

import (
	"context"
	"encoding/json"

	appwrite "github.com/mr-destructive/appwrite-go-sdk"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableDeployments(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "appwrite_deployments",
		Description: "Query deployment information of a function for an appwrite project",
		List: &plugin.ListConfig{
			Hydrate: deployments,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "search_query", Require: plugin.Optional},
				{Name: "query", Require: plugin.Optional},
				{Name: "settings", Require: plugin.Optional},
			},
		},
		Columns: []*plugin.Column{
			// Result columns
			{Name: "id", Type: proto.ColumnType_STRING, Transform: transform.FromField("Deployment.Id"), Description: "id"},
			{Name: "name", Type: proto.ColumnType_STRING, Transform: transform.FromField("Deployment.Name"), Description: "Name"},
			{Name: "created_at", Type: proto.ColumnType_STRING, Transform: transform.FromField("Deployment.CreatedAt"), Description: "created at"},
			{Name: "updated_at", Type: proto.ColumnType_STRING, Transform: transform.FromField("Deployment.UpdatedAt"), Description: "updated at"},
			{Name: "resource_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("Deployment.ResourceId"), Description: "resource id"},
			{Name: "resource_type", Type: proto.ColumnType_STRING, Transform: transform.FromField("Deployment.ResourceType"), Description: "resource type"},
			{Name: "entry_point", Type: proto.ColumnType_STRING, Transform: transform.FromField("Deployment.EntryPoint"), Description: "entry point"},
			{Name: "size", Type: proto.ColumnType_INT, Transform: transform.FromField("Deployment.Size"), Description: "size"},
			{Name: "build_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("Deployment.BuildId"), Description: "build id"},
			{Name: "activate", Type: proto.ColumnType_BOOL, Transform: transform.FromField("Deployment.Activate"), Description: "activate"},
			{Name: "status", Type: proto.ColumnType_STRING, Transform: transform.FromField("Deployment.Status"), Description: "status"},
			{Name: "build_stdout", Type: proto.ColumnType_STRING, Transform: transform.FromField("Deployment.BuildStdout"), Description: "build stdout"},
			{Name: "build_stderr", Type: proto.ColumnType_STRING, Transform: transform.FromField("Deployment.BuildStderr"), Description: "build stderr"},
			{Name: "build_time", Type: proto.ColumnType_STRING, Transform: transform.FromField("Deployment.BuildTime"), Description: "build time"},

			// Input Columns
			{Name: "function_id", Type: proto.ColumnType_STRING, Transform: transform.FromQual("function_id"), Description: "function id"},
			{Name: "search_query", Type: proto.ColumnType_STRING, Transform: transform.FromField("Search")},
			{Name: "query", Type: proto.ColumnType_STRING, Transform: transform.FromField("Query")},
			{Name: "settings", Type: proto.ColumnType_JSON, Transform: transform.FromQual("settings"), Description: "Settings is a JSONB object that accepts any of the completion API request parameters."},
		},
	}
}

type deploymentsRequestQual struct {
	FunctionId *string   `json:"function_id"`
	Search     *string   `json:"search_query"`
	Query      *[]string `json:"query"`
}

type deploymentsRow struct {
	Deployment appwrite.DeploymentObject
	Search     string
	Query      []string
}

func deployments(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("appwrite.deployments", "connection_error", err)
		return nil, err
	}
	queryString := d.EqualsQuals["query"].GetStringValue()
	var query []string
	if queryString != "" {
		err := json.Unmarshal([]byte(queryString), &query)
		if err != nil {
			plugin.Logger(ctx).Error("appwrite_functions.listFunctions", "connection_error", err)
			return nil, err
		}
	}
	if len(query) == 0 {
		query = []string{}
	}

	function_id := d.EqualsQuals["function_id"].GetStringValue()
	search := d.EqualsQuals["search_query"].GetStringValue()

	settingsString := d.EqualsQuals["settings"].GetJsonbValue()
	if settingsString != "" {
		var crQual deploymentsRequestQual
		err := json.Unmarshal([]byte(settingsString), &crQual)
		if err != nil {
			plugin.Logger(ctx).Error("appwrite.deployments", "unmarshal_error", err)
			return nil, err
		}
		if crQual.Query != nil {
			query = *crQual.Query
		}
		if crQual.Search != nil {
			search = *crQual.Search
		}
	}

	functions := appwrite.Function{
		Client: *conn,
	}
	deploymentsList, err := functions.ListDeployments(function_id, "", query)
	if err != nil {
		plugin.Logger(ctx).Error("appwrite.deployments", "api_error", err)
		return nil, err
	}
	plugin.Logger(ctx).Trace("appwrite.deployments", "response", deploymentsList)
	deployments := *deploymentsList
	for _, deployment := range deployments.Deployments {
		row := deploymentsRow{
			Deployment: deployment,
			Search:     search,
			Query:      query,
		}
		d.StreamListItem(ctx, row)
	}
	return nil, nil
}
