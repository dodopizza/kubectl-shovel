## kubectl-shovel gcdump

Get dotnet-gcdump results

### Synopsis

This subcommand will run dotnet-gcdump tool for running in k8s appplication.
Result will be saved locally so you'll be able to analyze it with appropriate tools.
You can find more info about dotnet-gcdump tool by the following links:

	* https://devblogs.microsoft.com/dotnet/collecting-and-analyzing-memory-dumps/
	* https://docs.microsoft.com/en-us/dotnet/core/diagnostics/dotnet-gcdump

```
kubectl-shovel gcdump [flags]
```

### Examples

```
The only required flag is `--pod-name`. So you can use it like this:

	kubectl shovel gcdump --pod-name my-app-65c4fc589c-gznql

Use `-o`/`--output` to define name of dump file:

	kubectl shovel gcdump --pod-name my-app-65c4fc589c-gznql -o ./myapp.gcdump

Also use `-n`/`--namespace` if your pod is not in current context's namespace:

	kubectl shovel gcdump --pod-name my-app-65c4fc589c-gznql -n default
```

### Options

```
      --as string                      Username to impersonate for the operation
      --as-group stringArray           Group to impersonate for the operation, this flag can be repeated to specify multiple groups.
      --cache-dir string               Default cache directory (default "/Users/signal/.kube/cache")
      --certificate-authority string   Path to a cert file for the certificate authority
      --client-certificate string      Path to a client certificate file for TLS
      --client-key string              Path to a client key file for TLS
      --cluster string                 The name of the kubeconfig cluster to use
  -c, --container string               Target container in pod. Required if pod run multiple containers
      --context string                 The name of the kubeconfig context to use
  -h, --help                           help for gcdump
      --image string                   Image of dumper to use for job (default "dodopizza/kubectl-shovel-dumper:undefined")
      --insecure-skip-tls-verify       If true, the server's certificate will not be checked for validity. This will make your HTTPS connections insecure
      --kubeconfig string              Path to the kubeconfig file to use for CLI requests.
  -n, --namespace string               If present, the namespace scope for this CLI request
  -o, --output string                  Output file (default "./output.gcdump")
      --pod-name string                Target pod
  -p, --process-id int                 The process ID to collect the trace from (default 1)
      --request-timeout string         The length of time to wait before giving up on a single server request. Non-zero values should contain a corresponding time unit (e.g. 1s, 2m, 3h). A value of zero means don't timeout requests. (default "0")
  -s, --server string                  The address and port of the Kubernetes API server
      --timeout timeout                Give up on collecting the GC dump if it takes longer than this many seconds.
                                       Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
                                       Will be rounded to seconds. If no unit provided defaults to seconds.
                                       (default 30 sec)
      --tls-server-name string         Server name to use for server certificate validation. If it is not provided, the hostname used to contact the server is used
      --token string                   Bearer token for authentication to the API server
      --user string                    The name of the kubeconfig user to use
```

### SEE ALSO

* [kubectl-shovel](kubectl-shovel.md)	 - Get diagnostics from running in k8s dotnet application

