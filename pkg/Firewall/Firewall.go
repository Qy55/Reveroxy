package Firewall

var (
	strAllow     string = "Allowed"
	strDenied    string = "Denied"
	strNotAllowd string = "NotAllowed"
)

type Firewall struct {
	ndRootAllow *Node
	ndRootDeny  *Node
}

func (fw Firewall) FilterByIP(_ipInput []int) (string, bool) {
	ndInput := NewNode(32, _ipInput)
	if fw.ndRootAllow.Check(ndInput) {
		if fw.ndRootDeny.Check(ndInput) {
			return strDenied, false
		} else {
			return strAllow, true
		}
	} else {
		return strNotAllowd, false
	}
}

func (fw *Firewall) AddRule(_strRule string, _ndRule *Node) bool {
	if _strRule == RuleAllow {
		fw.ndRootAllow.Append(_ndRule)
		return true

	} else if _strRule == RuleDeny {
		fw.ndRootDeny.Append(_ndRule)
		return true

	}

	return false

}
