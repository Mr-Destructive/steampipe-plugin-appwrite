package appwrite

import (
	"context"
	"encoding/json"

	appwrite "github.com/mr-destructive/appwrite-go-sdk"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableDatabases(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "appwrite_databases",
		Description: "Query database meta information in an appwrite project",
		List: &plugin.ListConfig{
			Hydrate: databases,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "search_query", Require: plugin.Optional},
				{Name: "query", Require: plugin.Optional},
				{Name: "settings", Require: plugin.Optional},
			},
		},
		Columns: []*plugin.Column{
			// Result columns
			{Name: "id", Type: proto.ColumnType_STRING, Transform: transform.FromField("Database.Id"), Description: "id"},
			{Name: "name", Type: proto.ColumnType_STRING, Transform: transform.FromField("Database.Name"), Description: "Name"},

			// Input Columns
			{Name: "search_query", Type: proto.ColumnType_STRING, Transform: transform.FromField("Search")},
			{Name: "query", Type: proto.ColumnType_STRING, Transform: transform.FromField("Query")},
			{Name: "settings", Type: proto.ColumnType_JSON, Transform: transform.FromQual("settings"), Description: "Settings is a JSONB object that accepts any of the completion API request parameters."},
		},
	}
}

type databasessRequestQual struct {
	Search *string   `json:"search_query"`
	Query  *[]string `json:"query"`
}

type databasesRow struct {
	Database appwrite.DatabaseObject
	Search   string
	Query    []string
}

func databases(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("appwrite.databases", "connection_error", err)
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
		var crQual databasessRequestQual
		err := json.Unmarshal([]byte(settingsString), &crQual)
		if err != nil {
			plugin.Logger(ctx).Error("appwrite.databases", "unmarshal_error", err)
			return nil, err
		}
		if crQual.Query != nil {
			query = *crQual.Query
		}
		if crQual.Search != nil {
			search = *crQual.Search
		}
	}

	database := appwrite.Database{
		Client: *conn,
	}

	databasesList, err := database.ListDatabases(search, query)
	if err != nil {
		plugin.Logger(ctx).Error("appwrite.databases", "api_error", err)
		return nil, err
	}
	plugin.Logger(ctx).Trace("appwrite.databases", "response", databasesList)
	databases := *databasesList
	for _, database := range databases.Databases {
		row := databasesRow{
			Database: database,
			Search:   search,
			Query:    query,
		}
		d.StreamListItem(ctx, row)
	}
	return nil, nil
}
