package output

import (
	"bytes"
	"testing"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestRenderTunnels(t *testing.T) {
	tunnels := []*models.Tunnel{
		{
			ID:          "id22",
			Lhost:       "lhost",
			Lport:       "123",
			Rhost:       "rhost",
			Rport:       "124",
			LportRandom: false,
			Scheme:      "ssh",
			ACL:         "0.0.0.0",
		},
	}

	tr := &TunnelRenderer{
		ColCountCalculator: func() int {
			return 150
		},
	}

	buf := &bytes.Buffer{}
	err := tr.RenderTunnels(buf, tunnels)
	assert.NoError(t, err)
	if err != nil {
		return
	}

	actualRenderResult := RemoveEmptySpaces(buf.String())
	assert.Equal(
		t,
		"Tunnels ID CLIENT LHOST LPORT RHOST RPORT LPORTRAND SCHEME ACL id22 lhost 123 rhost 124 false ssh 0.0.0.0 ",
		actualRenderResult,
	)
}
func TestRenderTunnel(t *testing.T) {
	tunnel := &models.Tunnel{
		ID:          "id22",
		Lhost:       "lhost",
		Lport:       "123",
		Rhost:       "rhost",
		Rport:       "124",
		LportRandom: false,
		Scheme:      "ssh",
		ACL:         "0.0.0.0",
	}

	tr := &TunnelRenderer{
		ColCountCalculator: func() int {
			return 150
		},
	}

	buf := &bytes.Buffer{}
	err := tr.RenderTunnel(buf, tunnel)
	assert.NoError(t, err)
	if err != nil {
		return
	}

	actualRenderResult := RemoveEmptySpaces(buf.String())
	assert.Equal(
		t,
		"Tunnel KEY VALUE ID: id22 CLIENT: LHOST: lhost LPORT: 123 RHOST: rhost RPORT: 124 LPORT RANDOM: false SCHEME: ssh ACL: 0.0.0.0 ",
		actualRenderResult,
	)
}
