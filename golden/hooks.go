package golden

type (
	PreSaveHookVars struct {
		// Current is a present content of golden file.
		Current []byte
	}

	// PreSaveHook is a hook calling before golden file write.
	PreSaveHook func(TestingT, []byte, PreSaveHookVars) []byte

	// Hooks is a compose of arbitrary hooks for golden file handle phases.
	Hooks struct {
		preSave PreSaveHook
	}
)

func NewHooksDefault() Hooks {
	return Hooks{
		preSave: NewPreSaveHookDefault(),
	}
}

func NewPreSaveHookDefault() PreSaveHook {
	return func(_ TestingT, b []byte, _ PreSaveHookVars) []byte {
		return b
	}
}
