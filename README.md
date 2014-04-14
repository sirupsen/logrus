# Logrus <img src="http://i.imgur.com/hTeVwmJ.png" width="40" height="40" alt=":walrus:" class="emoji" title=":walrus:"/>&nbsp;[![Build Status](https://travis-ci.org/Sirupsen/logrus.svg?branch=master)](https://travis-ci.org/Sirupsen/logrus)

Logrus is a structured logger for Go (golang), completely API compatible with
the standard library logger. [Godoc][godoc].

Nicely color-coded in development (when a TTY is attached, otherwise just
plain text):

![Colored](http://i.imgur.com/PY7qMwd.png)

With `log.Formatter = new(logrus.JSONFormatter)`, for easy parsing by logstash
or Splunk:

```json
{"animal":"walrus","level":"info","msg":"A group of walrus emerges from the
ocean","size":10,"time":"2014-03-10 19:57:38.562264131 -0400 EDT"}

{"level":"warning","msg":"The group's number increased tremendously!",
"number":122,"omg":true,"time":"2014-03-10 19:57:38.562471297 -0400 EDT"}

{"animal":"walrus","level":"info","msg":"A giant walrus appears!",
"size":10,"time":"2014-03-10 19:57:38.562500591 -0400 EDT"}

{"animal":"walrus","level":"info","msg":"Tremendously sized cow enters the ocean.",
"size":9,"time":"2014-03-10 19:57:38.562527896 -0400 EDT"}

{"level":"fatal","msg":"The ice breaks!","number":100,"omg":true,
"time":"2014-03-10 19:57:38.562543128 -0400 EDT"}
```

With the default `log.Formatter = new(logrus.TextFormatter)` when a TTY is not
attached:

```text
time='2014-03-14 13:00:31.751756799 -0400 EDT' level='info' msg='A group of walrus emerges from the ocean' animal='walrus' size=10
time='2014-03-14 13:00:31.751994265 -0400 EDT' level='warning' msg='The group's number increased tremendously!' omg=true number=122
time='2014-03-14 13:00:31.752018319 -0400 EDT' level='info' msg='A giant walrus appears!' animal='walrus' size=10
time='2014-03-14 13:00:31.752034139 -0400 EDT' level='info' msg='Tremendously sized cow enters the ocean.' animal='walrus' size=9
time='2014-03-14 13:00:31.752048504 -0400 EDT' level='fatal' msg='The ice breaks!' omg=true number=100
```

#### Example

Note again that Logrus is API compatible with the standardlib logger, so if you
remove the `log` import and create a global `log` variable as below it will just
work.

```go
package main

import (
  "github.com/Sirupsen/logrus"
)

var log = logrus.New()

func init() {
  log.Formatter = new(logrus.JSONFormatter)
  log.Formatter = new(logrus.TextFormatter) // default
}

func main() {
  log.WithFields(logrus.Fields{
    "animal": "walrus",
    "size":   10,
  }).Info("A group of walrus emerges from the ocean")

  log.WithFields(logrus.Fields{
    "omg":    true,
    "number": 122,
  }).Warn("The group's number increased tremendously!")

  log.WithFields(logrus.Fields{
    "omg":    true,
    "number": 100,
  }).Fatal("The ice breaks!")
}
```

#### Fields

Logrus encourages careful, structured logging. It encourages the use of logging
fields instead of long, unparseable error messages. For example, instead of:
`log.Fatalf("Failed to send event %s to topic %s with key %d")`, you should log
the much more discoverable:

```go
log = logrus.New()

log.WithFields(logrus.Fields{
  "event": event,
  "topic": topic,
  "key": key
}).Fatal("Failed to send event")
```

We've found this API forces you to think about logging in a way that produces
much more useful logging messages. We've been in countless situations where just
a single added field to a log statement that was already there would've saved us
hours. The `WithFields` call is optional.

In general, with Logrus using any of the `printf`-family functions should be
seen as a hint you want to add a field, however, you can still use the
`printf`-family functions with Logrus.

#### Hooks

You can add hooks for logging levels. For example to send errors to an exception
tracking service on `Error`, `Fatal` and `Panic` or info to StatsD.

```go
log = logrus.New()
log.Hooks.Add(new(AirbrakeHook))

type AirbrakeHook struct{}

// `Fire()` takes the entry that the hook is fired for. `entry.Data[]` contains
// the fields for the entry. See the Fields section of the README.
func (hook *AirbrakeHook) Fire(entry *logrus.Entry) error {
  err := airbrake.Notify(entry.Data["error"].(error))
  if err != nil {
    log.WithFields(logrus.Fields{
      "source":   "airbrake",
      "endpoint": airbrake.Endpoint,
    }).Info("Failed to send error to Airbrake")
  }

  return nil
}

// `Levels()` returns a slice of `Levels` the hook is fired for.
func (hook *AirbrakeHook) Levels() []logrus.Level {
  return []logrus.Level{
    logrus.Error,
    logrus.Fatal,
    logrus.Panic,
  }
}
```

#### Level logging

Logrus has six logging levels: Debug, Info, Warning, Error, Fatal and Panic.

```go
log.Debug("Useful debugging information.")
log.Info("Something noteworthy happened!")
log.Warn("You should probably take a look at this.")
log.Error("Something failed but I'm not quitting.")
// Calls os.Exit(1) after logging
log.Fatal("Bye.")
// Calls panic() after logging
log.Panic("I'm bailing.")
```

You can set the logging level on a `Logger`, then it will only log entries with
that severity or anything above it:

```go
// Will log anything that is info or above (warn, error, fatal, panic). Default.
log.Level = logrus.Info
```

It may be useful to set `log.Level = logrus.Debug` in a debug or verbose
environment if your application has that.

#### Entries

Besides the fields added with `WithField` or `WithFields` some fields are
automatically added to all logging events:

1. `time`. The timestamp when the entry was created.
2. `msg`. The logging message passed to `{Info,Warn,Error,Fatal,Panic}` after
   the `AddFields` call. E.g. `Failed to send event.`
3. `level`. The logging level. E.g. `info`.

#### Environments

Logrus has no notion of environment.

If you wish for hooks and formatters to only be used in specific environments,
you should handle that yourself. For example, if your application has a global
variable `Environment`, which is a string representation of the environment you
could do:

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

This configuration is how `logrus` was intended to be used, but JSON in
production is mostly only useful if you do log aggregation with tools like
Splunk or Logstash.

#### Formatters

The built-in logging formatters are:

* `logrus.TextFormatter`. Logs the event in colors if stdout is a tty, otherwise
  without colors.
  * *Note:* to force colored output when there is no TTY, set the `ForceColors`
    field to `true`.
* `logrus.JSONFormatter`. Logs fields as JSON.

Third party logging formatters:

* [`zalgo`](https://github.com/aybabtme/logzalgo): invoking the P͉̫o̳̼̊w̖͈̰͎e̬͔̭͂r͚̼̹̲ ̫͓͉̳͈ō̠͕͖̚f̝͍̠ ͕̲̞͖͑Z̖̫̤̫ͪa͉̬͈̗l͖͎g̳̥o̰̥̅!̣͔̲̻͊̄ ̙̘̦̹̦.
* [`l2met`](https://github.com/meatballhat/logrus-formatters): log in [l2met](http://r.32k.io/l2met-introduction) format.

You can define your formatter by implementing the `Formatter` interface,
requiring a `Format` method. `Format` takes an `*Entry`. `entry.Data` is a
`Fields` type (`map[string]interface{}`) with all your fields as well as the
default ones (see Entries section above):

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

#### TODO

* Performance
* Default fields for an instance and inheritance
* Default available hooks (airbrake, statsd, dump cores)

[godoc]: https://godoc.org/github.com/Sirupsen/logrus
