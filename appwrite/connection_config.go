package appwrite

import (
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/schema"
)

type appwriteConfig struct {
	ProjectID *string `cty:"project_id" hcl:"project_id"`
	SecretKey *string `cty:"secret_key" hcl:"secret_key"`
}

var ConfigSchema = map[string]*schema.Attribute{
	"project_id": {
		Type: schema.TypeString,
	},
	"secret_key": {
		Type: schema.TypeString,
	},
}

func ConfigInstance() interface{} {
	return &appwriteConfig{}
}

func GetConfig(connection *plugin.Connection) appwriteConfig {
	if connection == nil || connection.Config == nil {
		return appwriteConfig{}
	}
	config, _ := connection.Config.(appwriteConfig)
	return config
}
