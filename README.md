# Learning bot

## Purpose

The aim of this project is to create a GitLab bot which would download
students code and generate a report containing advice on repairing
issues found in the source code.  

### Plan

Currently, the plan is to use [checkstyle] to list issues in the code,
which is parsed by this program, detecting which type of error it is
and on which line and file.  

Then an issue will be created on the project's issue tracker with a
link to the report. The report would contain a list of all the issues
in the source code, each of which would give a preview of the code and
a suggestion on how to fix it.

Alternatively, a CI-based addon would be built into GitLab which would
add a badge to each commit, which would redirect to the report page.

## Requirements

The following packages must be installed on your system.

- git
- go (tested with 1.12)

The program supports MySQL, but it is optional requirement as an SQLite
option is available.

## Installation from source

### Downloading

Since this is currently a private repository, using `go get` will not
suffice. Make sure you have configured git to authenticate with GitLab
via ssh.

```
$ mkdir -p ~/go/gitlab.com/gitedulab
$ git clone git@gitlab.com:gitedulab/learning-bot.git ~/go/gitlab.com/gitedulab/learning-bot
```

### Building

```
$ cd ~/go/gitlab.com/gitedulab/learning-bot
$ go get -u
$ go build
```

**Note:** You have to run the binary in the directory, so the program
is able to find web files required to render the web pages.

### Configuring and running

Before running the web server, the program has to be configured using
the `config` command-line option.

```
$ ./learning-bot config
```

The program will guide you through the process on creating a new configuration
file. You'll have to specify the SQL driver/server, checkstyle jar and configuration
location, bot private token, and so on...  

Once configured, the program would generate a `config.toml` file (which can be lated
edited, if required). And the web server (and bot) can start.

```
$ ./learning-bot run
```

## Development board

Development can be tracked at our [Trello board](https://trello.com/b/tTjkyF73/learning-bot).


[checkstyle]: https://checkstyle.org/
