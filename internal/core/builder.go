package core

import (
	"net/netip"

	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common/json/badoption"

	"go-vless-client/internal/config"
)

// buildOptions формирует конфигурацию sing-box из настроек сервера и приложения.
func buildOptions(srv config.ServerConfig, settings config.AppSettings) (option.Options, error) {
	opts := option.Options{
		Log: &option.LogOptions{
			Disabled: true,
		},
	}

	vlessOut, err := buildVLESSOutbound(srv)
	if err != nil {
		return option.Options{}, err
	}

	opts.Outbounds = []option.Outbound{
		vlessOut,
		{Type: "direct", Tag: "direct"},
	}

	switch settings.Mode {
	case config.ModeProxy:
		opts.Inbounds = buildProxyInbounds(settings)
	case config.ModeVPN:
		opts.Inbounds = []option.Inbound{buildTUNInbound()}
	case config.ModeBoth:
		opts.Inbounds = append(buildProxyInbounds(settings), buildTUNInbound())
	}

	return opts, nil
}

func buildVLESSOutbound(srv config.ServerConfig) (option.Outbound, error) {
	vlessOpts := option.VLESSOutboundOptions{
		ServerOptions: option.ServerOptions{
			Server:     srv.Address,
			ServerPort: uint16(srv.Port),
		},
		UUID: srv.UUID,
		Flow: srv.Flow,
	}

	if srv.TLS {
		sni := srv.SNI
		if sni == "" {
			sni = srv.Address
		}
		vlessOpts.TLS = &option.OutboundTLSOptions{
			Enabled:    true,
			ServerName: sni,
		}
	}

	if srv.Transport == config.TransportWS {
		vlessOpts.Transport = &option.V2RayTransportOptions{
			Type: "ws",
			WebsocketOptions: option.V2RayWebsocketOptions{
				Path: srv.Path,
			},
		}
	}

	return option.Outbound{
		Type:    "vless",
		Tag:     "vless-out",
		Options: vlessOpts,
	}, nil
}

func buildProxyInbounds(settings config.AppSettings) []option.Inbound {
	localhost := listenAddr("127.0.0.1")

	socks := option.Inbound{
		Type: "socks",
		Tag:  "socks-in",
		Options: option.SocksInboundOptions{
			ListenOptions: option.ListenOptions{
				Listen:     localhost,
				ListenPort: uint16(settings.Socks5Port),
			},
		},
	}

	http := option.Inbound{
		Type: "http",
		Tag:  "http-in",
		Options: option.HTTPMixedInboundOptions{
			ListenOptions: option.ListenOptions{
				Listen:     localhost,
				ListenPort: uint16(settings.HTTPPort),
			},
		},
	}

	return []option.Inbound{socks, http}
}

func buildTUNInbound() option.Inbound {
	return option.Inbound{
		Type: "tun",
		Tag:  "tun-in",
		Options: option.TunInboundOptions{
			InterfaceName: "tun0",
			Address: badoption.Listable[netip.Prefix]{
				netip.MustParsePrefix("172.19.0.1/30"),
			},
			AutoRoute:   true,
			StrictRoute: true,
			Stack:       "system",
		},
	}
}

func listenAddr(addr string) *badoption.Addr {
	a := badoption.Addr(netip.MustParseAddr(addr))
	return &a
}