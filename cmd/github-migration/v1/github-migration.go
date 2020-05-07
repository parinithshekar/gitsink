package v1

import (
	// "encoding/json"
	"fmt"
	// "io/ioutil"
	"os"
	"time"

	// "github.com/go-openapi/strfmt"
	// "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"

	// pkg "github.com/parinithshekar/github-migration-cli/pkg/v1"
	// logger "github.com/parinithshekar/github-migration-cli/wrap/logrus/v1"
	profile "github.com/parinithshekar/github-migration-cli/wrap/profile/v1"
	// runtime "github.com/go-openapi/runtime"
	// httptransport "github.com/go-openapi/runtime/client"
)

// TimeFormat
const (
	TimeFormat = time.RFC3339
)

// Execute called from main.go
func Execute() { // hello

	// Start the profiler and defer stopping it until the program exits.
	defer profile.Start().Stop()

	var (
		// Main git-migration command
		app         = kingpin.New("github-migration", "The Github-Migration CLI")
		appLogLevel = app.Flag("log-level", "Set log-level (trace|debug|info|warn|error|fatal|panic).").Default("info").OverrideDefaultFromEnvar("MERAKI_LOG_LEVEL").String()

		//////////
		// sync
		appSync        = app.Command("sync", "Sync Bitbucket and GitHub repositories")
		appSyncRunOnce = appSync.Flag("run-once", "Syncs the repositories once").Bool()
		// appSyncPersonalAccount    = appSync.Flag("personal-account", "Migrates/Syncs the repositories to personal GitHub account").Bool()
		appSyncBlockNewMigrations = appSync.Flag("block-new-migrations", "Block new migrations and sync only existing repos on GitHub").Bool()

		/////////
		// interactive - Can leave it out, does not make sense if supporting multiple sources for integrations
		appInteractive = app.Command("interactive", "Select the projects and repositories to migrate/sync")
	)

	p := kingpin.MustParse(app.Parse(os.Args[1:]))

	switch p {

	case appSync.FullCommand():
		fmt.Println("SYNC")
		fmt.Printf("App Log Level: %v\n", *appLogLevel)
		fmt.Printf("Run Once: %v\n", *appSyncRunOnce)
		// fmt.Printf("Personal Account: %v\n", *appSyncPersonalAccount)
		fmt.Printf("Block New: %v\n", *appSyncBlockNewMigrations)

	case appInteractive.FullCommand():
		fmt.Println("INTERACTIVE")
		fmt.Printf("App Log Level: %v\n", *appLogLevel)

	}
}
