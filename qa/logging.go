package qa

import (
	"bytes"
	"github.com/uber-go/zap"
	"sync"
)

type buffer struct {
	bytes.Buffer
	mu sync.Mutex
}

type BufferLogger struct {
	zap.Logger
	buf buffer
}

func (b *buffer) Write(p []byte) (n int, err error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.Buffer.Write(p)
}

func (b *buffer) Sync() error {
	return nil
}

func (b *buffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.Buffer.String()
}

// Logger creates new test logger
func Logger() *BufferLogger {

	dynamicLevel := zap.DynamicLevel()
	dynamicLevel.SetLevel(zap.DebugLevel)

	var b = &BufferLogger{}

	logger := zap.New(
		zap.NewJSONEncoder(),
		zap.AddCaller(),
		zap.Output(&b.buf),
		dynamicLevel,
	)

	b.Logger = logger
	return b
}

func (bl *BufferLogger) String() string {
	return bl.buf.String()
}
