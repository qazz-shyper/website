package core_test

import (
	"testing"

	"github.com/golang/protobuf/proto"

	"github.com/qazz-shyper/website/app/dispatcher"
	"github.com/qazz-shyper/website/app/proxyman"
	"github.com/qazz-shyper/website/common"
	"github.com/qazz-shyper/website/common/net"
	"github.com/qazz-shyper/website/common/protocol"
	"github.com/qazz-shyper/website/common/serial"
	"github.com/qazz-shyper/website/common/uuid"
	. "github.com/qazz-shyper/website/core"
	"github.com/qazz-shyper/website/features/dns"
	"github.com/qazz-shyper/website/features/dns/localdns"
	_ "github.com/qazz-shyper/website/main/distro/all"
	"github.com/qazz-shyper/website/proxy/dokodemo"
	"github.com/qazz-shyper/website/proxy/vmess"
	"github.com/qazz-shyper/website/proxy/vmess/outbound"
	"github.com/qazz-shyper/website/testing/servers/tcp"
)

func TestWebsiteDependency(t *testing.T) {
	instance := new(Instance)

	wait := make(chan bool, 1)
	instance.RequireFeatures(func(d dns.Client) {
		if d == nil {
			t.Error("expected dns client fulfilled, but actually nil")
		}
		wait <- true
	})
	instance.AddFeature(localdns.New())
	<-wait
}

func TestWebsiteClose(t *testing.T) {
	port := tcp.PickPort()

	userID := uuid.New()
	config := &Config{
		App: []*serial.TypedMessage{
			serial.ToTypedMessage(&dispatcher.Config{}),
			serial.ToTypedMessage(&proxyman.InboundConfig{}),
			serial.ToTypedMessage(&proxyman.OutboundConfig{}),
		},
		Inbound: []*InboundHandlerConfig{
			{
				ReceiverSettings: serial.ToTypedMessage(&proxyman.ReceiverConfig{
					PortRange: net.SinglePortRange(port),
					Listen:    net.NewIPOrDomain(net.LocalHostIP),
				}),
				ProxySettings: serial.ToTypedMessage(&dokodemo.Config{
					Address: net.NewIPOrDomain(net.LocalHostIP),
					Port:    uint32(0),
					NetworkList: &net.NetworkList{
						Network: []net.Network{net.Network_TCP, net.Network_UDP},
					},
				}),
			},
		},
		Outbound: []*OutboundHandlerConfig{
			{
				ProxySettings: serial.ToTypedMessage(&outbound.Config{
					Receiver: []*protocol.ServerEndpoint{
						{
							Address: net.NewIPOrDomain(net.LocalHostIP),
							Port:    uint32(0),
							User: []*protocol.User{
								{
									Account: serial.ToTypedMessage(&vmess.Account{
										Id: userID.String(),
									}),
								},
							},
						},
					},
				}),
			},
		},
	}

	cfgBytes, err := proto.Marshal(config)
	common.Must(err)

	server, err := StartInstance("protobuf", cfgBytes)
	common.Must(err)
	server.Close()
}
