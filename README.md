# Scaffold
A command-line tool for initializing new software projects, with an emphasis on go-kit based services.

## Getting Started
If you are routinely setting certain values on (nearly) every use of `scaffold`, it makes sense to configure them as personalized defaults.  This can
 be accomplished through either the configuration file or environment variables.  In either case, values supplied at the command-line override
these personalized defaults.  Below is a list of all available config settings and environment variables.

### Configuration File
`$HOME/.scaffold.yaml`:
```yaml
---
create-repo: true
disable-modules: false
github-auth-token: 8e...f3
github-enterprise-server-url: http(s)://fq.d.n/
github-organization: exampleorg
license: mit
```

### Environment Variables
  - `CREATE_REPO`
  - `DISABLE_MODULES`
  - `GITHUB_AUTH_TOKEN`
  - `GITHUB_ENTERPRISE_SERVER`
  - `GITHUB_ORGANIZATION`
  - `LICENSE`

### License Notes
Valid options for a license are any exact license keyword supported by the github api. This list can be found [here](https://docs.github.com/en/github/creating-cloning-and-archiving-repositories/licensing-a-repository#searching-github-by-license-type).

## Usage
```
scaffold -h
A workflow utility that streamlines creating standardized go projects, generating much of the "boilerplate" structure and layout.

This tool is focused on go, go-kit, jaeger opentracing, and github (enterprise)

Usage:
  scaffold [command]

Available Commands:
  help        Help about any command
  init        Initialize project
  new         Create a new project resource

Flags:
      --config string   config file (default is "$HOME/.scaffold.yaml)"
  -h, --help            help for scaffold

Use "scaffold [command] --help" for more information about a command.
```

### scaffold init
The `init` subcommand is used to initialize a new project directory and (optionally) the remote git repo.

```
scaffold init -h
Initialize a new go project

Usage:
  scaffold init [flags]

Flags:
      --create-repo                           Create a new remote/origin for this project (default is false)
      --disable-modules                       Disable go modules for this project (default is false)
      --github-auth-token string              Github (Enterprise) Personal Auth Token for creating remote/origin (Ignored unless --create-repo used)
      --github-enterprise-server-url string   Base URL of Github Enterprise server if needed (ex: https://ghe.example.com) (Ignored unless --create-repo used)
      --github-organization string            Org owner of remote/origin if not the user (Ignored unless --create-repo used)
      --github-repo-name string               Repository name to create (default is the project name) (Ignored unless --create-repo used)
  -h, --help                                  help for init
      --license string                        Project license (Ignored unless --create-repo used) (default "MIT")
      --module-name string                    Name of the go module for this project (default is the project name)
      --project-dir string                    Directory in which to initialize new project (default is project name)

Global Flags:
      --config string   config file (default is "$HOME/.scaffold.yaml)"
```

**Notes for `--create-repo`**: After passing the `--create-repo` flag, `scaffold` will attempt to create a repo via the github api, using a
 personal access token, then clone it.  The user must have permissions to perform the requested actions,.  The personal access token requires only
  the `repo` scope.  Thus, if a user specifies a `--github-organization` they do not have permission to create new repos in, the operation will fail.