# go-openfga

[![Language: Go](https://img.shields.io/badge/lang-Go-6ad7e5.svg?style=flat-square&logo=go)](https://golang.org/)
[![Go Reference](https://pkg.go.dev/badge/go.linka.cloud/go-openfga.svg)](https://pkg.go.dev/go.linka.cloud/go-openfga)
[![Go Report Card](https://goreportcard.com/badge/go.linka.cloud/go-openfga)](https://goreportcard.com/report/go.linka.cloud/go-openfga)

go-openfga is a simple implementation of an easy to use (and currently partial) client for the [OpenFGA](https://github.com/openfga/openfga) gRPC API.

It also provides a simple way to run openfga in process.

**Project status: *alpha***

Not all planned features are completed.
The API, spec, status and other user facing objects are subject to change.
We do not support backward-compatibility for the alpha releases.

## Overview

### Import

```go
import openfga "go.linka.cloud/go-openfga"
```
