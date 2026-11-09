package state

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/packetflinger/libquetoo/message"
)

type ServerInfo struct {
	Server  map[string]string
	Players []struct {
		Name  string
		Score int
		Ping  int
	}
}

// Request the current serverinfo from a Quetoo server. This will include all
// CVARs marked as serverinfo and the ping/score for all connected players.
func (s *Server) FetchInfo() (ServerInfo, error) {
	p := message.ConnectionlessPacket{Data: "status"}
	out, err := p.Send(s.Address, s.Port)
	if err != nil {
		return ServerInfo{}, err
	}
	if len(out.Data) < 5 {
		return ServerInfo{}, fmt.Errorf("invalid serverinfo response, is server running?")
	}
	info, err := ParseServerinfoResponse(out.Data)
	if err != nil {
		return ServerInfo{}, fmt.Errorf("%v", err)
	}
	return info, nil
}

// Create a ServerInfo struct from the data sent back from the server
func ParseServerinfoResponse(b []byte) (ServerInfo, error) {
	lines := strings.Split(strings.Trim(string(b[4:]), " \n\t"), "\n")
	info, err := parseServerinfo(lines)
	if err != nil {
		return ServerInfo{}, err
	}
	return info, nil
}

func parseServerinfo(s []string) (ServerInfo, error) {
	si := ServerInfo{}
	info := map[string]string{}

	if len(s) < 2 {
		return ServerInfo{}, fmt.Errorf("invalid server info response")
	}

	if s[0] != "status" {
		return ServerInfo{}, fmt.Errorf("response missing 'status' prefix")
	}

	info = ParseInfoString(s[1])
	if info == nil {
		return ServerInfo{}, fmt.Errorf("nil info map")
	}

	if len(s) > 2 {
		players := s[2:]
		if len(players) > 0 {
			playernames := ""
			for _, p := range players {
				name, score, ping, ok := parsePlayer(p)
				if !ok {
					continue
				}
				playernames = fmt.Sprintf("%s,%s", playernames, name)
				si.Players = append(si.Players, struct {
					Name  string
					Score int
					Ping  int
				}{
					Name:  name,
					Score: score,
					Ping:  ping,
				})
			}
			if playernames != "" {
				info["players"] = playernames[1:]
			}
		}
		info["player_count"] = fmt.Sprintf("%d", len(si.Players))
	} else {
		info["player_count"] = "0"
	}
	si.Server = info
	return si, nil
}

// Parse and info string (in the format of:
// "key1\val1\key2\val2\key3\val3" into a key/value map)
func ParseInfoString(info string) map[string]string {
	if info == "" {
		return nil
	}
	out := make(map[string]string)
	vars := strings.Split(strings.TrimPrefix(info, "\\"), "\\")

	for i := 0; i+1 < len(vars); i += 2 {
		out[strings.ToLower(vars[i])] = vars[i+1]
	}
	return out
}

// Parse one player line from a status response. Newer servers send an info
// string ("\score\9\ping\0\name\foo\ai\1"), older ones send
// `score ping "name"`.
func parsePlayer(line string) (name string, score, ping int, ok bool) {
	if strings.HasPrefix(line, "\\") {
		kv := ParseInfoString(line)
		name, ok = kv["name"]
		score, _ = strconv.Atoi(kv["score"])
		ping, _ = strconv.Atoi(kv["ping"])
		return name, score, ping, ok
	}
	fields := strings.SplitN(line, " ", 3)
	if len(fields) < 3 {
		return "", 0, 0, false
	}
	score, _ = strconv.Atoi(fields[0])
	ping, _ = strconv.Atoi(fields[1])
	return strings.Trim(fields[2], "\""), score, ping, true
}
