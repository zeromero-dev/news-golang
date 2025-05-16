# Project test-news

One Paragraph of project description goes here

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

## MakeFile

Note: If templ is not installing through the script, you should manually add the GOPATH.
You can do this by adding the following line to your `.zshrc` or `.bashrc` file: `export PATH="$PATH:$HOME/go/bin"`
Then run `source ~/.zshrc` or `source ~/.bashrc` to apply the changes.

````bash

Run build make command with tests
```bash
make all
````

Build the application

```bash
make build
```

Run the application

```bash
make run
```

Create DB container

```bash
make docker-run
```

Shutdown DB Container

```bash
make docker-down
```

DB Integrations Test:

```bash
make itest
```

Live reload the application:

```bash
make watch
```

Run the test suite:

```bash
make test
```

Clean up binary from the last build:

```bash
make clean
```
