package exitcode

const (
	OK             = 0 // success, no issues
	ErrUsage       = 1 // operational failure or invalid usage
	ErrPolicy      = 2 // policy file invalid
	ErrDiagnostics = 3 // diagnostics at error threshold
	ErrBlockers    = 4 // diagnostics at blocker threshold
)
