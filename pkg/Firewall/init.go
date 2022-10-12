package Firewall

import (
	"io"
	"os"
)

var (
	FW Firewall

	RuleAllow string = "allow"
	RuleDeny  string = "deny"

	strAllowFile string = "./config/allow.conf"
	strDenyFile  string = "./config/deny.conf"
)

func init() {
	FW = Firewall{
		ndRootAllow: NewNode(8, []int{0, 0, 0, 0}),
		ndRootDeny:  NewNode(8, []int{0, 0, 0, 0}),
	}

	configFirewall(RuleAllow, strAllowFile)
	configFirewall(RuleDeny, strDenyFile)
}

func configFirewall(_strRule string, _strFileInput string) {

	fileInput, err := os.Open(_strFileInput)
	if err != nil {
		panic(err)
	}
	defer fileInput.Close()
	arrBuff := make([]byte, 1<<30)
	for {
		nCnt, err := fileInput.Read(arrBuff)
		if err != nil && err != io.EOF {
			panic(err)
		}

		if nCnt == 0 {
			break
		}

		strCIDR := ""
		nIdx := 0
		for {
			bResult := false
			if arrBuff[nIdx] == 0 {
				break
			}
			if arrBuff[nIdx] == byte('\n') || arrBuff[nIdx] == byte('\r') {
				arrBuff = arrBuff[nIdx+1:]
				nIdx = 0

				ndInput := CIDRtoNode(strCIDR)
				if ndInput != nil {
					bResult = FW.AddRule(_strRule, ndInput)
				}
				strCIDR = ""
			}
			if !bResult {
				strCIDR += string(arrBuff[nIdx])
				nIdx += 1
			}
		}

	}
}
