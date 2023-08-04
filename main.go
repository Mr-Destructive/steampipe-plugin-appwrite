package main

import (
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
    "github.com/turbot/steampipe-plugin-appwrite/appwrite"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{PluginFunc: appwrite.Plugin})
}
