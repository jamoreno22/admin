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

func main() {

	var brokerIp string
	brokerIp = "10.10.28.20:8000"
	var conn *grpc.ClientConn

	conn, err := grpc.Dial(brokerIp, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}

	ac := l3.NewBrokerClient

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
			comm = l3.Command{action: 1, name: split2[0], domain: split2[1], option: "", parameter: ""}
		case "Update":
			comm = l3.Command{action: 2, name: split2[0], domain: split2[1], option: split[2], parameter: split[3]}
		case "Delete":
			comm = l3.Command{action: 3, name: split2[0], domain: split2[1], option: "", parameter: ""}
		default:
			log.Println("Ingrese un comando válido")
			break
		}
		runDNSIsAvailable(ac, command)
	}

}

func runDNSIsAvailable(ac l3.BrokerClient, comm string) error {
	msg := l3.Message{text: comm}
	_, err := l3.DNSIsAvailable(context.Background(), &msg)
	return err
}

func pingDataNode(ip string) bool {
	timeOut := time.Duration(10 * time.Second)
	_, err := net.DialTimeout("tcp", ip, timeOut)
	if err != nil {
		return false
	}
	return true
}
