package listener

import (
	"context"
	"fmt"
	"github.com/smallnest/rpcx/client"
	"sabathe/client/grumble"
	"sabathe/server/register"
)

// khls new a kcp-http listener

func InitKhlsListner(Salient client.XClient, c *grumble.Context) error {
	var resp register.Response
	var req register.Requests
	req.Listener.ListenerKCPAddress = c.Flags.String("kcp")
	req.Listener.ListenerHTTPAddress = c.Flags.String("http")
	req.Listener.Name = c.Flags.String("name")
	err := Salient.Call(context.Background(), "GetHttpListener", &req, &resp)
	if err != nil {
		fmt.Println("😰", err)
		return err
	}
	if resp.Message == "success" {
		fmt.Println(fmt.Sprintf("💀 start kcp in %v http in %v", c.Flags.String("kcp"), c.Flags.String("http")))
	}
	if resp.Error != nil {
		fmt.Println("😰", resp.Error)
		return nil
	}
	return nil
}
