package diagnostic

import "fmt"

type Severity int

const (
	Info Severity = iota
	Warning
	Error
	Blocker
)

func (s Severity) String() string {
	switch s {
	case Info:
		return "INFO"
	case Warning:
		return "WARNING"
	case Error:
		return "ERROR"
	case Blocker:
		return "BLOCKER"
	default:
		return "UNKNOWN"
	}
}

// Diagnostic represents a single finding from an analyzer.
type Diagnostic struct {
	Path        string   `json:"path"`
	Severity    Severity `json:"severity"`
	Rule        string   `json:"rule"`
	Message     string   `json:"message"`
	Explanation string   `json:"explanation,omitempty"`
	CanAutoFix  bool     `json:"can_auto_fix,omitempty"`
}

func (d Diagnostic) String() string {
	return fmt.Sprintf("%-7s %s: %s", d.Severity, d.Path, d.Message)
}

// Set collects diagnostics and provides summary methods.
type Set struct {
	Items []Diagnostic
}

func (s *Set) Add(d Diagnostic) {
	s.Items = append(s.Items, d)
}

func (s *Set) AddAll(other *Set) {
	if other != nil {
		s.Items = append(s.Items, other.Items...)
	}
}

func (s *Set) Count(sev Severity) int {
	n := 0
	for _, d := range s.Items {
		if d.Severity == sev {
			n++
		}
	}
	return n
}

func (s *Set) HasBlockers() bool { return s.Count(Blocker) > 0 }
func (s *Set) HasErrors() bool   { return s.Count(Error) > 0 }

func (s *Set) Summary() string {
	return fmt.Sprintf("%d info, %d warnings, %d errors, %d blockers",
		s.Count(Info), s.Count(Warning), s.Count(Error), s.Count(Blocker))
}
