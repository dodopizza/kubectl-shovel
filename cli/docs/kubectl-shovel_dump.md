## kubectl-shovel dump

Get dotnet-dump results

### Synopsis

This subcommand will run dump tool for running in k8s application.
Result will be saved locally (or on host) so you'll be able to analyze it with appropriate instruments.
Tool specific additional arguments are also supported.
You can find more info about this tool by the following links:

	* https://docs.microsoft.com/en-us/dotnet/core/diagnostics/dotnet-dump
	* https://docs.microsoft.com/en-us/dotnet/core/diagnostics/debug-linux-dumps

```
kubectl-shovel dump [flags]
```

### Examples

```
The only required flag is `--pod-name`. So you can use it like this:

	kubectl shovel dump --pod-name my-app-65c4fc589c-gznql

Use `-o`/`--output` to define name of dump file:

	kubectl shovel dump --pod-name my-app-65c4fc589c-gznql -o ./myapp.dump

Also use `-n`/`--namespace` if your pod is not in current context's namespace:

	kubectl shovel dump --pod-name my-app-65c4fc589c-gznql -n default
```

### Options

```
      --as string                      Username to impersonate for the operation. User could be a regular user or a service account in a namespace.
      --as-group stringArray           Group to impersonate for the operation, this flag can be repeated to specify multiple groups.
      --as-uid string                  UID to impersonate for the operation.
      --cache-dir string               Default cache directory (default "/home/user/.kube/cache")
      --certificate-authority string   Path to a cert file for the certificate authority
      --client-certificate string      Path to a client certificate file for TLS
      --client-key string              Path to a client key file for TLS
      --cluster string                 The name of the kubeconfig cluster to use
  -c, --container string               Target container in pod. Required if pod run multiple containers
      --context string                 The name of the kubeconfig context to use
      --diag                           Enable dump collection diagnostic logging
  -h, --help                           help for dump
      --image string                   Image of dumper to use for job (default "dodopizza/kubectl-shovel-dumper:undefined")
      --insecure-skip-tls-verify       If true, the server's certificate will not be checked for validity. This will make your HTTPS connections insecure
      --kubeconfig string              Path to the kubeconfig file to use for CLI requests.
      --limit-cpu string               Limit maximal consumptions cpu for the executing job
      --limit-memory string            Limit maximal consumptions memory for the executing job
  -n, --namespace string               If present, the namespace scope for this CLI request
  -o, --output string                  Output file (default "./output.dump")
      --output-host-path string        Host folder, where will be stored artifact (default "/tmp/kubectl-shovel")
      --pod-name string                Target pod
  -p, --process-id int                 The process ID to collect the trace from (default 1)
      --request-timeout string         The length of time to wait before giving up on a single server request. Non-zero values should contain a corresponding time unit (e.g. 1s, 2m, 3h). A value of zero means don't timeout requests. (default "0")
  -s, --server string                  The address and port of the Kubernetes API server
  -t, --store-output-on-host           Store output on node instead of downloading it locally
      --tls-server-name string         Server name to use for server certificate validation. If it is not provided, the hostname used to contact the server is used
      --token string                   Bearer token for authentication to the API server
      --type type                      The kinds of information that are collected from process. Supported types:
                                       Full, Heap, Mini, Triage
                                       Full - The largest dump containing all memory including the module images
                                       Heap - A large and relatively comprehensive dump containing module lists, thread lists, all stacks, exception information and all memory except for mapped images
                                       Mini - A small dump containing module lists, thread lists, exception information and all stacks
                                       Triage - A small dump containing minimal information (default Full)
      --user string                    The name of the kubeconfig user to use
```

### SEE ALSO

* [kubectl-shovel](kubectl-shovel.md)	 - Get diagnostics from running in k8s dotnet application

