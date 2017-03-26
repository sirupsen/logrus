# Mail Hook for Logrus 

In some deployments, you'll want to report errors by email. If you add this hook, an email will send for the following levels:

* Error
* Fatal
* Panic

The subject is of the form `APPLICATION_NAME - LEVEL` and the body contains the timestamp and the message.

## Usage

For `APPLICATION_NAME`, substitute a short string that will identify your application or service in the logs.

```go
import (
  "log/syslog"
  "github.com/Sirupsen/logrus"
  "github.com/Sirupsen/logrus/hooks/mail"
)

func main() {
  log       := logrus.New()
  // if you do not need authentication for your smtp host
  hook, err := logrus_mail.NewMailHook("APPLICATION_NAME", "HOST", PORT, "FROM", "TO")

  if err == nil {
    log.Hooks.Add(hook)
  }
}
```

Example with authentication:
```go
  // if you need authentication for your smtp host
  hook, err := logrus_mail.NewMailAuthHook("APPLICATION_NAME", "HOST", PORT, "FROM", "TO", "USERNAME", "PASSWORD)
```

If you want to send mails with gmail:
```go
 hook, err := logrus_mail.NewMailAuthHook("testapp", "smtp.gmail.com", 587, "user.name@gmail.com", "user.name@gmail.com", "user.name", "password")
```
