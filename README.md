# ZapDog

[Zap](https://pkg.go.dev/go.uber.org/zap) custom core for logging to DataDog.

## Usage

```go
package main

import (
	"context"
	"github.com/sevco/zapdog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	datadog, _ := zapdog.NewDataDogLogger(context.TODO(), "DD-API-KEY", zapdog.Options{
		Host:     "",
		Source:   "",
		Service:  "",
		Hostname: "",
		Tags:     []string{},
	})
	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		datadog,
		zap.NewAtomicLevel(),
	))
	defer logger.Sync()
	
	logger.Info("constructed a logger")
}
```