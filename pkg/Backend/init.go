package Backend

import (
	"os/exec"
	"strings"
)

var (
	strDomainName     string = ".c0ff33.kro.kr"
	strServicePort    string = ""
	strDomainNamePort string = strDomainName + strServicePort

	// [prefix + strDomainName] = Server
	mapNameServer map[string]*App = map[string]*App{}
)

func GetAddrWSL() string {
	cmd := exec.Command("wsl", "hostname", "-I")
	out, err := cmd.Output()

	if err != nil {
		return "Fail"
	}
	ret := strings.TrimSpace(string(out))
	return ret
}

func AddApp(_srvBackend *App) error {
	strHostName := _srvBackend.strServicePrefix + strDomainNamePort
	_, bExist := mapNameServer[strHostName]
	if !bExist {
		mapNameServer[strHostName] = _srvBackend
	} else {
		println("Already Registered Service")
	}

	return nil
}
