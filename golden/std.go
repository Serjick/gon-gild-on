package golden

type tNamer interface {
	// Name returns the name of the running (sub-) test or benchmark.
	//
	// The name will include the name of the test along with the names of
	// any nested sub-tests. If two sibling sub-tests have the same name,
	// Name will append a suffix to guarantee the returned name is unique.
	Name() string
}

type tLogger interface {
	// Logf formats its arguments according to the format, analogous to Printf, and
	// records the text in the error log. A final newline is added if not provided. For
	// tests, the text will be printed only if the test fails or the -test.v flag is
	// set. For benchmarks, the text is always printed to avoid having performance
	// depend on the value of the -test.v flag.
	Logf(msg string, args ...any)
}

type tFailer interface {
	// Errorf is equivalent to Logf followed by Fail.
	// Fail marks the function as having failed but continues execution.
	Errorf(format string, args ...any)

	// FailNow marks the function as having failed and stops its execution
	// by calling runtime.Goexit (which then runs all deferred calls in the
	// current goroutine).
	// Execution will continue at the next test or benchmark.
	// FailNow must be called from the goroutine running the
	// test or benchmark function, not from other goroutines
	// created during the test. Calling FailNow does not stop
	// those other goroutines.
	FailNow()
}

// TestingT is a interface compartible with standart *testing.T.
type TestingT interface {
	tNamer
	tLogger
	tFailer
}
