package appwrite

import (
	"context"
	"encoding/json"

	appwrite "github.com/mr-destructive/appwrite-go-sdk"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableExecutions(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "appwrite_executions",
		Description: "",
		List: &plugin.ListConfig{
			Hydrate: executions,
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

type executionssRequestQual struct {
	Search *string
	Order  *string
}

type executionsRow struct {
	appwrite.ExecutionObject
	search string
}

func executions(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("appwrite.executions", "connection_error", err)
		return nil, err
	}

	settingsString := d.EqualsQuals["settings"].GetJsonbValue()
	if settingsString != "" {
		var crQual executionssRequestQual
		err := json.Unmarshal([]byte(settingsString), &crQual)
		if err != nil {
			plugin.Logger(ctx).Error("appwrite.executions", "unmarshal_error", err)
			return nil, err
		}
	}

	functions := appwrite.Function{
		Client: *conn,
	}
	search := d.EqualsQuals["search"].GetStringValue()

	executionsList, err := functions.ListExecutions(search, "", []string{})
	if err != nil {
		plugin.Logger(ctx).Error("appwrite.executions", "api_error", err)
		return nil, err
	}
	plugin.Logger(ctx).Trace("appwrite.executions", "response", executionsList)
	executions := *executionsList
	for _, execution := range executions.Executions {
		row := executionsRow{execution, search}
		d.StreamListItem(ctx, row)
	}
	return nil, nil
}
