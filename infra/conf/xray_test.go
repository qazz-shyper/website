package conf_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/google/go-cmp/cmp"

	"github.com/qazz-shyper/website/app/dispatcher"
	"github.com/qazz-shyper/website/app/log"
	"github.com/qazz-shyper/website/app/proxyman"
	"github.com/qazz-shyper/website/app/router"
	"github.com/qazz-shyper/website/common"
	clog "github.com/qazz-shyper/website/common/log"
	"github.com/qazz-shyper/website/common/net"
	"github.com/qazz-shyper/website/common/protocol"
	"github.com/qazz-shyper/website/common/serial"
	core "github.com/qazz-shyper/website/core"
	. "github.com/qazz-shyper/website/infra/conf"
	"github.com/qazz-shyper/website/proxy/blackhole"
	dns_proxy "github.com/qazz-shyper/website/proxy/dns"
	"github.com/qazz-shyper/website/proxy/freedom"
	"github.com/qazz-shyper/website/proxy/vmess"
	"github.com/qazz-shyper/website/proxy/vmess/inbound"
	"github.com/qazz-shyper/website/transport/internet"
	"github.com/qazz-shyper/website/transport/internet/http"
	"github.com/qazz-shyper/website/transport/internet/tls"
	"github.com/qazz-shyper/website/transport/internet/websocket"
)

func TestWebsiteConfig(t *testing.T) {
	createParser := func() func(string) (proto.Message, error) {
		return func(s string) (proto.Message, error) {
			config := new(Config)
			if err := json.Unmarshal([]byte(s), config); err != nil {
				return nil, err
			}
			return config.Build()
		}
	}

	runMultiTestCase(t, []TestCase{
		{
			Input: `{
				"outbound": {
					"protocol": "freedom",
					"settings": {}
				},
				"log": {
					"access": "/var/log/xray/access.log",
					"loglevel": "error",
					"error": "/var/log/xray/error.log"
				},
				"inbound": {
					"streamSettings": {
						"network": "ws",
						"wsSettings": {
							"headers": {
								"host": "example.domain"
							},
							"path": ""
						},
						"tlsSettings": {
							"alpn": "h2"
						},
						"security": "tls"
					},
					"protocol": "vmess",
					"port": 443,
					"settings": {
						"clients": [
							{
								"security": "aes-128-gcm",
								"id": "0cdf8a45-303d-4fed-9780-29aa7f54175e"
							}
						]
					}
				},
				"inbounds": [{
					"streamSettings": {
						"network": "ws",
						"wsSettings": {
							"headers": {
								"host": "example.domain"
							},
							"path": ""
						},
						"tlsSettings": {
							"alpn": "h2"
						},
						"security": "tls"
					},
					"protocol": "vmess",
					"port": "443-500",
					"allocate": {
						"strategy": "random",
						"concurrency": 3
					},
					"settings": {
						"clients": [
							{
								"security": "aes-128-gcm",
								"id": "0cdf8a45-303d-4fed-9780-29aa7f54175e"
							}
						]
					}
				}],
				"outboundDetour": [
					{
						"tag": "blocked",
						"protocol": "blackhole"
					},
					{
						"protocol": "dns"
					}
				],
				"routing": {
					"strategy": "rules",
					"settings": {
						"rules": [
							{
								"ip": [
									"10.0.0.0/8"
								],
								"type": "field",
								"outboundTag": "blocked"
							}
						]
					}
				},
				"transport": {
					"httpSettings": {
						"path": "/test"
					}
				}
			}`,
			Parser: createParser(),
			Output: &core.Config{
				App: []*serial.TypedMessage{
					serial.ToTypedMessage(&log.Config{
						ErrorLogType:  log.LogType_File,
						ErrorLogPath:  "/var/log/xray/error.log",
						ErrorLogLevel: clog.Severity_Error,
						AccessLogType: log.LogType_File,
						AccessLogPath: "/var/log/xray/access.log",
					}),
					serial.ToTypedMessage(&dispatcher.Config{}),
					serial.ToTypedMessage(&proxyman.InboundConfig{}),
					serial.ToTypedMessage(&proxyman.OutboundConfig{}),
					serial.ToTypedMessage(&router.Config{
						DomainStrategy: router.Config_AsIs,
						Rule: []*router.RoutingRule{
							{
								Geoip: []*router.GeoIP{
									{
										Cidr: []*router.CIDR{
											{
												Ip:     []byte{10, 0, 0, 0},
												Prefix: 8,
											},
										},
									},
								},
								TargetTag: &router.RoutingRule_Tag{
									Tag: "blocked",
								},
							},
						},
					}),
				},
				Outbound: []*core.OutboundHandlerConfig{
					{
						SenderSettings: serial.ToTypedMessage(&proxyman.SenderConfig{
							StreamSettings: &internet.StreamConfig{
								ProtocolName: "tcp",
								TransportSettings: []*internet.TransportConfig{
									{
										ProtocolName: "http",
										Settings: serial.ToTypedMessage(&http.Config{
											Path: "/test",
										}),
									},
								},
							},
						}),
						ProxySettings: serial.ToTypedMessage(&freedom.Config{
							DomainStrategy: freedom.Config_AS_IS,
							UserLevel:      0,
						}),
					},
					{
						Tag: "blocked",
						SenderSettings: serial.ToTypedMessage(&proxyman.SenderConfig{
							StreamSettings: &internet.StreamConfig{
								ProtocolName: "tcp",
								TransportSettings: []*internet.TransportConfig{
									{
										ProtocolName: "http",
										Settings: serial.ToTypedMessage(&http.Config{
											Path: "/test",
										}),
									},
								},
							},
						}),
						ProxySettings: serial.ToTypedMessage(&blackhole.Config{}),
					},
					{
						SenderSettings: serial.ToTypedMessage(&proxyman.SenderConfig{
							StreamSettings: &internet.StreamConfig{
								ProtocolName: "tcp",
								TransportSettings: []*internet.TransportConfig{
									{
										ProtocolName: "http",
										Settings: serial.ToTypedMessage(&http.Config{
											Path: "/test",
										}),
									},
								},
							},
						}),
						ProxySettings: serial.ToTypedMessage(&dns_proxy.Config{
							Server: &net.Endpoint{},
						}),
					},
				},
				Inbound: []*core.InboundHandlerConfig{
					{
						ReceiverSettings: serial.ToTypedMessage(&proxyman.ReceiverConfig{
							PortList: &net.PortList{Range: []*net.PortRange{net.SinglePortRange(443)}},
							StreamSettings: &internet.StreamConfig{
								ProtocolName: "websocket",
								TransportSettings: []*internet.TransportConfig{
									{
										ProtocolName: "websocket",
										Settings: serial.ToTypedMessage(&websocket.Config{
											Header: []*websocket.Header{
												{
													Key:   "host",
													Value: "example.domain",
												},
											},
										}),
									},
									{
										ProtocolName: "http",
										Settings: serial.ToTypedMessage(&http.Config{
											Path: "/test",
										}),
									},
								},
								SecurityType: "xray.transport.internet.tls.Config",
								SecuritySettings: []*serial.TypedMessage{
									serial.ToTypedMessage(&tls.Config{
										NextProtocol: []string{"h2"},
									}),
								},
							},
						}),
						ProxySettings: serial.ToTypedMessage(&inbound.Config{
							User: []*protocol.User{
								{
									Level: 0,
									Account: serial.ToTypedMessage(&vmess.Account{
										Id: "0cdf8a45-303d-4fed-9780-29aa7f54175e",
										SecuritySettings: &protocol.SecurityConfig{
											Type: protocol.SecurityType_AES128_GCM,
										},
									}),
								},
							},
						}),
					},
					{
						ReceiverSettings: serial.ToTypedMessage(&proxyman.ReceiverConfig{
							PortList: &net.PortList{Range: []*net.PortRange{{
								From: 443,
								To:   500,
							}}},
							AllocationStrategy: &proxyman.AllocationStrategy{
								Type: proxyman.AllocationStrategy_Random,
								Concurrency: &proxyman.AllocationStrategy_AllocationStrategyConcurrency{
									Value: 3,
								},
							},
							StreamSettings: &internet.StreamConfig{
								ProtocolName: "websocket",
								TransportSettings: []*internet.TransportConfig{
									{
										ProtocolName: "websocket",
										Settings: serial.ToTypedMessage(&websocket.Config{
											Header: []*websocket.Header{
												{
													Key:   "host",
													Value: "example.domain",
												},
											},
										}),
									},
									{
										ProtocolName: "http",
										Settings: serial.ToTypedMessage(&http.Config{
											Path: "/test",
										}),
									},
								},
								SecurityType: "xray.transport.internet.tls.Config",
								SecuritySettings: []*serial.TypedMessage{
									serial.ToTypedMessage(&tls.Config{
										NextProtocol: []string{"h2"},
									}),
								},
							},
						}),
						ProxySettings: serial.ToTypedMessage(&inbound.Config{
							User: []*protocol.User{
								{
									Level: 0,
									Account: serial.ToTypedMessage(&vmess.Account{
										Id: "0cdf8a45-303d-4fed-9780-29aa7f54175e",
										SecuritySettings: &protocol.SecurityConfig{
											Type: protocol.SecurityType_AES128_GCM,
										},
									}),
								},
							},
						}),
					},
				},
			},
		},
	})
}

func TestMuxConfig_Build(t *testing.T) {
	tests := []struct {
		name   string
		fields string
		want   *proxyman.MultiplexingConfig
	}{
		{"default", `{"enabled": true, "concurrency": 16}`, &proxyman.MultiplexingConfig{
			Enabled:     true,
			Concurrency: 16,
		}},
		{"empty def", `{}`, &proxyman.MultiplexingConfig{
			Enabled:     false,
			Concurrency: 8,
		}},
		{"not enable", `{"enabled": false, "concurrency": 4}`, &proxyman.MultiplexingConfig{
			Enabled:     false,
			Concurrency: 4,
		}},
		{"forbidden", `{"enabled": false, "concurrency": -1}`, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MuxConfig{}
			common.Must(json.Unmarshal([]byte(tt.fields), m))
			if got := m.Build(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MuxConfig.Build() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_Override(t *testing.T) {
	tests := []struct {
		name string
		orig *Config
		over *Config
		fn   string
		want *Config
	}{
		{
			"combine/empty",
			&Config{},
			&Config{
				LogConfig:    &LogConfig{},
				RouterConfig: &RouterConfig{},
				DNSConfig:    &DNSConfig{},
				Transport:    &TransportConfig{},
				Policy:       &PolicyConfig{},
				API:          &APIConfig{},
				Stats:        &StatsConfig{},
				Reverse:      &ReverseConfig{},
			},
			"",
			&Config{
				LogConfig:    &LogConfig{},
				RouterConfig: &RouterConfig{},
				DNSConfig:    &DNSConfig{},
				Transport:    &TransportConfig{},
				Policy:       &PolicyConfig{},
				API:          &APIConfig{},
				Stats:        &StatsConfig{},
				Reverse:      &ReverseConfig{},
			},
		},
		{
			"combine/newattr",
			&Config{InboundConfigs: []InboundDetourConfig{{Tag: "old"}}},
			&Config{LogConfig: &LogConfig{}}, "",
			&Config{LogConfig: &LogConfig{}, InboundConfigs: []InboundDetourConfig{{Tag: "old"}}},
		},
		{
			"replace/inbounds",
			&Config{InboundConfigs: []InboundDetourConfig{{Tag: "pos0"}, {Protocol: "vmess", Tag: "pos1"}}},
			&Config{InboundConfigs: []InboundDetourConfig{{Tag: "pos1", Protocol: "kcp"}}},
			"",
			&Config{InboundConfigs: []InboundDetourConfig{{Tag: "pos0"}, {Tag: "pos1", Protocol: "kcp"}}},
		},
		{
			"replace/inbounds-replaceall",
			&Config{InboundConfigs: []InboundDetourConfig{{Tag: "pos0"}, {Protocol: "vmess", Tag: "pos1"}}},
			&Config{InboundConfigs: []InboundDetourConfig{{Tag: "pos1", Protocol: "kcp"}, {Tag: "pos2", Protocol: "kcp"}}},
			"",
			&Config{InboundConfigs: []InboundDetourConfig{{Tag: "pos1", Protocol: "kcp"}, {Tag: "pos2", Protocol: "kcp"}}},
		},
		{
			"replace/notag-append",
			&Config{InboundConfigs: []InboundDetourConfig{{}, {Protocol: "vmess"}}},
			&Config{InboundConfigs: []InboundDetourConfig{{Tag: "pos1", Protocol: "kcp"}}},
			"",
			&Config{InboundConfigs: []InboundDetourConfig{{}, {Protocol: "vmess"}, {Tag: "pos1", Protocol: "kcp"}}},
		},
		{
			"replace/outbounds",
			&Config{OutboundConfigs: []OutboundDetourConfig{{Tag: "pos0"}, {Protocol: "vmess", Tag: "pos1"}}},
			&Config{OutboundConfigs: []OutboundDetourConfig{{Tag: "pos1", Protocol: "kcp"}}},
			"",
			&Config{OutboundConfigs: []OutboundDetourConfig{{Tag: "pos0"}, {Tag: "pos1", Protocol: "kcp"}}},
		},
		{
			"replace/outbounds-prepend",
			&Config{OutboundConfigs: []OutboundDetourConfig{{Tag: "pos0"}, {Protocol: "vmess", Tag: "pos1"}}},
			&Config{OutboundConfigs: []OutboundDetourConfig{{Tag: "pos1", Protocol: "kcp"}, {Tag: "pos2", Protocol: "kcp"}}},
			"config.json",
			&Config{OutboundConfigs: []OutboundDetourConfig{{Tag: "pos1", Protocol: "kcp"}, {Tag: "pos2", Protocol: "kcp"}}},
		},
		{
			"replace/outbounds-append",
			&Config{OutboundConfigs: []OutboundDetourConfig{{Tag: "pos0"}, {Protocol: "vmess", Tag: "pos1"}}},
			&Config{OutboundConfigs: []OutboundDetourConfig{{Tag: "pos2", Protocol: "kcp"}}},
			"config_tail.json",
			&Config{OutboundConfigs: []OutboundDetourConfig{{Tag: "pos0"}, {Protocol: "vmess", Tag: "pos1"}, {Tag: "pos2", Protocol: "kcp"}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.orig.Override(tt.over, tt.fn)
			if r := cmp.Diff(tt.orig, tt.want); r != "" {
				t.Error(r)
			}
		})
	}
}
