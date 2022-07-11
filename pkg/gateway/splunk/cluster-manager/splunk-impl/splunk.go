package splunk

import (
	"net/http"

	"github.com/go-logr/logr"
	"github.com/go-resty/resty/v2"
	gateway "github.com/splunk/splunk-operator/pkg/gateway/splunk/cluster-manager"
)

// splunkGateway implements the gateway.Gateway interface
// and uses gateway to manage the host.
type splunkGateway struct {
	// a logger configured for this host
	log logr.Logger
	// a debug logger configured for this host
	debugLog logr.Logger
	// an event publisher for recording significant events
	publisher gateway.EventPublisher
	// client for talking to splunk
	client *resty.Client
	// credentials
	credentials *gateway.SplunkCredentials
}

func (f splunkGateway) GetClusterConfig() error {
	url := "/cluster/config"
	resp, err := f.client.R().Get(url)
	if resp.StatusCode() == http.StatusOK {
		f.log.Info("response success set to", "result", resp.Result())
	} else {
		f.log.Info("response failure set to", "result", err)
	}
	return nil
}

func (f splunkGateway) GetIndexVolumes() error {
	url := "/data/index-volumes"
	resp, err := f.client.R().Get(url)
	if resp.StatusCode() == http.StatusOK {
		f.log.Info("response success set to", "result", resp.Result())
	} else {
		f.log.Info("response failure set to", "result", err)
	}
	return nil
}

func (f splunkGateway) GetDeploymentHealth() error {
	url := "/server/health/deployment"
	resp, err := f.client.R().Get(url)
	if resp.StatusCode() == http.StatusOK {
		f.log.Info("response success set to", "result", resp.Result())
	} else {
		f.log.Info("response failure set to", "result", err)
	}
	return nil
}