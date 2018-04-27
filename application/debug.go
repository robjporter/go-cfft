package application

import (
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Sirupsen/logrus"
)

func (a *Application) DEBUGOverrideLocalHXServer(server string) {
	a.addToLogDebug(a.Stats.GetCounter("tasks"), nil, "Overriding DB loaded information.")
	a.HX.Credentials.Url = server
	override = true
}

func (a *Application) addToLogDebug(task int64, fields map[string]interface{}, message string) {
	if fields == nil {
		fields = make(map[string]interface{})
	}
	pc, file, line, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc)

	dotName := filepath.Ext(fn.Name())
	fnName := strings.TrimLeft(dotName, ".") + "()"

	fields[".Task Number"] = task
	fields["File"] = filepath.Base(file)
	fields["Function"] = fnName
	fields["Line"] = line
	a.Logger.WithFields(logrus.Fields(fields)).Debug(message)
}

func (a *Application) addToLogWarning(task int64, fields map[string]interface{}, message string) {
	if fields == nil {
		fields = make(map[string]interface{})
	}
	pc, file, line, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc)

	dotName := filepath.Ext(fn.Name())
	fnName := strings.TrimLeft(dotName, ".") + "()"

	fields[".Task Number"] = task
	fields["File"] = filepath.Base(file)
	fields["Function"] = fnName
	fields["Line"] = line
	a.Logger.WithFields(logrus.Fields(fields)).Warning(message)
}

func (a *Application) addToLogFatal(task int64, fields map[string]interface{}, message string) {
	if fields == nil {
		fields = make(map[string]interface{})
	}
	pc, file, line, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc)

	dotName := filepath.Ext(fn.Name())
	fnName := strings.TrimLeft(dotName, ".") + "()"

	fields[".Task Number"] = task
	fields["File"] = filepath.Base(file)
	fields["Function"] = fnName
	fields["Line"] = line
	a.Logger.WithFields(logrus.Fields(fields)).Fatal(message)
}
