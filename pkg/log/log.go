package log

import (
	"io"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
)

// New returns a logger implemented using the logrus package.
func New(wr io.Writer, level, dir string) *logrus.Logger {
	if wr == nil {
		wr = os.Stderr
	}

	lr := logrus.New()
	lr.Out = wr

	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		lvl = logrus.WarnLevel
		lr.Warnf("failed to parse log-level '%s', defaulting to 'warning'", level)
	}
	lr.SetLevel(lvl)
	lr.SetFormatter(getFormatter(false))

	if dir != "" {
		_ = ensureDir(dir)
		// nolint:gocritic
		fileHook, err := NewLogrusFileHook(dir+"/app.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
		if err == nil {
			lr.Hooks.Add(fileHook)
		} else {
			lr.Warnf("Failed to open logfile, using standard out: %v", err)
		}
	}

	return lr
}

func NewNullLogger() *logrus.Logger {
	lr := logrus.New()
	lr.SetOutput(ioutil.Discard)
	return lr
}

func NewFromWriter(w io.Writer) *logrus.Logger {
	lr := logrus.New()
	lr.SetOutput(w)
	lr.SetLevel(logrus.DebugLevel)
	return lr
}

// getFormatter returns the default log formatter.
func getFormatter(disableColors bool) *textFormatter {
	return &textFormatter{
		DisableColors:    disableColors,
		ForceFormatting:  true,
		ForceColors:      true,
		DisableTimestamp: false,
		FullTimestamp:    true,
		DisableSorting:   true,
		TimestampFormat:  "2006-01-02 15:04:05.000000",
		SpacePadding:     45,
	}
}

func ensureDir(path string) error {
	isExisting, err := exists(path)
	if err != nil {
		return err
	}
	if !isExisting {
		return os.MkdirAll(path, os.ModePerm)
	}
	return nil
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}
