package appwrite

import (
	"context"
	"encoding/json"

	appwrite "github.com/mr-destructive/appwrite-go-sdk"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableBuckets(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "appwrite_buckets",
		Description: "Query buckets meta information from an appwrite project",
		List: &plugin.ListConfig{
			Hydrate: buckets,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "search_query", Require: plugin.Optional},
				{Name: "settings", Require: plugin.Optional},
			},
		},
		Columns: []*plugin.Column{
			// Result columns
			{Name: "id", Type: proto.ColumnType_STRING, Transform: transform.FromField("Id"), Description: "id"},
			{Name: "name", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name"), Description: "Name"},
			{Name: "file_extensions", Type: proto.ColumnType_STRING, Transform: transform.FromField("AllowedFileExtensions"), Description: "file extensions"},
			{Name: "created_at", Type: proto.ColumnType_STRING, Transform: transform.FromField("CreatedAt"), Description: "created at"},
			{Name: "updated_at", Type: proto.ColumnType_STRING, Transform: transform.FromField("UpdatedAt"), Description: "updated at"},
			{Name: "permissions", Type: proto.ColumnType_JSON, Transform: transform.FromField("Permissions"), Description: "permissions"},
			{Name: "file_security", Type: proto.ColumnType_BOOL, Transform: transform.FromField("FileSecurity"), Description: "file security"},
			{Name: "enabled", Type: proto.ColumnType_BOOL, Transform: transform.FromField("Enabled"), Description: "enabled"},
			{Name: "maximum_file_size", Type: proto.ColumnType_INT, Transform: transform.FromField("MaximumFileSize"), Description: "maximum file size"},
			{Name: "compression_type", Type: proto.ColumnType_STRING, Transform: transform.FromField("CompressionType"), Description: "compression type"},
			{Name: "encryption", Type: proto.ColumnType_BOOL, Transform: transform.FromField("Encryption"), Description: "encryption"},
			{Name: "antivirus", Type: proto.ColumnType_BOOL, Transform: transform.FromField("Antivirus"), Description: "antivirus"},

			// Input Columns
			{Name: "search_query", Type: proto.ColumnType_STRING, Transform: transform.FromField("Search")},
			{Name: "settings", Type: proto.ColumnType_JSON, Transform: transform.FromQual("settings"), Description: "Settings is a JSONB object that accepts any of the completion API request parameters."},
		},
	}
}

type bucketsRequestQual struct {
	Search *string
	Order  *string
}

type bucketsRow struct {
	appwrite.Bucket
	Search string
}

func buckets(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// NOTE: IMPLEMENT THIS based on the service/provider
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("appwrite.buckets", "connection_error", err)
		return nil, err
	}

	settingsString := d.EqualsQuals["settings"].GetJsonbValue()
	if settingsString != "" {
		// Overwrite any settings provided in the settings qual. If a field
		// is not passed in the settings, then default to the settings above.
		var crQual accountsRequestQual
		err := json.Unmarshal([]byte(settingsString), &crQual)
		if err != nil {
			plugin.Logger(ctx).Error("appwrite.buckets", "unmarshal_error", err)
			return nil, err
		}
	}

	// Query the sdk with appropriate methods and serialize the response
	storage := appwrite.Storage{
		Client: *conn,
	}
	search := d.EqualsQuals["search_query"].GetStringValue()

	bucketsList, err := storage.ListBuckets(search, 0, 0, "")
	if err != nil {
		plugin.Logger(ctx).Error("appwrite.buckets", "api_error", err)
		return nil, err
	}
	plugin.Logger(ctx).Trace("appwrite.buckets", "response", bucketsList)
	buckers := *bucketsList
	for _, bucket := range buckers.Buckets {
		row := bucketsRow{bucket, search}
		d.StreamListItem(ctx, row)
	}

	return nil, nil
}
