package mthod

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/xtaci/kcp-go"
	"math/rand"
	"net"
	"strings"
	"time"
)

func RandomString(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := []byte(str)
	result := []byte{}
	rand.Seed(time.Now().UnixNano() + int64(rand.Intn(100)))
	for i := 0; i < length; i++ {
		result = append(result, b[rand.Intn(len(b))])
	}
	return string(result)
}

func Package(message []byte) ([]byte, error) {
	// message length
	var length = int32(len(message))
	var pkg = new(bytes.Buffer)
	// write length to buffer
	err := binary.Write(pkg, binary.BigEndian, length)
	if err != nil {
		return nil, err
	}
	// write message in buffer
	err = binary.Write(pkg, binary.BigEndian, message)
	if err != nil {
		return nil, err
	}
	// return buffer as bytes
	return pkg.Bytes(), nil
	//return message,nil
}

func UnKcpPackage(conn *kcp.UDPSession) []byte {
	reader := bufio.NewReader(conn)
	peek, err := reader.Peek(4)
	if err != nil {
		return nil
	}
	buffer := bytes.NewBuffer(peek)
	//读取数据长度
	var length int32
	err = binary.Read(buffer, binary.BigEndian, &length)
	if err != nil {
		return nil
	}
	data := make([]byte, length+4)
	_, err = reader.Read(data)
	if err != nil {
		return nil
	}
	return data[4:]
}

func IsOpen(address string) bool {
	fmt.Println(address)
	port := strings.Split(address, ":")[1]
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("0.0.0.0:%v", port), 1*time.Second)
	if err != nil {
		// port is not open
		return true
	}
	_ = conn.Close()
	return false
}

func MessageJson(v interface{}) ([]byte, error) {
	message, err := json.Marshal(v)
	if err != nil {
		return []byte{}, err
	}
	packs, err := Package(message)
	if err != nil {
		return []byte{}, err
	}
	return packs, nil
}
