package lib

import (
	"net"

	"github.com/pkg/errors"
)

type GuacamoleTunnel struct {
	GuacamoleAddr, Protocol, Host, Port, User, Password, Uuid string
	W, H, Dpi                                                 int
}

func NewGuacamoleTunnel(tunnelConfig *GuacamoleTunnel) (*SimpleTunnel, error) {
	config := NewGuacamoleConfig(
		WithConnectionID(tunnelConfig.Uuid),
		WithProtocol(tunnelConfig.Protocol),
		WithOptimalScreenHeight(tunnelConfig.H),
		WithOptimalScreenWidth(tunnelConfig.W),
		WithOptimalResolution(tunnelConfig.Dpi),
		WithAudioMimetypes([]string{"audio/L16", "rate=44100", "channels=2"}),
		WithParameters(map[string]string{
			"scheme":      tunnelConfig.Protocol,
			"hostname":    tunnelConfig.Host,
			"port":        tunnelConfig.Port,
			"ignore-cert": "true",
			"security":    "",
			"username":    tunnelConfig.User,
			"password":    tunnelConfig.Password,
		}),
	)

	conn, err := net.Dial("tcp", tunnelConfig.GuacamoleAddr)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	stream, err := NewStream(conn, SocketTimeout).Handshake(config)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	tunnel := NewSimpleTunnel(stream)
	return tunnel, nil
}
