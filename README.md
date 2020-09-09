# kubectl shovel

Plugin for kubectl that will help you to get diagnostic info from dotnet application.
At the moment there is support for:

* `dotnet-gcdump`
* `dotnet-trace`

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
