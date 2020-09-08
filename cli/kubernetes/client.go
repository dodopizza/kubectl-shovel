package kubernetes

import (
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"

	// import auth plugins to make oidc auth work
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

// Client is kubernetes client set wrapper
type Client struct {
	Namespace string
	*kubernetes.Clientset
}

// NewClient return new kubernetes client
func NewClient(clientGetter genericclioptions.RESTClientGetter) (*Client, error) {
	restConfig, err := clientGetter.ToRESTConfig()
	if err != nil {
		return nil, err
	}

	clientSet, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}

	ns, _, err := clientGetter.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		return nil, err
	}

	return &Client{ns, clientSet}, nil
}
