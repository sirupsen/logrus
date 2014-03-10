package logrus

type Formatter interface {
	Format(*Entry) ([]byte, error)
}
