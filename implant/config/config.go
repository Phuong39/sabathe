package config

// implant struct

type Implant struct {
	AgentId      string   `json:"agent_id"`     // every implant agent id
	ResponseUrl  string   `json:"response_url"` // every implant submit result path
	HostName     string   `json:"host_name"`    // hostname obtained by the implant
	UserName     string   `json:"user_name"`    // username obtained by the implant
	PID          int      `json:"pid"`          // process pid obtained by the implant
	Platform     string   `json:"platform"`     // the obtained system platform
	Architecture string   `json:"architecture"` // the obtained system architecture
	Protocol     string   `json:"protocol"`     // protocol used to send data back to implants
	IPs          []string `json:"ips"`          // the network information obtained by the implant
}

type ImplantConfig struct {
	MsgReturnAddress   string // Message return address
	CmdDeliveryAddress string // Command delivery address
}

type ImplantCmd struct {
	Cmd  string
	Args []string
	body []byte
}

type Response struct {
	Body []byte
}
