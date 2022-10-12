package Backend

import (
	"net"
	"net/http"
	"net/http/httputil"
)

type ServiceHttp struct {
	rpBackend *httputil.ReverseProxy
}
type ServiceTcp struct {
	net.Conn
	net.TCPConn
}

func (proxy *ServiceHttp) SetBackend(_rpBackend *httputil.ReverseProxy) {
	proxy.rpBackend = _rpBackend
}

func (proxy *ServiceHttp) ConfigProxy(_tpConfig *http.Transport) {
	proxy.rpBackend.Transport = _tpConfig
}

func (proxy *ServiceHttp) SetModResp(_funcTarget func(*http.Response) error) {
	proxy.rpBackend.ModifyResponse = _funcTarget
}

func (proxy *ServiceHttp) Serve(_rw http.ResponseWriter, _req *http.Request) {
	proxy.rpBackend.ServeHTTP(_rw, _req)
}
