# Base API Example

This API example is a framework for a REST API

## Compiling
This is designed as a go module aware program and thus requires go 1.11 or better
You can clone it anywhere, just run `make` inside the cloned directory to build

## Requirements
This does require a postgres database to be setup and reachable. It will attempt to create and migrate the database upon starting.

## Configuration
The configuration can be specified in a number of ways. By default you can create a json file and call it with the -c option
you can also specify environment variables that align with the config file values.

Example:
```json
{
	"logger": {
        "level": "debug"
	}
}
```
Can be set via an environment variable:
```
LOGGER_LEVEL=debug
```

### Options:
| Setting                        | Description                                                 | Default      |
|--------------------------------|-------------------------------------------------------------|--------------|
| logger.level                   | The default logging level                                   | "info"       |
| logger.encoding                | Logging format (console or json)                            | "console"    |
| logger.color                   | Enable color in console mode                                | true         |
| logger.disable_caller          | Hide the caller source file and line number                 | false        |
| logger.disable_stacktrace      | Hide a stacktrace on debug logs                             | true         |
| ---                            | ---                                                         | ---          |
| server.host                    | The host address to listen on (blank=all addresses)         | ""           |
| server.port                    | The port number to listen on                                | 8900         |
| server.tls                     | Enable https/tls                                            | false        |
| server.devcert                 | Generate a development cert                                 | false        |
| server.certfile                | The HTTPS/TLS server certificate                            | "server.crt" |
| server.keyfile                 | The HTTPS/TLS server key file                               | "server.key" |
| server.log_requests            | Log API requests                                            | true         |
| server.profiler_enabled        | Enable the profiler                                         | false        |
| server.profiler_path           | Where should the profiler be available                      | "/debug"     |
| ---                            | ---                                                         | ---          |
| storage.type                   | The database type (supports postgres)                       | "postgres"   |
| storage.username               | The database username                                       | "postgres"   |
| storage.password               | The database password                                       | "password"   |
| storage.host                   | Thos hostname for the database                              | "postgres"   |
| storage.port                   | The port for the database                                   | 5432         |
| storage.database               | The database                                                | "gorestapi"  |
| storage.sslmode                | The postgres sslmode to use                                 | "disable"    |
| storage.retries                | How many times to try to reconnect to the database on start | 5            |
| storage.sleep_between_retriews | How long to sleep between retries                           | "7s"         |
| storage.max_connections        | How many pooled connections to have                         | 80           |
| storage.wipe_confirm           | Wipe the database during start                              | false        |
| ---                            | ---                                                         | ---          |
| pidfile                        | Write a pidfile (only if specified)                         | ""           |
| profiler.enabled               | Enable the debug pprof interface                            | "false"      |
| profiler.host                  | The profiler host address to listen on                      | ""           |
| profiler.port                  | The profiler port to listen on                              | "6060"       |


## Data Storage
Data is stored in a postgres database

## TLS/HTTPS
You can enable https by setting the config option server.tls = true and pointing it to your keyfile and certfile.
To create a self-signed cert: `openssl req -new -newkey rsa:2048 -days 3650 -nodes -x509 -keyout server.key -out server.crt`
