# gitsink

Work in progress CLI to sync repositories from other providers to GitHub Enterprise or github.com

## Summary

Sync your repositories from any of the repository management providers (BitBucket, GitLab, GitHub Public) to GitHub Enterprise or the public github.com.

You can move repositories between different organizations, projects, teams, personal user accounts, etc. This tool ensures all branches and all tags are synced, with the complete commit history.

## Usage

```
$ ./gitsink
usage: gitsink [<flags>] <command> [<args> ...]

The Github-Sync CLI

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
$ ./gitsink sync --help
usage: gitsink sync [<flags>]

Sync Bitbucket and GitHub repositories

Flags:
  --help                  Show context-sensitive help (also try --help-long and --help-man)
  --log-level="info"      Set log-level (trace|debug|info|warn|error|fatal|panic)
  --run-once              Syncs the repositories once
  --block-new-migrations  Block new migrations and sync only existing repos on GitHub
```

## Installing

```
go get github.com/parishekar/gitsink
```

## Building

```
make all
```
