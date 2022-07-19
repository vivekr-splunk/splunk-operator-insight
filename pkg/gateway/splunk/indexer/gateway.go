package indexer

import (
	"context"
	clustermodel "github.com/splunk/splunk-operator/pkg/gateway/splunk/model/services/cluster"
	managermodel "github.com/splunk/splunk-operator/pkg/gateway/splunk/model/services/cluster/manager"
	peermodel "github.com/splunk/splunk-operator/pkg/gateway/splunk/model/services/cluster/peer"
)

// SplunkCredentials contains the information necessary to communicate with
// the Splunk service
type SplunkCredentials struct {

	// Address holds the URL for splunk service
	Address string `json:"address"`

	//Port port to connect
	Port int32 `json:"port"`

	//ServicesNamespace  optional for services endpoints
	ServicesNamespace string `json:"servicesNs,omitempty"`

	//User optional for services endpoints
	User string `json:"user,omitempty"`

	//App optional for services endpoints
	App string `json:"app,omitempty"`

	//CredentialsName The name of the secret containing the Splunk credentials (requires
	// keys "username" and "password").
	CredentialsName string `json:"credentialsName"`

	//TrustedCAFile Server trusted CA file
	TrustedCAFile string `json:"trustedCAFile,omitempty"`

	//ClientCertificateFile client certification if we are using to connect to server
	ClientCertificateFile string `json:"clientCertificationFile,omitempty"`

	//ClientPrivateKeyFile client private key if we are using to connect to server
	ClientPrivateKeyFile string `json:"clientPrivateKeyFile,omitempty"`

	// DisableCertificateVerification disables verification of splunk
	// certificates when using HTTPS to connect to the Splunk.
	DisableCertificateVerification bool `json:"disableCertificateVerification,omitempty"`
}

// EventPublisher is a function type for publishing events associated
// with gateway functions.
type EventPublisher func(reason, message string)

// Factory is the interface for creating new Gateway objects.
type Factory interface {
	NewGateway(ctx context.Context, sad *SplunkCredentials, publisher EventPublisher) (Gateway, error)
}

// Gateway holds the state information for talking to
// splunk gateway backend.
type Gateway interface {

	// Access cluster node configuration details.
	// endpoint: https://<host>:<mPort>/services/cluster/config
	GetClusterConfig(context context.Context) (*[]clustermodel.ClusterConfigContent, error)

	// Provides bucket configuration information for a cluster manager node.
	// endpoint: https://<host>:<mPort>/services/cluster/manager/buckets
	GetClusterManagerBuckets(context context.Context) (*[]managermodel.ClusterManagerBucketContent, error)

	// Access current generation cluster manager information and create a cluster generation.
	// List peer nodes participating in the current generation for this manager.
	// endpoint: https://<host>:<mPort>/services/cluster/manager/generation
	GetClusterManagerGeneration(context context.Context) (*[]managermodel.ClusterManagerGenerationContent, error)

	// Used by the load balancers to check the high availability mode of a given cluster manager.
	// The active cluster manager will return "HTTP 200", denoting "healthy", and a startup or standby cluster manager will return "HTTP 503".
	// Authentication and authorization:
	// 	This endpoint is unauthenticated because some load balancers don't support authentication on a health check endpoint.
	// endpoint: https://<host>:<mPort>/services/cluster/manager/ha_active_status
	// FIXME TODO, not sure how the structure looks
	GetClusterManagerHAActiveStatus(context context.Context) error

	// Performs health checks to determine the cluster health and search impact, prior to a rolling upgrade of the indexer cluster.
	// Authentication and Authorization:
	// 		Requires the admin role or list_indexer_cluster capability.
	// endpoint: https://<host>:<mPort>/services/cluster/manager/health
	GetClusterManagerHealth(context context.Context) (*[]managermodel.ClusterManagerHealthContent, error)

	// Access cluster index information.
	// endpoint: https://<host>:<mPort>/services/cluster/manager/indexes
	GetClusterManagerIndexes(context context.Context) (*[]managermodel.ClusterManagerIndexesContent, error)

	// Access information about cluster manager node.
	// get List cluster manager node details.
	// endpoint: https://<host>:<mPort>/services/cluster/manager/info
	GetClusterManagerInfo(context context.Context) (*[]managermodel.ClusterManagerInfoContent, error)

	// Access cluster manager peers.
	// endpoint: https://<host>:<mPort>/services/cluster/manager/peers
	GetClusterManagerPeers(context context.Context) (*[]managermodel.ClusterManagerPeerContent, error)

	// Display the details of all cluster managers participating in cluster manager redundancy, and switch the HA state of the cluster managers.
	// Authentication and authorization
	//		The GET on this endpoint needs the capability list_indexer_cluster, and the POST on this endpoint needs the capability edit_indexer_cluster.
	// GET Display the details of all cluster managers participating in cluster manager redundancy
	// endpoint: https://<host>:<mPort>/services/cluster/manager/redundancy
	GetClusterManagerRedundancy(context context.Context) (*[]managermodel.ClusterManagerRedundancyContent, error)

	// Access cluster site information.
	// list List available cluster sites.
	// endpoint: https://<host>:<mPort>/services/cluster/manager/sites
	GetClusterManagerSites(context context.Context) (*[]managermodel.ClusterManagerSiteContent, error)

	// Endpoint to get the status of a rolling restart.
	// GET the status of a rolling restart.
	// endpoint: https://<host>:<mPort>/services/cluster/manager/status
	GetClusterManagerStatus(context context.Context) (*[]managermodel.ClusterManagerStatusContent, error)

	// Access cluster peers bucket configuration.
	// GET
	// List cluster peers bucket configuration.
	// endpoint: https://<host>:<mPort>/services/cluster/peer/buckets
	GetClusterPeerBuckets(context context.Context) (*[]peermodel.ClusterPeerBucket, error)

	// Manage peer buckets.
	// GET
	// List peer specified bucket information.
	// endpoint: https://<host>:<mPort>/services/cluster/peer/buckets/	}
	GetClusterPeerInfo(context context.Context) (*[]peermodel.ClusterPeerInfo, error)
}
