# kubectl shovel

[![Go Report Card](https://goreportcard.com/badge/github.com/dodopizza/kubectl-shovel)](https://goreportcard.com/report/github.com/dodopizza/kubectl-shovel)

Plugin for kubectl that will help you to gather diagnostic info from dotnet application.
It can work with .Net Core 3.0+ applications and Kubernetes clusters with docker runtime.

At the moment the following diagnostic tools are supported:

* `dotnet-gcdump`
* `dotnet-trace`

Inspired by [`kubectl-flame`](https://github.com/VerizonMedia/kubectl-flame).

## Installation

You can find latest release on repository [release page](https://github.com/dodopizza/kubectl-shovel/releases).
Once you download compatible with your OS binary, move it to any directory specified in your `$PATH`.

## Usage

Feel free to use it as a kubectl plugin or standalone executable (`kubectl shovel`/`kubectl-shovel`)

Get gcdump:

```shell
kubectl shovel gcdump --pod-name pod-name-74df554df7-qldq7 -o ./dump.gcdump
```

Or trace:

```shell
kubectl shovel trace --pod-name pod-name-74df554df7-qldq7 -o ./trace.nettrace
```

You can find more info in [cli documentation](./cli/docs/kubectl-shovel.md) or by using `-h/--help` flag.

# How it works

It will run the job on the specified pod's node and mount its `/tmp` folder with dotnet-diagnostic socket. After this specified diagnostic tool will be launched.
