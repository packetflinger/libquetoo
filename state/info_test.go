package state

import (
	"encoding/hex"
	"maps"
	"testing"
)

func TestServerinfoResponse(t *testing.T) {
	tests := []struct {
		name string
		data string
		want ServerInfo
	}{
		{
			name: "test1",
			data: "ffffffff7374617475730a6275696c645c7838365f36342d70632d6c696e75782d676e755c656469746f725c300a",
			want: ServerInfo{
				Server: map[string]string{
					"build":        "x86_64-pc-linux-gnu",
					"editor":       "0",
					"player_count": "0",
				},
				Players: nil,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b, err := hex.DecodeString(tc.data)
			if err != nil {
				t.Error(err)
			}
			si, err := ParseServerinfoResponse(b)
			if err != nil {
				t.Error(err)
			}
			if !maps.Equal(tc.want.Server, si.Server) {
				t.Errorf("info maps not equal:\ngot: %v\nwant: %v\n", si.Server, tc.want.Server)
			}
		})
	}
}
