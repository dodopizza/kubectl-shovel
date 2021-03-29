package kubernetes

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/kubectl/pkg/scheme"

	// import auth plugins to make oidc auth work
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

// Client is kubernetes client set wrapper
type Client struct {
	Namespace string
	*kubernetes.Clientset
	*rest.Config
}

// NewClient return new kubernetes client
func NewClient(clientGetter genericclioptions.RESTClientGetter) (*Client, error) {
	restConfig, err := clientGetter.ToRESTConfig()
	if err != nil {
		return nil, err
	}
	restConfig.GroupVersion = &schema.GroupVersion{Group: "", Version: "v1"}
	if restConfig.APIPath == "" {
		restConfig.APIPath = "/api"
	}
	if restConfig.NegotiatedSerializer == nil {
		restConfig.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	}
	if len(restConfig.UserAgent) == 0 {
		restConfig.UserAgent = rest.DefaultKubernetesUserAgent()
	}

	clientSet, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}

	ns, _, err := clientGetter.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		return nil, err
	}

	return &Client{ns, clientSet, restConfig}, nil
}
