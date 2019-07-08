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
- Go *(tested with 1.12)*
- OpenJDK Runtime Environment *(tested with 1.8.0)*
- unzip *(Unix utility)*
- checkstyle *(jar, tested with 8.22)*

The program supports MySQL, but it is optional requirement as an SQLite
option is available. No additional packages are needed to be installed for
database functionality.

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

## Further configuration

There are some default configuration values which aren't setup in the wizard. Here
is the documentation for `config.toml`.

- `site_title`: Title of the website and bot.
- `site_url`: The URL of the bot website, including port, used in issue tracker link.
- `bot_private_access_token`: Bot account's private access token.
- `checkstyle_jar_path`: Path for the checkstyle JAR file, must be downloaded separately.
- `checkstyle_config_path`: Path for the checkstyle configuration file, provided by the project.
- `gitlab_instance_url`: The URL of the GitLab instance for the API, links and integration.
- `lms_title`: Optional LMS link title.
- `lms_url`: Optional LMS link URL.
- `check_active_repositories_cron`: Cron job schedule interval. Learn more about format [here](https://godoc.org/github.com/robfig/cron#hdr-Predefined_schedules).
- `timezone`: Timezone used for report and database.
- `code_snippet_include_previous_lines`: Maximum number of lines to include before troubled line in report code snippet.
- `database_configuration`: Database configuration field.
  - `type`: Driver type; 0 for SQLite, 1 for MySQL.
  - `host`: [MySQL] The host of the MySQL server.
  - `name`: [MySQL] The name of the MySQL user.
  - `ssl_mode`: [MySQL] The SSL/TLS mode of the MySQL connection.
  - `path`: [SQLite] The path for the SQLite database file.


## Installing the Project Service (GitLab)

The project service can be installed by coping the file `learning_bot_service.rb` in
the `project_services_gitlab` folder of this project to the GitLab-CE source at
`app/models/project_services/[learning_bot_service.rb]`.  

After copying the file, you will have to change the `learning_bot_service_url` constant
to be the bot's instance URL (including port, if non-default).  

Then make sure to include the service models in `app/models/project.rb` and
`spec/models/project_spec.rb` to make sure that GitLab recognises and loads
the new service.

## Development board

Development can be tracked at our [Trello board](https://trello.com/b/tTjkyF73/learning-bot).


[checkstyle]: https://checkstyle.org/
