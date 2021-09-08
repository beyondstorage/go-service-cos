[![Build Status](https://github.com/beyondstorage/go-service-cos/workflows/Unit%20Test/badge.svg?branch=master)](https://github.com/beyondstorage/go-service-cos/actions?query=workflow%3A%22Unit+Test%22)
[![Integration Tests](https://teamcity.beyondstorage.io/app/rest/builds/buildType:(id:Services_Cos_IntegrationTests)/statusIcon)](https://teamcity.beyondstorage.io/buildConfiguration/Services_Cos_IntegrationTests)
[![License](https://img.shields.io/badge/license-apache%20v2-blue.svg)](https://github.com/Xuanwo/storage/blob/master/LICENSE)
[![](https://img.shields.io/matrix/beyondstorage@go-storage:matrix.org.svg?logo=matrix)](https://matrix.to/#/#beyondstorage@go-storage:matrix.org)

# go-services-cos

[COS(Cloud Object Storage)](https://cloud.tencent.com/product/cos) service support for [go-storage](https://github.com/beyondstorage/go-storage).

## Install

```go
go get github.com/beyondstorage/go-service-cos/v2
```

## Usage

```go
import (
	"log"

	_ "github.com/beyondstorage/go-service-cos/v2"
	"github.com/beyondstorage/go-storage/v4/services"
)

func main() {
	store, err := services.NewStoragerFromString("cos://bucket_name/path/to/workdir?credential=hmac:<account_name>:<account_key>")
	if err != nil {
		log.Fatal(err)
	}
	
	// Write data from io.Reader into hello.txt
	n, err := store.Write("hello.txt", r, length)
}
```

- See more examples in [go-storage-example](https://github.com/beyondstorage/go-storage-example).
- Read [more docs](https://beyondstorage.io/docs/go-storage/services/cos) about go-service-cos.
