package all

import (
	// The following are necessary as they register handlers in their init functions.

	// Required features. Can't remove unless there is replacements.
	_ "github.com/qazz-shyper/website/app/dispatcher"
	_ "github.com/qazz-shyper/website/app/proxyman/inbound"
	_ "github.com/qazz-shyper/website/app/proxyman/outbound"

	// Default commander and all its services. This is an optional feature.
	_ "github.com/qazz-shyper/website/app/commander"
	_ "github.com/qazz-shyper/website/app/log/command"
	_ "github.com/qazz-shyper/website/app/proxyman/command"
	_ "github.com/qazz-shyper/website/app/stats/command"

	// Other optional features.
	_ "github.com/qazz-shyper/website/app/dns"
	_ "github.com/qazz-shyper/website/app/dns/fakedns"
	_ "github.com/qazz-shyper/website/app/log"
	_ "github.com/qazz-shyper/website/app/policy"
	_ "github.com/qazz-shyper/website/app/reverse"
	_ "github.com/qazz-shyper/website/app/router"
	_ "github.com/qazz-shyper/website/app/stats"

	// Inbound and outbound proxies.
	_ "github.com/qazz-shyper/website/proxy/blackhole"
	_ "github.com/qazz-shyper/website/proxy/dns"
	_ "github.com/qazz-shyper/website/proxy/dokodemo"
	_ "github.com/qazz-shyper/website/proxy/freedom"
	_ "github.com/qazz-shyper/website/proxy/http"
	_ "github.com/qazz-shyper/website/proxy/mtproto"
	_ "github.com/qazz-shyper/website/proxy/shadowsocks"
	_ "github.com/qazz-shyper/website/proxy/socks"
	_ "github.com/qazz-shyper/website/proxy/trojan"
	_ "github.com/qazz-shyper/website/proxy/vless/inbound"
	_ "github.com/qazz-shyper/website/proxy/vless/outbound"
	_ "github.com/qazz-shyper/website/proxy/vmess/inbound"
	_ "github.com/qazz-shyper/website/proxy/vmess/outbound"

	// Transports
	_ "github.com/qazz-shyper/website/transport/internet/domainsocket"
	_ "github.com/qazz-shyper/website/transport/internet/grpc"
	_ "github.com/qazz-shyper/website/transport/internet/http"
	_ "github.com/qazz-shyper/website/transport/internet/kcp"
	_ "github.com/qazz-shyper/website/transport/internet/quic"
	_ "github.com/qazz-shyper/website/transport/internet/tcp"
	_ "github.com/qazz-shyper/website/transport/internet/tls"
	_ "github.com/qazz-shyper/website/transport/internet/udp"
	_ "github.com/qazz-shyper/website/transport/internet/websocket"
	_ "github.com/qazz-shyper/website/transport/internet/xtls"

	// Transport headers
	_ "github.com/qazz-shyper/website/transport/internet/headers/http"
	_ "github.com/qazz-shyper/website/transport/internet/headers/noop"
	_ "github.com/qazz-shyper/website/transport/internet/headers/srtp"
	_ "github.com/qazz-shyper/website/transport/internet/headers/tls"
	_ "github.com/qazz-shyper/website/transport/internet/headers/utp"
	_ "github.com/qazz-shyper/website/transport/internet/headers/wechat"
	_ "github.com/qazz-shyper/website/transport/internet/headers/wireguard"

	// JSON & TOML & YAML
	_ "github.com/qazz-shyper/website/main/json"
	_ "github.com/qazz-shyper/website/main/toml"
	_ "github.com/qazz-shyper/website/main/yaml"

	// Load config from file or http(s)
	_ "github.com/qazz-shyper/website/main/confloader/external"

	// Commands
	_ "github.com/qazz-shyper/website/main/commands/all"
)
