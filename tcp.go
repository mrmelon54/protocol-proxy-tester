package main

import (
	"bytes"
	"fmt"
	"github.com/google/uuid"
	"net"
	"net/netip"
	"time"
)

func testTcp(listenAddr, connectAddr netip.AddrPort) {
	lAddr := net.TCPAddrFromAddrPort(listenAddr)
	cAddr := net.TCPAddrFromAddrPort(connectAddr)

	// initiate connection
	fmt.Println("Connecting")
	l, err := net.ListenTCP("tcp", lAddr)
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer l.Close()
	c, err := net.DialTCP("tcp", nil, cAddr)
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer c.Close()
	innerConn, err := l.AcceptTCP()
	if err != nil {
		fmt.Println("Error accepting connection:", err)
		return
	}
	defer innerConn.Close()

	// test writing from connect port
	fmt.Println("Test 1->2")
	var b [1024]byte
	var rn int

	u := uuid.New()
	_, err = c.Write(u[:])
	if err != nil {
		fmt.Println("Error writing to connect port:", err)
		return
	}
	innerConn.SetDeadline(time.Now().Add(10 * time.Second))
	rn, err = innerConn.Read(b[:])
	if err != nil {
		fmt.Println("Error reading from connection:", err)
		return
	}
	if !bytes.Equal(b[:rn], u[:]) {
		fmt.Println("Differing response value (expecting 16 length and matching byte values):", rn)
		return
	}

	// test writing from the other side
	fmt.Println("Test 2->1")

	u = uuid.New()
	innerConn.SetDeadline(time.Now().Add(10 * time.Second))
	_, err = innerConn.Write(u[:])
	if err != nil {
		fmt.Println("Error writing through the connection:", err)
		return
	}
	c.SetDeadline(time.Now().Add(10 * time.Second))
	rn, err = c.Read(b[:rn])
	if err != nil {
		fmt.Println("Error reading from connection:", err)
		return
	}
	if !bytes.Equal(b[:rn], u[:]) {
		fmt.Println("Differing response value (expecting 16 length and matching byte values):", rn)
		return
	}

	// done
	fmt.Println("Test success")
}
