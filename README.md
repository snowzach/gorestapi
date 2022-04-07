# Base API Example

This API example is a basic framework for a REST API

## Compiling
This is designed as a go module aware program and thus requires go 1.11 or better
You can clone it anywhere, just run `make` inside the cloned directory to build

## Requirements
This does require a postgres database to be setup and reachable. It will attempt to create and migrate the database upon starting.

## Configuration
The configuration is designed to be specified with environment variables in all caps with underscores instead of periods. 
```
Example:
LOGGER_LEVEL=debug
```

### Options:
| Setting                         | Description                                                 | Default                 |
| ------------------------------- | ----------------------------------------------------------- | ----------------------- |
| logger.level                    | The default logging level                                   | "info"                  |
| logger.encoding                 | Logging format (console, json or stackdriver)               | "console"               |
| logger.color                    | Enable color in console mode                                | true                    |
| logger.dev_mode                 | Dump additional information as part of log messages         | true                    |
| logger.disable_caller           | Hide the caller source file and line number                 | false                   |
| logger.disable_stacktrace       | Hide a stacktrace on debug logs                             | true                    |
| ---                             | ---                                                         | ---                     |
| metrics.enabled                 | Enable metrics server                                       | true                    |
| metrics.host                    | Host/IP to listen on for metrics server                     | ""                      |
| metrics.port                    | Port to listen on for metrics server                        | 6060                    |
| profiler.enabled                | Enable go profiler on metrics server under /debug/pprof/    | true                    |
| pidfile                         | If set, creates a pidfile at the given path                 | ""                      |
| ---                             | ---                                                         | ---                     |
| server.host                     | The host address to listen on (blank=all addresses)         | ""                      |
| server.port                     | The port number to listen on                                | 8900                    |
| server.tls                      | Enable https/tls                                            | false                   |
| server.devcert                  | Generate a development cert                                 | false                   |
| server.certfile                 | The HTTPS/TLS server certificate                            | "server.crt"            |
| server.keyfile                  | The HTTPS/TLS server key file                               | "server.key"            |
| server.log.enabled              | Log server requests                                         | true                    |
| server.log.level                | Log level for server requests                               | "info                   |
| server.log.request_body         | Log the request body                                        | false                   |
| server.log.response_body        | Log the response body                                       | false                   |
| server.log.ignore_paths         | The endpoint prefixes to not log                            | []string{"/version"}    |
| server.cors.enabled             | Enable CORS middleware                                      | false                   |
| server.cors.allowed_origins     | CORS Allowed origins                                        | []string{"*"}           |
| server.cors.allowed_methods     | CORS Allowed methods                                        | []string{...everything} |
| server.cors.allowed_headers     | CORS Allowed headers                                        | []string{"*"}           |
| server.cors.allowed_credentials | CORS Allowed credentials                                    | false                   |
| server.cors.max_age             | CORS Max Age                                                | 300                     |
| server.metrics.enabled          | Enable metrics on server endpoints                          | true                    |
| server.metrics.ignore_paths     | The endpoint prefixes to not capture metrics on             | []string{"/version"}    |
| ---                             | ---                                                         | ---                     |
| database.username               | The database username                                       | "postgres"              |
| database.password               | The database password                                       | "password"              |
| database.host                   | Thos hostname for the database                              | "postgres"              |
| database.port                   | The port for the database                                   | 5432                    |
| database.database               | The database                                                | "gorestapi"             |
| database.auto_create            | Automatically create database                               | true                    |
| database.search_path            | Set the search path                                         | ""                      |
| database.sslmode                | The postgres sslmode to use                                 | "disable"               |
| database.sslcert                | The postgres sslcert file                                   | ""                      |
| database.sslkey                 | The postgres sslkey file                                    | ""                      |
| database.sslrootcert            | The postgres sslrootcert file                               | ""                      |
| database.retries                | How many times to try to reconnect to the database on start | 7                       |
| database.sleep_between_retries  | How long to sleep between retries                           | "7s"                    |
| database.max_connections        | How many pooled connections to have                         | 40                      |
| database.loq_queries            | Log queries (must set logging.level=debug)                  | false                   |
| database.wipe_confirm           | Wipe the database during start                              | false                   |


## Data Storage
Data is stored in a postgres database by default.

## Query Logic
Find requests `GET /api/things` and `GET /api/widgets` uses a url query parser to allow very complex logic including AND, OR and precedence operators. 
For the documentation on how to use this format see https://github.com/snowzach/queryp

## Swagger Documentation
When you run the API it has built in Swagger documentation available at `/api/api-docs/` (trailing slash required)
The documentation is automatically generated.

## TLS/HTTPS
You can enable https by setting the config option server.tls = true and pointing it to your keyfile and certfile.
To create a self-signed cert: `openssl req -new -newkey rsa:2048 -days 3650 -nodes -x509 -keyout server.key -out server.crt`
It also has the option to automatically generate a development cert every time it runs using the server.devcert option.

## Relocation
If you want to start with this as boilerplate for your project, you can clone this repo and use the `make relocate` option to rename the package.
`make relocate TARGET=github.com/myname/mycoolproject`
