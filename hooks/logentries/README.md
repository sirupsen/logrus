# Logentries Hook for Logrus <img src="http://i.imgur.com/hTeVwmJ.png" width="40" height="40" alt=":walrus:" class="emoji" title=":walrus:" />

[Logentries](https://logentries.com/) provides Log Management and Analytics in Real-time.

This implementation uses the [token-based logging](https://logentries.com/doc/input-token/).

## Usage

For `LE_TOKEN`, substitute the log token for the log set you want to send log entries to

For `APP_NAME`, substitute a short string that will readily identify your application or service in the logs.

```go
import (
  log "github.com/Sirupsen/logrus"
  "github.com/Sirupsen/logrus/hooks/logentries"
)

func main() {
	log.AddHook(logrus_logentries.NewLogentriesHook(APP_NAME, LE_TOKEN))
}
```
