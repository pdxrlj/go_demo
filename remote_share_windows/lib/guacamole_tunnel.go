package lib

import (
	"fmt"
	"net"

	"github.com/pkg/errors"
)

type GuacamoleTunnel struct {
	GuacamoleAddr, Protocol, Host, Port, User, Password, Uuid string
	W, H, Dpi                                                 int
}

func NewGuacamoleTunnel(tunnelConfig *GuacamoleTunnel) error {
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
		return errors.WithStack(err)
	}
	stream := NewStream(conn, SocketTimeout).Handshake(config)
	fmt.Println("stream", stream)
	return nil
}
