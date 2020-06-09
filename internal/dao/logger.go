package dao

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"regexp"
	"time"
	"unicode"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	xlog "xorm.io/xorm/log"

	"github.com/duchenhao/backend-demo/internal/model"
)

var (
	sqlRegexp                = regexp.MustCompile(`\?`)
	numericPlaceHolderRegexp = regexp.MustCompile(`\$\d+`)
)

type Logger struct {
	logger  *zap.Logger
	level   zapcore.Level
	showSQL bool
}

func newLogger(logger *zap.Logger) xlog.ContextLogger {
	return &Logger{
		logger: logger,
	}
}

func (l *Logger) BeforeSQL(xlog.LogContext) {}

func (l *Logger) AfterSQL(ctx xlog.LogContext) {
	logger := l.logger
	if reqCtx, ok := ctx.Ctx.(*model.ReqContext); ok {
		logger = logger.With(zap.String("request_id", reqCtx.RequestId))
	}

	if ctx.ExecuteTime > 0 {
		logger.Info("[SQL]",
			zap.String("sql", formatSql(ctx.SQL, ctx.Args)),
			zap.Int("time_ms", int(ctx.ExecuteTime/time.Millisecond)),
		)
	} else {
		logger.Info("[SQL]",
			zap.String("sql", formatSql(ctx.SQL, ctx.Args)),
		)
	}
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.logger.Debug(fmt.Sprintf(format, v...))
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.logger.Error(fmt.Sprintf(format, v...))
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.logger.Info(fmt.Sprintf(format, v...))
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	l.logger.Warn(fmt.Sprintf(format, v...))
}

func (l *Logger) Level() xlog.LogLevel {
	switch l.level {
	case zapcore.ErrorLevel:
		return xlog.LOG_ERR
	case zapcore.WarnLevel:
		return xlog.LOG_WARNING
	case zapcore.InfoLevel:
		return xlog.LOG_INFO
	case zapcore.DebugLevel:
		return xlog.LOG_DEBUG
	default:
		return xlog.LOG_ERR
	}
}

func (l *Logger) SetLevel(xlog.LogLevel) {}

func (l *Logger) ShowSQL(show ...bool) {
	if len(show) == 0 {
		l.showSQL = true
		return
	}
	l.showSQL = show[0]
}

func (l *Logger) IsShowSQL() bool {
	return l.showSQL
}

func formatSql(sql string, args []interface{}) string {
	formattedValues := make([]string, 0)
	var realSql string
	for _, value := range args {
		indirectValue := reflect.Indirect(reflect.ValueOf(value))
		if indirectValue.IsValid() {
			value = indirectValue.Interface()
			if t, ok := value.(time.Time); ok {
				if t.IsZero() {
					formattedValues = append(formattedValues, fmt.Sprintf("'%v'", "0000-00-00 00:00:00"))
				} else {
					formattedValues = append(formattedValues, fmt.Sprintf("'%v'", t.Format("2006-01-02 15:04:05")))
				}
			} else if b, ok := value.([]byte); ok {
				if str := string(b); isPrintable(str) {
					formattedValues = append(formattedValues, fmt.Sprintf("'%v'", str))
				} else {
					formattedValues = append(formattedValues, "'<binary>'")
				}
			} else if r, ok := value.(driver.Valuer); ok {
				if value, err := r.Value(); err == nil && value != nil {
					formattedValues = append(formattedValues, fmt.Sprintf("'%v'", value))
				} else {
					formattedValues = append(formattedValues, "NULL")
				}
			} else {
				switch value.(type) {
				case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, bool:
					formattedValues = append(formattedValues, fmt.Sprintf("%v", value))
				default:
					formattedValues = append(formattedValues, fmt.Sprintf("'%v'", value))
				}
			}
		} else {
			formattedValues = append(formattedValues, "NULL")
		}
	}

	// differentiate between $n placeholders or else treat like ?
	if numericPlaceHolderRegexp.MatchString(sql) {
		for index, value := range formattedValues {
			placeholder := fmt.Sprintf(`\$%d([^\d]|$)`, index+1)
			realSql = regexp.MustCompile(placeholder).ReplaceAllString(sql, value+"$1")
		}
	} else {
		formattedValuesLength := len(formattedValues)
		for index, value := range sqlRegexp.Split(sql, -1) {
			realSql += value
			if index < formattedValuesLength {
				realSql += formattedValues[index]
			}
		}
	}

	return realSql
}

func isPrintable(s string) bool {
	for _, r := range s {
		if !unicode.IsPrint(r) {
			return false
		}
	}
	return true
}
