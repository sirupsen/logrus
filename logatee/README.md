# Logatee
A timid logrus. Useful for testing code that uses logrus instances.

# Usage
Logatee will only help with code that injects a `*logrus.Logger`. To test with logatee, simply call `logatee.New(t)` to get a `*logrus.Logger` instance, and then pass it to your service under test. Logatee will group your log output with each test, and allow you to inspect what your service logged from your tests.

```golang
func TestMyService(t *testing.T) {
    logger := logatee.New(t)
    svc := NewMyService(logger)
    svc.DoSomeThings()

    // All logs from svc.DoSomeThings() will be written to
    // test output. Use go test -v to see logs for passing tests.

    logs := logatee.Logs(logger) 
    // logs has each logrus.Entry that was sent to logger. inspect it
    // to check that your error handling is doing the right thing. 

    // Before testing anything else that re-uses logger, reset its log 
    // count.
    logatee.Reset(logger)
}
```

# Test Suites
Logatee is designed to play nicely with test [suites](https://godoc.org/github.com/stretchr/testify/suite). Use `logatee.NewFunc` to ensure that tests are always grouped correctly with suite-based tests.

```golang
func (suite *TestSuite) SetupSuite() {
    // re-use logger throughout all tests in suite
    logger := logatee.NewFunc(suite.T)
    suite.logger = logger
}

func (suite *TestSuite) SetupTest() {
    // reset the logs for each test
    logatee.Reset(suite.logger)
}
```
