package test

import (
	"fmt"
	"net"
	"testing"
)

func TestAAA(t *testing.T) {
	fmt.Println("testing")
}

func Test_array(t *testing.T) {
	var list [123]int
	fmt.Println(len(list), cap(list))
}

func Test_mac(t *testing.T) {
	macAddressList := getMacAddrs()
	for i, mac := range macAddressList {
		fmt.Println(i, mac)

	}
}

func getMacAddrs() (macAddrs []string) {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		fmt.Printf("fail to get net interfaces: %v", err)
		return macAddrs
	}
	// fmt.Println(netInterfaces)
	for _, netInterface := range netInterfaces {
		macAddr := netInterface.HardwareAddr.String()

		if len(macAddr) == 0 {
			continue
		}
		macAddrs = append(macAddrs, macAddr)
	}
	return
}
