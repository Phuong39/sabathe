package goroutes

import (
	"github.com/xtaci/kcp-go"
	"net/http"
	"sync"
)

var jobID = 0

// 定义作业所需要结构体

var Jobs = &jobs{
	active: map[int]*Job{},
	mutex:  &sync.RWMutex{},
}

type Job struct {
	ID               int     `json:"id,omitempty"`            // 作业id
	Name             string  `json:"name,omitempty"`          // 作业名称
	Description      string  `json:"description,omitempty"`   // 作业描述
	PersistentID     string  `json:"persistent_id,omitempty"` // 当前作业id
	ReturnMsgAddress string  `json:"return_msg_address,omitempty"`
	SendMsgAddress   string  `json:"send_msg_address,omitempty"`
	Servers          *Server `json:"servers,omitempty"`
}

// 用于打印作业信息

type JobInfo struct {
	ID               int    `json:"id,omitempty"`                 // 作业id
	Name             string `json:"name,omitempty"`               // 作业名称
	Description      string `json:"description,omitempty"`        // 作业描述
	PersistentID     string `json:"persistent_id,omitempty"`      // 当前作业id
	ReturnMsgAddress string `json:"return_msg_address,omitempty"` // 数据返回消息地址
	SendMsgAddress   string `json:"send_msg_address,omitempty"`   // 命令接收消息地址
}

type Server struct {
	KcpServer *kcp.Listener
	//HttpServer *http.Server
	HttpServer *http.Server
}

type jobs struct {
	active map[int]*Job  // 活跃的会话
	mutex  *sync.RWMutex // 读写锁
}

// 获取所有的活跃作业

func (j *jobs) All() []*Job {
	j.mutex.RLock()         // 加读锁
	defer j.mutex.RUnlock() // 解读锁
	var all []*Job
	// 遍历所有的活跃作业，并将作业添加到all切片当中
	for _, job := range j.active {
		all = append(all, job)
	}
	return all
}

// 添加一个活跃作业

func (j *jobs) Add(job *Job) {
	j.mutex.Lock()         // 加写锁
	defer j.mutex.Unlock() // 解写锁
	j.active[job.ID] = job // 添加一个作业到active中
	EventBroker.Publish(Event{
		Job:       job,
		EventType: "start-job",
	})
}

func (j *jobs) Remove(job *Job) {
	j.mutex.Lock()         // 加写锁
	defer j.mutex.Unlock() // 解写锁
	delete(j.active, job.ID)
	EventBroker.Publish(Event{
		Job:       job,
		EventType: "stop-job",
	})
}

func (j *jobs) Get(jobID int) *Job {
	if jobID <= 0 {
		return nil
	}
	j.mutex.RLock()
	defer j.mutex.RUnlock()
	return j.active[jobID]
}

func NextJobID() int {
	newID := jobID + 1
	jobID++
	return newID
}
