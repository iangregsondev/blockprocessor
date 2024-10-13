# Documentation

## Pre-requisites

Where possible it is recommended to install any tools via a package manager, i.e. Brew
Also a good recommendation to install GO is using gobrew (https://github.com/kevincobain2000/gobrew)

* Go 1.23.2 (although a previous may work)
* Python 3.8 or higher
* Task (https://taskfile.dev/#/installation)
* pipx (https://pipx.pypa.io/stable/installation/)
* jq (for mac / linux) - usually installed via your package manager i.e. ```brew install jq```

## Installation of tooling

We have a number of tools that need to be installed, the following commands will install them.
Please run this from the root of the project.

```bash
task provision:setup
```

## Configuration of tooling

You can test this by running the following command:

```bash
task common:precommit
```

## Add the entry to /etc/host file

Add the following entry to your /etc/host file

```bash
127.0.0.1 kafka
```

This will ensure that all services run inside or outside of docker.


## Summary of the repository

This repository contains the following services.

* Bitcoin Block Processor
* Ethereum Block Processor
* Solana Block Processor
* Bitcoin Transaction Processor
* Ethereum Transaction Processor
* Solana Transaction Processor
* Multichain Transaction Observer


## How to run the services

All services require some infrastructure to be running in docker, first an external network needs to be created, please run the following command (this only needs to be done once):

```bash
task infra:init-local-dev-env
```

Then to start the services, please run the following command:

```bash
task infra:docker-up
```

The above docker file includes a kafka topic viewer and can be viewed in the browser at http://localhost:8080/

The above docker file will automatically create the required topics for the services to run.


If you prefer to run these in docker (which uses the same external network) you can run the following command:

```bash
task build:docker:docker-up
```

## Summary of configuration options

We are using consumer groups in kafka, this is to ensure that we can run multiple instances of the same service and they will share the load of the messages.

With consumer groups we also ensure that the last message read is always remembered, this can be disabled (while developing) in the configuration file.
Removing the consumer groups will cause the same messages to be consumed when the services are restarted.


## Summary of all available tasks

```bash
task: Available tasks for this project:
* app:build:bitcoin-block-processor:                  Build the Bitcoin block processor
* app:build:bitcoin-transaction-processor:            Build the Bitcoin transaction processor
* app:build:ethereum-block-processor:                 Build the Ethereum block processor
* app:build:ethereum-transaction-processor:           Build the Ethereum transaction processor
* app:build:multichain-transaction-observer:          Build the Multichain transaction observer
* app:build:solana-block-processor:                   Build the Solana block processor
* app:build:solana-transaction-processor:             Build the Solana transaction processor
* app:delete:database:all:                            Delete all databases
* app:delete:database:bitcoin-block-processor:        Delete the Bitcoin block processor database
* app:delete:database:ethereum-block-processor:       Delete the Ethereum block processor database
* app:delete:database:solana-block-processor:         Delete the Solana block processor database
* app:start:all:                                      Start all processors and observer
* app:start:all:except-observer:                      Start all processors but not observer
* app:start:bitcoin-block-processor:                  Start the Bitcoin block processor
* app:start:bitcoin-transaction-processor:            Start the Bitcoin transaction processor
* app:start:bitcoin:all:                              Start all Bitcoin processors
* app:start:bitcoin:all:except-observer:              Start all Bitcoin processors but not observer
* app:start:ethereum-block-processor:                 Start the Ethereum block processor
* app:start:ethereum-transaction-processor:           Start the Ethereum transaction processor
* app:start:ethereum:all:                             Start all Ethereum processors
* app:start:ethereum:all:except-observer:             Start all Ethereum processors but not observer
* app:start:multichain-transaction-observer:          Start the Multichain transaction observer
* app:start:solana-block-processor:                   Start the Solana block processor
* app:start:solana-transaction-processor:             Start the Solana transaction processor
* app:start:solana:all:                               Start all Solana processors
* app:start:solana:all:except-observer:               Start all Solana processors but not observer
* build:docker:docker-build:                          Build docker services
* build:docker:docker-down:                           Take docker services down
* build:docker:docker-up:                             Bring docker services up and detach
* build:docker:docker-up:no-detach:                   Bring docker services up in foreground
* common:clean-mocks:                                 Remove mock files
* common:gen-mocks:                                   Generate mock files
* common:lint:                                        Run linters for changed files
* common:lint:allfiles:                               Run linters for all files
* common:precommit:                                   Run pre-commit checks
* common:test:                                        Run tests
* infra:docker-down:                                  Take docker services down
* infra:docker-up:                                    Bring docker services up and detach
* infra:docker-up:no-detach:                          Bring docker services up in foreground
* infra:init-local-dev-env:                           Initialize local development environment
* provision:configure-precommit:                      Setup git hook scripts for pre-commit
* provision:install-golangci-lint:                    Install golangci-lint
* provision:install-mockery:                          Install mockery
* provision:install-precommit:                        Install pre-commit
* provision:setup:                                    Setup development environment
```

## Other documentation

Please check under the directory /docs for more documentation on the services, as well has some diagrams.

There is also a discussion on other ideas.
