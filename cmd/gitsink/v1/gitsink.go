package v1

import (
	"os"
	"fmt"
	// "io/ioutil"
	// "encoding/json"

	// "github.com/go-openapi/strfmt"
	// "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"

	// pkg "github.com/parinithshekar/gitsink/pkg/v1"
	bbcloud "github.com/parinithshekar/gitsink/plugins/input/bitbucket/cloud"
	git "github.com/parinithshekar/gitsink/plugins/output/git"
	ghpublic "github.com/parinithshekar/gitsink/plugins/output/github/public"
	logger "github.com/parinithshekar/gitsink/wrap/logrus/v1"
	profile "github.com/parinithshekar/gitsink/wrap/profile/v1"
	// bbserver "github.com/parinithshekar/gitsink/plugins/input/bitbucket/server"
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
	fmt.Printf("%+v\n", config)

	switch p {

	case appSync.FullCommand():
		log.Infof("SYNC")
		log.Infof("App Log Level: %v\n", *appLogLevel)
		log.Infof("Run Once: %v\n", *appSyncRunOnce)
		fmt.Printf("Personal Account: %v\n", *appSyncPersonalAccount)
		log.Infof("Block New: %v\n", *appSyncBlockNewMigrations)

	case appInteractive.FullCommand():
		log.Infof("INTERACTIVE")
		log.Infof("App Log Level: %v\n", *appLogLevel)

	case appTest.FullCommand():
		log.Infof("TEST")
		var input plugins.Input
		var output plugins.Output
		for _, integration := range config.Integrations {
			fmt.Println(integration.Name)

			// INPUT PLUGIN
			// get input plugin based on input type
			switch integration.Source.Type {
			case "bitbucket-cloud":
				input = bbcloud.New(integration.Source)

			default:
				log.Errorf("Unsupported source type: %v", integration.Source.Type)
			}
			// Authenticate credentials for reading from input
			_, err := input.Authenticate()
			if err != nil {
				log.Errorf(err.Error())
			}
			// Get repositories to sync
			repos, err := input.Repositories(true)
			if err != nil {
				log.Errorf(err.Error())
			}

			// OUTPUT PLUGIN
			// get output plugin based on output type
			switch integration.Target.Type {
			case "github-public":
				output = ghpublic.New(integration.Target)

			default:
				log.Errorf("Unsupported target type: %v", integration.Target.Type)
			}
			// Authenticate credentials for pushing to output
			_, err = output.Authenticate()
			if err != nil {
				log.Errorf(err.Error())
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
