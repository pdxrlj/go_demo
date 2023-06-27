package lib

import (
	"testing"
)

func TestNewGuacamoleTunnel(t *testing.T) {
	if err := NewGuacamoleTunnel(&GuacamoleTunnel{
		GuacamoleAddr: "192.168.1.223:4822",
		Protocol:      "vnc",
		Host:          "192.168.1.223",
		Port:          "5901",
		User:          "",
		Password:      "vncpassword",
		Uuid:          "",
		W:             1024,
		H:             768,
		Dpi:           128,
	}); err != nil {
		t.Errorf("NewGuacamoleTunnel() error = %v", err)
	}
}
