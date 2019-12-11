## qencode-api-go-client

## Install
```sh
go get github.com/mavvverick/qencode
```

## Usage

```go
import "github.com/mavvverick/qencode"
```


## Initialize

To get task token

```
package main
import (
	"context"
	"github.com/mavvverick/qencode"
)


func main() {
	client, err := qencode.NewClient(nil, accessKey)
    if err != nil {
		return client, err
	}
}
```