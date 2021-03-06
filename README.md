# kubectl shovel

[![Testing and publishing](https://github.com/dodopizza/kubectl-shovel/workflows/Testing%20and%20publishing/badge.svg)](https://github.com/dodopizza/kubectl-shovel/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/dodopizza/kubectl-shovel)](https://goreportcard.com/report/github.com/dodopizza/kubectl-shovel)
[![FOSSA Status](https://app.fossa.com/api/projects/custom%2B20998%2Fgit%40github.com%3Adodopizza%2Fkubectl-shovel.git.svg?type=shield)](https://app.fossa.com/projects/custom%2B20998%2Fgit%40github.com%3Adodopizza%2Fkubectl-shovel.git?ref=badge_shield)
[![GitHub Release](https://img.shields.io/github/release/dodopizza/kubectl-shovel.svg?style=flat)](https://github.com/dodopizza/kubectl-shovel/releases)

Plugin for kubectl that will help you to gather diagnostic info from running in Kubernetes dotnet applications.
It can work with .NET Core 3.0+ applications and Kubernetes clusters with docker or containerd runtime.

At the moment the following diagnostic tools are supported:

* `dotnet-gcdump`
* `dotnet-trace`

Inspired by [`kubectl-flame`](https://github.com/VerizonMedia/kubectl-flame).

## Installation

### Krew

You can install `kubectl shovel` via [`krew`](https://krew.sigs.k8s.io/).
At first install `krew` if you don't have it yet following the guide - [Installing](https://krew.sigs.k8s.io/docs/user-guide/setup/install/).
Then you will be able to install `shovel` plugin:

```
kubectl krew install shovel
```

### Precompiled binaries

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

Most of dotnet tools flags supported as well to use, e.g `--duration` and `--format` for `trace`.
You can find more info and examples in [cli documentation](./cli/docs/kubectl-shovel.md) or by using `-h/--help` flag.

## How it works

It runs the job with specified tool on the specified pod's node and mount its `/tmp` folder with dotnet-diagnostic socket.
So it requires permissions to get pods and create jobs and allowance to mount `/var/lib/docker` path for docker runtime and `/run/containerd` for containerd from a host in read-only mode.

## Development

To run all kinds of checks and generators please use:

```
make prepare
```

### Prerequisites

* golang
* docker
* kind

### Testing

#### Unit tests:

```
make test-unit
```

#### Integration tests

> kind-clusters use containerd as container runtime, so functionality with docker-runtime won't be covered.

* Integration tests require running kind-cluster. You can create it with `kind create cluster`. Also you can specify some version for cluster: `kind create cluster --image=kindest/node:<version>`, e.g v1.19.1 version.
* Then run integration tests with `make test-integration`. It will:
  * Build docker image for dumper
  * Upload it to kind-cluster
  * Build CLI
  * Run tests using these artifacts

#### All in one

```
make test
```
