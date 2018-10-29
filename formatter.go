package logrus

import "time"

// Default key names for the default fields
const (
	defaultTimestampFormat = time.RFC3339
	FieldKeyMsg            = "msg"
	FieldKeyLevel          = "level"
	FieldKeyTime           = "time"
	FieldKeyLogrusError    = "logrus_error"
	FieldKeyFunc           = "func"
	FieldKeyFile           = "file"
)

// The Formatter interface is used to implement a custom Formatter. It takes an
// `Entry`. It exposes all the fields, including the default ones:
//
// * `entry.Data["msg"]`. The message passed from Info, Warn, Error ..
// * `entry.Data["time"]`. The timestamp.
// * `entry.Data["level"]. The level the entry was logged at.
//
// Any additional fields added with `WithField` or `WithFields` are also in
// `entry.Data`. Format is expected to return an array of bytes which are then
// logged to `logger.Out`.
type Formatter interface {
	Format(*Entry) ([]byte, error)
}

// This is to not silently overwrite `time`, `msg`, `func` and `level` fields when
// dumping it. If this code wasn't there doing:
//
//  logrus.WithField("level", 1).Info("hello")
//
// Would just silently drop the user provided level. Instead with this code
// it'll logged as:
//
//  {"level": "info", "fields.level": 1, "msg": "hello", "time": "..."}
//
// It's not exported because it's still using Data in an opinionated way. It's to
// avoid code duplication between the two default formatters.
func prefixFieldClashes(data Fields, fieldMap FieldMap, reportCaller bool) {
	reservedKeyNames := []fieldKey{
		FieldKeyTime,
		FieldKeyMsg,
		FieldKeyLevel,
		FieldKeyLogrusError,
	}
	// If reportCaller is not set, 'func' will not conflict.
	if reportCaller {
		reservedKeyNames = append(reservedKeyNames, FieldKeyFunc, FieldKeyFile)
	}
	for _, srcKeyName := range reservedKeyNames {
		destKeyName := fieldMap.resolve(srcKeyName)
		if itemValue, ok := data[destKeyName]; ok {
			data["fields."+destKeyName] = itemValue
			delete(data, destKeyName)
		}
	}
}
