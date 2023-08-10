package appwrite

import (
	"context"
	"encoding/json"

	appwrite "github.com/mr-destructive/appwrite-go-sdk"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableFunctions(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "appwrite_functions",
		Description: "Query meta information of a function for an appwrite project",
		List: &plugin.ListConfig{
			Hydrate: functions,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "search_query", Require: plugin.Optional},
				{Name: "query", Require: plugin.Optional},
				{Name: "settings", Require: plugin.Optional},
			},
		},
		Columns: []*plugin.Column{
			// Result columns
			{Name: "id", Type: proto.ColumnType_STRING, Transform: transform.FromField("Function.Id"), Description: "id"},
			{Name: "name", Type: proto.ColumnType_STRING, Transform: transform.FromField("Function.Name"), Description: "Name"},
			{Name: "runtime", Type: proto.ColumnType_STRING, Transform: transform.FromField("Function.Runtime"), Description: "Runtime"},

			// Input Columns
			{Name: "search_query", Type: proto.ColumnType_STRING, Transform: transform.FromField("Search")},
			{Name: "query", Type: proto.ColumnType_STRING, Transform: transform.FromField("Query")},
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

func functions(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("appwrite_functions.functions", "connection_error", err)
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
	search := d.EqualsQuals["search_query"].GetStringValue()

	settingsString := d.EqualsQuals["settings"].GetJsonbValue()
	if settingsString != "" {
		var crQual functionssRequestQual
		err := json.Unmarshal([]byte(settingsString), &crQual)
		if err != nil {
			plugin.Logger(ctx).Error("appwrite.functions", "unmarshal_error", err)
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
		plugin.Logger(ctx).Error("appwrite.functions", "api_error", err)
		return nil, err
	}
	plugin.Logger(ctx).Trace("appwrite.functions", "response", functionsList)
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
