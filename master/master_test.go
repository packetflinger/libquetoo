package master

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/packetflinger/libquetoo/message"
)

func TestFetchServers(t *testing.T) {
	oob := message.ConnectionlessPacket{
		Data: "getservers",
	}
	msg, err := oob.Send("master.quetoo.org", 1996)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%s\n", hex.Dump(msg.Data))
	t.Error()
}
