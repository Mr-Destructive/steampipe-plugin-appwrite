package appwrite

import (
	"context"
	"encoding/json"

	appwrite "github.com/mr-destructive/appwrite-go-sdk"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableCollections(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "appwrite_collections",
		Description: "Query collections of an appwrite database.",
		List: &plugin.ListConfig{
			Hydrate: collections,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "database_id", Require: plugin.Optional},
				{Name: "search", Require: plugin.Optional},
				{Name: "settings", Require: plugin.Optional},
			},
		},
		Columns: []*plugin.Column{
			// Result columns
			{Name: "id", Type: proto.ColumnType_STRING, Transform: transform.FromField("Id"), Description: "id"},
			{Name: "name", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name"), Description: "Name"},

			// Input Columns
			{Name: "database_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("DatabaseId"), Description: "DatabaseId"},
			{Name: "search", Type: proto.ColumnType_STRING, Transform: transform.FromField("Search")},
			{Name: "offset", Type: proto.ColumnType_INT, Transform: transform.FromField("Offset")},
			{Name: "settings", Type: proto.ColumnType_JSON, Transform: transform.FromQual("settings"), Description: "Settings is a JSONB object that accepts any of the completion API request parameters."},
		},
	}
}

type collectionsRequestQual struct {
	DatabaseId *string `json:"database_id"`
	Search     *string
	Order      *string
}

type collectionRow struct {
	appwrite.Collection
}

func collections(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("appwrite.collections", "connection_error", err)
		return nil, err
	}

	settingsString := d.EqualsQuals["settings"].GetJsonbValue()
	if settingsString != "" {
		var crQual collectionsRequestQual
		err := json.Unmarshal([]byte(settingsString), &crQual)
		if err != nil {
			plugin.Logger(ctx).Error("appwrite.collections", "unmarshal_error", err)
			return nil, err
		}
	}

	database := appwrite.Database{
		Client: *conn,
	}
	databaseId := d.EqualsQuals["database_id"].GetStringValue()
	search := d.EqualsQuals["search"].GetStringValue()

	collectionList, err := database.ListCollections(databaseId, search, []string{})
	if err != nil {
		plugin.Logger(ctx).Error("appwrite.collections", "api_error", err)
		return nil, err
	}
	plugin.Logger(ctx).Trace("appwrite.collections", "response", collectionList)
	collections := *collectionList
	for _, collection := range collections.Collections {
		row := collectionRow{collection}
		d.StreamListItem(ctx, row)
	}
	return nil, nil
}
