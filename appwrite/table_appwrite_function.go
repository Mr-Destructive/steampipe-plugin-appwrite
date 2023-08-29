package appwrite

import (
	"context"
	"encoding/json"

	appwrite "github.com/mr-destructive/appwrite-go-sdk"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAppwriteFunction(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "appwrite_function",
		Description: "Query meta information of a function for an appwrite project",
		List: &plugin.ListConfig{
			Hydrate: listFunctions,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "search_query", Require: plugin.Optional},
				{Name: "query", Require: plugin.Optional},
				{Name: "settings", Require: plugin.Optional},
			},
		},
		Columns: []*plugin.Column{
			// Result columns
			{Name: "id", Type: proto.ColumnType_STRING, Transform: transform.FromField("Function.Id"), Description: "The unique ID for the function."},
			{Name: "name", Type: proto.ColumnType_STRING, Transform: transform.FromField("Function.Name"), Description: "The Name of the function."},
			{Name: "created_at", Type: proto.ColumnType_STRING, Transform: transform.FromField("Function.CreatedAt"), Description: "Function creation date in ISO 8601 format."},
			{Name: "updated_at", Type: proto.ColumnType_STRING, Transform: transform.FromField("Function.UpdatedAt"), Description: "Function updation date in ISO 8601 format."},
			{Name: "execute", Type: proto.ColumnType_STRING, Transform: transform.FromField("Function.Execute"), Description: "A list of string as permissions for the execution of the function."},
			{Name: "enabled", Type: proto.ColumnType_BOOL, Transform: transform.FromField("Function.Enabled"), Description: "A boolean flag to indicate if the function is enabled."},
			{Name: "variable", Type: proto.ColumnType_STRING, Transform: transform.FromField("Function.Variable"), Description: "The list of variables for the function."},
			{Name: "runtime", Type: proto.ColumnType_STRING, Transform: transform.FromField("Function.Runtime"), Description: "The runtime for the function execution."},
			{Name: "deployment", Type: proto.ColumnType_STRING, Transform: transform.FromField("Function.Deployment"), Description: "Function's active deployment ID."},
			{Name: "events", Type: proto.ColumnType_STRING, Transform: transform.FromField("Function.Events"), Description: "The list of trigger events for the function."},
			{Name: "schedule", Type: proto.ColumnType_STRING, Transform: transform.FromField("Function.Schedule"), Description: "The schedule for the function execution in CRON format."},
			{Name: "schedule_next", Type: proto.ColumnType_STRING, Transform: transform.FromField("Function.ScheduleNext"), Description: "The next scheduled execution time of function in ISO 8601 format."},
			{Name: "schedule_previous", Type: proto.ColumnType_STRING, Transform: transform.FromField("Function.SchedulePrevious"), Description: "The previous scheduled execution time of function in ISO 8601 format."},
			{Name: "timeout", Type: proto.ColumnType_STRING, Transform: transform.FromField("Function.Timeout"), Description: "The execution time of the function in seconds."},

			// Input Columns
			{Name: "search_query", Type: proto.ColumnType_STRING, Transform: transform.FromField("Search"), Description: "The string as a search filter the results from the request."},
			{Name: "query", Type: proto.ColumnType_STRING, Transform: transform.FromField("Query"), Description: "The string of query type to filter the results from the request."},
			{Name: "settings", Type: proto.ColumnType_JSON, Transform: transform.FromQual("settings"), Description: "Settings is a JSONB object that accepts any of the completion API request parameters."},
		},
	}
}

type functionssRequestQual struct {
	Search *string   `json:"search_query"`
	Query  *[]string `json:"query"`
}

type functionsRow struct {
	Function appwrite.FunctionObject
	Query    []string
	Search   string
}

func listFunctions(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("appwrite_function.functions", "connection_error", err)
		return nil, err
	}

	queryString := d.EqualsQuals["query"].GetStringValue()
	var query []string

	if queryString != "" {
		err := json.Unmarshal([]byte(queryString), &query)
		if err != nil {
			plugin.Logger(ctx).Error("appwrite_function.listFunctions", "connection_error", err)
			return nil, err
		}
	}
	if len(query) == 0 {
		query = []string{}
	}
	search := d.EqualsQuals["search_query"].GetStringValue()

	settingsString := d.EqualsQuals["settings"].GetJsonbValue()
	if settingsString != "" {
		var crQual functionssRequestQual
		err := json.Unmarshal([]byte(settingsString), &crQual)
		if err != nil {
			plugin.Logger(ctx).Error("appwrite_function.listFunctions", "unmarshal_error", err)
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

	functionsList, err := functions.ListFunctions(search, query)
	if err != nil {
		plugin.Logger(ctx).Error("appwrite_function.listFunctions", "api_error", err)
		return nil, err
	}
	plugin.Logger(ctx).Trace("appwrite_function.listFunctions", "response", functionsList)
	funcs := *functionsList
	for _, f := range funcs.Functions {
		row := functionsRow{
			Function: f,
			Query:    query,
			Search:   search,
		}
		d.StreamListItem(ctx, row)
	}
	return nil, nil
}
