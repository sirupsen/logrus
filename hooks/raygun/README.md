# [logrus-raygun-hook](https://github.com/squirkle/logrus-raygun-hook)
A Raygun.io hook for logrus

## Usage

```go
import (
  log "github.com/Sirupsen/logrus"
  "github.com/squirkle/logrus-raygun-hook"
)

func init() {
  log.AddHook(raygun.NewHook("https://api.raygun.io/entries", "yourApiKey", true))
}
```

## Project status
Both this logrus hook and the [goraygun](https://github.com/SDITools/goraygun) library are **works in progress**.  Be aware of the possibility of upcoming improvements/API changes.
