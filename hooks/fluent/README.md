# Fluentd Hook for Logrus <img src="http://i.imgur.com/hTeVwmJ.png" width="40" height="40" alt=":walrus:" class="emoji" title=":walrus:"/>

## Usage

```go
import (
	"github.com/Sirupsen/logrus"
	"github.com/Sirupsen/logrus/hooks/fluent"
)

func main() {
	hook := fluent.NewHook("localhost", 24224)
	hook.SetLevels([]logrus.Level{
		logrus.PanicLevel,
		logrus.ErrorLevel,
	})

	logrus.AddHook(hook)
}
```


## Special fields

Some logrus fields have a special meaning in this hook.

- `tag` is used as a fluentd tag. (if `tag` is omitted, Entry.Message is used as a fluentd tag)
