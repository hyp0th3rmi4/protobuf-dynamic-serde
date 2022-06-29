package logging

import "go.uber.org/zap"

var Log *zap.Logger
var SugarLog *zap.SugaredLogger

func init() {

	Log, _ = zap.NewProduction()
	defer Log.Sync() // flushes buffer, if any
	SugarLog = Log.Sugar()
}
