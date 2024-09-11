# Wingman

Wingman is a command-line tool written in Go that provides developers with an easy and efficient workflow for developing projects consisted of multiple services.

There are couple of features in wingman that differentiate it from other similar tools:

1. Can work with/watch more that one service at a time;
2. When it detects a file change it will not restart all services, just the ones affected by the code change;
3. You can inject environment variables to each service separately, or to all of them globally;
4. It packs a simple reverse-proxy that allows developers to "unify" services under a single port, based on a unique service "proxy handle";
5. When using the proxy you can also define static, storage or SPA routes, which help you run frontend apps on the same (proxy) port with your services.

## Installation
To install Wingman, just use the following command:

```sh
$ go install github.com/oblaxio/wingman@latest
```

## Usage
To use Wingman, navigate to your go project directory and run the following command:

```sh
$ wingman init
```
This will initialize a Wingman config file called `wingman.yaml` which you can later edit and customize according to your project needs.

> **Warning**
> Wingman relies on the go.mod, so your project needs to be created using [`go mod init`](https://go.dev/ref/mod#go-mod-init)

This is how a wingman config file looks like:

```yaml
version: 1 # The config version number. For now it's 1
module: oblax.io # The name of the go module
build_dir: bin # The build directory for the services
watchers:
  include_dirs: ["pkg", "services"] # Directories to be watched
  exclude_dirs: ["vendor", "modules"] # Directories to be excluded from watching
  include_files: ["*.go"] # Types of tiles to be watched
  exclude_files: ["test_*.go"] # Types of files not to be watched

env: # Environment variables available to all services at start
  GODEBUG: 'x509sha1=1' 
  OBLAX_REST_ERROR_MODE: 'development'

service_groups: 
  testing: ["obx-test-service-one", "obx-test-service-two"] # A service list

services: 
  obx-test-service-one: # An example of a GRPC/Protobuf service
    entrypoint: services/obx-test-service-one # Service location directory
    executable: obx-test-service-one # Name of the built service
    ldflags: # Build flags
      oblax.io/services/obx-test-service-one.Version: 'v0.1'
      oblax.io/services/obx-test-service-one.Build: 'dev-build'
      oblax.io/services/obx-test-service-one.Name: 'obx-test-service-one'
    env: # Service specific environment variables
      PORT: 10001

  obx-test-service-two: # An example of a REST service with reverse proxy
    entrypoint: services/obx-test-service-two
    executable: obx-test-service-two
    proxy_type: service 
    proxy_handle: /api/v1/test-service-two # When someone asks for this route
    proxy_address: 127.0.0.1 # ...proxy to this address
    proxy_port: 10002 # ...and this port
    ldflags:
      oblax.io/services/obx-test-service-two.Version: 'v0.1'
      oblax.io/services/obx-test-service-two.Build: 'dev-build'
      oblax.io/services/obx-test-service-two.Name: 'obx-test-service-two'
    env:
      PORT: 10002

  obx-test-web-storage: # An example of a static file handler
    proxy_type: static
    proxy_handle: /public/platform # Whenever someone asks for this route ...
    proxy_static_dir: services/obx-test-service-two/public # ...static files
```

The config file in the example above is suited for a project with the following structure:

```
.
├── bin
├── go.mod
├── go.sum
├── pkg
│   ├── libone
│   │   └── libone.go
│   ├── libtwo
│   │   └── libtwo.go
│   └── shared
│       └── shared.go
├── public
├── services
│   ├── obx-test-service-one
│   │   ├── handlers
│   │   │   └── handler.go
│   │   └── main.go
│   └── obx-test-service-two
│       ├── handlers
│       │   └── handler.go
│       ├── public
│       │   ├── image.jpg
│       │   ├── script.js
│       │   └── style.css
│       └── main.go
└── wingman.yaml
```
## Running wingman
After your wingman configuration is all set and done the next step is running it. For this you'll have to execute the following command inside your project directory:

```sh
$ wingman start
```
or if you want to run a group
```sh
$ wingman start testing
```

## Things you might find... interesting?
Wingman was created by [Beyond Basics](https://beyondbasics.co) as a tool to help with the development of the [Oblax](https://oblax.io) platform and we've been dogfooding it since it's inception. 
