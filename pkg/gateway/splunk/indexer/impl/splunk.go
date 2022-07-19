package impl

import (
	"context"
	"net/http"
	//"encoding/json"
	"github.com/go-logr/logr"
	"github.com/go-resty/resty/v2"

	gateway "github.com/splunk/splunk-operator/pkg/gateway/splunk/indexer"
	model "github.com/splunk/splunk-operator/pkg/gateway/splunk/model"
	clustermodel "github.com/splunk/splunk-operator/pkg/gateway/splunk/model/services/cluster"
	managermodel "github.com/splunk/splunk-operator/pkg/gateway/splunk/model/services/cluster/manager"
	peermodel "github.com/splunk/splunk-operator/pkg/gateway/splunk/model/services/cluster/peer"
	commonmodel "github.com/splunk/splunk-operator/pkg/gateway/splunk/model/services/common"
	//logz "sigs.k8s.io/controller-runtime/pkg/log/zap"
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

// Access cluster node configuration details.
// endpoint: https://<host>:<mPort>/services/cluster/config
func (p *splunkGateway) GetClusterConfig(context context.Context) (*[]clustermodel.ClusterConfigContent, error) {
	url := clustermodel.GetClusterConfigUrl

	// featch the configheader into struct
	splunkError := &model.SplunkError{}
	envelop := &clustermodel.ClusterConfigHeader{}
	resp, err := p.client.R().
		SetResult(envelop).
		SetError(&splunkError).
		ForceContentType("application/json").
		Get(url)
	if err != nil {
		p.log.Error(err, "get cluster configuration failed")
	}
	if resp.StatusCode() != http.StatusOK {
		p.log.Info("response failure set to", "result", err)
	}
	if resp.StatusCode() > 400 {
		p.log.Info("response failure set to", "result", splunkError.Messages[0].Text)
		return nil, splunkError
	}

	contentList := []clustermodel.ClusterConfigContent{}
	for _, entry := range envelop.Entry {
		contentList = append(contentList, entry.Content)
	}
	return &contentList, nil
}

// Provides bucket configuration information for a cluster manager node.
// endpoint: https://<host>:<mPort>/services/cluster/manager/buckets
func (p *splunkGateway) GetClusterManagerBuckets(context context.Context) (*[]managermodel.ClusterManagerBucketContent, error) {
	url := clustermodel.GetClusterManagerBucketUrl

	// featch the configheader into struct
	splunkError := &model.SplunkError{}
	envelop := &managermodel.ClusterManagerBucketHeader{}
	resp, err := p.client.R().
		SetResult(envelop).
		SetError(&splunkError).
		ForceContentType("application/json").
		Get(url)
	if err != nil {
		p.log.Error(err, "get cluster manager buckets failed")
	}
	if resp.StatusCode() != http.StatusOK {
		p.log.Info("response failure set to", "result", err)
	}
	if resp.StatusCode() > 400 {
		p.log.Info("response failure set to", "result", splunkError.Messages[0].Text)
		return nil, splunkError
	}

	contentList := []managermodel.ClusterManagerBucketContent{}
	for _, entry := range envelop.Entry {
		contentList = append(contentList, entry.Content)
	}
	return &contentList, nil
}

// Access current generation cluster manager information and create a cluster generation.
// List peer nodes participating in the current generation for this manager.
// endpoint: https://<host>:<mPort>/services/cluster/manager/generation
func (p *splunkGateway) GetClusterManagerGeneration(context context.Context) (*[]managermodel.ClusterManagerGenerationContent, error) {
	url := clustermodel.GetClusterManagerGenerationUrl

	// featch the configheader into struct
	splunkError := &model.SplunkError{}
	envelop := &managermodel.ClusterManagerGenerationHeader{}
	resp, err := p.client.R().
		SetResult(envelop).
		SetError(&splunkError).
		ForceContentType("application/json").
		Get(url)
	if err != nil {
		p.log.Error(err, "get cluster manager buckets failed")
	}
	if resp.StatusCode() != http.StatusOK {
		p.log.Info("response failure set to", "result", err)
	}
	if resp.StatusCode() > 400 {
		p.log.Info("response failure set to", "result", splunkError.Messages[0].Text)
		return nil, splunkError
	}

	contentList := []managermodel.ClusterManagerGenerationContent{}
	for _, entry := range envelop.Entry {
		contentList = append(contentList, entry.Content)
	}
	return &contentList, nil
}

// Used by the load balancers to check the high availability mode of a given cluster manager.
// The active cluster manager will return "HTTP 200", denoting "healthy", and a startup or standby cluster manager will return "HTTP 503".
// Authentication and authorization:
// 	This endpoint is unauthenticated because some load balancers don't support authentication on a health check endpoint.
// endpoint: https://<host>:<mPort>/services/cluster/manager/ha_active_status
// FIXME TODO, not sure how the structure looks
func (p *splunkGateway) GetClusterManagerHAActiveStatus(context context.Context) error {
	return nil
}

// Performs health checks to determine the cluster health and search impact, prior to a rolling upgrade of the indexer cluster.
// Authentication and Authorization:
// 		Requires the admin role or list_indexer_cluster capability.
// endpoint: https://<host>:<mPort>/services/cluster/manager/health
func (p *splunkGateway) GetClusterManagerHealth(context context.Context) (*[]managermodel.ClusterManagerHealthContent, error) {
	url := clustermodel.GetClusterManagerHealthUrl

	// featch the configheader into struct
	splunkError := &model.SplunkError{}
	envelop := &managermodel.ClusterManagerHealthHeader{}
	resp, err := p.client.R().
		SetResult(envelop).
		SetError(&splunkError).
		ForceContentType("application/json").
		Get(url)
	if err != nil {
		p.log.Error(err, "get cluster manager buckets failed")
	}
	if resp.StatusCode() != http.StatusOK {
		p.log.Info("response failure set to", "result", err)
	}
	if resp.StatusCode() > 400 {
		p.log.Info("response failure set to", "result", splunkError.Messages[0].Text)
		return nil, splunkError
	}

	contentList := []managermodel.ClusterManagerHealthContent{}
	for _, entry := range envelop.Entry {
		contentList = append(contentList, entry.Content)
	}
	return &contentList, nil
}

// Access cluster index information.
// endpoint: https://<host>:<mPort>/services/cluster/manager/indexes
func (p *splunkGateway) GetClusterManagerIndexes(context context.Context) (*[]managermodel.ClusterManagerIndexesContent, error) {
	url := clustermodel.GetClusterManagerIndexesUrl

	// featch the configheader into struct
	splunkError := &model.SplunkError{}
	envelop := &managermodel.ClusterManagerIndexesHeader{}
	resp, err := p.client.R().
		SetResult(envelop).
		SetError(&splunkError).
		ForceContentType("application/json").
		Get(url)
	if err != nil {
		p.log.Error(err, "get cluster manager buckets failed")
	}
	if resp.StatusCode() != http.StatusOK {
		p.log.Info("response failure set to", "result", err)
	}
	if resp.StatusCode() > 400 {
		p.log.Info("response failure set to", "result", splunkError.Messages[0].Text)
		return nil, splunkError
	}

	contentList := []managermodel.ClusterManagerIndexesContent{}
	for _, entry := range envelop.Entry {
		contentList = append(contentList, entry.Content)
	}
	return &contentList, nil
}

// Access information about cluster manager node.
// get List cluster manager node details.
// endpoint: https://<host>:<mPort>/services/cluster/manager/info
func (p *splunkGateway) GetClusterManagerInfo(context context.Context) (*[]managermodel.ClusterManagerInfoContent, error) {
	url := clustermodel.GetClusterManagerInfoUrl

	// featch the configheader into struct
	splunkError := &model.SplunkError{}
	envelop := &managermodel.ClusterManagerInfoHeader{}
	resp, err := p.client.R().
		SetResult(envelop).
		SetError(&splunkError).
		ForceContentType("application/json").
		Get(url)
	if err != nil {
		p.log.Error(err, "get cluster manager buckets failed")
	}
	if resp.StatusCode() != http.StatusOK {
		p.log.Info("response failure set to", "result", err)
	}
	if resp.StatusCode() > 400 {
		p.log.Info("response failure set to", "result", splunkError.Messages[0].Text)
		return nil, splunkError
	}

	contentList := []managermodel.ClusterManagerInfoContent{}
	for _, entry := range envelop.Entry {
		contentList = append(contentList, entry.Content)
	}
	return &contentList, nil
}

// Access cluster manager peers.
// endpoint: https://<host>:<mPort>/services/cluster/manager/peers
func (p *splunkGateway) GetClusterManagerPeers(context context.Context) (*[]managermodel.ClusterManagerPeerContent, error) {
	url := clustermodel.GetClusterManagerPeersUrl

	// featch the configheader into struct
	splunkError := &model.SplunkError{}
	envelop := &managermodel.ClusterManagerPeerHeader{}
	resp, err := p.client.R().
		SetResult(envelop).
		SetError(&splunkError).
		ForceContentType("application/json").
		Get(url)
	if err != nil {
		p.log.Error(err, "get cluster manager buckets failed")
	}
	if resp.StatusCode() != http.StatusOK {
		p.log.Info("response failure set to", "result", err)
	}
	if resp.StatusCode() > 400 {
		p.log.Info("response failure set to", "result", splunkError.Messages[0].Text)
		return nil, splunkError
	}

	contentList := []managermodel.ClusterManagerPeerContent{}
	for _, entry := range envelop.Entry {
		contentList = append(contentList, entry.Content)
	}
	return &contentList, nil
}

// Display the details of all cluster managers participating in cluster manager redundancy, and switch the HA state of the cluster managers.
// Authentication and authorization
//		The GET on this endpoint needs the capability list_indexer_cluster, and the POST on this endpoint needs the capability edit_indexer_cluster.
// GET Display the details of all cluster managers participating in cluster manager redundancy
// endpoint: https://<host>:<mPort>/services/cluster/manager/redundancy
func (p *splunkGateway) GetClusterManagerRedundancy(context context.Context) (*[]managermodel.ClusterManagerRedundancyContent, error) {
	url := clustermodel.GetClusterManagerRedundancyUrl

	// featch the configheader into struct
	splunkError := &model.SplunkError{}
	envelop := &managermodel.ClusterManagerRedundancyHeader{}
	resp, err := p.client.R().
		SetResult(envelop).
		SetError(&splunkError).
		ForceContentType("application/json").
		Get(url)
	if err != nil {
		p.log.Error(err, "get cluster manager buckets failed")
	}
	if resp.StatusCode() != http.StatusOK {
		p.log.Info("response failure set to", "result", err)
	}
	if resp.StatusCode() > 400 {
		p.log.Info("response failure set to", "result", splunkError.Messages[0].Text)
		return nil, splunkError
	}

	contentList := []managermodel.ClusterManagerRedundancyContent{}
	for _, entry := range envelop.Entry {
		contentList = append(contentList, entry.Content)
	}
	return &contentList, nil
}

// Access cluster site information.
// list List available cluster sites.
// endpoint: https://<host>:<mPort>/services/cluster/manager/sites
func (p *splunkGateway) GetClusterManagerSites(context context.Context) (*[]managermodel.ClusterManagerSiteContent, error) {
	url := clustermodel.GetClusterManagerSitesUrl

	// featch the configheader into struct
	splunkError := &model.SplunkError{}
	envelop := &managermodel.ClusterManagerSiteHeader{}
	resp, err := p.client.R().
		SetResult(envelop).
		SetError(&splunkError).
		ForceContentType("application/json").
		Get(url)
	if err != nil {
		p.log.Error(err, "get cluster manager buckets failed")
	}
	if resp.StatusCode() != http.StatusOK {
		p.log.Info("response failure set to", "result", err)
	}
	if resp.StatusCode() > 400 {
		p.log.Info("response failure set to", "result", splunkError.Messages[0].Text)
		return nil, splunkError
	}

	contentList := []managermodel.ClusterManagerSiteContent{}
	for _, entry := range envelop.Entry {
		contentList = append(contentList, entry.Content)
	}
	return &contentList, nil
}

// Endpoint to get the status of a rolling restart.
// GET the status of a rolling restart.
// endpoint: https://<host>:<mPort>/services/cluster/manager/status
func (p *splunkGateway) GetClusterManagerStatus(context context.Context) (*[]managermodel.ClusterManagerStatusContent, error) {
	url := clustermodel.GetClusterManagerStatusUrl

	// featch the configheader into struct
	splunkError := &model.SplunkError{}
	envelop := &managermodel.ClusterManagerStatusHeader{}
	resp, err := p.client.R().
		SetResult(envelop).
		SetError(&splunkError).
		ForceContentType("application/json").
		Get(url)
	if err != nil {
		p.log.Error(err, "get cluster manager buckets failed")
	}
	if resp.StatusCode() != http.StatusOK {
		p.log.Info("response failure set to", "result", err)
	}
	if resp.StatusCode() > 400 {
		p.log.Info("response failure set to", "result", splunkError.Messages[0].Text)
		return nil, splunkError
	}

	contentList := []managermodel.ClusterManagerStatusContent{}
	for _, entry := range envelop.Entry {
		contentList = append(contentList, entry.Content)
	}
	return &contentList, nil
}

// Access cluster peers bucket configuration.
// GET
// List cluster peers bucket configuration.
// endpoint: https://<host>:<mPort>/services/cluster/peer/buckets
func (p *splunkGateway) GetClusterPeerBuckets(context context.Context) (*[]peermodel.ClusterPeerBucket, error) {
	url := clustermodel.GetClusterPeerBucketsUrl

	// featch the configheader into struct
	splunkError := &model.SplunkError{}
	envelop := &commonmodel.Header{}
	resp, err := p.client.R().
		SetResult(envelop).
		SetError(&splunkError).
		ForceContentType("application/json").
		Get(url)
	if err != nil {
		p.log.Error(err, "get cluster manager buckets failed")
	}
	if resp.StatusCode() != http.StatusOK {
		p.log.Info("response failure set to", "result", err)
	}
	if resp.StatusCode() > 400 {
		p.log.Info("response failure set to", "result", splunkError.Messages[0].Text)
		return nil, splunkError
	}

	contentList := []peermodel.ClusterPeerBucket{}
	for _, entry := range envelop.Entry {
		contentList = append(contentList, entry.Content.(peermodel.ClusterPeerBucket))
	}
	return &contentList, nil
}

// Manage peer buckets.
// GET
// List peer specified bucket information.
// endpoint: https://<host>:<mPort>/services/cluster/peer/buckets/{name}
func (p *splunkGateway) GetClusterPeerInfo(context context.Context) (*[]peermodel.ClusterPeerInfo, error) {
	url := clustermodel.GetClusterPeerInfoUrl

	// featch the configheader into struct
	splunkError := &model.SplunkError{}
	envelop := &commonmodel.Header{}
	resp, err := p.client.R().
		SetResult(envelop).
		SetError(&splunkError).
		ForceContentType("application/json").
		Get(url)
	if err != nil {
		p.log.Error(err, "get cluster manager buckets failed")
	}
	if resp.StatusCode() != http.StatusOK {
		p.log.Info("response failure set to", "result", err)
	}
	if resp.StatusCode() > 400 {
		p.log.Info("response failure set to", "result", splunkError.Messages[0].Text)
		return nil, splunkError
	}

	contentList := []peermodel.ClusterPeerInfo{}
	for _, entry := range envelop.Entry {
		contentList = append(contentList, entry.Content.(peermodel.ClusterPeerInfo))
	}
	return &contentList, nil
}
