package static_msg

import (
	"context"
	"fmt"
	"log/slog"
)

const constMsg = "msg"

var varMsg = "msg"

func tests() {
	ctx := context.Background()

	slog.Info("msg")
	slog.InfoContext(ctx, "msg")
	slog.Log(ctx, slog.LevelInfo, "msg")
	slog.With("key", "value").Info("msg")

	slog.Info(constMsg)
	slog.InfoContext(ctx, constMsg)
	slog.Log(ctx, slog.LevelInfo, constMsg)
	slog.With("key", "value").Info(constMsg)

	slog.Info(varMsg)                      // want `message should be a string literal or a constant`
	slog.InfoContext(ctx, varMsg)          // want `message should be a string literal or a constant`
	slog.Log(ctx, slog.LevelInfo, varMsg)  // want `message should be a string literal or a constant`
	slog.With("key", "value").Info(varMsg) // want `message should be a string literal or a constant`

	slog.Info(fmt.Sprintf("msg"))                      // want `message should be a string literal or a constant`
	slog.InfoContext(ctx, fmt.Sprintf("msg"))          // want `message should be a string literal or a constant`
	slog.Log(ctx, slog.LevelInfo, fmt.Sprintf("msg"))  // want `message should be a string literal or a constant`
	slog.With("key", "value").Info(fmt.Sprintf("msg")) // want `message should be a string literal or a constant`

	// binary expressions:
	slog.Info("msg" + "msg")
	slog.Info("msg" + "msg" + "msg")
	slog.Info("msg" + constMsg)
	slog.Info("msg" + varMsg)                                 // want `message should be a string literal or a constant`
	slog.Info("msg" + fmt.Sprintf("msg"))                     // want `message should be a string literal or a constant`
	slog.Info("msg" + constMsg + varMsg + fmt.Sprintf("msg")) // want `message should be a string literal or a constant`
}

func issue92() {
	slog.Info(constMsg)
	slog.Info(anotherConstMsg)
}
