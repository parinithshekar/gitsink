# github-migration-cli

Wrok in progress CLI to migrate repositories from other providers to GitHub Enterprise or github.com

## Summary

Migrate your  repositories from any of the repository management providers (BitBucket, GitLab, GitHub Public) to GitHub Enterprise or the public github.com.

You can move repositories between different organizations, projects, teams, personal user accounts, etc. This tool ensures all branches and all tags are synced, with the complete commit history.

## Usage

```
$ ./github-migration-cli
usage: git-migration-cli [<flags>] <command> [<args> ...]

The Github-Migration CLI

Flags:
  --help              Show context-sensitive help (also try --help-long and --help-man)
  --log-level="info"  Set log-level (trace|debug|info|warn|error|fatal|panic)

Commands:
  help [<command>...]
    Show help

  sync [<flags>]
    Sync Bitbucket and GitHub repositories

  interactive
    Select the projects and repositories to migrate/sync
```

```
$ ./github-migration-cli sync --help
usage: github-migration-cli sync [<flags>]

Sync Bitbucket and GitHub repositories

Flags:
  --help                  Show context-sensitive help (also try --help-long and --help-man)
  --log-level="info"      Set log-level (trace|debug|info|warn|error|fatal|panic)
  --run-once              Syncs the repositories once
  --block-new-migrations  Block new migrations and sync only existing repos on GitHub
```



## Installing

```
go get github.com/parishekar/github-migration-cli
```


## Building

```
make all
```
