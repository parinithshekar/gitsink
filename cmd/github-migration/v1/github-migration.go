package v1

import (
	// "encoding/json"
	"fmt"
	// "io/ioutil"
	"os"

	// "github.com/go-openapi/strfmt"
	// "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"

	// pkg "github.com/parinithshekar/github-migration-cli/pkg/v1"
	logger "github.com/parinithshekar/github-migration-cli/wrap/logrus/v1"
	profile "github.com/parinithshekar/github-migration-cli/wrap/profile/v1"
	bbcloud "github.com/parinithshekar/github-migration-cli/plugins/input/bitbucket/cloud"
	ghpublic "github.com/parinithshekar/github-migration-cli/plugins/output/github/public"
	git "github.com/parinithshekar/github-migration-cli/plugins/output/git"
	// bbserver "github.com/parinithshekar/github-migration-cli/plugins/input/bitbucket/server"
	plugins "github.com/parinithshekar/github-migration-cli/plugins/interfaces"
	config "github.com/parinithshekar/github-migration-cli/common/config"
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
		appSync        = app.Command("sync", "Sync Bitbucket and GitHub repositories")
		appSyncRunOnce = appSync.Flag("run-once", "Syncs the repositories once").Bool()
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


			// get input plugin based on input type
			switch integration.Source.Type {
			case "bitbucket-cloud":
				input = bbcloud.New(integration.Source)
			}
			// Authenticate credentials for reading from input
			_, err := input.Authenticate()
			if err != nil {
				log.Errorf(err.Error())
			}
			// Get repositories to sync
			repos, err := input.Repositories(true)
			if err != nil { log.Errorf(err.Error()) }


			// get output plugin based on output type
			switch integration.Target.Type {
			case "github-public":
				output = ghpublic.New(integration.Target)
			}
			// Authenticate credentials for pushing to output
			_, err = output.Authenticate()
			if err != nil {
				log.Errorf(err.Error())
			} else {
				log.Infof("GITHUB PUBLIC SUCCESS")
			}

			// Start syncing repos
			gitClient := git.New(&output)
			gitClient.SyncRepos(repos)
		}
	}
}
