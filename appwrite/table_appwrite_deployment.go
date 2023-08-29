package appwrite

import (
	"context"
	"encoding/json"

	appwrite "github.com/mr-destructive/appwrite-go-sdk"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAppwriteDeployment(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "appwrite_deployment",
		Description: "Query deployment information of a function for an appwrite project",
		List: &plugin.ListConfig{
			Hydrate: listDeployments,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "search_query", Require: plugin.Optional},
				{Name: "query", Require: plugin.Optional},
				{Name: "settings", Require: plugin.Optional},
			},
		},
		Columns: []*plugin.Column{
			// Result columns
			{Name: "id", Type: proto.ColumnType_STRING, Transform: transform.FromField("Deployment.Id"), Description: "The unique ID for the deployment."},
			{Name: "created_at", Type: proto.ColumnType_STRING, Transform: transform.FromField("Deployment.CreatedAt"), Description: "Deployment creation date in ISO 8601 format."},
			{Name: "updated_at", Type: proto.ColumnType_STRING, Transform: transform.FromField("Deployment.UpdatedAt"), Description: "Deployment updation date in ISO 8601 format."},
			{Name: "resource_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("Deployment.ResourceId"), Description: "The unique ID for the resource in the deployment."},
			{Name: "resource_type", Type: proto.ColumnType_STRING, Transform: transform.FromField("Deployment.ResourceType"), Description: "The type of resource in the deployment."},
			{Name: "entry_point", Type: proto.ColumnType_STRING, Transform: transform.FromField("Deployment.EntryPoint"), Description: "The entrypoint file to use to execute the deployment code."},
			{Name: "size", Type: proto.ColumnType_INT, Transform: transform.FromField("Deployment.Size"), Description: "The size of code for deployment function in bytes."},
			{Name: "build_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("Deployment.BuildId"), Description: "The unique ID for the current build of the deployment."},
			{Name: "activate", Type: proto.ColumnType_BOOL, Transform: transform.FromField("Deployment.Activate"), Description: "A boolean to indicate whether the deployment should be automatically activated."},
			{Name: "status", Type: proto.ColumnType_STRING, Transform: transform.FromField("Deployment.Status"), Description: "The deployment status as either processing, building, pending, ready, or failed"},
			{Name: "build_stdout", Type: proto.ColumnType_STRING, Transform: transform.FromField("Deployment.BuildStdout"), Description: "The standard output for the current build of deployment."},
			{Name: "build_stderr", Type: proto.ColumnType_STRING, Transform: transform.FromField("Deployment.BuildStderr"), Description: "The standard error for the current build of deployment."},
			{Name: "build_time", Type: proto.ColumnType_STRING, Transform: transform.FromField("Deployment.BuildTime"), Description: "The time taken for the current build in seconds."},

			// Input Columns
			{Name: "function_id", Type: proto.ColumnType_STRING, Transform: transform.FromQual("function_id"), Description: "The unique ID for the function to fetch the deployments from."},
			{Name: "search_query", Type: proto.ColumnType_STRING, Transform: transform.FromField("Search"), Description: "The string to filter the results from the request."},
			{Name: "query", Type: proto.ColumnType_STRING, Transform: transform.FromField("Query"), Description: "A string of query type as filter for the request."},
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

func listDeployments(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("appwrite_deployment.listDeployments", "connection_error", err)
		return nil, err
	}
	queryString := d.EqualsQuals["query"].GetStringValue()
	var query []string
	if queryString != "" {
		err := json.Unmarshal([]byte(queryString), &query)
		if err != nil {
			plugin.Logger(ctx).Error("appwrite_deployment.listDeployments", "connection_error", err)
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
			plugin.Logger(ctx).Error("appwrite_deployment.listDeployments", "unmarshal_error", err)
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
		plugin.Logger(ctx).Error("appwrite_deployment.listDeployments", "api_error", err)
		return nil, err
	}
	plugin.Logger(ctx).Trace("appwrite_deployment.listDeployments", "response", deploymentsList)
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
