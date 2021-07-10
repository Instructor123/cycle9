package wireGuard

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"strings"
)

func checkErr(e error){
	if e != nil{
		panic(e)
	}
}
// True means there is a new key file, false means there was previously a keyfile and we're currently just reusing it.
func genKey(keyName string)bool{

	retValue := true

	//determine if file already exists
	cwd, err := exec.Command("pwd").Output()
	checkErr(err)
	fileName :=  strings.TrimRight(string(cwd), "\n")+"/"+keyName
	fileInfo, err := os.Stat(fileName)

	if 0 == fileInfo.Size(){
		//create the private key:					wg genkey > <location>
		output, err := exec.Command("wg", "genkey").Output()
		checkErr(err)

		err = ioutil.WriteFile(fileName,output,0744)
		checkErr(err)
	} else {
		fmt.Println("Key already exists, not generating a new key.")
		retValue = false
	}

	return retValue
}

//Creates the interface wg0 using the provided IP and then creates the wireguard client listening on the provided port.
func Initialize(vpnIP, listenPort string)bool{
	retValue := false

	iFace, err := net.InterfaceByName("wg0")

	if nil != iFace {
		fmt.Println("Interface 'wg0' already created, skipping")
	} else {
		//create the link:							ip link add wg0 type wireguard
		cmd := exec.Command("ip", "link", "add", "wg0", "type", "wireguard")
		_, err = cmd.Output()
		checkErr(err)

		//add the ip address:						ip addr add <address>/24 dev wg0
		cmd = exec.Command("ip", "addr", "add", vpnIP+"/24", "dev", "wg0")
		_, err = cmd.Output()
		checkErr(err)

		retValue = genKey("private")

		//set the private key						wg set wg0 private-key <private key location>
		cmd = exec.Command("wg", "set", "wg0", "private-key", "private", "listen-port", listenPort)
		_, err = cmd.Output()
		checkErr(err)
		//set the link to up						ip link set wg0 up
		cmd = exec.Command("ip", "link", "set", "wg0", "up")
		_, err = cmd.Output()
		checkErr(err)
	}

	return retValue
}

//Creates a peer
func ConfigPeer(key, port, remoteAddr, vpn, listenPort string){
	cmd := exec.Command("wg", "set", "wg0", "listen-port", listenPort, "peer", key,
		"allowed-ips", vpn, "endpoint", remoteAddr+":"+port)
	_, err := cmd.Output()
	checkErr(err)
}
