package main

import (
	"bytes"
	"fmt"
	"github.com/google/uuid"
	"net"
	"net/netip"
	"time"
)

func testUdp(listenAddr, connectAddr netip.AddrPort) {
	lAddr := net.UDPAddrFromAddrPort(listenAddr)
	cAddr := net.UDPAddrFromAddrPort(connectAddr)

	// initiate connection
	fmt.Println("Connecting")
	l, err := net.ListenUDP("udp", lAddr)
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer l.Close()
	c, err := net.DialUDP("udp", nil, cAddr)
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer c.Close()

	// test writing from connect port
	fmt.Println("Test 1->2")
	var b [1024]byte
	var rn int
	var recvAddr *net.UDPAddr

	u := uuid.New()
	_, err = c.Write(u[:])
	if err != nil {
		fmt.Println("Error writing to connect port:", err)
		return
	}
	l.SetDeadline(time.Now().Add(10 * time.Second))
	rn, recvAddr, err = l.ReadFromUDP(b[:])
	if err != nil {
		fmt.Println("Error reading from listen port:", err)
		return
	}
	if !bytes.Equal(b[:rn], u[:]) {
		fmt.Println("Differing response value (expecting 16 length and matching byte values):", rn)
		return
	}

	// test writing from the other side
	fmt.Println("Test 2->1")

	u = uuid.New()
	l.SetDeadline(time.Now().Add(10 * time.Second))
	_, err = l.WriteToUDP(u[:], recvAddr)
	if err != nil {
		fmt.Println("Error writing to listen port:", err)
		return
	}
	l.SetDeadline(time.Now().Add(10 * time.Second))
	rn, err = c.Read(b[:])
	if err != nil {
		fmt.Println("Error reading from connect port:", err)
		return
	}
	if !bytes.Equal(b[:rn], u[:]) {
		fmt.Println("Differing response value (expecting 16 length and matching byte values):", rn)
		return
	}

	// done
	fmt.Println("Test success")
}
