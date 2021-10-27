package retailcrm

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type wrappedLogger struct {
	lastMessage string
}

func (w *wrappedLogger) Debugf(msg string, v ...interface{}) {
	w.lastMessage = fmt.Sprintf(msg, v...)
}

func TestDebugLoggerAdapter_Printf(t *testing.T) {
	wrapped := &wrappedLogger{}
	logger := DebugLoggerAdapter(wrapped)
	logger.Printf("Test message #%d", 1)

	assert.Equal(t, "Test message #1", wrapped.lastMessage)
}
