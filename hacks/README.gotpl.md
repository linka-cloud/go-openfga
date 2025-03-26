# go-openfga & protoc-gen-go-openfga

[![Language: Go](https://img.shields.io/badge/lang-Go-6ad7e5.svg?style=flat-square&logo=go)](https://golang.org/)
[![Go Reference](https://pkg.go.dev/badge/go.linka.cloud/go-openfga.svg)](https://pkg.go.dev/go.linka.cloud/go-openfga)
[![Go Report Card](https://goreportcard.com/badge/go.linka.cloud/go-openfga)](https://goreportcard.com/report/go.linka.cloud/go-openfga)

**Project status: *alpha***

Not all planned features are completed.
The API, spec, status and other user facing objects are subject to change.
We do not support backward-compatibility for the alpha releases.


## `go-openfga`

### Overview

`go-openfga` is a simple implementation of an easy to use (and currently partial) client for the [OpenFGA](https://github.com/openfga/openfga) gRPC API.

It also provides a simple way to run openfga in process.

#### Import

```go
import openfga "go.linka.cloud/go-openfga"
```

## `protoc-gen-go-openfga`

### Overview

`protoc-gen-go-openfga` is a protoc plugin that generates openfga schema and go code to register access checks into the interceptor.


### Installation

```shell
go install go.linka.cloud/go-openfga/cmd/protoc-gen-go-openfga
```

### Usage

Use the plugin as any other protoc plugins.

### Generated code

For a given base [`openfga` module](example/base.fga):

```openfga
{{ file "example/base.fga" }}
```

For a given [`resource.proto`](example/pb/resource.proto):

```protobuf
{{ file "example/pb/resource.proto" }}
```

The following [`resource.fga`](example/pb/resource.fga) `openfga` module will be generated:

```openfga
{{ file "example/pb/resource.fga" }}
```

And following code will be generated:

```go
{{ file "example/pb/resource.pb.fga.go" }}
```

### Usage

See the [example](example) directory for complete example.

```go
{{ file "example/main.go" }}
```
