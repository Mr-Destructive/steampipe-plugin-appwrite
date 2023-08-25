package appwrite

import (
	"context"
	"encoding/json"

	appwrite "github.com/mr-destructive/appwrite-go-sdk"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAppwriteCollection(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "appwrite_collection",
		Description: "Query collections of an appwrite database.",
		List: &plugin.ListConfig{
			Hydrate: listCollections,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "database_id", Require: plugin.Optional},
				{Name: "search_query", Require: plugin.Optional},
				{Name: "query", Require: plugin.Optional},
				{Name: "settings", Require: plugin.Optional},
			},
		},
		Columns: []*plugin.Column{
			// Result columns
			{Name: "id", Type: proto.ColumnType_STRING, Transform: transform.FromField("Collection.Id"), Description: "id"},
			{Name: "name", Type: proto.ColumnType_STRING, Transform: transform.FromField("Collection.Name"), Description: "Name"},
			{Name: "created_at", Type: proto.ColumnType_STRING, Transform: transform.FromField("Collection.CreatedAt"), Description: "created at"},
			{Name: "updated_at", Type: proto.ColumnType_STRING, Transform: transform.FromField("Collection.UpdatedAt"), Description: "updated at"},
			{Name: "document_security", Type: proto.ColumnType_BOOL, Transform: transform.FromField("Collection.DocumentSecurity"), Description: "document security"},
			{Name: "permissions", Type: proto.ColumnType_JSON, Transform: transform.FromField("Collection.Permissions"), Description: "permissions"},
			{Name: "enabled", Type: proto.ColumnType_BOOL, Transform: transform.FromField("Collection.Enabled"), Description: "enabled"},
			{Name: "attributes", Type: proto.ColumnType_JSON, Transform: transform.FromField("Collection.Attributes"), Description: "attributes"},
			{Name: "indexes", Type: proto.ColumnType_JSON, Transform: transform.FromField("Collection.Indexes"), Description: "indexes"},

			// Input Columns
			{Name: "database_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("DatabaseId"), Description: "DatabaseId"},
			{Name: "search_query", Type: proto.ColumnType_STRING, Transform: transform.FromField("Search")},
			{Name: "query", Type: proto.ColumnType_STRING, Transform: transform.FromField("Query")},
			{Name: "settings", Type: proto.ColumnType_JSON, Transform: transform.FromQual("settings"), Description: "Settings is a JSONB object that accepts any of the completion API request parameters."},
		},
	}
}

type collectionsRequestQual struct {
	DatabaseId *string   `json:"database_id"`
	Search     *string   `json:"search_query"`
	Query      *[]string `json:"query"`
}

type collectionRow struct {
	Collection appwrite.Collection
	Search     string
	Query      []string
}

func listCollections(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("appwrite_collection.listCollections", "connection_error", err)
		return nil, err
	}

	queryString := d.EqualsQuals["query"].GetStringValue()
	var query []string
	if queryString != "" {
		err := json.Unmarshal([]byte(queryString), &query)
		if err != nil {
			plugin.Logger(ctx).Error("appwrite_collection.listCollections", "connection_error", err)
			return nil, err
		}
	}
	search := d.EqualsQuals["search_query"].GetStringValue()

	settingsString := d.EqualsQuals["settings"].GetJsonbValue()
	if settingsString != "" {
		var crQual collectionsRequestQual
		err := json.Unmarshal([]byte(settingsString), &crQual)
		if err != nil {
			plugin.Logger(ctx).Error("appwrite_collection.listCollections", "unmarshal_error", err)
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
	databaseId := d.EqualsQuals["database_id"].GetStringValue()

	collectionList, err := database.ListCollections(databaseId, search, query)
	if err != nil {
		plugin.Logger(ctx).Error("appwrite_collection.listCollections", "api_error", err)
		return nil, err
	}
	plugin.Logger(ctx).Trace("appwrite_collection.listCollections", "response", collectionList)
	collections := *collectionList
	for _, collection := range collections.Collections {
		row := collectionRow{
			Collection: collection,
			Search:     search,
			Query:      query,
		}
		d.StreamListItem(ctx, row)
	}
	return nil, nil
}
