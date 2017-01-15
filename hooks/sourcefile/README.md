# Source File Hooks for Logrus <img src="http://i.imgur.com/hTeVwmJ.png" width="40" height="40" alt=":walrus:" class="emoji" title=":walrus:"/>

Add `source_file` field to entries with predefined log levels. Format of the field is `<package>/<sourcefile.go>:<line_number>`

## Usage

```go
package main

import (
  "github.com/Sirupsen/logrus"
  "github.com/Sirupsen/logrus/hooks/sourcefile"
)

func main() {
  log := logrus.New()
  log.Hooks.Add(&logrus_sourcefile.SourceFileHook{LogLevel: logrus.InfoLevel})

  log.Info("Hello World")
}
```
