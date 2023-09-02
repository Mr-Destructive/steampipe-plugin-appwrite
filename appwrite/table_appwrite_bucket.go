package appwrite

import (
	"context"
	"encoding/json"

	appwrite "github.com/mr-destructive/appwrite-go-sdk"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAppwriteBucket(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "appwrite_bucket",
		Description: "Query buckets meta information from an appwrite project",
		List: &plugin.ListConfig{
			Hydrate: listBuckets,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "search_query", Require: plugin.Optional},
				{Name: "settings", Require: plugin.Optional},
			},
		},
		Columns: []*plugin.Column{
			// Result columns
			{Name: "id", Type: proto.ColumnType_STRING, Transform: transform.FromField("Id"), Description: "The ID of the bucket."},
			{Name: "title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Id"), Description: "The Name or ID of the bucket."},
			{Name: "name", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name"), Description: "The name of the bucket."},
			{Name: "file_extensions", Type: proto.ColumnType_STRING, Transform: transform.FromField("AllowedFileExtensions"), Description: "The allowed file extensions for the bucket."},
			{Name: "created_at", Type: proto.ColumnType_STRING, Transform: transform.FromField("CreatedAt"), Description: "Bucket creation time in ISO 8601 format."},
			{Name: "updated_at", Type: proto.ColumnType_STRING, Transform: transform.FromField("UpdatedAt"), Description: "Bucket update time in ISO 8601 format"},
			{Name: "permissions", Type: proto.ColumnType_JSON, Transform: transform.FromField("Permissions"), Description: "The permission setting(list of strings) for the bucket."},
			{Name: "file_security", Type: proto.ColumnType_BOOL, Transform: transform.FromField("FileSecurity"), Description: "A boolean value/flag for file-level security is enabled on the bucket."},
			{Name: "enabled", Type: proto.ColumnType_BOOL, Transform: transform.FromField("Enabled"), Description: "Flag for checking if the bucket is enabled or disabled as storage."},
			{Name: "maximum_file_size", Type: proto.ColumnType_INT, Transform: transform.FromField("MaximumFileSize"), Description: "Maximum file size supported in the bucket."},
			{Name: "compression_type", Type: proto.ColumnType_STRING, Transform: transform.FromField("CompressionType"), Description: "Compression algorithm choosen for compression. Will be one of none, gzip, or zstd"},
			{Name: "encryption", Type: proto.ColumnType_BOOL, Transform: transform.FromField("Encryption"), Description: "A boolean value for chacking if encryption is enabled in the bucket or not."},
			{Name: "antivirus", Type: proto.ColumnType_BOOL, Transform: transform.FromField("Antivirus"), Description: "A boolean value for chacking if the virus scanning is enabled in the bucket or not."},

			// Input Columns
			{Name: "search_query", Type: proto.ColumnType_STRING, Transform: transform.FromField("Search"), Description: "The string value to filter the results from the request."},
			{Name: "settings", Type: proto.ColumnType_JSON, Transform: transform.FromQual("settings"), Description: "Settings is a JSONB object that accepts any of the completion API request parameters."},
		},
	}
}

type bucketsRequestQual struct {
	Search *string
}

type bucketsRow struct {
	appwrite.Bucket
	Search string
}

func listBuckets(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("appwrite_bucket.listBuckets", "connection_error", err)
		return nil, err
	}

	settingsString := d.EqualsQuals["settings"].GetJsonbValue()
	if settingsString != "" {
		// Overwrite any settings provided in the settings qual. If a field
		// is not passed in the settings, then default to the settings above.
		var crQual bucketsRequestQual
		err := json.Unmarshal([]byte(settingsString), &crQual)
		if err != nil {
			plugin.Logger(ctx).Error("appwrite_bucket.listBuckets", "unmarshal_error", err)
			return nil, err
		}
	}

	storage := appwrite.Storage{
		Client: *conn,
	}
	search := d.EqualsQuals["search_query"].GetStringValue()

	bucketsList, err := storage.ListBuckets(search, 0, 0, "")
	if err != nil {
		plugin.Logger(ctx).Error("appwrite_bucket.listBuckets", "api_error", err)
		return nil, err
	}
	plugin.Logger(ctx).Trace("appwrite_bucket.listBuckets", "response", bucketsList)
	buckers := *bucketsList
	for _, bucket := range buckers.Buckets {
		row := bucketsRow{bucket, search}
		d.StreamListItem(ctx, row)
	}

	return nil, nil
}
