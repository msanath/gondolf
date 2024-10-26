package ctxslog_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/msanath/gondolf/pkg/ctxslog"
	"github.com/stretchr/testify/require"
)

// captureStdout captures the stdout output for testing.
func captureStdout(f func()) string {
	// Keep backup of the real stdout
	old := os.Stderr
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Execute the function that will write to stdout
	f()

	// Close the write end of the pipe so we can read from it
	w.Close()
	os.Stdout = old // Restore original stdout

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestSlog(t *testing.T) {
	testMessage := "Hello Gondolf"

	capturedOutput := captureStdout(func() {
		logger := slog.New(
			slog.NewJSONHandler(os.Stdout, ctxslog.NewCustomHandler(slog.LevelInfo)),
		).With(slog.Int("processId", os.Getpid()))
		ctx := ctxslog.NewContext(context.Background(), logger)

		// Fetch the logger from the context and log a message.
		ctxlogger := ctxslog.FromContext(ctx)
		ctxlogger.Info(testMessage)
	})

	type logMessage struct {
		Time      string `json:"time"`
		Level     string `json:"level"`
		Msg       string `json:"msg"`
		Source    string `json:"source"`
		ProcessID int    `json:"processId"`
	}

	var message logMessage
	err := json.Unmarshal([]byte(capturedOutput), &message)
	require.NoError(t, err)
	require.Equal(t, message.Level, "INFO")
	require.Equal(t, message.Msg, testMessage)
	require.Equal(t, message.Source, "ctxslog/new_test.go:47") // This should match the line number above. If the line number changes, this test will fail.
	require.Equal(t, message.ProcessID, os.Getpid())

	// Parse the time string
	parsedTime, err := time.Parse(time.RFC3339Nano, message.Time)
	require.NoError(t, err)
	require.WithinDuration(t, parsedTime, time.Now(), time.Second)
}
