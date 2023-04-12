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
go install github.com/oblax/wingman@latest
```

## Usage
To use Wingman, navigate to your go project directory and run the following command:

```sh
wingman init
```
This will initialize a Wingman config file called `wingman.yaml` which you can later edit and customize according to your project needs.

> **Warning**
> Wingman relies on the go.mod, so your project needs to be created using [`go mod init`](https://go.dev/ref/mod#go-mod-init)

This is how a wingman config file looks like:

```yaml
version: 1.0 
module: github.com/oblaxio/wingman # the module path as in your project's go.mod file
env: 
  MASTER_KEY: key123431ee
build_dir: bin # the directory where the built services reside
watchers:
  include_dirs: ["pkg", "services"]
  exclude_dirs: ["vendor", "modules"]  
  include_files: ["*.go"]
  exclude_files: ["test_*.go"]
proxy:
  enabled: true
  port: 8080
  address: 127.0.0.1
  api_prefix: api # needed to differentiate wheter it's an api request or a request for the static/frontend assets
  log_requests: true # whether to log the api requests in the terminal
  storage: # usually used to give access to a file-storage service (s3, minio, etc.)
    enabled: true
    prefix: storage # an endpoint prefix to distinct the sorage route from the api, static or SPA routes
    bucket: bucketname # name of the storage bucket
    service: minio
    address: 127.0.0.1
    port: 9000
  spa: # either SPA or static can be enabled because they both rely on the same routing parameters
    enabled: true
    address: 127.0.0.1
    port: 3000
  static:
    enabled: true
    dir: /public
    index: index.html
services:
  svc-one:
    entrypoint: services/svc-one # the path where the services main file is located
    executable: svc-one # name of the executable generated (usually the name of it's directory)
    proxy_handle: /api/v1/svc-one # the handle used by the proxy to identify to which service the request goes
    proxy_address: 127.0.0.1 # service address used by the proxy to redirect the request to
    proxy_port: 10001 # service port used by the proxy to redirect the request to
    env: 
      PORT: 10001
  svc-two:
    entrypoint: services/svc-two
    executable: svc-two
    proxy_handle: /api/v1/svc-two
    proxy_address: 127.0.0.1
    proxy_port: 10002
    env: 
      PORT: 10002
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
│   ├── svc-one
│   │   ├── handlers
│   │   │   └── handler.go
│   │   └── main.go
│   └── svc-two
│       ├── handlers
│       │   └── handler.go
│       └── main.go
└── wingman.yaml
```

## Things you might find interesting
Wingman was created to help with the development of the [Oblax](https://oblax.io) platform and we've been dogfooding it since it's inception. 