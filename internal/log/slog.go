package log

import (
	"context"
	"io"
	"os"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/exp/slog"
	"gopkg.in/natefinch/lumberjack.v2"
)

func SetSlog(filepath string) {
	w := newCombinedWriter(filepath)
	opts := slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	handler := newTextHandler(w, &opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)
	gin.DefaultWriter = w
}

// Combined Writer
type combinedWriter struct {
	stdout *os.File
	logger *lumberjack.Logger
}

func newCombinedWriter(filepath string) *combinedWriter {
	logger := &lumberjack.Logger{
		Filename:  filepath,
		MaxSize:   1, // MB
		LocalTime: true,
	}
	return &combinedWriter{
		stdout: os.Stdout,
		logger: logger,
	}
}

func (w *combinedWriter) Write(p []byte) (n int, err error) {
	_, _ = w.stdout.Write(p)
	return w.logger.Write(p)
}

// textHandler proxy
type textHandler struct {
	*slog.TextHandler
}

func newTextHandler(w io.Writer, opts *slog.HandlerOptions) *textHandler {
	return &textHandler{
		TextHandler: slog.NewTextHandler(w, opts),
	}
}

func (h *textHandler) Handle(ctx context.Context, r slog.Record) error {
	spanCtx := trace.SpanContextFromContext(ctx)
	r.Add(slog.String("traceId", spanCtx.TraceID().String()))
	r.Add(slog.String("spanId", spanCtx.SpanID().String()))
	return h.TextHandler.Handle(ctx, r)
}
