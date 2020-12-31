package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	l3 "github.com/jamoreno22/admin/pkg/proto"
	"google.golang.org/grpc"
)

type consistency struct {
	zfName string
	rv     l3.VectorClock
	ip     string
}

var cons consistency

func main() {

	var brokerIP string
	brokerIP = "10.10.28.20:8000"
	var conn *grpc.ClientConn

	conn, err := grpc.Dial(brokerIP, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}

	bc := l3.NewBrokerClient(conn)

	var command string
	var comm l3.Command

	defer conn.Close()
	for {
		fmt.Println("Ingrese comando")
		fmt.Scanln(&command)
		split := strings.Split(command, " ")
		split2 := strings.Split(split[0], ".")
		switch split[0] {
		case "Create":
			comm = l3.Command{Action: 1, Name: split2[0], Domain: split2[1], Option: "", Parameter: "", Ip: "10.10.10.10"}
		case "Update":
			comm = l3.Command{Action: 2, Name: split2[0], Domain: split2[1], Option: split[2], Parameter: split[3]}
		case "Delete":
			comm = l3.Command{Action: 3, Name: split2[0], Domain: split2[1], Option: "", Parameter: ""}
		default:
			log.Println("Ingrese un comando válido")
			continue
		}

		var dnsIP string
		dnsIP = runDNSIsAvailable(bc, command)
		var conn1 *grpc.ClientConn

		conn1, err1 := grpc.Dial(dnsIP, grpc.WithInsecure())
		if err1 != nil {
			log.Fatalf("did not connect: %s", err)
		}

		dnsc := l3.NewDNSClient(conn1)

		dnsc.Action(context.Background(), &comm)
	}

}

func runDNSIsAvailable(bc l3.BrokerClient, comm string) string {
	msg := l3.Message{Text: comm}
	state, err := bc.DNSIsAvailable(context.Background(), &msg)
	if err != nil {
		fmt.Println("DNSIsAvailable error")
	}
	dnsIps := []string{"10.10.28.17:8000", "10.10.28.18:8000", "10.10.28.19:8000"}
	if state.Dns1 == true {
		return dnsIps[0]
	} else if state.Dns2 == true {
		return dnsIps[1]
	} else if state.Dns3 == true {
		return dnsIps[2]
	}
	log.Fatalln("Dns servers not available")
	return "Dns not available"
}

func pingDataNode(ip string) bool {
	timeOut := time.Duration(10 * time.Second)
	_, err := net.DialTimeout("tcp", ip, timeOut)
	if err != nil {
		return false
	}
	return true
}
