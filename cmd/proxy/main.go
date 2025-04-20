package main

import (
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/zekrotja/sshttproxy/pkg/proxy"
	"github.com/zekrotja/sshttproxy/pkg/stdioconn"
)

type Config struct {
	LogLevel slog.Level `env:"SSHTTPROXY_LOGLEVEL" envDefault:"info"`
	LogFile  string     `env:"SSHTTPROXY_LOGFILE"`
	Target   string     `env:"SSHTTPROXY_TARGET,notEmpty"`
}

func main() {
	if len(os.Args) > 1 {
		err := godotenv.Load(os.Args[1])
		if err != nil {
			slog.Error("failed loading env file", "path", os.Args[1], "err", err.Error())
			os.Exit(1)
		}
	}

	cfg, err := env.ParseAs[Config]()
	if err != nil {
		slog.Error("failed parsing config from env", "err", err.Error())
		os.Exit(1)
	}

	var logOutput io.Writer
	if cfg.LogFile != "" {
		f, err := os.OpenFile(cfg.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
		if err != nil {
			slog.Error("failed opening log file", "err", err.Error())
			os.Exit(1)
		}
		defer f.Close()
		logOutput = f
	} else {
		logOutput = os.Stderr
	}

	logger := slog.New(slog.NewTextHandler(logOutput, &slog.HandlerOptions{Level: cfg.LogLevel}))
	slog.SetDefault(logger)

	slog.Info("Server running ...")

	proxy := &proxy.Proxy{
		Upstream: cfg.Target,
	}

	listener := stdioconn.NewListenerFromStdInOut()

	err = http.Serve(listener, proxy)
	if err != nil {
		slog.Error("failed starting http server", "err", err.Error())
		os.Exit(1)
	}
}
