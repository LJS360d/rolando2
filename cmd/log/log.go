package log

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Log *zap.SugaredLogger
)

// ANSI escape codes for coloring console output
const (
	ColorGray   string = "\033[90m"
	ColorWhite  string = "\033[97m"
	ColorGreen  string = "\033[32m"
	ColorYellow string = "\033[33m"
	ColorRed    string = "\033[31m"
	ColorReset  string = "\033[0m"
)

type WebhookSyncer struct {
	url string
}

// Write sends the log entry to the webhook asynchronously.
func (ws *WebhookSyncer) Write(p []byte) (n int, err error) {
	go func() {
		payload := map[string]string{
			"content": removeANSICodes(string(p)),
		}
		body, _ := json.Marshal(payload)
		_, _ = http.Post(ws.url, "application/json", bytes.NewBuffer(body))
	}()
	return len(p), nil
}

func removeANSICodes(input string) string {
	// Regular expression to match ANSI escape codes
	re := regexp.MustCompile(`\x1b\[[0-9;]*m`)

	// Replace all matches with an empty string
	return re.ReplaceAllString(input, "")
}

// Sync is required to implement zapcore.WriteSyncer.
func (ws *WebhookSyncer) Sync() error {
	return nil
}

func init() {
	// Optionally set the webhook URL via an environment variable or pass as a flag
	webhookURL := os.Getenv("LOG_WEBHOOK")

	var writeSyncers []zapcore.WriteSyncer

	// Always log to stdout
	writeSyncers = append(writeSyncers, zapcore.AddSync(os.Stdout))

	// Add webhook syncer if URL is provided
	if webhookURL != "" {
		webhookSyncer := &WebhookSyncer{url: webhookURL}
		writeSyncers = append(writeSyncers, zapcore.AddSync(webhookSyncer))
	}

	// Custom time encoder for the format [dd/mm/YYYY HH:mm:ss]
	timeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(ColorGray + t.Format("[02/01/2006 15:04:05]") + ColorReset)
	}

	// Custom level encoder with colored output
	levelEncoder := func(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		var color string
		switch l {
		case zapcore.InfoLevel:
			color = ColorGreen
		case zapcore.WarnLevel:
			color = ColorYellow
		case zapcore.ErrorLevel:
		case zapcore.FatalLevel:
			color = ColorRed
		default:
			color = ColorWhite
		}
		enc.AppendString(string(color) + strings.ToUpper(l.String()) + string(ColorReset))
	}

	// Configure the encoder
	encoderCfg := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		EncodeTime:     timeEncoder,
		LevelKey:       "level",
		EncodeLevel:    levelEncoder,
		MessageKey:     "msg",
		CallerKey:      "caller",
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}

	// Use a console encoder with the custom settings
	encoder := zapcore.NewConsoleEncoder(encoderCfg)

	// Create a core that writes logs to all configured write syncers
	core := zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(writeSyncers...),
		zap.InfoLevel,
	)

	// Create the logger
	logger := zap.New(core)
	defer logger.Sync()
	Log = logger.Sugar()
}
