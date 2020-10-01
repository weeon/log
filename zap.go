package log

import (
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var LEVELS = map[string]zapcore.Level{
	"debug": zap.DebugLevel,
	"info":  zap.InfoLevel,
	"warn":  zap.WarnLevel,
	"err":   zap.ErrorLevel,
}

type Logger struct {
	path  string
	rlog  *lumberjack.Logger
	log   *zap.Logger
	sugar *zap.SugaredLogger

	level zapcore.Level
	// pid   []interface{}

	rolling        bool
	lastRotateTime time.Time
	lastRotateRW   sync.Mutex

	opt *Opt
}

type Opt struct {
	disableFile        bool
	service, namespace string
	writers            []zapcore.WriteSyncer
}

type LoggerOpt func(opt *Opt)

func SetWritersOpt(ws ...zapcore.WriteSyncer) LoggerOpt {
	return func(opt *Opt) {
		opt.writers = ws
	}
}

func SetServiceOpt(service string) LoggerOpt {
	return func(opt *Opt) {
		opt.service = service
	}
}

func SetNamespaceOpt(ns string) LoggerOpt {
	return func(opt *Opt) {
		opt.namespace = ns
	}
}

func DisableOpt() LoggerOpt {
	return func(opt *Opt) {
		opt.disableFile = true
	}
}

func NewLogger(path string, level zapcore.Level, opts ...LoggerOpt) (*Logger, error) {
	out := new(Logger)
	out.rlog = new(lumberjack.Logger)

	opt := new(Opt)
	for _, fn := range opts {
		fn(opt)
	}

	out.opt = opt

	out.path = path
	out.lastRotateTime = time.Now()
	out.level = level
	// out.pid = []interface{}{env.Pid}

	// config lumberjack
	out.rlog.Filename = path
	out.rlog.MaxSize = 0x1000 * 5 // automatic rolling file on it increment than 2GB
	out.rlog.LocalTime = true
	out.rlog.Compress = true
	out.rlog.MaxBackups = 60 // reserve last 60 day logs

	// config encoder config
	ec := zap.NewProductionEncoderConfig()
	ec.EncodeLevel = zapcore.CapitalLevelEncoder
	ec.TimeKey = "@timestamp"
	ec.EncodeTime = zapcore.ISO8601TimeEncoder

	// config core
	c := zapcore.AddSync(out.rlog)

	core := zapcore.NewCore(zapcore.NewJSONEncoder(ec), c, out.level)

	if len(opt.writers) != 0 {
		cs := []zapcore.Core{core}
		for _, w := range opt.writers {
			cr := zapcore.NewCore(zapcore.NewJSONEncoder(ec), w, out.level)
			cs = append(cs, cr)
		}
		core = zapcore.NewTee(cs...)
	}

	out.log = zap.New(
		core,
		zap.AddCaller(),
		zap.AddCallerSkip(2),
	).
		With(zap.String("service", opt.service), zap.String("namespace", opt.namespace))

	// default enable daily rotate
	out.rolling = true

	out.sugar = out.log.Sugar()
	return out, nil
}

func (tlog *Logger) checkRotate() {
	if !tlog.rolling {
		return
	}

	n := time.Now()

	tlog.lastRotateRW.Lock()
	defer tlog.lastRotateRW.Unlock()

	last := tlog.lastRotateTime
	y, m, d := last.Year(), last.Month(), last.Day()
	if y != n.Year() || m != n.Month() || d != n.Day() {
		go tlog.rlog.Rotate()
		tlog.lastRotateTime = n
	}
}

func (tlog *Logger) EnableDailyFile() {
	tlog.rolling = true
}

func (tlog *Logger) GetZapLogger() *zap.Logger {
	return tlog.log
}

func (tlog *Logger) Error(v ...interface{}) {
	tlog.checkRotate()
	if !tlog.level.Enabled(zap.ErrorLevel) {
		return
	}

	defer tlog.log.Sync()
	tlog.sugar.Error(v)
}

func (tlog *Logger) Errorf(format string, v ...interface{}) {
	tlog.checkRotate()
	if !tlog.level.Enabled(zap.ErrorLevel) {
		return
	}

	defer tlog.log.Sync()
	tlog.sugar.Errorf(format, v...)
}

func (tlog *Logger) Errorw(format string, v ...interface{}) {
	tlog.checkRotate()
	if !tlog.level.Enabled(zap.ErrorLevel) {
		return
	}

	defer tlog.log.Sync()
	tlog.sugar.Errorw(format, v...)
}

func (tlog *Logger) Warn(format string, v ...interface{}) {
	tlog.checkRotate()
	if !tlog.level.Enabled(zap.WarnLevel) {
		return
	}

	defer tlog.log.Sync()
	tlog.sugar.Warnf(format, v...)
}

func (tlog *Logger) Warnw(format string, v ...interface{}) {
	tlog.checkRotate()
	if !tlog.level.Enabled(zap.WarnLevel) {
		return
	}

	defer tlog.log.Sync()
	tlog.sugar.Warnw(format, v...)
}

func (tlog *Logger) Info(v ...interface{}) {
	tlog.checkRotate()
	if !tlog.level.Enabled(zap.InfoLevel) {
		return
	}

	defer tlog.log.Sync()
	tlog.sugar.Info(v)
}

func (tlog *Logger) Infof(format string, v ...interface{}) {
	tlog.checkRotate()
	if !tlog.level.Enabled(zap.InfoLevel) {
		return
	}

	defer tlog.log.Sync()
	tlog.sugar.Infof(format, v...)
}

func (tlog *Logger) Infow(format string, v ...interface{}) {
	tlog.checkRotate()
	if !tlog.level.Enabled(zap.InfoLevel) {
		return
	}

	defer tlog.log.Sync()
	tlog.sugar.Infow(format, v...)
}

func (tlog *Logger) Debug(v ...interface{}) {
	tlog.checkRotate()
	if !tlog.level.Enabled(zap.DebugLevel) {
		return
	}

	defer tlog.log.Sync()
	tlog.sugar.Debug(v)
}

func (tlog *Logger) Debugf(format string, v ...interface{}) {
	tlog.checkRotate()
	if !tlog.level.Enabled(zap.DebugLevel) {
		return
	}

	defer tlog.log.Sync()
	tlog.sugar.Debugf(format, v...)
}

func (tlog *Logger) Debugw(format string, v ...interface{}) {
	tlog.checkRotate()
	if !tlog.level.Enabled(zap.DebugLevel) {
		return
	}

	defer tlog.log.Sync()
	tlog.sugar.Debugw(format, v...)
}

var _logger *Logger

func Stdout() {
	l, _ := NewLogger("stdout", zapcore.DebugLevel)
	SetDefault(l)
}

// active

func Error(v ...interface{}) {
	_logger.Error(v...)
}

func Errorf(format string, v ...interface{}) {
	_logger.Errorf(format, v...)
}

func Warn(format string, v ...interface{}) {
	_logger.Warn(format, v...)
}

func Info(v ...interface{}) {
	_logger.Info(v...)
}

func Infof(format string, v ...interface{}) {
	_logger.Infof(format, v...)
}

func Debug(v ...interface{}) {
	_logger.Debug(v...)
}

func Debugf(format string, v ...interface{}) {
	_logger.Debugf(format, v...)
}

func Errorw(format string, v ...interface{}) {
	_logger.Errorw(format, v...)
}

func Warnw(format string, v ...interface{}) {
	_logger.Warnw(format, v...)
}

func Infow(format string, v ...interface{}) {
	_logger.Infow(format, v...)
}

func Debugw(format string, v ...interface{}) {
	_logger.Debugw(format, v...)
}
