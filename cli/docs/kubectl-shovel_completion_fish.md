## kubectl-shovel completion fish

generate the autocompletion script for fish

### Synopsis


Generate the autocompletion script for the fish shell.

To load completions in your current shell session:
$ kubectl-shovel completion fish | source

To load completions for every new session, execute once:
$ kubectl-shovel completion fish > ~/.config/fish/completions/kubectl-shovel.fish

You will need to start a new shell for this setup to take effect.


```
kubectl-shovel completion fish [flags]
```

### Options

```
  -h, --help              help for fish
      --no-descriptions   disable completion descriptions
```

### SEE ALSO

* [kubectl-shovel completion](kubectl-shovel_completion.md)	 - generate the autocompletion script for the specified shell

