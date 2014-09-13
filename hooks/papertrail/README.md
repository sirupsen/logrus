# Papertrail Hook for Logrus

## Usage

You can find your Papertrail UDP port on your [Papertrail account page](https://papertrailapp.com/account/destinations).

```go
import (
  "log/syslog"
  "github.com/Sirupsen/logrus"
  "github.com/Sirupsen/logrus/hooks/papertrail"
)

func main() {
  log       := logrus.New()
  hook, err := logrus_papertrail.NewPapertrailHook("logs.papertrailapp.com", YOUR_PAPERTRAIL_UDP_PORT, YOUR_APP_NAME)

  if err == nil {
    log.Hooks.Add(hook)
  }
}
```
