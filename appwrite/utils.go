package appwrite

import (
	"context"
	"errors"
	"os"
	"strings"

	appwrite "github.com/mr-destructive/appwrite-go-sdk"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func connect(ctx context.Context, d *plugin.QueryData) (*appwrite.Client, error) {

	cacheKey := "appwrite"
	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.(*appwrite.Client), nil
	}

	conn, err := connectCached(ctx, d, nil)
	if err != nil {
		return nil, err
	}
	return conn.(*appwrite.Client), nil
}

var connectCached = plugin.HydrateFunc(connectUncached).Memoize()

func connectUncached(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {

	conn := &appwrite.Client{}

	// Default to the env var settings
	secretKey := os.Getenv("APPWRITE_SECRET_KEY")
	projectID := os.Getenv("APPWRITE_PROJECT_ID")

	// Prefer config settings
	appwriteConfig := GetConfig(d.Connection)
	if appwriteConfig.SecretKey != nil || appwriteConfig.ProjectID != nil {
		secretKey = *appwriteConfig.SecretKey
		projectID = *appwriteConfig.ProjectID
	}

	// Error if the minimum config is not set
	if secretKey == "" || projectID == "" {
		return conn, errors.New("api_key must be configured")
	}

	conn.SetEndpoint("https://cloud.appwrite.io/v1")
	conn.SetKey(secretKey)
	conn.SetProject(projectID)

	return conn, nil
}

func isNotFoundError(err error) bool {
	return strings.Contains(err.Error(), "status code: 404")
}
