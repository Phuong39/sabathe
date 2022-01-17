package register

import (
	"context"
	"errors"
	"fmt"
	"github.com/smallnest/rpcx/protocol"
	"github.com/smallnest/rpcx/server"
	"log"
	"sabathe/client/color"
	"sabathe/server/transport"
)

// RegisterService register rpcx server
func RegisterService(address string) {
	s := server.NewServer()
	//s.RegisterName("sabathe", new(Service), "")
	err := s.Register(new(Service), "")
	if err != nil {
		log.Fatal("[x]", err)
	}
	s.AuthFunc = auth
	color.PrintlnYellow(fmt.Sprintf("[*] start a rpcx service in %v", address))
	transport.RpcServer = s
	err = s.Serve("tcp", address)
	if err != nil {
		log.Fatal("[x]", err)
	}


}

// auth use token auth rpc client
func auth(ctx context.Context, req *protocol.Message, token string) error {

	if token == "TGzv3JOkF0XG5Qx2TlKwi" {
		return nil
	}

	return errors.New("invalid token")
}
