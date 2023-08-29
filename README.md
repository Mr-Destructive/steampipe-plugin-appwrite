
# Appwrite Plugin for Steampipe

Use SQL to query account and more from [Appwrite](https://appwrite.io/).

- **[Get started →](https://hub.steampipe.io/plugins/mr-destructive/appwrite)**
- Documentation: [Table definitions & examples](https://hub.steampipe.io/plugins/mr-destructive/appwrite/tables)
- Community: [Slack Channel](https://steampipe.io/community/join)
- Get involved: [Issues](https://github.com/mr-destructive/steampipe-plugin-appwrite/issues)

## Quick start

Install the plugin with [Steampipe](https://steampipe.io):

```shell
steampipe plugin install mr-destructive/appwrite
```

Configure your plugin in `~/.steampipe/config/appwrite.spc`:

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
Or through environment variables:

```
export APPWRITE_SECRET_KEY="7a1f0d410a6ab90110232e3f9578a0e5ac33453493930e195c7"
export APPWRITE_PROJECT_ID="68a121f3e41164679a30"
```

Run steampipe:

```shell
steampipe query
```

Run a query:

```sql
select
  *
from
  account
where
  name = ''
```

## Developing

Prerequisites:

- [Steampipe](https://steampipe.io/downloads)
- [Golang](https://golang.org/doc/install)

Clone:

```sh
git clone https://github.com/mr-destructive/steampipe-plugin-appwrite.git
cd steampipe-plugin-appwrite
```

Build, which automatically installs the new version to your `~/.steampipe/plugins` directory:

```
make
```

Configure the plugin:

```
cp config/* ~/.steampipe/config
vi ~/.steampipe/config/appwrite.spc
```

Try it!

```
steampipe query
> .inspect appwrite
```

Further reading:

- [Writing plugins](https://steampipe.io/docs/develop/writing-plugins)
- [Writing your first table](https://steampipe.io/docs/develop/writing-your-first-table)

## Contributing

Please see the [contribution guidelines](https://github.com/turbot/steampipe/blob/main/CONTRIBUTING.md) and our [code of conduct](https://github.com/turbot/steampipe/blob/main/CODE_OF_CONDUCT.md). All contributions are subject to the [Apache 2.0 open source license](https://github.com/mr-destructive/steampipe-plugin-appwrite/blob/main/LICENSE).

`help wanted` issues:

- [Steampipe](https://github.com/turbot/steampipe/labels/help%20wanted)
- [AppwritePlugin](https://github.com/mr-destructive/steampipe-plugin-appwrite/labels/help%20wanted)
