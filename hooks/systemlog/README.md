# System log Hooks for Logrus <img src="http://i.imgur.com/hTeVwmJ.png" width="40" height="40" alt=":walrus:" class="emoji" title=":walrus:"/>

## Usage

```go
import (
  "log/syslog"
  "github.com/sirupsen/logrus"
  "github.com/sirupsen/logrus/hooks/systemlog"
)

func main() {
  log       := logrus.New()
  hook, err := systemlog.NewSystemlogHook("udp", "localhost:514", "")

  if err == nil {
    log.Hooks.Add(hook)
  }
}
```

If you want to connect to local syslog (Ex. "/dev/log" or "/var/run/syslog" or "/var/run/log"). Just assign empty string to the first two parameters of `NewSystemlogHook`. It should look like the following.

```go
import (
  "log/syslog"
  "github.com/sirupsen/logrus"
  "github.com/sirupsen/logrus/hooks/systemlog"
)

func main() {
  log       := logrus.New()
  hook, err := systemlog.NewSyslogHook("", "", "")

  if err == nil {
    log.Hooks.Add(hook)
  }
}
```

On Windows it connects to event log. You may use third parameter to specify your event source.
```go
import (
  "log/syslog"
  "github.com/sirupsen/logrus"
  "github.com/sirupsen/logrus/hooks/systemlog"
)

func main() {
  log       := logrus.New()
  hook, err := systemlog.NewSyslogHook("", "localhost", "MySource")

  if err == nil {
    log.Hooks.Add(hook)
  }
}
```
