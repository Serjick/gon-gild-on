package golden

// TestingT is a interface compartible with standart *testing.T.
type TestingT interface {
	// Name is a test fullname.
	Name() string
	// Logf to write log message into test output.
	Logf(msg string, args ...any)
}
