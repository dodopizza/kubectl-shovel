## kubectl-shovel trace

Get dotnet-trace results

### Synopsis

This subcommand will capture runtime events with dotnet-trace tool for running in k8s application.
Result will be saved locally in nettrace format so you'll be able to convert it and analyze with appropriate tools.
You can find more info about dotnet-trace tool by the following links:

	* https://github.com/dotnet/diagnostics/blob/master/documentation/dotnet-trace-instructions.md
	* https://docs.microsoft.com/en-us/dotnet/core/diagnostics/dotnet-trace

```
kubectl-shovel trace [flags]
```

### Examples

```
The only required flag is `--pod-name`. So you can use it like this:

	kubectl shovel trace --pod-name my-app-65c4fc589c-gznql

Use `-o`/`--output` to define name of dump file:

	kubectl shovel trace --pod-name my-app-65c4fc589c-gznql -o ./myapp.trace

Also use `-n`/`--namespace` if your pod is not in current context's namespace:

	kubectl shovel trace --pod-name my-app-65c4fc589c-gznql -n default

Use `--duration` to define duration of trace to 30 seconds:

	kubectl shovel trace --pod-name my-app-65c4fc589c-gznql -o ./myapp.trace --duration 30s

Use `--format` to specify Speedscope format:

	kubectl shovel trace --pod-name my-app-65c4fc589c-gznql -o ./myapp.trace --format Speedscope

And then you can analyze it with https://www.speedscope.app/
Or convert any other format to speedscope format with:

	dotnet trace convert myapp.trace --format Speedscope
```

### Options

```
      --as string                      Username to impersonate for the operation. User could be a regular user or a service account in a namespace.
      --as-group stringArray           Group to impersonate for the operation, this flag can be repeated to specify multiple groups.
      --as-uid string                  UID to impersonate for the operation.
      --buffersize int                 Sets the size of the in-memory circular buffer, in megabytes (default 256)
      --cache-dir string               Default cache directory (default "/home/user/.kube/cache")
      --certificate-authority string   Path to a cert file for the certificate authority
      --client-certificate string      Path to a client certificate file for TLS
      --client-key string              Path to a client key file for TLS
      --clreventlevel clreventlevel    Verbosity of CLR events to be emitted. Supported levels:
                                       logalways, critical, error, warning, informational, verbose
      --clrevents clrevents            A list of CLR runtime provider keywords to enable separated by "+" signs.
                                       This is a simple mapping that lets you specify event keywords via string aliases rather than their hex values.
                                       More info here: https://docs.microsoft.com/en-us/dotnet/core/diagnostics/dotnet-trace#options-1
      --cluster string                 The name of the kubeconfig cluster to use
  -c, --container string               Target container in pod. Required if pod run multiple containers
      --context string                 The name of the kubeconfig context to use
      --duration duration              Trace for the given timespan and then automatically stop the trace.Provided in the form of dd:hh:mm:ss or corresponding time unit representation (e.g. 1s, 2m, 3h) (default 10s) (default 00:00:00:10)
      --format format                  Sets the output format for the trace file conversion. Supported formats:
                                       NetTrace, Chromium, Speedscope (default "NetTrace")
  -h, --help                           help for trace
      --image string                   Image of dumper to use for job (default "dodopizza/kubectl-shovel-dumper:undefined")
      --insecure-skip-tls-verify       If true, the server's certificate will not be checked for validity. This will make your HTTPS connections insecure
      --kubeconfig string              Path to the kubeconfig file to use for CLI requests.
  -n, --namespace string               If present, the namespace scope for this CLI request
  -o, --output string                  Output file (default "./output.trace")
      --pod-name string                Target pod
  -p, --process-id int                 The process ID to collect the trace from (default 1)
      --profile profile                A named pre-defined set of provider configurations that allowscommon tracing scenarios to be specified succinctly.
                                       The following profiles are available:
                                       cpu-sampling, gc-verbose, gc-collect
      --providers providers            A comma-separated list of EventPipe providers to be enabled.
                                       These providers supplement any providers implied by --profile <profile-name>.
                                       If there's any inconsistency for a particular provider,
                                       this configuration takes precedence over the implicit configuration from the profile.
                                       
                                       This list of providers is in the form:
                                       
                                       * Provider[,Provider]
                                       * Provider is in the form: KnownProviderName[:Flags[:Level][:KeyValueArgs]].
                                       * KeyValueArgs is in the form: [key1=value1][;key2=value2].
                                       
                                       To learn more about some of the well-known providers in .NET, refer to:
                                       https://docs.microsoft.com/en-us/dotnet/core/diagnostics/well-known-event-providers
      --request-timeout string         The length of time to wait before giving up on a single server request. Non-zero values should contain a corresponding time unit (e.g. 1s, 2m, 3h). A value of zero means don't timeout requests. (default "0")
  -s, --server string                  The address and port of the Kubernetes API server
  -t, --store-output-on-host           Flag, indicating that output should be stored on host /tmp folder
      --tls-server-name string         Server name to use for server certificate validation. If it is not provided, the hostname used to contact the server is used
      --token string                   Bearer token for authentication to the API server
      --user string                    The name of the kubeconfig user to use
```

### SEE ALSO

* [kubectl-shovel](kubectl-shovel.md)	 - Get diagnostics from running in k8s dotnet application

