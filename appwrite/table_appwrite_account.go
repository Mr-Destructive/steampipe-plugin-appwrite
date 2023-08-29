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
			{Name: "id", Type: proto.ColumnType_STRING, Transform: transform.FromField("Id"), Description: "The unique ID of the user account."},
			{Name: "email", Type: proto.ColumnType_STRING, Transform: transform.FromField("Email"), Description: "The Email ID of the user account."},
			{Name: "name", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name"), Description: "The Name of the user account."},
			{Name: "status", Type: proto.ColumnType_STRING, Transform: transform.FromField("Status"), Description: "The active status of the user account."},
			{Name: "phone", Type: proto.ColumnType_STRING, Transform: transform.FromField("Phone"), Description: "The phone number of the user account."},
			{Name: "password", Type: proto.ColumnType_STRING, Transform: transform.FromField("Phone"), Description: "The password of the user account."},
			{Name: "email_verification", Type: proto.ColumnType_BOOL, Transform: transform.FromField("EmailVerification"), Description: "The status of the email verification of the user account."},
			{Name: "phone_verification", Type: proto.ColumnType_BOOL, Transform: transform.FromField("PhoneVerification"), Description: "The status of the phone verification of the user account."},

			// Input Columns
			{Name: "search_query", Type: proto.ColumnType_STRING, Transform: transform.FromField("Search"), Description: "The string as a search filter the results from the request."},
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
