package state

import (
	"encoding/hex"
	"fmt"
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
			// data: "ffffffff7374617475730a6275696c645c7838365f36342d70632d6c696e75782d676e755c656469746f725c300a",
			data: "status\nbuild\\x86_64-pc-linux-gnu\\editor\\0",
			want: ServerInfo{
				Server: map[string]string{
					"build":        "x86_64-pc-linux-gnu",
					"editor":       "0",
					"player_count": "0",
				},
				Players: nil,
			},
		},
		{
			name: "test2",
			data: "status\nvar1\\val1\\var2\\val2\\var3\\val3\n1 2 player1\n4 5 player2\n6 7 name with space",
			want: ServerInfo{
				Server: map[string]string{
					"player_count": "3",
					"players":      "player1,player2,name with space",
					"var1":         "val1",
					"var2":         "val2",
					"var3":         "val3",
				},
				Players: []struct {
					Name  string
					Score int
					Ping  int
				}{
					{
						Name:  "player1",
						Score: 1,
						Ping:  2,
					},
					{
						Name:  "player2",
						Score: 4,
						Ping:  5,
					},
					{
						Name:  "name with space",
						Score: 6,
						Ping:  7,
					},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b, err := hex.DecodeString(fmt.Sprintf("%x", tc.data))
			if err != nil {
				t.Error(err)
			}
			b = append([]byte{255, 255, 255, 255}, b...)

			si, err := ParseServerinfoResponse(b)
			if err != nil {
				t.Error(err)
			}
			if !maps.Equal(tc.want.Server, si.Server) {
				t.Errorf("info maps not equal:\ngot: %v\nwant: %v\n", si.Server, tc.want.Server)
			}
			if len(si.Players) > 0 {
				for i := range si.Players {
					if si.Players[i].Name != tc.want.Players[i].Name || si.Players[i].Ping != tc.want.Players[i].Ping || si.Players[i].Score != tc.want.Players[i].Score {
						t.Errorf("Player not equal:\ngot: %v\nwant: %v\n", si.Players[i], tc.want.Players[i])
					}
				}
			}
		})
	}
}
