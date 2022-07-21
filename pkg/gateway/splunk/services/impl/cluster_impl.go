package impl

import (
	"context"
	"net/http"

	splunkmodel "github.com/splunk/splunk-operator/pkg/gateway/splunk/model"
	clustermodel "github.com/splunk/splunk-operator/pkg/gateway/splunk/model/services/cluster"
	peermodel "github.com/splunk/splunk-operator/pkg/gateway/splunk/model/services/cluster/peer"
	commonmodel "github.com/splunk/splunk-operator/pkg/gateway/splunk/model/services/common"
	//logz "sigs.k8s.io/controller-runtime/pkg/log/zap"
)

// GetClusterPeerBuckets Access cluster peers bucket configuration.
// GET
// List cluster peers bucket configuration.
// endpoint: https://<host>:<mPort>/services/cluster/peer/buckets
func (p *splunkGateway) GetClusterPeerBuckets(context context.Context) (*[]peermodel.ClusterPeerBucket, error) {
	url := clustermodel.GetClusterPeerBucketsUrl

	// featch the configheader into struct
	splunkError := &splunkmodel.SplunkError{}
	envelop := &commonmodel.Header{}
	resp, err := p.client.R().
		SetResult(envelop).
		SetError(&splunkError).
		ForceContentType("application/json").
		SetQueryParams(map[string]string{"output_mode": "json", "count": "0"}).
		Get(url)
	if err != nil {
		p.log.Error(err, "get cluster manager buckets failed")
	}
	if resp.StatusCode() != http.StatusOK {
		p.log.Info("response failure set to", "result", err)
	}
	if resp.StatusCode() > 400 {
		if len(splunkError.Messages) > 0 {
			p.log.Info("response failure set to", "result", splunkError.Messages[0].Text)
		}
		return nil, splunkError
	}

	contentList := []peermodel.ClusterPeerBucket{}
	for _, entry := range envelop.Entry {
		contentList = append(contentList, entry.Content.(peermodel.ClusterPeerBucket))
	}
	return &contentList, err
}

// GetClusterPeerInfo Manage peer buckets.
// GET
// List peer specified bucket information.
// endpoint: https://<host>:<mPort>/services/cluster/peer/buckets/{name}
func (p *splunkGateway) GetClusterPeerInfo(context context.Context) (*[]peermodel.ClusterPeerInfo, error) {
	url := clustermodel.GetClusterPeerInfoUrl

	// featch the configheader into struct
	splunkError := &splunkmodel.SplunkError{}
	envelop := &commonmodel.Header{}
	resp, err := p.client.R().
		SetResult(envelop).
		SetError(&splunkError).
		ForceContentType("application/json").
		SetQueryParams(map[string]string{"output_mode": "json", "count": "0"}).
		Get(url)
	if err != nil {
		p.log.Error(err, "get cluster manager info failed")
	}
	if resp.StatusCode() != http.StatusOK {
		p.log.Info("response failure set to", "result", err)
	}
	if resp.StatusCode() > 400 {
		if len(splunkError.Messages) > 0 {
			p.log.Info("response failure set to", "result", splunkError.Messages[0].Text)
		}
		return nil, splunkError
	}

	contentList := []peermodel.ClusterPeerInfo{}
	for _, entry := range envelop.Entry {
		contentList = append(contentList, entry.Content.(peermodel.ClusterPeerInfo))
	}
	return &contentList, err
}
