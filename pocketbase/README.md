# Pocketbase

There are two flavors of the backend:

1. standard release downloaded from https://github.com/pocketbase/pocketbase/releases. This one is a good start, but most real-world applications would require more (see next).
2. custom compiled (`go build`), possibly with my customizations and perhaps yours too.

Out of the box, the project assumes #2 (custom compiled with my customizations).

## standard (official) release of pocketbase

Download from release archive from https://github.com/pocketbase/pocketbase/releases/latest, unzip it and place the `pocketbase` binary in this folder, and you're done.

## custom build

If you would like to extend PocketBase and use it as a framework then there is a `main.go` in this folder that you can customize and build using `go build` or do live development using `modd`.

See https://pocketbase.io/docs/use-as-framework/ for details.

# Setup

## Architecture

> **Note:** For optimal set up, ensure you are using a standard distribution of Linux. For other operating systems, you may run into issues, or need additional configuration.
> A docker-compose setup is included with the project, which can be used on any OS.

## Build

Assuming you have Go language tools installed ...

`go build`

If you don't have Go and don't want to install it, you can use docker-compose setup. Otherwise, your only choice is to download the binary from https://github.com/pocketbase/pocketbase/releases/latest, and placing it in this folder. But then you will not be able to use any of the custom code (such as "config-driven hooks")

## Run migrations

Before you can run the actual backend, you must run the migrations using `./pocketbase migrate up` in the current directory. It will create appropriate schema tables/collections.

## Run the backend

You can run the PocketBase backend direct with `./pocketbase serve` or using `npm run backend` in the `sk` directory. Note that if you want the backend to also serve the frontend assets, then you must add the `--publicDir ../frontend/build` option.

## Docker

A highly recommended option is to run it inside a Docker container. A `Dockerfile` is included that builds a production Docker image. Also, a `docker-compose.yml` along with an _override_ file example are included, which should be used during development.

## Active development with `modd`

Finally, if you are going to actively develop using Go using PocketBase as a framework, then you probably want to use [modd](https://github.com/cortesi/modd), a development tool that rebuilds and restarts your Go binary everytime a source file changes (live reload on change). An basic `modd.conf` config file is included in this setup. You can run it by installing `modd` (`go install github.com/cortesi/modd/cmd/modd@latest`) and then running `modd`. All this is done automatically for you if you are using Docker.

# Schema (Collections)

With the 0.9 version of PocketBase, JavaScript auto-migrations as implemented. The JS files in `pb_migrations` can create/drop/modify collections and data. These are executed automatically by PocketBase on startup.

Not only that, they are also generated automatically whenever you change the schema! So go ahead and make changes to the schema and watch new JS files generated in the `pb_migrations` folder. Just remember to commit them to version control.

## Generated Types

The file `generated-types.ts` contains TypeScript definitions of `Record` types mirroring the fields in your database collections. But it needs to be regenerated every time you modify the schema. This can be done by simply running the `typegen` script in the frontend's `package.json`. So remember to do that.

# Hooks

PocketBase provides API's like .OnModelBefore* and .OnModelAfter* to run
callbacks when records change. This app builds on top of that by providing
a "hooks" table that drives those hooks using configuration. It has the
following fields:

- collection: name of the collection that triggers an action
- event: insert/update/delete event that triggers the action
- action_type: "command" if you want to run a program/script or "post" if
  you want to POST to a webhook endpoint. The record will be marshaled to
  JSON and passed to the command as STDIN or to the webhook POST as
  request body (with header 'content-type: application/json')
- action: path to the command/script or URL of the webhook to POST to
- action_params: a string that will be passed as argument to the action

So now by configuring the above table, you can execute external commands/scripts
and POST data to external webhooks in reaction to insert/update/delete of
records.

Most web services these days provide webhook endpoints (e.g. sendgrid, mailchimp, stripe, etc) which you can POST directly to. But if you need special
processing then you can write a script that receives changed data as JSON, parses and manipulates it using [`jq`](https://github.com/stedolan/jq) before
sending it on its way.

See `example-hook-script.sh` for a demonstration.

Possible use cases:

- Clone git repo when a record is inserted to "repositories" table
- Execute a terraform script when a new cluster is inserted to "clusters" table
- Send an acknowledgement email when a "contact" form table is inserted to.
- Charge a credit card when payment_token table is inserted to and then
  send email upon success/failure
- Recalculate inventory levels as "orders" table is inserted to, and then
  send notifications when inventory becomes low.
