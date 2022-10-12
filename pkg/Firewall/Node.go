package Firewall

import (
	"strconv"
	"strings"
)

var (
	strSepCIDR string = "/"
	strSepIP   string = "."
)

type Node struct {
	nValue int
	nMask  int

	ndBig       *Node
	ndSmall     *Node
	ndNextLayer *Node
}

func NewNode(_nMask int, _lstValues []int) *Node {
	if _nMask <= 0 {
		return nil
	}

	ndNew := Node{
		nValue:      _lstValues[0],
		nMask:       _nMask,
		ndBig:       nil,
		ndSmall:     nil,
		ndNextLayer: NewNode(_nMask-8, _lstValues[1:])}
	if _nMask > 8 {
		ndNew.nMask = 8
	}
	return &ndNew
}

func (_nd *Node) Append(_ndNew *Node) {
	if _nd.nValue == _ndNew.nValue {
		if _nd.ndNextLayer == nil {
			_nd.ndNextLayer = _ndNew.ndNextLayer
		} else {
			_nd.ndNextLayer.Append(_ndNew.ndNextLayer)
		}
	} else {
		if _nd.nValue > _ndNew.nValue {
			if _nd.ndSmall == nil {
				_nd.ndSmall = _ndNew
			} else {
				_nd.ndSmall.Append(_ndNew)
			}
		} else if _nd.nValue < _ndNew.nValue {
			if _nd.ndBig == nil {
				_nd.ndBig = _ndNew
			} else {

				_nd.ndBig.Append(_ndNew)
			}
		}
	}
}
func (_nd Node) isContain(_nValue int) bool {
	nRange := (1 << (8 - _nd.nMask)) - 1
	if _nd.nValue <= _nValue && (_nd.nValue+nRange) >= _nValue {
		return true
	}
	return false
}

func (_nd *Node) Check(_ndInput *Node) bool {
	if _ndInput == nil {
		return false
	}

	if _nd.isContain(_ndInput.nValue) {
		if _nd.ndNextLayer == nil {
			return true
		} else {
			return _nd.ndNextLayer.Check(_ndInput.ndNextLayer)
		}
	} else {
		if _nd.nValue > _ndInput.nValue {
			if _nd.ndSmall == nil {
				return false
			}
			return _nd.ndSmall.Check(_ndInput)
		} else if _nd.nValue < _ndInput.nValue {
			if _nd.ndBig == nil {
				return false
			}
			return _nd.ndBig.Check(_ndInput)
		}
	}

	return false
}

func CIDRtoNode(_strInput string) *Node {
	var (
		lstIntIP []int = []int{}
		nMask    int
	)
	strInput := strings.Trim(_strInput, "\n")
	lstInput := strings.Split(strInput, strSepCIDR)
	if len(lstInput) != 2 {
		return nil
	}

	lstIP := strings.Split(lstInput[0], strSepIP)
	if len(lstIP) != 4 {
		return nil
	}

	nMask64, err := strconv.ParseInt(lstInput[1], 10, 16)
	if err != nil {
		return nil
	}

	for nIdx := range lstIP {
		nConv, err := strconv.ParseInt(lstIP[nIdx], 10, 16)
		if err != nil {
			return nil
		}
		nByte := int(nConv)
		lstIntIP = append(lstIntIP, nByte)
	}
	nMask = int(nMask64)
	return NewNode(nMask, lstIntIP)

}
