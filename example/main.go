package main

import (
	"context"
	"flag"
	camect_go "github.com/kaiiorg/camect-go"
	"log/slog"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"time"
)

var (
	logLevels = map[string]slog.Level{
		"DEBUG": slog.LevelDebug,
		"INFO":  slog.LevelInfo,
		"WARN":  slog.LevelWarn,
		"ERROR": slog.LevelError,
	}
)

var (
	json  = flag.Bool("json", false, "print JSON logs")
	level = flag.String("level", "INFO", "log level: DEBUG, INFO, WARN, ERROR")

	ip       = flag.String("ip", "0.0.0.0", "ip address of hub")
	username = flag.String("username", "admin", "username of hub local admin")
	password = flag.String("password", "this isn't a real password, provide your own", "password of hub local admin")
)

func main() {
	flag.Parse()
	initSlog()

	slog.Info("camect-go", "goVersion", runtime.Version())

	hub := camect_go.New(
		*ip,
		*username,
		*password,
		nil, // logger, nil means use global default
	)

	info, err := hub.Info()
	if err != nil {
		panic(err)
	}

	slog.Info(
		"got hub info",
		"hubName", info.Name,
		"hubId", info.Id,
	)

	ctx, ctxCancel := context.WithCancel(context.Background())
	eventsChan, err := hub.Events(ctx, 1)
	if err != nil {
		slog.Error("failed to listen for events", "error", err.Error())
		return
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	for {
		select {
		case sig := <-signalChan:
			slog.Info("got signal to exit", "signal", sig)
			ctxCancel()
			time.Sleep(time.Second)
		case e := <-eventsChan:
			slog.Info("got event", "event", e)
		}

		if ctx.Err() != nil {
			break
		}
	}
}

func initSlog() {
	level, ok := logLevels[strings.ToUpper(*level)]
	if !ok {
		panic("invalid log level")
	}

	options := &slog.HandlerOptions{
		Level: level,
	}

	var handler slog.Handler
	if *json {
		handler = slog.NewJSONHandler(os.Stdout, options)
	} else {
		handler = slog.NewTextHandler(os.Stdout, options)
	}

	slog.SetDefault(slog.New(handler))
}
