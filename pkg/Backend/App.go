package Backend

import (
	Logger "Reveroxy/pkg/Logger"
	"errors"
	"net/http"
	"strings"
)

type App struct {
	strServicePrefix string
	strBaseDomain    string
	mapNameProxy     map[string]*ServiceHttp
	mapAddHeader     map[string]string
	GetName          func(*http.Request) string
}

func (serv *App) Init(_strPrefix string) {
	Logger.WriteLog("initializing for " + _strPrefix)
	serv.strBaseDomain = strDomainNamePort
	serv.strServicePrefix = _strPrefix
	serv.mapNameProxy = make(map[string]*ServiceHttp)
	serv.mapAddHeader = make(map[string]string)
}

func (serv *App) SetDomainName(_strDomain string, _strPort string) {
	serv.strBaseDomain = "." + _strDomain + ":" + _strPort
}

func (serv *App) AddProxy(_strName string, _prTarget *ServiceHttp) {
	serv.mapNameProxy[_strName] = _prTarget
}

func (serv *App) ConfigProxy(_strName string, _tpConfig *http.Transport) error {
	prTarget, bExist := serv.mapNameProxy[_strName]
	if !bExist {
		return errors.New("NoProxyFound" + _strName)
	}
	prTarget.ConfigProxy(_tpConfig)
	return nil
}

func (serv *App) SetModResp(_strName string, _funcTarget func(*http.Response) error) error {

	prTarget, bExist := serv.mapNameProxy[_strName]
	if !bExist {
		return errors.New("NoProxyFound" + _strName)
	}
	prTarget.SetModResp(_funcTarget)
	return nil
}

func (serv App) Serve(_strProxyName string, _rw http.ResponseWriter, _req *http.Request) {
	strProxyName := serv.GetName(_req)
	Logger.WriteLog("[", strings.ReplaceAll(_req.Host, strDomainNamePort, ""), ":", strProxyName, "]", _req.RemoteAddr, ":", _req.RequestURI)
	proxyTarget, bExist := serv.mapNameProxy[strProxyName]
	if bExist {
		proxyTarget.Serve(_rw, _req)
	} else {
		println("Not defined proxy")
	}
}
