package appwrite

import (
	"context"
	"encoding/json"

	appwrite "github.com/mr-destructive/appwrite-go-sdk"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAppwriteDocument(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "appwrite_document",
		Description: "Query documents of a collection from an appwrite database",
		List: &plugin.ListConfig{
			Hydrate: listDocuments,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "database_id", Require: plugin.Optional},
				{Name: "collection_id", Require: plugin.Optional},
				{Name: "search_query", Require: plugin.Optional},
				{Name: "settings", Require: plugin.Optional},
			},
		},
		Columns: []*plugin.Column{
			// Result columns
			{Name: "id", Type: proto.ColumnType_STRING, Transform: transform.FromField("Document.Id"), Description: "The unique ID for the document."},
			{Name: "name", Type: proto.ColumnType_STRING, Transform: transform.FromField("Document.Name"), Description: "The Name of the document."},
			{Name: "title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Id"), Description: "The Name or ID of the document."},
			{Name: "fields", Type: proto.ColumnType_JSON, Transform: transform.FromField("Document.Fields"), Description: "The fields in the document."},
			{Name: "created_at", Type: proto.ColumnType_STRING, Transform: transform.FromField("Document.CreatedAt"), Description: "created at"},
			{Name: "updated_at", Type: proto.ColumnType_STRING, Transform: transform.FromField("Document.UpdatedAt"), Description: "updated at"},
			{Name: "permissions", Type: proto.ColumnType_JSON, Transform: transform.FromField("Document.Permissions"), Description: "permissions"},

			// Input Columns
			{Name: "database_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("DatabaseId"), Description: "DatabaseId"},
			{Name: "collection_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("CollectionId"), Description: "CollectionId"},
			{Name: "search_query", Type: proto.ColumnType_STRING, Transform: transform.FromField("Search"), Description: "The string to filter the results from the request."},
			{Name: "settings", Type: proto.ColumnType_JSON, Transform: transform.FromQual("settings"), Description: "Settings is a JSONB object that accepts any of the completion API request parameters."},
		},
	}
}

type documentsRequestQual struct {
	DatabaseId   *string `json:"database_id"`
	CollectionId *string `json:"collection_id"`
	Search       *string `json:"search_query"`
}

type documentRow struct {
	Document appwrite.Document
	Search   string
}

func listDocuments(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("appwrite_document.listDocuments", "connection_error", err)
		return nil, err
	}

	settingsString := d.EqualsQuals["settings"].GetJsonbValue()
	if settingsString != "" {
		var crQual documentsRequestQual
		err := json.Unmarshal([]byte(settingsString), &crQual)
		if err != nil {
			plugin.Logger(ctx).Error("appwrite_document.listDocuments", "unmarshal_error", err)
			return nil, err
		}
	}

	database := appwrite.Database{
		Client: *conn,
	}
	databaseId := d.EqualsQuals["database_id"].GetStringValue()
	collectionId := d.EqualsQuals["collection_id"].GetStringValue()
	search := d.EqualsQuals["search_query"].GetStringValue()

	documentList, err := database.ListDocuments(databaseId, collectionId, []interface{}{}, 0, 0, "", "", "", search, 0, 0)
	if err != nil {
		plugin.Logger(ctx).Error("appwrite_document.listDocuments", "api_error", err)
		return nil, err
	}
	plugin.Logger(ctx).Trace("appwrite_document.listDocuments", "response", documentList)
	documents := *documentList
	for _, document := range documents.Documents {
		row := documentRow{
			Document: document,
			Search:   search,
		}
		d.StreamListItem(ctx, row)
	}
	return nil, nil
}
