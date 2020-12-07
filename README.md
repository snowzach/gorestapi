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
| logger.disable_caller           | Hide the caller source file and line number                 | false                   |
| logger.disable_stacktrace       | Hide a stacktrace on debug logs                             | true                    |
| ---                             | ---                                                         | ---                     |
| server.host                     | The host address to listen on (blank=all addresses)         | ""                      |
| server.port                     | The port number to listen on                                | 8900                    |
| server.tls                      | Enable https/tls                                            | false                   |
| server.devcert                  | Generate a development cert                                 | false                   |
| server.certfile                 | The HTTPS/TLS server certificate                            | "server.crt"            |
| server.keyfile                  | The HTTPS/TLS server key file                               | "server.key"            |
| server.log_requests             | Log API requests                                            | true                    |
| server.log_requests_body        | Log the requests with the body                              | false                   |
| server.log_disabled_http        | The endpoints to not log                                    | []string{"/version"}    |
| server.profiler_enabled         | Enable the profiler                                         | false                   |
| server.profiler_path            | Where should the profiler be available                      | "/debug"                |
| server.cors.allowed_origins     | CORS Allowed origins                                        | []string{"*"}           |
| server.cors.allowed_methods     | CORS Allowed methods                                        | []string{...everything} |
| server.cors.allowed_headers     | CORS Allowed headers                                        | []string{"*"}           |
| server.cors.allowed_credentials | CORS Allowed credentials                                    | false                   |
| server.cors.max_age             | CORS Max Age                                                | 300                     |
| ---                             | ---                                                         | ---                     |
| database.type                   | The database type (supports postgres)                       | "postgres"              |
| database.username               | The database username                                       | "postgres"              |
| database.password               | The database password                                       | "password"              |
| database.host                   | Thos hostname for the database                              | "postgres"              |
| database.port                   | The port for the database                                   | 5432                    |
| database.database               | The database                                                | "gorestapi"             |
| database.sslmode                | The postgres sslmode to use                                 | "disable"               |
| database.retries                | How many times to try to reconnect to the database on start | 7                       |
| database.sleep_between_retriews | How long to sleep between retries                           | "7s"                    |
| database.max_connections        | How many pooled connections to have                         | 40                      |
| database.wipe_confirm           | Wipe the database during start                              | false                   |
| database.loq_queries            | Log queries (must set logging.level=debug)                  | false                   |
| ---                             | ---                                                         | ---                     |
| pidfile                         | Write a pidfile (only if specified)                         | ""                      |
| profiler.enabled                | Enable the debug pprof interface                            | "false"                 |
| profiler.host                   | The profiler host address to listen on                      | ""                      |
| profiler.port                   | The profiler port to listen on                              | "6060"                  |


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
