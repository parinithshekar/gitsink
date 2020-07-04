package v1

import (
	"os"
	"fmt"
	// "io/ioutil"
	// "encoding/json"

	// "github.com/go-openapi/strfmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"

	// pkg "github.com/parinithshekar/gitsink/pkg/v1"
	bbcloud "github.com/parinithshekar/gitsink/plugins/input/bitbucket/cloud"
	bbserver "github.com/parinithshekar/gitsink/plugins/input/bitbucket/server"
	git "github.com/parinithshekar/gitsink/plugins/output/git"
	ghpublic "github.com/parinithshekar/gitsink/plugins/output/github/public"
	logger "github.com/parinithshekar/gitsink/wrap/logrus/v1"
	profile "github.com/parinithshekar/gitsink/wrap/profile/v1"
	config "github.com/parinithshekar/gitsink/common/config"
	plugins "github.com/parinithshekar/gitsink/plugins/interfaces"
	// runtime "github.com/go-openapi/runtime"
	// httptransport "github.com/go-openapi/runtime/client"
)

// Execute Runs the core of the CLI
func Execute() { // hello

	// Start the profiler and defer stopping it until the program exits.
	defer profile.Start().Stop()

	// Logger object
	log := logger.New()

	var (
		// Main git-migration command
		app         = kingpin.New("github-migration", "The Github-Migration CLI")
		appLogLevel = app.Flag("log-level", "Set log-level (trace|debug|info|warn|error|fatal|panic).").Default("info").OverrideDefaultFromEnvar("MERAKI_LOG_LEVEL").String()

		//////////
		// sync
		appSync                   = app.Command("sync", "Sync Bitbucket and GitHub repositories")
		appSyncRunOnce            = appSync.Flag("run-once", "Syncs the repositories once").Bool()
		appSyncPersonalAccount    = appSync.Flag("personal-account", "Migrates/Syncs the repositories to personal GitHub account").Bool()
		appSyncBlockNewMigrations = appSync.Flag("block-new-migrations", "Block new migrations and sync only existing repos on GitHub").Bool()

		/////////
		// interactive - Can leave it out, does not make sense if supporting multiple sources for integrations
		appInteractive = app.Command("interactive", "Select the projects and repositories to migrate/sync")

		/////////
		// test
		appTest = app.Command("test", "Test out new features")
	)

	p := kingpin.MustParse(app.Parse(os.Args[1:]))

	config := config.Parse()

	switch p {

	case appSync.FullCommand():
		fmt.Println("SYNC")
		fmt.Printf("App Log Level: %v\n", *appLogLevel)
		fmt.Printf("Run Once: %v\n", *appSyncRunOnce)
		fmt.Printf("Personal Account: %v\n", *appSyncPersonalAccount)
		fmt.Printf("Block New: %v\n", *appSyncBlockNewMigrations)

	case appInteractive.FullCommand():
		fmt.Println("INTERACTIVE")
		fmt.Printf("App Log Level: %v\n", *appLogLevel)

	case appTest.FullCommand():
		fmt.Printf("TEST")
		var input plugins.Input
		var output plugins.Output
		for _, integration := range config.Integrations {
			fmt.Println(integration.Name)

			var err error
			// INPUT PLUGIN
			// get input plugin based on input type
			switch integration.Source.Type {
			case "bitbucket-cloud":
				input, err = bbcloud.New(integration.Source)
				
			case "bitbucket-server":
				input, err = bbserver.New(integration.Source)
				
			default:
				log.WithFields(logrus.Fields{
					"integration": integration.Name,
					"source": integration.Source.Type,
					}).Errorf("Unsupported source type")
			}
			// Error during initializing source input plugin
			if err != nil {
				log.WithFields(logrus.Fields{
					"error": err.Error(),
					"integration": integration.Name,
					"targetType": integration.Source.Type,
				}).Errorf("Initializing source failed")
			}
			
			// Authenticate credentials for reading from input
			_, err = input.Authenticate()
			if err != nil {
				log.WithFields(logrus.Fields{
					"error": err.Error(),
					"integration": integration.Name,
					"source": integration.Source.Type,
				}).Errorf("Source authentication failed")
			}
			// Get repositories to sync
			repos, err := input.Repositories(true)
			if err != nil {
				log.WithFields(logrus.Fields{
					"error": err.Error(),
					"integration": integration.Name,
					"source": integration.Source.Type,
				}).Errorf("Fetching repository list failed")
			}

			// OUTPUT PLUGIN
			// get output plugin based on output type
			switch integration.Target.Type {
			case "github-public":
				output, err = ghpublic.New(integration.Target)
				if err != nil {
					log.WithFields(logrus.Fields{
						"error": err.Error(),
						"integration": integration.Name,
						"targetType": integration.Target.Type,
					}).Errorf("Initializing target failed")
				}

			default:
				log.WithFields(logrus.Fields{
					"integration": integration.Name,
					"targetType": integration.Target.Type,
				}).Errorf("Unsupported target type")
			}
			// Authenticate credentials for pushing to output
			_, err = output.Authenticate()
			if err != nil {
				log.WithFields(logrus.Fields{
					"error": err.Error(),
					"integration": integration.Name,
					"source": integration.Target.Type,
				}).Errorf("Target authentication failed")
			} else {
				log.Infof("GITHUB PUBLIC SUCCESS")
			}

			// SYNC REPOS
			// Check if repos need to by synced or migrated
			// Makes new repo on target if there doesn't already exist one
			repos = output.SyncCheck(repos)
			fmt.Println(repos)

			// Start syncing repos
			gitClient := git.New(input, output, integration.Name)
			gitClient.SyncRepos(repos)
		}
	}
}
