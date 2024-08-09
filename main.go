package main

import (
	"flag"
	"fmt"
	"net/netip"
	"strconv"
)

var udpFlag = flag.Bool("u", false, "Enable UDP protocol")
var addrFlag = flag.String("addr", "127.0.0.1", "Listen address")
var listenFlag = flag.String("l", "", "Address to listen on")
var connectFlag = flag.String("c", "", "Address to connect to")

func main() {
	flag.Parse()

	if *listenFlag == "" || *connectFlag == "" {
		flag.Usage()
		return
	}

	addr, err := netip.ParseAddr(*addrFlag)
	if err != nil {
		fmt.Println("Invalid address:", *addrFlag, err)
		return
	}
	listenPort, err := strconv.ParseUint(*listenFlag, 10, 16)
	if err != nil {
		fmt.Println("Invalid port:", *listenFlag, err)
		return
	}
	connectPort, err := strconv.ParseUint(*connectFlag, 10, 16)
	if err != nil {
		fmt.Println("Invalid port:", *connectFlag, err)
		return
	}

	listenAddr := netip.AddrPortFrom(addr, uint16(listenPort))
	connectAddr := netip.AddrPortFrom(addr, uint16(connectPort))

	if *udpFlag {
		testUdp(listenAddr, connectAddr)
	} else {
		testTcp(listenAddr, connectAddr)
	}
}
