# Learning bot

## Description

The aim of this project is to create a GitLab bot which would download
students code and generate a report containing advice on repairing
issues found in the source code.  

Currently, the plan is to use [checkstyle] to list issues in the code,
which is parsed by this program, detecting which type of error it is
and on which line and file.  

Then an issue will be created on the project's issue tracker with a
link to the report. The report would contain a list of all the issues
in the source code, each of which would give a preview of the code and
a suggestion on how to fix it.

## Requirements

This program is written in Go, and is tested with Go 1.12. A database
is optional, as SQLite is an option.

## Download and Build

```
$ mkdir -p ~/go/gitlab.com/gitedulab
$ git clone git@gitlab.com:gitedulab/learning-bot.git ~/go/gitlab.com/gitedulab/learning-bot
$ cd ~/go/gitlab.com/gitedulab/learning-bot
$ go get -u
$ go build
$ ./learning-bot
```

## Development board

Development can be tracked at our [Trello board](https://trello.com/b/tTjkyF73/learning-bot).


[checkstyle]: https://checkstyle.org/
