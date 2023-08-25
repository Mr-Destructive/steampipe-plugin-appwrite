package appwrite

import (
	"context"
	"encoding/json"

	appwrite "github.com/mr-destructive/appwrite-go-sdk"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAppwriteExecution(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "appwrite_execution",
		Description: "Query executions meta information of a function deployment in an appwrite project",
		List: &plugin.ListConfig{
			Hydrate: listExecutions,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "search_query", Require: plugin.Optional},
				{Name: "query", Require: plugin.Optional},
				{Name: "settings", Require: plugin.Optional},
			},
		},
		Columns: []*plugin.Column{
			// Result columns
			{Name: "id", Type: proto.ColumnType_STRING, Transform: transform.FromField("Execution.Id"), Description: "id"},
			{Name: "name", Type: proto.ColumnType_STRING, Transform: transform.FromField("Execution.Name"), Description: "Name"},
			{Name: "created_at", Type: proto.ColumnType_STRING, Transform: transform.FromField("Execution.CreatedAt"), Description: "created at"},
			{Name: "updated_at", Type: proto.ColumnType_STRING, Transform: transform.FromField("Execution.UpdatedAt"), Description: "updated at"},
			{Name: "permissions", Type: proto.ColumnType_JSON, Transform: transform.FromField("Execution.Permissions"), Description: "permissions"},
			{Name: "function_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("Execution.FunctionId"), Description: "function id"},
			{Name: "trigger", Type: proto.ColumnType_STRING, Transform: transform.FromField("Execution.Trigger"), Description: "trigger"},
			{Name: "status", Type: proto.ColumnType_STRING, Transform: transform.FromField("Execution.Status"), Description: "status"},
			{Name: "status_code", Type: proto.ColumnType_INT, Transform: transform.FromField("Execution.StatusCode"), Description: "status code"},
			{Name: "response", Type: proto.ColumnType_STRING, Transform: transform.FromField("Execution.Response"), Description: "response"},
			{Name: "stdout", Type: proto.ColumnType_STRING, Transform: transform.FromField("Execution.Stdout"), Description: "stdout"},
			{Name: "stderr", Type: proto.ColumnType_STRING, Transform: transform.FromField("Execution.Stderr"), Description: "stderr"},
			{Name: "duration", Type: proto.ColumnType_STRING, Transform: transform.FromField("Execution.Duration"), Description: "duration"},

			// Input Columns
			{Name: "function_id", Type: proto.ColumnType_STRING, Transform: transform.FromQual("function_id"), Description: "function id"},
			{Name: "search_query", Type: proto.ColumnType_STRING, Transform: transform.FromField("Search")},
			{Name: "query", Type: proto.ColumnType_STRING, Transform: transform.FromField("Query")},
			{Name: "settings", Type: proto.ColumnType_JSON, Transform: transform.FromQual("settings"), Description: "Settings is a JSONB object that accepts any of the completion API request parameters."},
		},
	}
}

type executionssRequestQual struct {
	Search *string   `json:"search_query"`
	Query  *[]string `json:"query"`
}

type executionsRow struct {
	Execution appwrite.ExecutionObject
	Search    string
	Query     []string
}

func listExecutions(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("appwrite_execution.listExecutions", "connection_error", err)
		return nil, err
	}
	queryString := d.EqualsQuals["query"].GetStringValue()
	var query []string
	if queryString != "" {
		err := json.Unmarshal([]byte(queryString), &query)
		if err != nil {
			plugin.Logger(ctx).Error("appwrite_execution.listExecutions", "connection_error", err)
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
		var crQual executionssRequestQual
		err := json.Unmarshal([]byte(settingsString), &crQual)
		if err != nil {
			plugin.Logger(ctx).Error("appwrite_execution.listExecutions", "unmarshal_error", err)
			return nil, err
		}
	}

	functions := appwrite.Function{
		Client: *conn,
	}

	executionsList, err := functions.ListExecutions(function_id, "", query)
	if err != nil {
		plugin.Logger(ctx).Error("appwrite_execution.listExecutions", "api_error", err)
		return nil, err
	}
	plugin.Logger(ctx).Trace("appwrite_execution.listExecutions", "response", executionsList)
	executions := *executionsList
	for _, execution := range executions.Executions {
		row := executionsRow{
			Execution: execution,
			Search:    search,
			Query:     query,
		}
		d.StreamListItem(ctx, row)
	}
	return nil, nil
}