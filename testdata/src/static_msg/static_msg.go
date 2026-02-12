package static_msg

import (
	"fmt"
	"log/slog"
)

const constMsg = "msg"

var varMsg = "msg"

func _() {
	slog.Info("msg")
	slog.Info(constMsg)
	slog.Info(anotherConstMsg)
	slog.Info(varMsg)             // want `message should be a string literal or a constant`
	slog.Info(fmt.Sprintf("msg")) // want `message should be a string literal or a constant`

	slog.Info("msg" + "msg")
	slog.Info("msg" + constMsg)
	slog.Info("msg" + anotherConstMsg)
	slog.Info("msg" + varMsg)             // want `message should be a string literal or a constant`
	slog.Info("msg" + fmt.Sprintf("msg")) // want `message should be a string literal or a constant`
}
