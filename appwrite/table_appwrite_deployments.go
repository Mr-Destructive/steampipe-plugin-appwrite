package appwrite

import (
	"context"
	"encoding/json"

	appwrite "github.com/mr-destructive/appwrite-go-sdk"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableDeployments(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "appwrite_deployments",
		Description: "",
		List: &plugin.ListConfig{
			Hydrate: accounts,
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

type deploymentssRequestQual struct {
	Search *string
	Order  *string
}

type deploymentsRow struct {
	appwrite.DeploymentListResponse
}

func deployments(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("appwrite.deployments", "connection_error", err)
		return nil, err
	}

	settingsString := d.EqualsQuals["settings"].GetJsonbValue()
	if settingsString != "" {
		var crQual accountsRequestQual
		err := json.Unmarshal([]byte(settingsString), &crQual)
		if err != nil {
			plugin.Logger(ctx).Error("appwrite.deployments", "unmarshal_error", err)
			return nil, err
		}
	}

	storage := appwrite.Storage{
		Client: *conn,
	}
	search := d.EqualsQuals["search"].GetStringValue()
	limit := d.EqualsQuals["limit"].GetInt64Value()
	offset := d.EqualsQuals["offset"].GetInt64Value()
	order := d.EqualsQuals["order"].GetStringValue()

	deploymentsList, err := storage.ListDeployments(search, int(limit), int(offset), order)
	if err != nil {
		plugin.Logger(ctx).Error("appwrite.deployments", "api_error", err)
		return nil, err
	}
	plugin.Logger(ctx).Trace("appwrite.deployments", "response", deploymentsList)
	for _, f := range deploymentsList {
		row := deploymentsRow{f}
		d.StreamListItem(ctx, row)
	}
	return nil, nil
}
