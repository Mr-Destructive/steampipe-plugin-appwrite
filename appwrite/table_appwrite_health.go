package appwrite

import (
	"context"
	"encoding/json"

	appwrite "github.com/mr-destructive/appwrite-go-sdk"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAppwriteHealth(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "appwrite_health",
		Description: "Get health of various services in your appwrite project",
		List: &plugin.ListConfig{
			Hydrate: health,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "service", Require: plugin.Optional},
				{Name: "settings", Require: plugin.Optional},
			},
		},
		Columns: []*plugin.Column{
			// Result columns
			{Name: "ping", Type: proto.ColumnType_INT, Transform: transform.FromField("Status.Ping")},
			{Name: "status", Type: proto.ColumnType_STRING, Transform: transform.FromField("Status.Status")},
			{Name: "real_time", Type: proto.ColumnType_INT, Transform: transform.FromField("Time.RealTime")},
			{Name: "local_time", Type: proto.ColumnType_INT, Transform: transform.FromField("Time.LocalTime")},
			{Name: "diff", Type: proto.ColumnType_INT, Transform: transform.FromField("Time.Diff")},
			{Name: "size", Type: proto.ColumnType_INT, Transform: transform.FromField("Queue.Size")},

			// Input Columns
			{Name: "service", Type: proto.ColumnType_STRING, Transform: transform.FromField("Service")},
			{Name: "settings", Type: proto.ColumnType_JSON, Transform: transform.FromQual("settings"), Description: "Settings is a JSONB object that accepts any of the completion API request parameters."},
		},
	}
}

type healthsRequestQual struct {
	Service string `json:"service"`
}

type healthRow struct {
	Status appwrite.HealthStatus
	Queue  appwrite.HealthQueue
	Time   appwrite.HealthTime
}

func health(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("appwrite_health.health", "connection_error", err)
		return nil, err
	}
	service := d.EqualsQuals["service"].GetStringValue()

	settingsString := d.EqualsQuals["settings"].GetJsonbValue()
	if settingsString != "" {
		var crQual healthsRequestQual
		err := json.Unmarshal([]byte(settingsString), &crQual)
		if err != nil {
			plugin.Logger(ctx).Error("appwrite_health.health", "unmarshal_error", err)
			return nil, err
		}
	}
	client := *conn
	var row interface{}

	switch service {
	case "http":
		response, _ := client.Health()
		row = healthRow{Status: *response}
	case "db":
		response, _ := client.DBHealth()
		row = healthRow{Status: *response}
	case "cache":
		response, _ := client.CacheHealth()
		row = healthRow{Status: *response}
	case "local-storage":
		response, _ := client.LocalStorageHealth()
		row = healthRow{Status: *response}
	case "function-queue":
		response, _ := client.FunctionsQueue()
		row = healthRow{Queue: *response}
	case "logs-queue":
		response, _ := client.LogsQueue()
		row = healthRow{Queue: *response}
	case "webhooks-queue":
		response, _ := client.WebHooksQueue()
		row = healthRow{Queue: *response}
	case "time":
		response, _ := client.TimeHealth()
		row = healthRow{Time: *response}
	default:
		response, _ := client.TimeHealth()
		row = healthRow{Time: *response}
	}
	if err != nil {
		plugin.Logger(ctx).Error("appwrite_health.health", "api_error", err)
		return nil, err
	}
	d.StreamListItem(ctx, row)
	return nil, nil
}
