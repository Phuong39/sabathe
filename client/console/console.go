package console

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/protocol"
	"os"
	color2 "sabathe/client/color"
	"sabathe/client/console/listener"
	"sabathe/client/console/session"
	"sabathe/client/grumble"
	"sabathe/server/register"
	"strings"
)

// console and submit command

var App = grumble.New(&grumble.Config{
	Name:                  "sabathe",
	Description:           "",
	PromptColor:           color.New(),
	HelpSubCommands:       true,
	HelpHeadlineUnderline: true,
	HelpHeadlineColor:     color.New(),
})

var MsgChan = make(chan *protocol.Message)

// use tcp to connect rpc server

func ServiceStart(addr string) client.XClient {
	c, _ := client.NewPeer2PeerDiscovery(fmt.Sprintf("tcp@%v", addr), "")
	Salient := client.NewBidirectionalXClient("Service", client.Failtry, client.RandomSelect, c, client.DefaultOption, MsgChan)
	Salient.Auth("TGzv3JOkF0XG5Qx2TlKwi") //token
	var request register.Requests
	var response register.Response
	request.AuthRequest = "kali123"
	err := Salient.Call(context.Background(), "AuthUserLogin", &request, &response)
	if err != nil || response.AuthResponse == false {
		fmt.Println("[!] Auth User Login is failed")
		os.Exit(1)
	}
	fmt.Println("[+] Auth User Login is success")
	return Salient
}

// create shell server

func ServiceConsole(address string) {
	Salient := ServiceStart(address)
	var debugSessionID = 0 // use set debug session id
	// listen lost message
	go func() {
		for msg := range MsgChan {
			if (strings.Contains(string(msg.Payload), "lost a session") || strings.Contains(string(msg.Payload), "session is close")) && strings.Contains(string(msg.Payload), fmt.Sprintf("%v", debugSessionID)) {
				debugSessionID = 0
				App.Commands().Del("background")
				App.Commands().Del("shell")
				App.SetDefaultPrompt()
			}
			fmt.Println(color2.Clearln + fmt.Sprintf("%v", string(msg.Payload)))
		}
	}()

	// ------[about session]-----
	App.AddCommand(&grumble.Command{
		Name: "session",
		Help: "Operation session Settings",
		Flags: func(f *grumble.Flags) {
			f.BoolL("list", false, "list all alive sessions")
			f.IntL("kill", 0, "kill a alive session by id")
		},
		Args: func(a *grumble.Args) {
			a.Int("id", "session id", grumble.Default(0))
		},
		Run: func(c *grumble.Context) error {
			if c.Flags.Bool("list") == true {
				fmt.Println(debugSessionID)
				return session.ListSessions(Salient)
			}
			if c.Args["id"].Value.(int) != 0 {
				defaultId := debugSessionID
				err, id := session.UserSession(Salient, c.Args["id"].Value.(int), c, defaultId)
				debugSessionID = id
				if debugSessionID != 0 {
					session.SessionConsole(App, &debugSessionID, Salient)
				}
				return err
			}
			if c.Flags.Int("kill") != 0 {
				return session.KillSession(Salient, c.Flags.Int("kill"))
			}
			return nil
		},
	})
	// ------[about listener]-----
	App.AddCommand(&grumble.Command{
		Name: "khls",
		Help: "New create kcp-http listener",
		Flags: func(f *grumble.Flags) {
			f.StringL("http", "0.0.0.0:8080", "Set the listening address for data return")
			f.StringL("name", "", "Set the listener alias name")
			f.StringL("kcp", "0.0.0.0:8081", "Set the listening address for command delivery")
			f.BoolL("start", false, "Verify that these addresses create listeners")
		},
		Run: func(c *grumble.Context) error {
			if c.Flags.Bool("start") == true {
				_ = listener.InitKhlsListner(Salient, c)
			}
			return nil
		},
	})
	// list all listen
	App.AddCommand(&grumble.Command{
		Name: "listen",
		Help: "Operation listener Settings",
		Flags: func(f *grumble.Flags) {
			f.BoolL("all", true, "list all already set listener")
			f.IntL("kill", 0, "kill a running listener")
		},
		Run: func(c *grumble.Context) error {
			var resp register.Response
			var req register.Requests
			if c.Flags.Bool("all") == true && c.Flags.Int("kill") == 0 {
				err := Salient.Call(context.Background(), "GetListeners", nil, &resp)
				if err != nil {
					return err
				}
				if resp.Error != nil {
					return resp.Error
				} else {
					if len(resp.ActiveJobMsg) > 0 {
						t := table.NewWriter()
						t.SetOutputMirror(os.Stdout)
						t.AppendHeader(table.Row{"ID", "Listener Type", "Alias Name", "Command delivery address", "Message return address"})
						for _, l := range resp.ActiveJobMsg {
							t.AppendRow([]interface{}{l.ID, l.Description, l.Name, l.SendMsgAddress, l.ReturnMsgAddress})
						}
						t.Render()
					} else {
						fmt.Println("😊 no set any listener")
					}
				}
			}
			if c.Flags.Int("kill") != 0 {
				req.Listener.JobID = c.Flags.Int("kill")
				err := Salient.Call(context.Background(), "KillListener", &req, &resp)
				if err != nil {
					fmt.Println("😰", err)
				}
				if resp.Error != nil {
					fmt.Println("😰", resp.Error)
				} else {
					fmt.Println(fmt.Sprintf("😈 kill listener id %v", c.Flags.Int("kill")))
				}
			}
			return nil
		},
	})
	err := App.Run()
	if err != nil {
		return
	}
}
