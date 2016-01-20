package logrus

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntryCaller(t *testing.T) {
	logger := New()
	entry := NewEntry(logger)

	entry.Debug("c", caller(0))
	assert.Equal(t, entry.Data["caller"], entry.Data["c"])

	entry.Debugf("c", caller(0))
	assert.Equal(t, entry.Data["caller"], entry.Data["c"])

	entry.Debugln("c", caller(0))
	assert.Equal(t, entry.Data["caller"], entry.Data["c"])

	entry.WithError(nil).Debug("c", caller(0))
	assert.Equal(t, entry.Data["caller"], entry.Data["c"])

	entry.WithField("first", 1).Debug("c", caller(0))
	assert.Equal(t, entry.Data["caller"], entry.Data["c"])

	entry.WithFields(Fields{"first": 1}).Debug("c", caller(0))
	assert.Equal(t, entry.Data["caller"], entry.Data["c"])
}
