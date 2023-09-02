package appwrite

import (
	"context"
	"encoding/json"

	appwrite "github.com/mr-destructive/appwrite-go-sdk"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAppwriteUser(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "appwrite_user",
		Description: "Query users in an appwrite project",
		List: &plugin.ListConfig{
			Hydrate: listUsers,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "search_query", Require: plugin.Optional},
				{Name: "offset", Require: plugin.Optional},
				{Name: "settings", Require: plugin.Optional},
			},
		},
		Columns: []*plugin.Column{
			// Result columns
			{Name: "id", Type: proto.ColumnType_STRING, Transform: transform.FromField("Id"), Description: "The unique ID of the account user."},
			{Name: "title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Id"), Description: "The Name or ID of the user."},
			{Name: "email", Type: proto.ColumnType_STRING, Transform: transform.FromField("Email"), Description: "The Email ID of the account user."},
			{Name: "name", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name"), Description: "The Name of the account user."},
			{Name: "status", Type: proto.ColumnType_BOOL, Transform: transform.FromField("Status"), Description: "The active status of the account user."},
			{Name: "phone", Type: proto.ColumnType_STRING, Transform: transform.FromField("Phone"), Description: "The phone number of the account user."},
			{Name: "password", Type: proto.ColumnType_STRING, Transform: transform.FromField("Phone"), Description: "The password of the account user."},
			{Name: "email_verification", Type: proto.ColumnType_BOOL, Transform: transform.FromField("EmailVerification"), Description: "The status of the email verification of the account user."},
			{Name: "phone_verification", Type: proto.ColumnType_BOOL, Transform: transform.FromField("PhoneVerification"), Description: "The status of the phone verification of the account user."},

			// Input Columns
			{Name: "search_query", Type: proto.ColumnType_STRING, Transform: transform.FromField("Search"), Description: "The string as a search filter the results from the request."},
			{Name: "settings", Type: proto.ColumnType_JSON, Transform: transform.FromQual("settings"), Description: "Settings is a JSONB object that accepts any of the completion API request parameters."},
		},
	}
}

type usersRequestQual struct {
	Search *string `json:"search_query"`
	Offset *int    `json:"offset"`
	Limit  *int
	Order  *string
}

type usersRow struct {
	appwrite.UserObject
	Search string
}

func listUsers(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("appwrite_user.listUsers", "connection_error", err)
		return nil, err
	}

	settingsString := d.EqualsQuals["settings"].GetJsonbValue()
	if settingsString != "" {
		// Overwrite any settings provided in the settings qual. If a field
		// is not passed in the settings, then default to the settings above.
		var crQual usersRequestQual
		err := json.Unmarshal([]byte(settingsString), &crQual)
		if err != nil {
			plugin.Logger(ctx).Error("appwrite_user.listUsers", "unmarshal_error", err)
			return nil, err
		}
	}

	users := appwrite.Users{
		Client: *conn,
	}
	search := d.EqualsQuals["search_query"].GetStringValue()

	userList, err := users.List(search, 0, 0, "")
	if err != nil {
		plugin.Logger(ctx).Error("appwrite_user.listUsers", "api_error", err)
		return nil, err
	}
	plugin.Logger(ctx).Trace("appwrite_user.listUsers", "response", userList)
	for _, u := range userList {
		row := usersRow{u, search}
		d.StreamListItem(ctx, row)
	}

	return nil, nil
}
