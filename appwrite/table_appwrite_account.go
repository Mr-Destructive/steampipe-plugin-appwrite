package appwrite

import (
	"context"
	"encoding/json"

	appwrite "github.com/mr-destructive/appwrite-go-sdk"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAppwriteAccount(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "appwrite_account",
		Description: "Query users in an appwrite project",
		List: &plugin.ListConfig{
			Hydrate: listAccounts,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "search_query", Require: plugin.Optional},
				{Name: "offset", Require: plugin.Optional},
				{Name: "settings", Require: plugin.Optional},
			},
		},
		Columns: []*plugin.Column{
			// Result columns
			{Name: "id", Type: proto.ColumnType_STRING, Transform: transform.FromField("Id"), Description: "id"},
			{Name: "email", Type: proto.ColumnType_STRING, Transform: transform.FromField("Email"), Description: "Email"},
			{Name: "name", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name"), Description: "Name"},
			{Name: "status", Type: proto.ColumnType_STRING, Transform: transform.FromField("Status"), Description: "status"},
			{Name: "phone", Type: proto.ColumnType_STRING, Transform: transform.FromField("Phone"), Description: "phone"},
			{Name: "password", Type: proto.ColumnType_STRING, Transform: transform.FromField("Phone"), Description: "phone"},
			{Name: "email_verification", Type: proto.ColumnType_BOOL, Transform: transform.FromField("EmailVerification"), Description: "email verification"},
			{Name: "phone_verification", Type: proto.ColumnType_BOOL, Transform: transform.FromField("PhoneVerification"), Description: "phone verification"},

			// Input Columns
			{Name: "search_query", Type: proto.ColumnType_STRING, Transform: transform.FromField("Search")},
			{Name: "offset", Type: proto.ColumnType_INT, Transform: transform.FromField("Offset")},
			{Name: "settings", Type: proto.ColumnType_JSON, Transform: transform.FromQual("settings"), Description: "Settings is a JSONB object that accepts any of the completion API request parameters."},
		},
	}
}

type accountsRequestQual struct {
	Search *string `json:"search_query"`
	Offset *int    `json:"offset"`
	Limit  *int
	Order  *string
}

type accountsRow struct {
	appwrite.UserObject
	Search string
}

func listAccounts(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("appwrite_appwrite_account.listAccounts", "connection_error", err)
		return nil, err
	}

	settingsString := d.EqualsQuals["settings"].GetJsonbValue()
	if settingsString != "" {
		// Overwrite any settings provided in the settings qual. If a field
		// is not passed in the settings, then default to the settings above.
		var crQual accountsRequestQual
		err := json.Unmarshal([]byte(settingsString), &crQual)
		if err != nil {
			plugin.Logger(ctx).Error("appwrite_account.listAccounts", "unmarshal_error", err)
			return nil, err
		}
	}

	// Query the sdk with appropriate methods and serialize the response
	users := appwrite.Users{
		Client: *conn,
	}
	search := d.EqualsQuals["search_query"].GetStringValue()

	userList, err := users.List(search, 0, 0, "")
	if err != nil {
		plugin.Logger(ctx).Error("appwrite_account.listAccounts", "api_error", err)
		return nil, err
	}
	plugin.Logger(ctx).Trace("appwrite_account.listAccounts", "response", userList)
	for _, u := range userList {
		row := accountsRow{u, search}
		d.StreamListItem(ctx, row)
	}

	return nil, nil
}
