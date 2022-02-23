package types

import (
	"fmt"
	"strings"
)

type Providers []string

const (
	ProvidersDelimiter = ","
)

func (p *Providers) String() string {
	return strings.Join(*p, ProvidersDelimiter)
}

func (p *Providers) Set(str string) error {
	if strings.TrimSpace(str) == "" {
		return fmt.Errorf("empty providers given")
	}
	providers := strings.Split(str, ProvidersDelimiter)
	*p = providers
	return nil
}

func (*Providers) Description() string {
	return "A comma-separated list of EventPipe providers to be enabled.\n" +
		"These providers supplement any providers implied by --profile <profile-name>.\n" +
		"If there's any inconsistency for a particular provider,\n" +
		"this configuration takes precedence over the implicit configuration from the profile.\n\n" +
		"This list of providers is in the form:\n\n" +
		"* Provider[,Provider]\n" +
		"* Provider is in the form: KnownProviderName[:Flags[:Level][:KeyValueArgs]].\n" +
		"* KeyValueArgs is in the form: [key1=value1][;key2=value2].\n\n" +
		"To learn more about some of the well-known providers in .NET, refer to:\n" +
		"https://docs.microsoft.com/en-us/dotnet/core/diagnostics/well-known-event-providers"
}

func (*Providers) Type() string {
	return "providers"
}
