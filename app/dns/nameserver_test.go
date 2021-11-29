package dns_test

import (
	"context"
	"testing"
	"time"

	. "github.com/qazz-shyper/website/app/dns"
	"github.com/qazz-shyper/website/common"
	dns_feature "github.com/qazz-shyper/website/features/dns"
)

func TestLocalNameServer(t *testing.T) {
	s := NewLocalNameServer()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	ips, err := s.QueryIP(ctx, "google.com", dns_feature.IPOption{
		IPv4Enable: true,
		IPv6Enable: true,
		FakeEnable: false,
	})
	cancel()
	common.Must(err)
	if len(ips) == 0 {
		t.Error("expect some ips, but got 0")
	}
}
