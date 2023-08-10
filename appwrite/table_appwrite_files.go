package appwrite

import (
	"context"
	"encoding/json"

	appwrite "github.com/mr-destructive/appwrite-go-sdk"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableFiles(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "appwrite_files",
		Description: "Query files meta information in a bucket for an appwrite project",
		List: &plugin.ListConfig{
			Hydrate: files,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "bucket_id", Require: plugin.Optional},
				{Name: "search_query", Require: plugin.Optional},
				{Name: "settings", Require: plugin.Optional},
			},
		},
		Columns: []*plugin.Column{
			// Result columns
			{Name: "id", Type: proto.ColumnType_STRING, Transform: transform.FromField("Id"), Description: "id"},
			{Name: "name", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name"), Description: "Name"},

			// Input Columns
			{Name: "bucket_id", Type: proto.ColumnType_STRING, Transform: transform.FromQual("bucket_id"), Description: "bucket id"},
			{Name: "search_query", Type: proto.ColumnType_STRING, Transform: transform.FromField("Search")},
			{Name: "offset", Type: proto.ColumnType_INT, Transform: transform.FromField("Offset")},
			{Name: "settings", Type: proto.ColumnType_JSON, Transform: transform.FromQual("settings"), Description: "Settings is a JSONB object that accepts any of the completion API request parameters."},
		},
	}
}

type filesRequestQual struct {
	BucketId *string `json:"bucket_id"`
	Search   *string
	Order    *string
}

type filesRow struct {
	appwrite.File
	BucketId string
	Search   string
}

func files(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("appwrite.files", "connection_error", err)
		return nil, err
	}

	settingsString := d.EqualsQuals["settings"].GetJsonbValue()
	if settingsString != "" {
		var crQual filesRequestQual
		err := json.Unmarshal([]byte(settingsString), &crQual)
		if err != nil {
			plugin.Logger(ctx).Error("appwrite.files", "unmarshal_error", err)
			return nil, err
		}
	}

	storage := appwrite.Storage{
		Client: *conn,
	}
	bucketId := d.EqualsQuals["bucket_id"].GetStringValue()
	search := d.EqualsQuals["search_query"].GetStringValue()

	filesList, err := storage.ListFiles(bucketId, search, 0, 0, "")
	if err != nil {
		plugin.Logger(ctx).Error("appwrite.files", "api_error", err)
		return nil, err
	}
	plugin.Logger(ctx).Trace("appwrite.files", "response", filesList)
	files := *filesList
	for _, f := range files.Files {
		row := filesRow{f, bucketId, search}
		d.StreamListItem(ctx, row)
	}
	return nil, nil
}
