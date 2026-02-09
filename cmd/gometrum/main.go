package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/Miklakapi/gometrum/internal/agent"
	"github.com/Miklakapi/gometrum/internal/cli"
	"github.com/Miklakapi/gometrum/internal/config"
	"github.com/Miklakapi/gometrum/internal/logger"
	"github.com/Miklakapi/gometrum/internal/mqtt"
	"github.com/Miklakapi/gometrum/internal/sensors"
	"github.com/Miklakapi/gometrum/internal/service"
	"github.com/Miklakapi/gometrum/internal/version"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	appCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	flags, err := cli.ParseFlags()
	if err != nil {
		printErrorAndExit(err, 2)
	}

	switch {
	case flags.GenerateConfig:
		if flags.ConfigPath == "" {
			fmt.Print(config.ExampleYAML)
			return
		}

		if err = config.SaveExample(flags.ConfigPath); err != nil {
			printErrorAndExit(err, 1)
		}
		return
	case flags.GenerateService:
		if strings.TrimSpace(flags.ServicePath) == "" {
			fmt.Print(service.ExampleService)
			return
		}
		if err := service.SaveExample(flags.ServicePath); err != nil {
			printErrorAndExit(err, 1)
		}
		return
	case flags.Validate:
		if _, err = config.LoadAndValidate(flags.ConfigPath); err != nil {
			printErrorAndExit(err, 1)
		}
		return
	case flags.Version:
		fmt.Println(version.VersionString())
		return
	}

	logger.SetupLogger(flags.LogLevel)

	cfg, err := config.LoadAndValidate(flags.ConfigPath)
	if err != nil {
		slog.Error("failed to load configuration", "err", err)
		os.Exit(1)
	}

	err = sensors.Prepare(&cfg)
	if err != nil {
		slog.Error("failed to load sensors configuration", "err", err)
		os.Exit(1)
	}

	sens, err := sensors.Build(cfg)
	if err != nil {
		slog.Error("failed to create sensors from configuration", "err", err)
		os.Exit(1)
	}

	s := agent.Settings{
		DiscoveryPrefix: cfg.MQTT.DiscoveryPrefix,
		StatePrefix:     cfg.MQTT.StatePrefix,
		DeviceId:        cfg.Agent.DeviceID,
		DeviceName:      cfg.Agent.DeviceName,
		Manufacturer:    cfg.Agent.Manufacturer,
		Model:           cfg.Agent.Model,
		Once:            flags.Once,
	}

	var pub mqtt.Publisher
	if flags.DryRun {
		pub = mqtt.NewDryRun()
	} else {
		o := MQTT.NewClientOptions()

		addr := net.JoinHostPort(cfg.MQTT.Host, strconv.Itoa(cfg.MQTT.Port))
		o.AddBroker("tcp://" + addr)

		o.SetClientID(cfg.MQTT.ClientID)
		o.SetUsername(cfg.MQTT.Username)
		o.SetPassword(cfg.MQTT.Password)

		o.SetCleanSession(true)
		o.SetAutoReconnect(true)
		o.SetConnectRetry(true)

		o.SetConnectTimeout(10 * time.Second)
		o.SetKeepAlive(30 * time.Second)
		o.SetPingTimeout(10 * time.Second)

		availabilityTopic := cfg.MQTT.StatePrefix + "/" + cfg.Agent.DeviceID + "/availability"

		o.SetWill(availabilityTopic, "offline", 1, true)

		pub = mqtt.New(o)
	}

	a, err := agent.New(s, sens, pub)
	if err != nil {
		slog.Error("failed to initialize agent", "err", err)
		os.Exit(1)
	}

	if flags.Purge {
		if err := a.Purge(); err != nil {
			slog.Error("purge failed", "err", err)
			os.Exit(1)
		}
		slog.Info("purge completed")
		return
	}

	if err = a.Run(appCtx); err != nil {
		slog.Error("agent stopped with error", "err", err)
		os.Exit(1)
	}
}

func printErrorAndExit(err error, code int) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(code)
}
