package session

import (
	"context"
	"fmt"
	"github.com/smallnest/rpcx/client"
	"sabathe/client/grumble"
	"sabathe/server/register"
)

// session console

func SessionConsole(c *grumble.App, debugId *int, SaClient client.XClient) {
	c.AddCommand(&grumble.Command{
		Name: "background",
		Help: "set session in background",
		Run: func(c *grumble.Context) error {
			c.App.SetDefaultPrompt()
			c.App.Commands().Del("background")
			c.App.Commands().Del("shell")
			*debugId = 0
			return nil
		},
	})
	c.AddCommand(&grumble.Command{
		Name: "shell",
		Help: "exec system command",
		Args: func(a *grumble.Args) {
			a.String("cmd", "system command")
		},
		Run: func(c *grumble.Context) error {
			fmt.Println(c.Args["cmd"].Value)
			var req register.Requests
			req.Command = c.Args["cmd"].Value.(string)
			req.Session.SessionID = *debugId
			err := SaClient.Call(context.Background(), "RunSystemCmd", &req, nil)
			if err != nil {
				return err
			}
			return nil
		},
	})
}
