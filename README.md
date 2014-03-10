# Logrus

Logrus is a simple, opinionated structured logging package for Go which is
completely API compatible with the standard library logger.

#### Fields

Logrus encourages careful, structured logging. It encourages the use of logging
fields instead of long, unparseable error messages. For example, instead of:
`log.Fatalf("Failed to send event %s to topic %s with key %d")`, you should log
the much more discoverable:

```go
log = logrus.New()

log.WithFields(&logrus.Fields{
  "event": event,
  "topic": topic,
  "key": key
}).Fatal("Failed to send event")
```

We've found this API forces you to think about logging in a way that produces
much more useful logging messages. The `WithFields` call is optional.

In general, with Logrus using any of the `printf`-family functions should be
seen as a hint you want to add a field, however, you can still use the
`printf`-family functions with Logrus.

#### Hooks

You can add hooks for logging levels. For example to send errors to an exception
tracking service on `LevelError`, `LevelFatal` and `LevelPanic`.

```go
log = logrus.New()

type AirbrakeHook struct {}

// `Fire()` takes the entry that the hook is fired for. `entry.Data[]` contains
// the fields for the entry. See the Fields section of the README.
func (hook *AirbrakeHook) Fire(entry *Entry) (error) {
  err := airbrake.Notify(errors.New(entry.String()))
  if err != nil {
    log.WithFields(logrus.Fields{
      "source": "airbrake",
      "endpoint": airbrake.Endpoint,
    }).Info("Failed to send error to Airbrake")
  }

  return nil
}

// `Levels()` returns a slice of `LevelTypes` the hook is fired for.
func (hook *AirbrakeHook) Levels() []logrus.LevelType {
  return []logrus.LevelType{
    logrus.LevelError,
    logrus.LevelFatal,
    logrus.LevelPanic
  }
}
```

#### Level logging

Logrus has six levels: Debug, Info, Warning, Error, Fatal and Panic.

```go
log.Debug("Useful debugging information.")
log.Info("Something noteworthy happened!")
log.Warn("You should probably take a look at this.")
log.Error("Something failed but I'm not quitting.")
log.Fatal("Bye.")
log.Panic("I'm bailing.")
```

You can set the logging level:

```go
// Will log anything that is info or above, default.
logrus.Level = LevelInfo
```

#### Entries

Besides the fields added with `WithField` or `WithFields` some fields are
automatically added to all logging events:

1. `time`. The timestamp when the entry was created.
2. `msg`. The logging message passed to `{Info,Warn,Error,Fatal,Panic}` after
   the `AddFields` call. E.g. `Failed to send event.`
3. `level`. The logging level. E.g. `info`.

#### Environments

Logrus has no notion of environment. If you wish for hooks and formatters to
only be used in specific environments, you should handle that yourself. For
example, if your application has a global variable `Environment`, which is a
string representation of the environment you could do:

```go
init() {
  // do something here to set environment depending on an environment variable
  // or command-line flag
  log := logrus.New()

  if Environment == "production" {
    log.Formatter = new(logrus.JSONFormatter)
  } else {
    // The TextFormatter is default, you don't actually have to do this.
    log.Formatter = new(logrus.TextFormatter)
  }
}
```

#### Formats

The built in logging formatters are:

* `logrus.TextFormatter`. Logs the event in colors if stdout is a tty, otherwise
  without colors.
* `logrus.JSONFormatter`. Logs fields as JSON.

You can define your formatter taking an entry. `entry.Data` is a `Fields` type
which is a `map[string]interface{}` with all your fields as well as the default
ones (see Entries above):

```go
type MyJSONFormatter struct {
}

log.Formatter = new(MyJSONFormatter)

func (f *JSONFormatter) Format(entry *Entry) ([]byte, error) {
  serialized, err := json.Marshal(entry.Data)
    if err != nil {
      return nil, fmt.Errorf("Failed to marshal fields to JSON, %v", err)
    }
  return append(serialized, '\n'), nil
}
```
