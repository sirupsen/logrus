package logrus

type RedactorFun func(serialized []byte) []byte

func defaultRedactor(serialized []byte) []byte {
	return serialized
}
