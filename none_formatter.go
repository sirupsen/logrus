package logrus

// NoneFormatter has only the Format() member function, which returns a null formatter
type NoneFormatter struct { }

// Format for NoneFormatter returns nil,nil to ignore all by log.Fatal
func (f *NoneFormatter) Format(entry *Entry) ([]byte, error) {
	return nil, nil
}
