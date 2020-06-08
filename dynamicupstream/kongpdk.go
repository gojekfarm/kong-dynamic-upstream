package dynamicupstream

import (
	"github.com/Kong/go-pdk"
)

type KongPDK interface {
	RequestPath() (string, error)
	SetUpstreamTarget(string, int) error
	SetUpstreamTargetRequestPath(string) error
}

type KongPDKWrapper struct {
	Kong *pdk.PDK
}

func (r *KongPDKWrapper) RequestPath() (string, error) {
	return r.Kong.Request.GetPath()
}

func (r *KongPDKWrapper) SetUpstreamTarget(host string, port int) error {
	return r.Kong.Service.SetTarget(host, port)
}

func (r *KongPDKWrapper) SetUpstreamTargetRequestPath(path string) error {
	return r.Kong.ServiceRequest.SetPath(path)
}
