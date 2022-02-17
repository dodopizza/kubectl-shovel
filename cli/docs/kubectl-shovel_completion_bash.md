## kubectl-shovel completion bash

generate the autocompletion script for bash

### Synopsis


Generate the autocompletion script for the bash shell.

This script depends on the 'bash-completion' package.
If it is not installed already, you can install it via your OS's package manager.

To load completions in your current shell session:
$ source <(kubectl-shovel completion bash)

To load completions for every new session, execute once:
Linux:
  $ kubectl-shovel completion bash > /etc/bash_completion.d/kubectl-shovel
MacOS:
  $ kubectl-shovel completion bash > /usr/local/etc/bash_completion.d/kubectl-shovel

You will need to start a new shell for this setup to take effect.
  

```
kubectl-shovel completion bash
```

### Options

```
  -h, --help              help for bash
      --no-descriptions   disable completion descriptions
```

### SEE ALSO

* [kubectl-shovel completion](kubectl-shovel_completion.md)	 - generate the autocompletion script for the specified shell

