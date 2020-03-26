package v1

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/cisco-sso/meraki-cli/info"
	"github.com/go-openapi/strfmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"

	apiclient "github.com/cisco-sso/meraki-cli/client"
	api_clients "github.com/cisco-sso/meraki-cli/client/clients"
	api_devices "github.com/cisco-sso/meraki-cli/client/devices"
	api_events "github.com/cisco-sso/meraki-cli/client/events"
	api_networks "github.com/cisco-sso/meraki-cli/client/networks"
	api_organizations "github.com/cisco-sso/meraki-cli/client/organizations"
	api_ssids "github.com/cisco-sso/meraki-cli/client/s_s_i_ds"
	pkg "github.com/cisco-sso/meraki-cli/pkg/v1"
	logger "github.com/cisco-sso/meraki-cli/wrap/logrus/v1"
	profile "github.com/cisco-sso/meraki-cli/wrap/profile/v1"
	runtime "github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
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
		appAuthToken = app.Flag("auth-token", "Auth token. Export env var MERAKI_AUTH_TOKEN as an alternative").
				Default("").OverrideDefaultFromEnvar("MERAKI_AUTH_TOKEN").String()
		appLogLevel = app.Flag("log-level", "Set log-level (trace|debug|info|warn|error|fatal|panic).").
				Default("info").OverrideDefaultFromEnvar("MERAKI_LOG_LEVEL").String()
		appServerHost = app.Flag("server-host", "Server host.").
				Default(apiclient.DefaultHost).OverrideDefaultFromEnvar("MERAKI_SERVER_HOST").String()
		appServerPort = app.Flag("server-port", "Server port.").
				Default("443").OverrideDefaultFromEnvar("MERAKI_SERVER_PORT").Int()

		///////////////////////////////////////
		// meraki devices
		appDevices = app.Command("devices", "Device commands.")

		// meraki devices list
		appDevicesList               = appDevices.Command("list", "List all devices.") // GetNetworkDevices
		appDevicesListOrganizationId = appDevicesList.Flag("organization-id", "ID of the organization.").Required().String()

		///////////////////////////////////////
		// meraki events
		appEvents = app.Command("events", "Event commands.")

		// meraki events list
		appEventsList            = appEvents.Command("list", "List all events for a network.")
		appEventsListNetworkId   = appEventsList.Flag("network-id", "ID of the network.").Required().String()
		appEventsListProductType = appEventsList.Flag("product-type", "Type of the product. Valid types (wireless|appliance|switch|systemsManager|camera|cellularGateway)").Required().String()

		// meraki events types
		appEventsTypes          = appEvents.Command("types", "List all events for a network.")
		appEventsTypesNetworkId = appEventsTypes.Flag("network-id", "ID of the network.").Required().String()

		///////////////////////////////////////
		// meraki networks
		appNetworks = app.Command("networks", "Network commands.")

		// meraki networks list
		appNetworksList               = appNetworks.Command("list", "List all networks for an organization.") // GetNetworkNetworks
		appNetworksListOrganizationId = appNetworksList.Flag("organization-id", "ID of the organization.").Required().String()

		// meraki networks get
		appNetworksGet          = appNetworks.Command("get", "Get a network.")
		appNetworksGetNetworkId = appNetworksGet.Flag("network-id", "ID of the network.").Required().String()

		// meraki networks clients list
		appNetworksClients              = appNetworks.Command("clients", "Client commands.")
		appNetworksClientsList          = appNetworksClients.Command("list", "List all Clients for a network.")
		appNetworksClientsListNetworkId = appNetworksClientsList.Flag("network-id", "ID of the network.").Required().String()

		// meraki networks ssids list
		appNetworksSsids              = appNetworks.Command("ssids", "SSID commands.")
		appNetworksSsidsList          = appNetworksSsids.Command("list", "List all SSIDs for a network.")
		appNetworksSsidsListNetworkId = appNetworksSsidsList.Flag("network-id", "ID of the network.").Required().String()

		///////////////////////////////////////
		// meraki organizations
		appOrganizations = app.Command("organizations", "Organization commands.")

		// meraki organizations list
		appOrganizationsList = appOrganizations.Command("list", "List all organizations.")

		// meraki organizations get
		appOrganizationsGet               = appOrganizations.Command("get", "Get a organization.")
		appOrganizationsGetOrganizationId = appOrganizationsGet.Flag("organization-id", "ID of the organization.").Required().String()

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
	} else {
		// Check here that it is Required, so we don't affect help
		log.Fatalf("Merak Auth Token must be provided either by arg '--auth-token' or envvar 'MERAKI_AUTH_TOKEN'")
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

	// create the API client
	transport := httptransport.New(a.Config.ServerHost, apiclient.DefaultBasePath, apiclient.DefaultSchemes)
	client := apiclient.New(transport, strfmt.Default)
	authInfo := httptransport.APIKeyAuth("X-Cisco-Meraki-API-Key", "header", a.Secrets.AuthToken)

	switch p {

	case appDevicesList.FullCommand():
		f := func() (interface{}, error) {
			params := api_devices.NewGetOrganizationDevicesParams()
			params.OrganizationID = *appDevicesListOrganizationId
			return client.Devices.GetOrganizationDevices(params, authInfo)
		}
		printPayload(f, log)

	case appEventsList.FullCommand():
		f := func() (interface{}, error) {
			params := api_events.NewGetNetworkEventsParams()
			params.NetworkID = *appEventsListNetworkId
			params.ProductType = appEventsListProductType
			return client.Events.GetNetworkEvents(params, authInfo)
		}
		printPayload(f, log)

	case appEventsTypes.FullCommand():
		f := func() (interface{}, error) {
			params := api_events.NewGetNetworkEventsEventTypesParams()
			params.NetworkID = *appEventsTypesNetworkId
			return client.Events.GetNetworkEventsEventTypes(params, authInfo)
		}
		printPayload(f, log)

	case appNetworksList.FullCommand():
		f := func() (interface{}, error) {
			params := api_networks.NewGetOrganizationNetworksParams()
			params.OrganizationID = *appNetworksListOrganizationId
			return client.Networks.GetOrganizationNetworks(params, authInfo)
		}
		printPayload(f, log)

	case appNetworksGet.FullCommand():
		f := func() (interface{}, error) {
			params := api_networks.NewGetNetworkParams()
			params.NetworkID = *appNetworksGetNetworkId
			return client.Networks.GetNetwork(params, authInfo)
		}
		printPayload(f, log)

	case appNetworksClientsList.FullCommand():
		f := func() (interface{}, error) {
			params := api_clients.NewGetNetworkClientsParams()
			params.NetworkID = *appNetworksClientsListNetworkId
			return client.Clients.GetNetworkClients(params, authInfo)
		}
		printPayload(f, log)

	case appNetworksSsidsList.FullCommand():
		f := func() (interface{}, error) {
			params := api_ssids.NewGetNetworkSsidsParams()
			params.NetworkID = *appNetworksSsidsListNetworkId
			return client.SsiDs.GetNetworkSsids(params, authInfo)
		}
		printPayload(f, log)

	case appOrganizationsList.FullCommand():
		f := func() (interface{}, error) {
			params := api_organizations.NewGetOrganizationsParams()
			return client.Organizations.GetOrganizations(params, authInfo)
		}
		printPayload(f, log)

	case appOrganizationsGet.FullCommand():
		f := func() (interface{}, error) {
			params := api_organizations.NewGetOrganizationParams()
			params.OrganizationID = *appOrganizationsGetOrganizationId
			return client.Organizations.GetOrganization(params, authInfo)
		}
		printPayload(f, log)

	case appVersion.FullCommand():
		type Version struct {
			Program         string `json:"program"`
			License         string `json:"license"`
			URL             string `json:"url"`
			BuildUser       string `json:"build_user"`
			BuildDate       string `json:"build_date"`
			Language        string `json:"language"`
			LanguageVersion string `json:"language_version"`
			Version         string `json:"version"`
			Revision        string `json:"revision"`
			Branch          string `json:"branch"`
		}
		version := Version{
			Program:         info.Program,
			License:         info.License,
			URL:             info.URL,
			BuildUser:       info.BuildUser,
			BuildDate:       info.BuildDate,
			Language:        info.Language,
			LanguageVersion: info.LanguageVersion,
			Version:         info.Version,
			Revision:        info.Revision,
			Branch:          info.Branch,
		}
		versionBytes, _ := json.MarshalIndent(version, "", "  ")
		fmt.Println(string(versionBytes))
	}
}

func printPayload(f func() (interface{}, error), log *logger.Logger) {
	log.WithFields(logrus.Fields{"args": os.Args}).Tracef("called")

	type PayloadInterface interface {
		GetPayload() interface{}
	}

	resp, err := f()
	if err != nil {
		failMsg := "Operation Failed: To debug, run with environment variable DEBUG=1 set"

		apiError, ok := err.(*runtime.APIError)
		if !ok {
			log.WithFields(logrus.Fields{"err": err}).Errorf(failMsg)
		}
		response, ok := apiError.Response.(runtime.ClientResponse)
		if !ok {
			log.WithFields(logrus.Fields{"err": err}).Errorf(failMsg)
		}
		body, err := ioutil.ReadAll(response.Body())
		if err != nil {
			log.WithFields(logrus.Fields{"err": err}).Errorf(failMsg)
		}
		// Upon error 400, we cannot seem to read the body of the http response for printing.
		//    It errors out in the ioutil.ReadAll statement above.
		//    TODO: Fix this by tuning the swagger file.
		//    https://github.com/go-openapi/runtime/issues/121
		fmt.Println(string(body))
		os.Exit(1)
	}

	var str interface{}
	p, ok := resp.(PayloadInterface)
	if ok {
		str = p.GetPayload()
	} else {
		str = resp
	}

	json, err := json.MarshalIndent(str, "", "  ")
	if err != nil {
		log.WithFields(logrus.Fields{"err": err, "str": str}).Errorf("failed")
	}
	fmt.Println(string(json))
}
