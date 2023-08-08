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
		Description: "",
		List: &plugin.ListConfig{
			Hydrate: databases,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "search", Require: plugin.Optional},
				{Name: "settings", Require: plugin.Optional},
			},
		},
		Columns: []*plugin.Column{
			// Result columns
			{Name: "id", Type: proto.ColumnType_STRING, Transform: transform.FromField("Id"), Description: "id"},
			{Name: "name", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name"), Description: "Name"},

			// Input Columns
			{Name: "search", Type: proto.ColumnType_STRING, Transform: transform.FromField("Search")},
			{Name: "offset", Type: proto.ColumnType_INT, Transform: transform.FromField("Offset")},
			{Name: "settings", Type: proto.ColumnType_JSON, Transform: transform.FromQual("settings"), Description: "Settings is a JSONB object that accepts any of the completion API request parameters."},
		},
	}
}

type databasessRequestQual struct {
	Search *string
	Order  *string
}

type databasesRow struct {
	appwrite.DatabaseObject
}

func databases(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("appwrite.databases", "connection_error", err)
		return nil, err
	}

	settingsString := d.EqualsQuals["settings"].GetJsonbValue()
	if settingsString != "" {
		var crQual databasessRequestQual
		err := json.Unmarshal([]byte(settingsString), &crQual)
		if err != nil {
			plugin.Logger(ctx).Error("appwrite.databases", "unmarshal_error", err)
			return nil, err
		}
	}

	database := appwrite.Database{
		Client: *conn,
	}
	search := d.EqualsQuals["search"].GetStringValue()

	databasesList, err := database.ListDatabases(search, []string{})
	if err != nil {
		plugin.Logger(ctx).Error("appwrite.databases", "api_error", err)
		return nil, err
	}
	plugin.Logger(ctx).Trace("appwrite.databases", "response", databasesList)
	databases := *databasesList
	for _, database := range databases.Databases {
		row := databasesRow{database}
		d.StreamListItem(ctx, row)
	}
	return nil, nil
}
