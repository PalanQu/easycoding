package db

import (
	"context"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

type gormLogger struct {
	log                   *logrus.Logger
	SlowThreshold         time.Duration
	SourceField           string
	SkipErrRecordNotFound bool
}

func newGormLogger(l *logrus.Logger) *gormLogger {
	return &gormLogger{
		log:                   l,
		SkipErrRecordNotFound: false,
		SourceField:           "source",
	}
}

func (l *gormLogger) LogMode(gormlogger.LogLevel) gormlogger.Interface {
	return l
}

func (l *gormLogger) Info(ctx context.Context, s string, args ...interface{}) {
	l.log.WithContext(ctx).Infof(s, args)
}

func (l *gormLogger) Warn(ctx context.Context, s string, args ...interface{}) {
	l.log.WithContext(ctx).Warnf(s, args)
}

func (l *gormLogger) Error(ctx context.Context, s string, args ...interface{}) {
	l.log.WithContext(ctx).Errorf(s, args)
}

func (l *gormLogger) Trace(
	ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, _ := fc()
	fields := logrus.Fields{}
	if l.SourceField != "" {
		fields[l.SourceField] = utils.FileWithLineNum()
	}
	fields["duration_ms"] = elapsed.Milliseconds()
	fields["sql"] = sql
	fields["system"] = "mysql"
	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound) && l.SkipErrRecordNotFound) {
		fields[logrus.ErrorKey] = err
		l.log.WithContext(ctx).WithFields(fields).Error()
		return
	}

	if l.SlowThreshold != 0 && elapsed > l.SlowThreshold {
		l.log.WithContext(ctx).WithFields(fields).Warn()
		return
	}

	l.log.WithContext(ctx).WithFields(fields).Debug()
}
