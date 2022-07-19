package fixutre

import (
	"context"
	//"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/go-logr/logr"
	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	gateway "github.com/splunk/splunk-operator/pkg/gateway/splunk/indexer"
	model "github.com/splunk/splunk-operator/pkg/gateway/splunk/model"
	clustermodel "github.com/splunk/splunk-operator/pkg/gateway/splunk/model/services/cluster"
	managermodel "github.com/splunk/splunk-operator/pkg/gateway/splunk/model/services/cluster/manager"
	peermodel "github.com/splunk/splunk-operator/pkg/gateway/splunk/model/services/cluster/peer"
	commonmodel "github.com/splunk/splunk-operator/pkg/gateway/splunk/model/services/common"
	logz "sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var log = logz.New().WithName("gateway").WithName("fixture")

// fixtureGateway implements the gateway.fixtureGateway interface
// and uses splunk to manage the host.
type fixtureGateway struct {
	// client for talking to splunk
	client *resty.Client
	// the splunk credentials
	credentials gateway.SplunkCredentials
	// a logger configured for this host
	log logr.Logger
	// an event publisher for recording significant events
	publisher gateway.EventPublisher
	// state of the splunk
	state *Fixture
}

// Fixture contains persistent state for a particular splunk instance
type Fixture struct {
}

// NewGateway returns a new Fixture Gateway
func (f *Fixture) NewGateway(ctx context.Context, sad *gateway.SplunkCredentials, publisher gateway.EventPublisher) (gateway.Gateway, error) {
	p := &fixtureGateway{
		log:       log.WithValues("splunk", sad.Address),
		publisher: publisher,
		state:     f,
		client:    resty.New(),
	}
	return p, nil
}

// Access cluster node configuration details.
// endpoint: https://<host>:<mPort>/services/cluster/config
func (p *fixtureGateway) GetClusterConfig(context context.Context) (*[]clustermodel.ClusterConfigContent, error) {
	// Read entire file content, giving us little control but
	// making it very simple. No need to close the file.
	content, err := ioutil.ReadFile("cluster_manager_buckets.json")
	if err != nil {
		log.Error(err, "fixture: error in get cluster config")
		return nil, err
	}
	httpmock.ActivateNonDefault(p.client.GetClient())
	fixtureData := string(content)
	responder := httpmock.NewStringResponder(200, fixtureData)
	fakeUrl := clustermodel.GetClusterConfigUrl
	httpmock.RegisterResponder("GET", fakeUrl, responder)

	// featch the configheader into struct
	envelop := &clustermodel.ClusterConfigHeader{}
	_, err = p.client.R().SetResult(envelop).Get(fakeUrl)
	contentList := []clustermodel.ClusterConfigContent{}
	for _, entry := range envelop.Entry {
		contentList = append(contentList, entry.Content)
	}
	return &contentList, nil
}

// Provides bucket configuration information for a cluster manager node.
// endpoint: https://<host>:<mPort>/services/cluster/manager/buckets
func (p *fixtureGateway) GetClusterManagerBuckets(context context.Context) (*[]managermodel.ClusterManagerBucketContent, error) {
	// Read entire file content, giving us little control but
	// making it very simple. No need to close the file.
	content, err := ioutil.ReadFile("cluster_config.json")
	if err != nil {
		log.Error(err, "fixture: error in get cluster config")
		return nil, err
	}
	httpmock.ActivateNonDefault(p.client.GetClient())
	fixtureData := string(content)
	responder := httpmock.NewStringResponder(200, fixtureData)
	fakeUrl := clustermodel.GetClusterManagerBucketUrl
	httpmock.RegisterResponder("GET", fakeUrl, responder)
	// featch the configheader into struct
	splunkError := &model.SplunkError{}
	envelop := &managermodel.ClusterManagerBucketHeader{}
	resp, err := p.client.R().
		SetResult(envelop).
		SetError(&splunkError).
		ForceContentType("application/json").
		Get(fakeUrl)
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
func (p *fixtureGateway) GetClusterManagerGeneration(context context.Context) (*[]managermodel.ClusterManagerGenerationContent, error) {
	// Read entire file content, giving us little control but
	// making it very simple. No need to close the file.
	content, err := ioutil.ReadFile("cluster_config.json")
	if err != nil {
		log.Error(err, "fixture: error in get cluster config")
		return nil, err
	}
	httpmock.ActivateNonDefault(p.client.GetClient())
	fixtureData := string(content)
	responder := httpmock.NewStringResponder(200, fixtureData)
	fakeUrl := clustermodel.GetClusterManagerGenerationUrl
	httpmock.RegisterResponder("GET", fakeUrl, responder)
	// featch the configheader into struct
	splunkError := &model.SplunkError{}
	envelop := &managermodel.ClusterManagerGenerationHeader{}
	resp, err := p.client.R().
		SetResult(envelop).
		SetError(&splunkError).
		ForceContentType("application/json").
		Get(fakeUrl)
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
func (p *fixtureGateway) GetClusterManagerHAActiveStatus(context context.Context) error {
	return nil
}

// Performs health checks to determine the cluster health and search impact, prior to a rolling upgrade of the indexer cluster.
// Authentication and Authorization:
// 		Requires the admin role or list_indexer_cluster capability.
// endpoint: https://<host>:<mPort>/services/cluster/manager/health
func (p *fixtureGateway) GetClusterManagerHealth(context context.Context) (*[]managermodel.ClusterManagerHealthContent, error) {
	// Read entire file content, giving us little control but
	// making it very simple. No need to close the file.
	content, err := ioutil.ReadFile("cluster_config.json")
	if err != nil {
		log.Error(err, "fixture: error in get cluster config")
		return nil, err
	}
	httpmock.ActivateNonDefault(p.client.GetClient())
	fixtureData := string(content)
	responder := httpmock.NewStringResponder(200, fixtureData)
	fakeUrl := clustermodel.GetClusterManagerHealthUrl
	httpmock.RegisterResponder("GET", fakeUrl, responder)

	// featch the configheader into struct
	splunkError := &model.SplunkError{}
	envelop := &managermodel.ClusterManagerHealthHeader{}
	resp, err := p.client.R().
		SetResult(envelop).
		SetError(&splunkError).
		ForceContentType("application/json").
		Get(fakeUrl)
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
func (p *fixtureGateway) GetClusterManagerIndexes(context context.Context) (*[]managermodel.ClusterManagerIndexesContent, error) {
	// Read entire file content, giving us little control but
	// making it very simple. No need to close the file.
	content, err := ioutil.ReadFile("cluster_config.json")
	if err != nil {
		log.Error(err, "fixture: error in get cluster config")
		return nil, err
	}
	httpmock.ActivateNonDefault(p.client.GetClient())
	fixtureData := string(content)
	responder := httpmock.NewStringResponder(200, fixtureData)
	fakeUrl := clustermodel.GetClusterManagerIndexesUrl
	httpmock.RegisterResponder("GET", fakeUrl, responder)
	// featch the configheader into struct
	splunkError := &model.SplunkError{}
	envelop := &managermodel.ClusterManagerIndexesHeader{}
	resp, err := p.client.R().
		SetResult(envelop).
		SetError(&splunkError).
		ForceContentType("application/json").
		Get(fakeUrl)
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
func (p *fixtureGateway) GetClusterManagerInfo(context context.Context) (*[]managermodel.ClusterManagerInfoContent, error) {
	// Read entire file content, giving us little control but
	// making it very simple. No need to close the file.
	content, err := ioutil.ReadFile("cluster_config.json")
	if err != nil {
		log.Error(err, "fixture: error in get cluster config")
		return nil, err
	}
	httpmock.ActivateNonDefault(p.client.GetClient())
	fixtureData := string(content)
	responder := httpmock.NewStringResponder(200, fixtureData)
	fakeUrl := clustermodel.GetClusterManagerInfoUrl
	httpmock.RegisterResponder("GET", fakeUrl, responder)
	// featch the configheader into struct
	splunkError := &model.SplunkError{}
	envelop := &managermodel.ClusterManagerInfoHeader{}
	resp, err := p.client.R().
		SetResult(envelop).
		SetError(&splunkError).
		ForceContentType("application/json").
		Get(fakeUrl)
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
func (p *fixtureGateway) GetClusterManagerPeers(context context.Context) (*[]managermodel.ClusterManagerPeerContent, error) {
	// Read entire file content, giving us little control but
	// making it very simple. No need to close the file.
	content, err := ioutil.ReadFile("cluster_config.json")
	if err != nil {
		log.Error(err, "fixture: error in get cluster config")
		return nil, err
	}
	httpmock.ActivateNonDefault(p.client.GetClient())
	fixtureData := string(content)
	responder := httpmock.NewStringResponder(200, fixtureData)
	fakeUrl := clustermodel.GetClusterManagerPeersUrl
	httpmock.RegisterResponder("GET", fakeUrl, responder)
	// featch the configheader into struct
	splunkError := &model.SplunkError{}
	envelop := &managermodel.ClusterManagerPeerHeader{}
	resp, err := p.client.R().
		SetResult(envelop).
		SetError(&splunkError).
		ForceContentType("application/json").
		Get(fakeUrl)
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
func (p *fixtureGateway) GetClusterManagerRedundancy(context context.Context) (*[]managermodel.ClusterManagerRedundancyContent, error) {
	// Read entire file content, giving us little control but
	// making it very simple. No need to close the file.
	content, err := ioutil.ReadFile("cluster_config.json")
	if err != nil {
		log.Error(err, "fixture: error in get cluster config")
		return nil, err
	}
	httpmock.ActivateNonDefault(p.client.GetClient())
	fixtureData := string(content)
	responder := httpmock.NewStringResponder(200, fixtureData)
	fakeUrl := clustermodel.GetClusterManagerRedundancyUrl
	httpmock.RegisterResponder("GET", fakeUrl, responder)
	// featch the configheader into struct
	splunkError := &model.SplunkError{}
	envelop := &managermodel.ClusterManagerRedundancyHeader{}
	resp, err := p.client.R().
		SetResult(envelop).
		SetError(&splunkError).
		ForceContentType("application/json").
		Get(fakeUrl)
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
func (p *fixtureGateway) GetClusterManagerSites(context context.Context) (*[]managermodel.ClusterManagerSiteContent, error) {
	// Read entire file content, giving us little control but
	// making it very simple. No need to close the file.
	content, err := ioutil.ReadFile("cluster_config.json")
	if err != nil {
		log.Error(err, "fixture: error in get cluster config")
		return nil, err
	}
	httpmock.ActivateNonDefault(p.client.GetClient())
	fixtureData := string(content)
	responder := httpmock.NewStringResponder(200, fixtureData)
	fakeUrl := clustermodel.GetClusterManagerSitesUrl
	httpmock.RegisterResponder("GET", fakeUrl, responder)
	// featch the configheader into struct
	splunkError := &model.SplunkError{}
	envelop := &managermodel.ClusterManagerSiteHeader{}
	resp, err := p.client.R().
		SetResult(envelop).
		SetError(&splunkError).
		ForceContentType("application/json").
		Get(fakeUrl)
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
func (p *fixtureGateway) GetClusterManagerStatus(context context.Context) (*[]managermodel.ClusterManagerStatusContent, error) {
	// Read entire file content, giving us little control but
	// making it very simple. No need to close the file.
	content, err := ioutil.ReadFile("cluster_config.json")
	if err != nil {
		log.Error(err, "fixture: error in get cluster config")
		return nil, err
	}
	httpmock.ActivateNonDefault(p.client.GetClient())
	fixtureData := string(content)
	responder := httpmock.NewStringResponder(200, fixtureData)
	fakeUrl := clustermodel.GetClusterManagerStatusUrl
	httpmock.RegisterResponder("GET", fakeUrl, responder)
	// featch the configheader into struct
	splunkError := &model.SplunkError{}
	envelop := &managermodel.ClusterManagerStatusHeader{}
	resp, err := p.client.R().
		SetResult(envelop).
		SetError(&splunkError).
		ForceContentType("application/json").
		Get(fakeUrl)
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
func (p *fixtureGateway) GetClusterPeerBuckets(context context.Context) (*[]peermodel.ClusterPeerBucket, error) {
	// Read entire file content, giving us little control but
	// making it very simple. No need to close the file.
	content, err := ioutil.ReadFile("cluster_config.json")
	if err != nil {
		log.Error(err, "fixture: error in get cluster config")
		return nil, err
	}
	httpmock.ActivateNonDefault(p.client.GetClient())
	fixtureData := string(content)
	responder := httpmock.NewStringResponder(200, fixtureData)
	fakeUrl := clustermodel.GetClusterPeerBucketsUrl
	httpmock.RegisterResponder("GET", fakeUrl, responder)
	// featch the configheader into struct
	splunkError := &model.SplunkError{}
	envelop := &commonmodel.Header{}
	resp, err := p.client.R().
		SetResult(envelop).
		SetError(&splunkError).
		ForceContentType("application/json").
		Get(fakeUrl)
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
func (p *fixtureGateway) GetClusterPeerInfo(context context.Context) (*[]peermodel.ClusterPeerInfo, error) {
	// Read entire file content, giving us little control but
	// making it very simple. No need to close the file.
	content, err := ioutil.ReadFile("cluster_config.json")
	if err != nil {
		log.Error(err, "fixture: error in get cluster config")
		return nil, err
	}
	httpmock.ActivateNonDefault(p.client.GetClient())
	fixtureData := string(content)
	responder := httpmock.NewStringResponder(200, fixtureData)
	fakeUrl := clustermodel.GetClusterPeerInfoUrl
	httpmock.RegisterResponder("GET", fakeUrl, responder)
	// featch the configheader into struct
	splunkError := &model.SplunkError{}
	envelop := &commonmodel.Header{}
	resp, err := p.client.R().
		SetResult(envelop).
		SetError(&splunkError).
		ForceContentType("application/json").
		Get(fakeUrl)
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
