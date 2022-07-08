package fixutre

import (
	"context"

	"github.com/go-logr/logr"
	gateway "github.com/splunk/splunk-operator/pkg/gateway/splunk/cluster-manager"
	logz "sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var log = logz.New().WithName("gateway").WithName("fixture")

// fixtureGateway implements the gateway.fixtureGateway interface
// and uses splunk to manage the host.
type fixtureGateway struct {
	// the splunk credentials
	credentials gateway.SplunkCredentials
	// a logger configured for this host
	log logr.Logger
	// an event publisher for recording significant events
	publisher gateway.EventPublisher
	// state of the splunk
	state *Fixture
}

// Fixture contains persistent state for a particular host
type Fixture struct {
}

// NewGateway returns a new Fixture Gateway
func (f *Fixture) NewGateway(ctx context.Context, sad *gateway.SplunkCredentials, publisher gateway.EventPublisher) (gateway.Gateway, error) {
	p := &fixtureGateway{
		log:       log.WithValues("splunk", sad.Address),
		publisher: publisher,
		state:     f,
	}
	return p, nil
}

func (p fixtureGateway) GetClusterConfig() error {
	return nil
}
