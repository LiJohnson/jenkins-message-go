package test

import (
	"fmt"
	"net"
	"regexp"
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

func Test_regx(t *testing.T) {
	urlReg, _ := regexp.Compile(`^>\s*http[\w:\?\.\-\=\#%&/]+`)

	var content string = `## 🌧【test2】build FAILURE  
	> **startAt** sd
	> **desc** roken since build 
	> Started by GitLab push by tjy312  
	> **params**  
	> testParam : fsdsf
	> **[build(#164)](http://localhost:8080/jenkins/job/test2/164/console)**  
> http://localhost:8080/jenkins/job/test2/164/console`

	content = urlReg.ReplaceAllString(content, "")
	fmt.Println(content)

}
