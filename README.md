# Logrus

Logrus is a simple, opinionated logging package for Go. It has three debugging
levels:

* `LevelDebug`: Debugging, usually turned off for deploys.
* `LevelInfo`: Info, useful for monitoring in production.
* `LevelWarning`: Warnings that should definitely be noted. These are sent to
  `airbrake`.
* `LevelFatal`: Fatal messages that causes the application to crash. These are
  sent to `airbrake`.

## Usage

The global logging level is set by: `logrus.Level = logrus.{LevelDebug,LevelWarning,LevelFatal}`.

Note that for `airbrake` to work, `airbrake.Endpoint` and `airbrake.ApiKey`
should be set.

There is a global logger, which new loggers inherit their settings from when
created (see example below), such as the place to redirect output. Logging can
be done with the global logging module:

```go
logrus.Debug("Something debugworthy happened: %s", importantStuff)
logrus.Info("Something infoworthy happened: %s", importantStuff)

logrus.Warning("Something bad happened: %s", importantStuff)
// Reports to Airbrake

logrus.Fatal("Something fatal happened: %s", importantStuff)
// Reports to Airbrake
// Then exits
```

Types are encouraged to include their own logging object. This allows to set a
context dependent prefix to know where a certain message is coming from, without
cluttering every single message with this.

```go
type Walrus struct {
  TuskSize uint64
  Sex      bool
  logger logrus.Logger
}

func NewWalrus(tuskSize uint64, sex bool) *Walrus {
  return &Walrus{
    TuskSize: tuskSize,
    Sex: bool,
    logger: logrus.NewLogger("Walrus"),
  }
}

func (walrus *Walrus) Mate(partner *Walrus) error {
  if walrus.Sex == partner.Sex {
    return errors.New("Incompatible mating partner.")
  }

  walrus.logger.Info("Walrus with tusk sizes %d and %d are mating!", walrus.TuskSize, partner.TuskSize)
  // Generates a logging message: <timestamp> [Info] [Walrus] Walrus with tusk sizes <int> and <int> are mating!

  // Walrus mating happens here

  return nil
}
```
