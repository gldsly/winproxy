package main

import (
	"bytes"
	"errors"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

var infoReg = regexp.MustCompile("\\d{0,3}\\.\\d{0,3}\\.\\d{0,3}\\.\\d{0,3}.*?\\n")
var spaceReg = regexp.MustCompile("\\s{2,}")
var ipRegStr = `^(([1-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.)(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){2}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`
func parseProxyInfo(row string) (*proxyItem, error) {
	row = spaceReg.ReplaceAllString(strings.TrimSpace(row), " ")
	rowSlice := strings.Split(row, " ")
	if len(rowSlice) < 4 {
		return nil, errors.New("数据错误")
	}
	return &proxyItem{
		ListenAddr: rowSlice[0],
		ListenPort: rowSlice[1],
		TargetAddr: rowSlice[2],
		TargetPort: rowSlice[3],
	}, nil
}

type proxyItem struct {
	ListenAddr string `json:"l_addr"`
	ListenPort string `json:"l_port"`
	TargetAddr string `json:"t_addr"`
	TargetPort string `json:"t_port"`
}

func checkIP(ip string) bool {
	addr := strings.Trim(ip, " ")

	if addr == "0.0.0.0" {
		return true
	}

	if match, _ := regexp.MatchString(ipRegStr, addr); match {
		return true
	}
	return false
}

func showProxyInfo(result *[]*proxyItem) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command("netsh", "interface", "portproxy", "show", "v4tov4")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("command run error: " + err.Error())
	}

	gbkDecoder := simplifiedchinese.GBK.NewDecoder()
	cmdRes, err := gbkDecoder.String(stdout.String())
	if err != nil {
		fmt.Println("decode command result error: " + err.Error())
		return
	}

	res := infoReg.FindAllString(cmdRes, -1)


	for _, item := range res {

		rowInfo, err := parseProxyInfo(item)
		if err != nil {
			fmt.Println(err)
			continue
		}
		*result = append(*result, rowInfo)
	}
}

func deleteProxy(listenAddr, listenPort string) error {


	// netsh interface portproxy delete v4tov4 listenaddress=0.0.0.0 listenport=10008

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	addr := "listenaddress=" + listenAddr
	port := "listenport=" + listenPort

	if !checkIP(listenAddr) {
		return errors.New("监听地址错误")
	}

	p, err := strconv.Atoi(listenPort)
	if err != nil {
		return err
	}
	if p < 0 || p > 65535 {
		return errors.New("端口超出范围")
	}

	cmd := exec.Command("netsh", "interface", "portproxy", "delete", "v4tov4", addr, port)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	return cmd.Run()
}

func createProxy(listenAddr, listenPort, targetAddr, targetPort string) error {

	// netsh interface portproxy add v4tov4 listenaddress=0.0.0.0 listenport=10008 connectaddress=192.168.233.3 connectport=10008

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	addr := "listenaddress=" + listenAddr
	port := "listenport=" + listenPort

	addr2 := "connectaddress=" + targetAddr
	port2 := "connectport=" + targetPort

	if !checkIP(listenAddr) {
		return errors.New("监听地址错误")
	}
	if !checkIP(targetAddr) {
		return errors.New("转发地址错误")
	}

	p, err := strconv.Atoi(listenPort)
	if err != nil {
		return err
	}
	if p < 0 || p > 65535 {
		return errors.New("监听端口超出范围")
	}

	p, err = strconv.Atoi(targetPort)
	if err != nil {
		return err
	}
	if p < 0 || p > 65535 {
		return errors.New("转发端口超出范围")
	}

	cmd := exec.Command("netsh", "interface", "portproxy", "add", "v4tov4", addr, port, addr2, port2)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	return cmd.Run()
}