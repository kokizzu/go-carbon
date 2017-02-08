package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"os/user"
	"path/filepath"
	"runtime"
	"strconv"
	"syscall"

	"github.com/lomik/go-carbon/carbon"
	daemon "github.com/sevlyar/go-daemon"
	"github.com/uber-go/zap"

	_ "net/http/pprof"
)

// Version of go-carbon
const Version = "0.9.0"

func httpServe(addr string) (func(), error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return nil, err
	}

	go http.Serve(listener, nil)
	return func() { listener.Close() }, nil
}

func loggingLevel(cfg *carbon.Config) (zap.Level, error) {
	var logLevel zap.Level

	if cfg.Common.LogLevel != "" {
		log.Println("[WARNING] `common.log-level` is DEPRICATED. Use `logging` config section")
		if err = logLevel.UnmarshalText([]byte(cfg.Common.LogLevel)); err != nil {
			return nil, err
		}
		return logLevel, nil
	}

	if err = logLevel.UnmarshalText([]byte(cfg.Logging.Level)); err != nil {
		return nil, err
	}
	return logLevel, nil
}

func loggingFile(cfg *carbon.Config) (string, error) {
	if cfg.Common.Logfile != "" {
		log.Println("[WARNING] `common.logfile` is DEPRICATED. Use `logging` config section")
		return cfg.Common.Logfile, nil
	}

	return cfg.Logging.File, nil
}

func loggingPrepare(filename string, owner *user.User) error {
	if filename == "" {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		return err
	}

	fd, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if fd != nil {
		fd.Close()
	}
	if err != nil {
		return err
	}
	if err := os.Chmod(filename, 0644); err != nil {
		return err
	}
	if owner != nil {

		uid, err := strconv.ParseInt(owner.Uid, 10, 0)
		if err != nil {
			return err
		}

		gid, err := strconv.ParseInt(owner.Gid, 10, 0)
		if err != nil {
			return err
		}

		if err := os.Chown(filename, int(uid), int(gid)); err != nil {
			return err
		}
	}

	return nil
}

func loggingInit(filename string, level zap.Level) (zap.Logger, error) {
	return nil, nil
}

func main() {
	var err error

	/* CONFIG start */

	configFile := flag.String("config", "", "Filename of config")
	printDefaultConfig := flag.Bool("config-print-default", false, "Print default config")
	checkConfig := flag.Bool("check-config", false, "Check config and exit")

	printVersion := flag.Bool("version", false, "Print version")

	isDaemon := flag.Bool("daemon", false, "Run in background")
	pidfile := flag.String("pidfile", "", "Pidfile path (only for daemon)")

	flag.Parse()

	if *printVersion {
		fmt.Print(Version)
		return
	}

	if *printDefaultConfig {
		if err = carbon.PrintConfig(carbon.NewConfig()); err != nil {
			log.Fatal(err)
		}
		return
	}

	app := carbon.New(*configFile)

	if err = app.ParseConfig(); err != nil {
		log.Fatal(err)
	}

	cfg := app.Config

	var runAsUser *user.User
	if cfg.Common.User != "" {
		runAsUser, err = user.Lookup(cfg.Common.User)
		if err != nil {
			log.Fatal(err)
		}
	}

	logLevel, err := loggingLevel(cfg)
	if err != nil {
		log.Fatal(err)
	}

	logFile, err := loggingFile(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// config parsed successfully. Exit in check-only mode
	if *checkConfig {
		return
	}

	if err := loggingPrepare(logFile, runAsUser); err != nil {
		log.Fatal(err)
	}

	logger, err := loggingInit(logFile, logLevel)
	if err != nil {
		log.Fatal(err)
	}

	if *isDaemon {
		runtime.LockOSThread()

		context := new(daemon.Context)
		if *pidfile != "" {
			context.PidFileName = *pidfile
			context.PidFilePerm = 0644
		}

		if runAsUser != nil {
			uid, err := strconv.ParseInt(runAsUser.Uid, 10, 0)
			if err != nil {
				log.Fatal(err)
			}

			gid, err := strconv.ParseInt(runAsUser.Gid, 10, 0)
			if err != nil {
				log.Fatal(err)
			}

			context.Credential = &syscall.Credential{
				Uid: uint32(uid),
				Gid: uint32(gid),
			}
		}

		child, _ := context.Reborn()

		if child != nil {
			return
		}
		defer context.Release()

		runtime.UnlockOSThread()
	}
	/* CONFIG end */

	// pprof
	// httpStop := func() {}
	if cfg.Pprof.Enabled {
		_, err = httpServe(cfg.Pprof.Listen)
		if err != nil {
			zap.Fatal(err)
		}
	}

	if err = app.Start(); err != nil {
		zap.Fatal(err)
	} else {
		zap.Info("started")
	}

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGUSR2)
		for {
			<-c
			app.DumpStop()
		}
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGHUP)
		for {
			<-c
			zap.Info("HUP received. Reload config")
			if err := app.ReloadConfig(); err != nil {
				zap.Error("config reload failed", zap.Error(err.Error))
			} else {
				zap.Info("config successfully reloaded")
			}
		}
	}()

	app.Loop()

	zap.Info("stopped")
}
