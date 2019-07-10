package logatee_test

import (
	"testing"

	"github.com/itzamna314/logrus/logatee"
	"github.com/sirupsen/logrus"
)

func TestLogatee(t *testing.T) {
	logger := logatee.New(t)
	logger.WithFields(logrus.Fields{
		"speed": "3 km/h",
		"depth": "6m",
	}).Trace("swimming")

	logger.WithFields(logrus.Fields{
		"animal":   "logatee",
		"diet":     "seaweed",
		"nickname": "sea cow",
	}).Info("hi guys! i'm the logatee!")

	logs := logatee.Logs(logger)
	if len(logs) != 2 {
		t.Fatalf("expected 2 logs, but found %d", len(logs))
	}

	ok := logs[0].Message == "swimming"
	ok = ok && logs[0].Level == logrus.TraceLevel
	ok = ok && len(logs[0].Data) == 2
	if ok {
		ok = ok && logs[0].Data["speed"] == "3 km/h"
		ok = ok && logs[0].Data["depth"] == "6m"
	}
	if !ok {
		t.Errorf("unexpected logs[0]: %#v", logs[0])
	}

	ok = logs[1].Message == "hi guys! i'm the logatee!"
	ok = ok && logs[1].Level == logrus.InfoLevel
	ok = ok && len(logs[1].Data) == 3
	if ok {
		ok = ok && logs[1].Data["animal"] == "logatee"
		ok = ok && logs[1].Data["diet"] == "seaweed"
		ok = ok && logs[1].Data["nickname"] == "sea cow"
	}
	if !ok {
		t.Errorf("unexpected logs[1]: %#v", logs[1])
	}

	logatee.Reset(logger)
	logs = logatee.Logs(logger)
	if len(logs) != 0 {
		t.Fatalf("expected 0 logs after reset, but found %d", len(logs))
	}

	logger.WithFields(logrus.Fields{
		"cause": "propeller",
	}).Error("ow")

	logs = logatee.Logs(logger)
	if len(logs) != 1 {
		t.Fatalf("expected 1 log after reset and log, but found %d", len(logs))
	}

	ok = logs[0].Message == "ow"
	ok = ok && logs[0].Level == logrus.ErrorLevel
	ok = ok && len(logs[0].Data) == 1
	if ok {
		ok = ok && logs[0].Data["cause"] == "propeller"
	}

	if !ok {
		t.Errorf("unexpected error log: %#v", logs[0])
	}
}
