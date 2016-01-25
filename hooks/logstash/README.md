# Logstash Hooks for Logrus

## Usage

```go
import (
  "github.com/Sirupsen/logrus"
  logrus_logstash "github.com/Sirupsen/logrus/hooks/logstash"
)

func main() {
  log := logrus.New()
  hook, err := logrus_logstash.NewLogstashHook("tcp", "localhost:9999")
  if err == nil {
    log.Hooks.Add(hook)
  }
}
```

