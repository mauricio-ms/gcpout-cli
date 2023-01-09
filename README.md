# gcpout-cli
In case of fire :fire:
```console
user:~$ git commit
user:~$ git push
user:~$ git out
```

Inspired by a nerd joke I created this tool to allow we as developers opening our PRs more easily, without leaving our terminals, opening Jira tickets on WEB, writing title, descriptions, filling the checklist boxes. Now, you can do all of that stuff in you terminal.

## Getting started
```console
user:~$ git clone <this-repo>
user:~$ cd this-repo
user:~$ go install
user:~$ gcpout-cli
```

## Running Commands

# Init
This command is mandatory in order to configure the Jira Server host.

```console
user:~$ gcpout-cli init
```

# OpenPr
This command allows you to:
- Select the project;
- Select the source and target branches;
- Inform the Jira Ticket;
- Inform the type of change (feature, bug, chore, etc);
- Answer the checklist questions.

Based on that this command will open the PR for you with a full description based on your answers.

```console
user:~$ gcpout-cli openPr
```

# Help
All comands has a help flag -h to get details about how run the command