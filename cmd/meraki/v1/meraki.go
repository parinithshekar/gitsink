package v1

import (
	"os"
	"time"

	"github.com/cisco-sso/meraki-cli/info"
	//"golang.org/x/net/context"
	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"

	pkg "github.com/cisco-sso/meraki-cli/pkg/v1"
	logger "github.com/cisco-sso/meraki-cli/wrap/logrus/v1"
	profile "github.com/cisco-sso/meraki-cli/wrap/profile/v1"
)

const (
	TimeFormat = time.RFC3339
)

func Execute() {

	// Start the profiler and defer stopping it until the program exits.
	defer profile.Start().Stop()

	var (
		// meraki
		app          = kingpin.New("meraki", "The Meraki CLI command.")
		appAuthToken = app.Flag("auth-token", "Auth token.").
				Default("").OverrideDefaultFromEnvar("MERAKI_AUTH_TOKEN").String()
		appLogLevel = app.Flag("log-level", "Set log-level (trace|debug|info|warn|error|fatal|panic).").
				Default("info").OverrideDefaultFromEnvar("MERAKI_LOG_LEVEL").String()
		appServerHost = app.Flag("server-host", "Server host.").
				Default("localhost").OverrideDefaultFromEnvar("MERAKI_SERVER_HOST").String()
		appServerPort = app.Flag("server-port", "Server port.").
				Default("10080").OverrideDefaultFromEnvar("MERAKI_SERVER_PORT").Int()

		///////////////////////////////////////
		// meraki devices
		appDevices = app.Command("devices", "Device commands.")

		// meraki devices list
		appDevicesList = appDevices.Command("list", "List all devices.") // GetNetworkDevices

		// meraki devices get
		appDevicesGet   = appDevices.Command("get", "Get a device.") // GetNetworkDevice
		appDevicesGetId = appDevicesGet.Arg("id", "ID of the device to get.").Required().String()

		///////////////////////////////////////
		// meraki events
		appEvents = app.Command("events", "Event commands.")

		// meraki events list
		appEventsList = appEvents.Command("list", "List all events.") // GetNetworkEvents

		///////////////////////////////////////
		// meraki version
		appVersion = app.Command("version", "Display version information.")
	)
	log := logger.New()

	a := &pkg.App{
		Config:  &pkg.AppConfig{},
		Secrets: &pkg.AppSecrets{},
	}

	p := kingpin.MustParse(app.Parse(os.Args[1:]))

	// meraki
	log.SetLevel(*appLogLevel)

	// TODO: Implement actual YAML files here. (rank: Config < Env Var < Flag)

	// Populate Secrets
	if *appAuthToken != "" {
		a.Secrets.AuthToken = *appAuthToken
	}

	// Populate Config
	if *appLogLevel != "" {
		a.Config.LogLevel = *appLogLevel
	}
	if *appServerHost != "" {
		a.Config.ServerHost = *appServerHost
	}
	if *appServerPort != 0 {
		a.Config.ServerPort = *appServerPort
	}

	// Initialize "just in case" vars.
	/*
		var ctx context.Context
		var err error
		var t string
		c := client.NewClient(a, log)
		if a.Secrets.AuthToken != "" {
			c.AuthToken = a.Secrets.AuthToken
		}
	*/

	switch p {

	case appDevicesList.FullCommand():
		log.WithField("args", "meraki devices list").Tracef("called")

	case appDevicesGet.FullCommand():
		log.WithFields(logrus.Fields{
			"args": "meraki devices get",
			"id":   *appDevicesGetId,
		}).Tracef("called")

	case appEventsList.FullCommand():
		log.WithField("args", "meraki events list").Tracef("called")

	case appVersion.FullCommand():
		log.WithFields(logrus.Fields{
			"program":          info.Program,
			"license":          info.License,
			"url":              info.URL,
			"build_user":       info.BuildUser,
			"build_date":       info.BuildDate,
			"language":         info.Language,
			"language_version": info.LanguageVersion,
			"version":          info.Version,
			"revision":         info.Revision,
			"branch":           info.Branch,
		}).Infof("version")
	}
}
