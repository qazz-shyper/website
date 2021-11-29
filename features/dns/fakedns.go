package dns

import (
	gonet "net"

	"github.com/qazz-shyper/website/common/net"
	"github.com/qazz-shyper/website/features"
)

type FakeDNSEngine interface {
	features.Feature
	GetFakeIPForDomain(domain string) []net.Address
	GetDomainFromFakeDNS(ip net.Address) string
	GetFakeIPRange() *gonet.IPNet
}

var FakeIPPool = "198.18.0.0/16"
