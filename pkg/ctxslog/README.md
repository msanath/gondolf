# Slog with FluentBit compatible handler.
slog is a wrapper around the log/slog package. It provides a way to create a logger and add it to the context, 
and a way to retrieve the logger from the context.

The format of the logger stems from the proposed format for inlogs v2 - https://go/inlogsV2-proposed-format

# Using context for passing the logger

This pkg provides functionality to pass the logger using context. To do so,
1. Create the logger and get a context for it
```go
import (
  "context"
  ligoslog "golnkd.in/nimbus/nimbus-pkg/v2/slog"
)

func main() {
	log := ligoslog.New()
	ctx := ligoslog.NewContext(context.Background(), log)
  ...
}
```

2. Fetch the logger from the context and log
```go
import (
  "context"
  ligoslog "golnkd.in/nimbus/nimbus-pkg/v2/slog"
)

func something(ctx context.Context) {
  log := ligoslog.FromContext(ctx)
  log.Info("Hello Nimbus!)
}
