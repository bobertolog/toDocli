package repository

type FakeLogger struct{}

func NewFakeLogger() *FakeLogger {
	return &FakeLogger{}
}

func (l *FakeLogger) Log(msg string) error {
	return nil // ничего не делает — заглушка
}
