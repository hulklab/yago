package preinit

import (
	"github.com/hulklab/yago/coms/orm"
	"github.com/hulklab/yago/example/app/libs/trace"
	xormLog "xorm.io/xorm/log"
)

func init() {
	orm.RegisterCtxLogger(new(CustomCtxLogger))
}

type CustomCtxLogger struct {
	level   xormLog.LogLevel
	showSQL bool
	Ctx     *trace.Context
}

func (l *CustomCtxLogger) Debugf(format string, v ...interface{}) {
	l.Ctx.Logger().Debugf(format, v...)
}

func (l *CustomCtxLogger) Errorf(format string, v ...interface{}) {
	l.Ctx.Logger().Errorf(format, v...)
}

func (l *CustomCtxLogger) Infof(format string, v ...interface{}) {
	l.Ctx.Logger().Infof(format, v...)
}

func (l *CustomCtxLogger) Warnf(format string, v ...interface{}) {
	l.Ctx.Logger().Warnf(format, v...)
}

// BeforeSQL implements ContextLogger
func (l *CustomCtxLogger) BeforeSQL(ctx xormLog.LogContext) {
	l.Ctx = trace.NewWithCtx(ctx.Ctx)
}

// AfterSQL implements ContextLogger
func (l *CustomCtxLogger) AfterSQL(ctx xormLog.LogContext) {
	if ctx.ExecuteTime > 0 {
		l.Ctx.Logger().Infof("[SQL]%s %v - %v", ctx.SQL, ctx.Args, ctx.ExecuteTime)
	} else {
		l.Ctx.Logger().Infof("[SQL]%s %v", ctx.SQL, ctx.Args)
	}
}

func (l *CustomCtxLogger) Level() xormLog.LogLevel {
	return l.level
}

func (l *CustomCtxLogger) SetLevel(c xormLog.LogLevel) {
	l.level = c
}

func (l *CustomCtxLogger) ShowSQL(show ...bool) {
	if len(show) > 0 {
		l.showSQL = show[0]
	}
}

func (l *CustomCtxLogger) IsShowSQL() bool {
	return l.showSQL
}
