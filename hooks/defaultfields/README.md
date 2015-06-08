# Default Fields Hooks for Logrus

Add default fields to be logged with every event. Useful when your logs are aggregated
and you want to identify where a particular event originated.

## Usage

```go
import (
	"log/syslog"

	"github.com/Sirupsen/logrus"
	"github.com/Sirupsen/logrus/hooks/defaultfields"
)

func main() {
	log  := logrus.New()
	hook := defaultfields.NewDefaultFields([]logrus.Level{
		logrus.DebugLevel,
		logrus.InfoLevel,
		logrus.WarnLevel,
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
	})
	hook.AddDefaultField("appId", "a1b4949a-0df9-11e5-8daa-5cf9dd6ef856")
	hook.AddDefaultField("appName", "MyApp")

    log.Hooks.Add(hook)
}
```