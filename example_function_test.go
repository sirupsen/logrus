package logrus_test

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"testing"
)

func TestLogger_LogFn(t *testing.T) {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.WarnLevel)

	log.InfoFn(func() []interface{} {
		fmt.Println("This is never run")
		return []interface{} {
			"Hello",
		}
	})

	log.ErrorFn(func() []interface{} {
		fmt.Println("This runs")
		return []interface{} {
			"Oopsi",
		}
	})
}
