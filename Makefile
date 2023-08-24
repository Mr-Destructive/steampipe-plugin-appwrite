install:
	go build -o ~/.steampipe/plugins/hub.steampipe.io/plugins/mr-destructive/appwrite@latest/steampipe-plugin-appwrite.plugin *.go
local:
	/home/meet/code/playground/github/steampipe/go/bin/go build -o ~/.steampipe/plugins/hub.steampipe.io/plugins/mr-destructive/appwrite@latest/steampipe-plugin-appwrite.plugin *.go
