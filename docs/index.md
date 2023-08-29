---
organization: Turbot
category: ["ai"]
icon_url: "/images/plugins/turbot/appwrite.svg"
brand_color: "#000000"
display_name: "Appwrite"
short_name: "appwrite"
description: "."
og_description: "Query Appwrite with SQL! Open source CLI. No DB required."
og_image: "/images/plugins/turbot/appwrite-social-graphic.png"
---

# Appwrite + Steampipe


[Steampipe](https://steampipe.io) is an open source CLI to instantly query cloud APIs using SQL.

Get generations for a given text prompt in your Appwrite account:

```sql
select
  
from
  
where
  
```

```
```

## Documentation

- **[Table definitions & examples →](/docs/tables)**

## Get started

### Install

Download and install the latest Appwrite plugin:

```bash
steampipe plugin install mr-destructive/appwrite
```

### Credentials

| Item        | Description                                                                                                                                                                                                                                                                                 |
|-------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Credentials |                                                                                                                                                                                  |
| Permissions | API Keys have the same permissions as the user who creates them, and if the user permissions change, the API key permissions also change.                                                                                                                                               |
| Radius      | Each connection represents a single appwrite Installation.                                                                                                                                                                                                                                   |
| Resolution  | 1. Credentials explicitly set in a steampipe config file (`~/.steampipe/config/appwrite.spc`)<br />2. Credentials specified in environment variables. |

### Configuration

Installing the latest comereai plugin will create a config file (`~/.steampipe/config/appwrite.spc`) with a single connection named `appwrite`:

```hcl
connection "appwrite" {
  plugin = "mr-destructive/appwrite"

  # Secret key for requests. Required.
  # This can also be set via the `APPWRITE_SECRET_KEY` environment variable.
  # secret_key = "7a1f0d410a6ab90110232e3f9578a0e5ac33453493930e195c7"

  # Project Id for specific appwrite project. Required
  # This can also be set via the `APPWRITE_PROJECT_ID` environment variable.
  # project_id = "68a121f3e41164679a30"
}
```

## Credentials from Environment Variables

The Appwrite plugin will use the standard Appwrite environment variables to obtain credentials only if other arguments (secret_key and project_id) are not specified in the connection:

```
export APPWRITE_SECRET_KEY="7a1f0d410a6ab90110232e3f9578a0e5ac33453493930e195c7"
export APPWRITE_PROJECT_ID="68a121f3e41164679a30"
```

## Get involved

- Open source: https://github.com/turbot/steampipe-plugin-appwrite
- Community: [Slack Channel](https://steampipe.io/community/join)
