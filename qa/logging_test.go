package qa

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {
	l := Logger()

	l.Info("test message")

	assert.Contains(t, l.String(), "test message")
}
