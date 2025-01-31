package retailcrm

import (
	"fmt"
	"strings"
	"time"
)

type SystemTime time.Time

const systemTimeLayout = "2006-01-02 15:04:05"

// UnmarshalJSON parses time.Time from system format.
func (st *SystemTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), `"`)
	nt, err := time.Parse(systemTimeLayout, s)
	*st = SystemTime(nt)
	return
}

// MarshalJSON will marshal time.Time to system format.
func (st SystemTime) MarshalJSON() ([]byte, error) {
	return []byte(st.String()), nil
}

// String returns the time in the custom format.
func (st *SystemTime) String() string {
	t := time.Time(*st)
	return fmt.Sprintf("%q", t.Format(systemTimeLayout))
}
