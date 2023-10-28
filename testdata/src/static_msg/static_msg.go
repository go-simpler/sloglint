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

	slog.Log(ctx, slog.LevelInfo, "msg")
	slog.Debug("msg")
	slog.DebugContext(ctx, "msg")

	slog.Log(ctx, slog.LevelInfo, constMsg)
	slog.Debug(constMsg)
	slog.DebugContext(ctx, constMsg)

	slog.Log(ctx, slog.LevelInfo, fmt.Sprintf("msg")) // want `messages should be string literals or constants`
	slog.Debug(fmt.Sprintf("msg"))                    // want `messages should be string literals or constants`
	slog.DebugContext(ctx, fmt.Sprintf("msg"))        // want `messages should be string literals or constants`

	slog.Log(ctx, slog.LevelInfo, varMsg) // want `messages should be string literals or constants`
	slog.Debug(varMsg)                    // want `messages should be string literals or constants`
	slog.DebugContext(ctx, varMsg)        // want `messages should be string literals or constants`
}
