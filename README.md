<img align="right" src="https://cloud.githubusercontent.com/assets/8041435/21744446/a5bb25b8-d4f4-11e6-9035-0a491d57f5ca.jpeg">
# Ghipu
**Ghipu** is an awesome tool to integrate your web service with Khipu's API.


## Overview

* Lightweight
* Extensible
* Zero dependencies

## TODO

**Ghipu** has the following server-2-khipu's functions to comunicate with its API.

- [x] `GET /banks`
- [x] `GET /payments`
- [x] `POST /Payments`
- [x] `GET /payments/{id}`
- [x] `DELETE /payments/{id}`
- [x] `POST /payments/{id}/refunds`
- [x] `POST /receivers`

## Instalation

```
go get github.com/mtavano/ghipu
```

## How to use

To use **Ghipu** you only need to get a new client and then start to gets through Khipu's API 

```go
import "github.com/mtavano/ghipu"

// ...

kc := ghipu.NewKhipuClient("secret", "receiver_id")

```

## Custom request

If you want to do a custom request to Khipu's API you can use **Ghipu** in the following way

```go
import (
    "http"

    "github.com/mtavano/ghipu"
)

// client definition ...

req, err := http.NewRequest("method", "urlStr", io.Reader)
if err != nil {
	// handle error properly
}

res, err := kc.DoRequest(req, "path")
if err != nil {
	// handle error properly
}
```
