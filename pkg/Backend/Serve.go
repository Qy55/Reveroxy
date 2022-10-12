package Backend

import (
	Firewall "Reveroxy/pkg/Firewall"
	Logger "Reveroxy/pkg/Logger"
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"strconv"
	"strings"
)

func chkFirewall(_strRemoteAddr string) (string, bool) {
	strRemoteIP := strings.Split(_strRemoteAddr, ":")[0]
	lstIP := strings.Split(strRemoteIP, ".")
	var lstInput []int
	for nInd := range lstIP {
		nDigit, err := strconv.ParseInt(lstIP[nInd], 10, 16)
		if err != nil {
			return "Wrong IP", false
		}
		lstInput = append(lstInput, int(nDigit))
	}
	return Firewall.FW.FilterByIP(lstInput)
}

func findServer(_strHostName string) (*App, string, bool) {
	if strings.HasSuffix(_strHostName, ":443") || strings.HasSuffix(_strHostName, ":80") {
		_strHostName = strings.Split(_strHostName, ":")[0]
	}
	for strHost, servTarget := range mapNameServer {
		if strings.HasSuffix(_strHostName, strHost) || strHost == _strHostName {
			return servTarget, strHost, true
		} else if strings.HasSuffix(_strHostName+strServicePort, strHost) {
			return servTarget, strHost, true
		}
	}
	return nil, "", false
}

func connCtx(ctx context.Context, c net.Conn) context.Context {
	strResult, bResult := chkFirewall(c.RemoteAddr().String())
	if bResult {
		return context.TODO()
	} else {
		ctxCancel, funcCancel := context.WithCancel(ctx)
		funcCancel()
		c.Close()
		Logger.WriteError(
			"[ Firewall :", strResult, "]",
			c.RemoteAddr().String(),
			": TCP Connection Closed")
		return ctxCancel
	}
}

func Serve(l net.Listener, handler http.Handler) error {
	srv := &http.Server{
		Handler:     handler,
		ConnContext: connCtx}
	return srv.Serve(l)
}

func ServeTLS(l net.Listener, handler http.Handler, certFile, keyFile string) error {
	tlsconfig := &tls.Config{}
	var err error
	// protos := []string{"http/1.0", "http/2.0"}
	// tlsconfig.NextProtos := []string{"http/1.0", "http/2.0"}

	tlsconfig.NextProtos = append(tlsconfig.NextProtos, "http/1.1")
	tlsconfig.NextProtos = append(tlsconfig.NextProtos, "h2")
	tlsconfig.Certificates = make([]tls.Certificate, 1)
	tlsconfig.Certificates[0], err = tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return err
	}
	srv := &http.Server{
		Handler:     handler,
		ConnContext: connCtx,
		TLSConfig:   tlsconfig}
	return srv.ServeTLS(l, certFile, keyFile)
}

func RedirectToTLS(_rw http.ResponseWriter, _req *http.Request) {

	strResult, bResult := chkFirewall(_req.RemoteAddr)

	if bResult {
		_, _, bExist := findServer(_req.Host)
		if bExist {
			Logger.WriteLog(
				"[ Reveroxy : NoTLS ]",
				_req.RemoteAddr, ">>", _req.Host,
				": HTTP Responsed with Redirect(307)")
			// http.Redirect(_rw, _req, "https://"+_req.Host+"/"+_req.RequestURI, http.StatusTemporaryRedirect)
			http.Redirect(_rw, _req, "https://"+_req.Host, http.StatusTemporaryRedirect)
		} else {
			strHeader := ""
			for key, val := range _req.Header {
				strHeader = strHeader + "\n\t" + key + ": " + strings.Join(val, ",")
			}
			Logger.WriteError(
				"[ Reveroxy : NoService ]",
				_req.RemoteAddr, ">>", _req.Host,
				": HTTP Responsed with Bad Request(400)",
				"\nRequested Header >>",
				strHeader)
			http.Error(_rw, "", http.StatusBadRequest)
		}
	} else {
		Logger.WriteError(
			"[ Firewall :", strResult, "]",
			_req.RemoteAddr,
			": HTTP Responsed with Bad Request(400)")
		http.Error(_rw, "", http.StatusBadRequest)
	}
}

func HandleTLS(_rw http.ResponseWriter, _req *http.Request) {

	strResult, bResult := chkFirewall(_req.RemoteAddr)

	if bResult {
		servTarget, strProxyName, bExist := findServer(_req.Host)
		if bExist {
			servTarget.Serve(strProxyName, _rw, _req)
		} else {
			strHeader := ""
			for key, val := range _req.Header {
				strHeader = strHeader + "\n\t" + key + ": " + strings.Join(val, ",")
			}
			Logger.WriteError(
				"[ Reveroxy : NoService ]",
				_req.RemoteAddr, ">>", _req.Host,
				": HTTPs Responsed with Bad Request(400)",
				"\nRequested Header >>",
				strHeader)
			http.Error(_rw, "", http.StatusBadRequest)
		}
	} else {
		Logger.WriteError(
			"[ Firewall :", strResult, "]",
			_req.RemoteAddr,
			": HTTPs Responsed with Bad Request(400)")
		http.Error(_rw, "", http.StatusBadRequest)

	}
}
