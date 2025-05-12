package reserv

import (
	"bytes"
	"log/slog"
	"math"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseFmt(t *testing.T) {
	handler, err := getHandlerFromFormat(FormatJSON, slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	})

	require.NoError(t, err)

	wantHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	})

	require.Equal(t, wantHandler, handler)

	handler, err = getHandlerFromFormat(FormatLogFmt, slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	})

	require.NoError(t, err)

	wantTextHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	})

	require.Equal(t, wantTextHandler, handler)

	handler, err = getHandlerFromFormat(FormatGCP, slog.HandlerOptions{
		Level:       slog.LevelInfo,
		AddSource:   true,
		ReplaceAttr: replaceSlogAttributesGCP,
	})

	require.NoError(t, err)

	require.IsType(t, &slog.JSONHandler{}, handler)

	// Create a buffer to capture the output
	var buf bytes.Buffer

	handlerWithBuffer := slog.NewJSONHandler(&buf, &slog.HandlerOptions{
		Level:       slog.LevelInfo,
		AddSource:   true,
		ReplaceAttr: replaceSlogAttributesGCP,
	})

	// Log a test message
	logger := slog.New(handlerWithBuffer)
	logger.Info("test message", "key", "value")

	// Check that the output contains "message" and "severity" keys. That are the must-have keys in GCP logging.
	output := buf.String()
	require.Contains(t, output, `"message":"test message"`)
	require.Contains(t, output, `"severity":"INFO"`)
	require.NotContains(t, output, `"level":"INFO"`)
}

func TestParseLevel(t *testing.T) {
	type args struct {
		lvl string
	}
	tests := []struct {
		name    string
		args    args
		want    slog.Level
		wantErr bool
	}{
		{
			name: "debug",
			args: args{
				lvl: LevelDebug,
			},
			want:    slog.LevelDebug,
			wantErr: false,
		},
		{
			name: "info",
			args: args{
				lvl: LevelInfo,
			},
			want:    slog.LevelInfo,
			wantErr: false,
		},
		{
			name: "warn",
			args: args{
				lvl: LevelWarn,
			},
			want:    slog.LevelWarn,
			wantErr: false,
		},
		{
			name: "error",
			args: args{
				lvl: LevelError,
			},
			want:    slog.LevelError,
			wantErr: false,
		},
		{
			name: "none",
			args: args{
				lvl: LevelNone,
			},
			want:    math.MaxInt,
			wantErr: false,
		},
		{
			name: "invalid",
			args: args{
				lvl: "invalid",
			},
			want:    -1,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseLevel(tt.args.lvl)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseLevel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.Equal(t, tt.want, got)
		})
	}
}
