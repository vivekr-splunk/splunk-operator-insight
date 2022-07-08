package splunk

import (
	"context"
	"fmt"
	"github.com/go-logr/logr"
	"github.com/go-resty/resty/v2"
	gateway "github.com/splunk/splunk-operator/pkg/gateway/splunk/cluster-manager"
	//cmmodel "github.com/splunk/splunk-operator/pkg/gateway/splunk/cluster-manager/model"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"time"
)

type splunkGatewayFactory struct {
	log logr.Logger
	//credentials to log on to splunk
	credentials *gateway.SplunkCredentials
	// client for talking to splunk
	client *resty.Client
}

func NewGatewayFactory(ctx context.Context, sad *gateway.SplunkCredentials) gateway.Factory {
	factory := splunkGatewayFactory{}
	factory.log = log.FromContext(ctx).WithName("gateway").WithName("splunk")
	err := factory.init(ctx, sad)
	if err != nil {
		factory.log.Error(err, "Cannot splunk gateway endpoint")
		return nil // FIXME we have to throw some kind of exception or error here
	}
	return factory
}

func (f *splunkGatewayFactory) init(ctx context.Context, sad *gateway.SplunkCredentials) error {
	f.log.Info("splunk settings",
		"endpoint", f.credentials.Address,
		"CACertFile", f.credentials.TrustedCAFile,
		"ClientCertFile", f.credentials.ClientCertificateFile,
		"ClientPrivKeyFile", f.credentials.ClientPrivateKeyFile,
		"TLSInsecure", f.credentials.DisableCertificateVerification,
	)

	client := resty.New()
	//splunkURL := fmt.Sprintf("https://%s:%d/%s", sad.Address, sad.Port, sad.ServicesNamespace)
	splunkURL := fmt.Sprintf("https://%s:%d/services", sad.Address, sad.Port)
	client.SetBaseURL(splunkURL)
	client.SetHeader("Content-Type", "application/json")
	client.SetHeader("Accept", "application/json")
	client.SetTimeout(time.Duration(60 * time.Minute))
	client.SetDebug(true)
	f.client = client

	return nil
}

func (f splunkGatewayFactory) splunkGateway(ctx context.Context, sad *gateway.SplunkCredentials, publisher gateway.EventPublisher) (*splunkGateway, error) {
	gatewayLogger := log.FromContext(ctx)

	gatewayLogger.Info("new splunk manager created to access rest endpoint")
	newGateway := &splunkGateway{
		credentials: sad,
		client:      f.client,
		publisher:   publisher,
	}
	return newGateway, nil
}

// NewGateway returns a new Splunk Gateway using the global
// configuration for finding the Splunk services.
func (f splunkGatewayFactory) NewGateway(ctx context.Context, sad *gateway.SplunkCredentials, publisher gateway.EventPublisher) (gateway.Gateway, error) {
	return f.splunkGateway(ctx, sad, publisher)
}
