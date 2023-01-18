# gcpout-cli
In case of fire :fire:
```console
user:~$ git commit
user:~$ git push
user:~$ git out
```

Inspired by a nerd joke, I created this tool to allow us developers to open our PRs more easily, without leaving our terminals, opening Jira tickets on the WEB, writing title, descriptions, filling in the checklist boxes. Now, you can do all this in your terminal.

# Getting started
```console
user:~$ git clone <this-repo>
user:~$ cd this-repo
user:~$ go install
user:~$ gcpout-cli
```

# Running Commands

## Init
This command is mandatory for configuring the Jira Server host.

```console
user:~$ gcpout-cli init
```

## OpenPr
This command allows you to:
- Select the project;
- Select the source and target branches;
- Inform the Jira Ticket;
- Inform the type of change (feature, bug, chore, etc.);
- Answer the checklist questions.

Based on that, this command will open the PR for you with a full description based on your answers.

```console
user:~$ gcpout-cli openPr
```

## Help
All commands have a help flag -h to get details about how to run the command.

## Running the commands without install it as a Go package

```console
user:~$ go run main.go <command>
```