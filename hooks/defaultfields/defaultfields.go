package defaultfields

import (
	"github.com/Sirupsen/logrus"
)

type DefaultFieldsHook struct {
	defaultFields map[string]interface{}
	levels        []logrus.Level
}

func (dfh DefaultFieldsHook) AddDefaultField(key string, val interface{}) {
	dfh.defaultFields[key] = val
}

func (dfh DefaultFieldsHook) Levels() []logrus.Level {
	return dfh.levels
}

func (dfh DefaultFieldsHook) Fire(entry *logrus.Entry) error {
	for key, val := range dfh.defaultFields {
		entry.Data[key] = val
	}
	return nil
}

func NewDefaultFieldsHook(levels []logrus.Level) *DefaultFieldsHook {
	fields := map[string]interface{}{}
	return &DefaultFieldsHook{
		defaultFields: fields,
		levels:        levels,
	}

}
