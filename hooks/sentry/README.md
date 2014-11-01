# Sentry Hook for Logrus <img src="http://i.imgur.com/hTeVwmJ.png" width="40" height="40" alt=":walrus:" class="emoji" title=":walrus:" />

[Sentry](https://getsentry.com) provides both self-hosted and hostes solutions for exception tracking.
Both client and server are [open source](https://github.com/getsentry/sentry).

## Usage

Every sentry application defined on the server gets a different [DNS](https://www.getsentry.com/docs/). In the example below replace `YOUR_DSN` with the one created for your application.


```go
import (
  "github.com/Sirupsen/logrus"
  "github.com/Sirupsen/logrus/hooks/sentry"
)

func main() {
  log       := logrus.New()
  hook, err := logrus_sentry.NewSentryHook(YOUR_DSN, []logrus.Level{
    logrus.PanicLevel,
    logrus.FatalLevel,
    logrus.ErrorLevel,
  })

  if err == nil {
    log.Hooks.Add(hook)
  }
}
```
