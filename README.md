![Baton Logo](./docs/images/baton-logo.png)

# `baton-jenkins` [![Go Reference](https://pkg.go.dev/badge/github.com/conductorone/baton-jenkins.svg)](https://pkg.go.dev/github.com/conductorone/baton-jenkins) ![main ci](https://github.com/conductorone/baton-jenkins/actions/workflows/main.yaml/badge.svg)

`baton-jenkins` is a connector for jenkins built using the [Baton SDK](https://github.com/conductorone/baton-sdk). It communicates with the jenkins API to sync data about users and roles.

Check out [Baton](https://github.com/conductorone/baton) to learn more the project in general.

# Getting Started

## brew

```
brew install conductorone/baton/baton conductorone/baton/baton-jenkins
baton-jenkins
baton resources
```

## docker

```
docker run --rm -v $(pwd):/out -e BATON_JENKINS_USERNAME=userID -e BATON_JENKINS_TOKEN=apiKey -e BATON_JENKINS_BASEURL=baseurl ghcr.io/conductorone/baton-jenkins:latest -f "/out/sync.c1z"
docker run --rm -v $(pwd):/out ghcr.io/conductorone/baton:latest -f "/out/sync.c1z" resources
```

## source

```
go install github.com/conductorone/baton/cmd/baton@main
go install github.com/conductorone/baton-jenkins/cmd/baton-jenkins@main

BATON_JENKINS_USERNAME=userID BATON_JENKINS_TOKEN=apiKey BATON_JENKINS_BASEURL=baseurl
baton resources
```

## How to test
You can use this docker-compose.yaml to launch an instance server to interact with Jenkins.

```
version: '3.7'
services:
  jenkins:
    image: jenkins/jenkins:lts
    privileged: true
    user: root
    ports:
      - 8080:8080
      - 50000:50000
    container_name: jenkins-lts
    volumes:
      - ~/jenkins_home:/var/jenkins_home
      - /var/run/docker.sock:/var/run/docker.sock
      - /usr/local/bin/docker:/usr/local/bin/docker
```

Run `docker-compose up -d` to launch the container. You will also need Jenkins’s initial admin password to perform the Jenkins web-based installation. You can get the Jenkins admin password with the following command.
```
docker exec jenkins-lts cat /var/jenkins_home/secrets/initialAdminPassword
```
You should see the password in the following output.
```
99b844a4ad13404796e1ab8bcf05edd1
```

## Access Jenkins Web UI 
At this point, Jenkins is installed and listens on port 8080. You can now access it using the URL [http://localhost:8080](http://localhost:8080). You should see the Jenkins initial setup password screen. 
- Provide your password and click on the Continue button. You should see the customized Jenkins screen. 
- Click on Install suggested plugins. You should see the Getting Started screen. 
- Create your new admin user and click on the Save and Continue buttons. You should see the instance configuration screen. 
- Click on the Save and Finish button

You can then access the Jenkins admin server at [http://localhost:8080](http://localhost:8080) and login with the admin credentials you provided in the previous step.

After you login you can create new resources to be synced by baton.

After creating new resources on the Jenkins server, use the baton-jenkins cli to sync the data from the server with the example command below. 
```
baton-jenkins --jenkins-username <user> --jenkins-password <password> --jenkins-baseurl <baseurl>
```
You can also use a token instead of user/password
```
baton-jenkins --jenkins-username <user> --jenkins-token <token> --jenkins-baseurl <baseurl>
```

After successfully syncing data, use the baton CLI to list the resources and see the synced data. baton resources baton stats

## Plugins
In order to manage users and roles within jenkins we need to install the followings plugins.

[People View Plugin](https://plugins.jenkins.io/people-view/)
Provides the "People" view and API that were part of Jenkins up to version 2.451.

`Installation options`
```
1) Using the GUI: From your Jenkins dashboard navigate to Manage Jenkins > Manage Plugins and select the Available tab. 
Locate this plugin by searching for people-view.
2) Using the CLI tool: jenkins-plugin-cli --plugins people-view:1.2
```

[Role-based Authorization Strategy](https://plugins.jenkins.io/role-strategy/)
Enables user authorization using a Role-Based strategy. Roles can be defined globally or for particular jobs or nodes selected by regular expressions.

`Installation options`
```
1) Using the GUI: From your Jenkins dashboard navigate to Manage Jenkins > Manage Plugins and select the Available tab. 
Locate this plugin by searching for role-strategy.
2) Using the CLI tool: jenkins-plugin-cli --plugins role-strategy:727.vd344b_eec783d
```

# Data Model

`baton-jenkins` will pull down information about the following jenkins resources:
- Users
- Roles
- Nodes
- Jobs 
- Views

# Contributing, Support and Issues

We started Baton because we were tired of taking screenshots and manually building spreadsheets. We welcome contributions, and ideas, no matter how small -- our goal is to make identity and permissions sprawl less painful for everyone. If you have questions, problems, or ideas: Please open a Github Issue!

See [CONTRIBUTING.md](https://github.com/ConductorOne/baton/blob/main/CONTRIBUTING.md) for more details.

# `baton-jenkins` Command Line Usage

```
baton-jenkins

Usage:
  baton-jenkins [flags]
  baton-jenkins [command]

Available Commands:
  capabilities       Get connector capabilities
  completion         Generate the autocompletion script for the specified shell
  help               Help about any command

Flags:
      --base-url string        required: Jenkins ($BATON_BASE_URL) (default "http://localhost:8080")
      --client-id string       The client ID used to authenticate with ConductorOne ($BATON_CLIENT_ID)
      --client-secret string   The client secret used to authenticate with ConductorOne ($BATON_CLIENT_SECRET)
  -f, --file string            The path to the c1z file to sync with ($BATON_FILE) (default "sync.c1z")
  -h, --help                   help for baton-jenkins
      --log-format string      The output format for logs: json, console ($BATON_LOG_FORMAT) (default "json")
      --log-level string       The log level: debug, info, warn, error ($BATON_LOG_LEVEL) (default "info")
      --password string        Application password used to connect to the Jenkins API ($BATON_PASSWORD)
  -p, --provisioning           This must be set in order for provisioning actions to be enabled ($BATON_PROVISIONING)
      --skip-full-sync         This must be set to skip a full sync ($BATON_SKIP_FULL_SYNC)
      --ticketing              This must be set to enable ticketing support ($BATON_TICKETING)
      --token string           HTTP access tokens in Jenkins ($BATON_TOKEN)
      --username string        required: Username of administrator used to connect to the Jenkins API ($BATON_USERNAME)
  -v, --version                version for baton-jenkins

Use "baton-jenkins [command] --help" for more information about a command.
```
