package cmd

var (
	examplesTemplate = "The only required flag is `--pod-name`. So you can use it like this:\n\n" +
		"\tkubectl shovel %[1]s --pod-name my-app-65c4fc589c-gznql\n\n" +
		"Use `-o`/`--output` to define name of dump file:\n\n" +
		"\tkubectl shovel %[1]s --pod-name my-app-65c4fc589c-gznql -o ./myapp.%[1]s\n\n" +
		"Also use `-n`/`--namespace` if your pod is not in current context's namespace:\n\n" +
		"\tkubectl shovel %[1]s --pod-name my-app-65c4fc589c-gznql -n default"
)
