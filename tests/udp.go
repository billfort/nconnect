package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/txthinking/socks5"
)

func StartUdpServer() error {
	a, err := net.ResolveUDPAddr("udp", udpPort)
	if err != nil {
		return err
	}
	udpServer, err := net.ListenUDP("udp", a)
	if err != nil {
		return err
	}

	defer udpServer.Close()
	fmt.Printf("UDP server is listening at %v\n", udpPort)

	b := make([]byte, 1024)
	for {
		n, addr, err := udpServer.ReadFromUDP(b)
		if err != nil {
			fmt.Printf("StartUdpServer.ReadFromUDP err: %v\n", err)
			return err
		}

		time.Sleep(100 * time.Millisecond)
		_, _, err = udpServer.WriteMsgUDP(b[:n], nil, addr)
		if err != nil {
			fmt.Printf("StartUdpServer.WriteMsgUDP err: %v\n", err)
			return err
		}
	}
}

func StartUDPClient(serverAddr string) error {
	proxyAddr := fmt.Sprintf("127.0.0.1:%v", port)
	s5c, err := socks5.NewClient(proxyAddr, "", "", 0, 60)
	if err != nil {
		return err
	}
	uc, err := s5c.Dial("udp", serverAddr)
	if err != nil {
		fmt.Println("StartUDPClient.s5c.Dial err: ", err)
		return err
	}
	defer uc.Close()

	user := &Person{Name: "udp_boy", Age: 0}
	for i := 0; i < numMsgs; i++ {
		user.Age++
		send, _ := json.Marshal(user)
		if _, err := uc.Write(send); err != nil {
			fmt.Println("StartUDPClient.Write err ", err)
			return err
		}

		recv := make([]byte, 512)
		n, err := uc.Read(recv)
		if err != nil {
			fmt.Println("StartUDPClient.Read err ", err)
			return err
		}
		if !bytes.Equal(recv[:n], send) {
			return fmt.Errorf("StartUDPClient.recv %v is not as same as sent %v", string(recv[:n]), string(send))
		} else {
			fmt.Printf("StartUDPClient got echo: %v\n", string(recv[:n]))
		}
	}

	return nil
}
