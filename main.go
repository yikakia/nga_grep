package main

import (
	"context"
	"log/slog"
	"runtime/debug"

	"github.com/yikakia/nga_grep/app/cmd"
	"github.com/yikakia/nga_grep/internal/observe"
)

func main() {
	observe.InitAll()

	defer func() {
		observe.FlushAll(context.Background())
	}()

	defer func() {
		if r := recover(); r != nil {
			slog.Error("recovered from panic", slog.Any("panic", r), slog.String("stack", string(debug.Stack())))
		}
	}()
	cmd.Execute()
}
