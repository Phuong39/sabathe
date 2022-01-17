package session

import (
	"context"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/smallnest/rpcx/client"
	"os"
	"sabathe/client/grumble"
	"sabathe/server/register"
)

// Operation session Settings

func ListSessions(Salient client.XClient) error {
	var resp register.Response
	err := Salient.Call(context.Background(), "GetSessions", nil, &resp)
	if err != nil {
		return err
	}
	if resp.Error != nil {
		return resp.Error
	} else {
		if len(resp.ActiveSessionMsg) > 0 {
			t := table.NewWriter()
			t.SetOutputMirror(os.Stdout)
			t.AppendHeader(table.Row{"ID", "Platform", "AgentId (unique uuid)", "Implant type", "current user"})
			for _, s := range resp.ActiveSessionMsg {
				t.AppendRow([]interface{}{s.SessionID, s.Description, s.Name, s.ImplantType, s.ImplantUser})
			}
			t.Render()
		}
	}
	return nil
}

func KillSession(Salient client.XClient, id int) error {
	var resp register.Response
	var req register.Requests
	req.Session.SessionID = id
	err := Salient.Call(context.Background(), "KillSession", &req, &resp)
	if err != nil {
		return err
	}
	if resp.Error != nil {
		return resp.Error
	}
	return nil
}

func UserSession(Salient client.XClient, id int, c *grumble.Context, defaultId int) (error, int) {
	var resp register.Response
	var req register.Requests
	req.Session.SessionID = id
	err := Salient.Call(context.Background(), "UseSession", &req, &resp)
	if err != nil {
		return err, defaultId
	}
	if resp.Error != nil {
		fmt.Println("😱", resp.Error)
		return nil, defaultId
	}
	c.App.SetPrompt(fmt.Sprintf("[%v] » ", resp.SessionInfo.Name))
	return nil, resp.SessionInfo.SessionID
}
