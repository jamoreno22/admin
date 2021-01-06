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

// Consistency struct
type Consistency struct {
	zfName string
	rv     l3.VectorClock
	ip     string
	com    l3.Command
}

var consList []*Consistency

var actionRv *l3.VectorClock

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
			comm = l3.Command{Action: 1, Name: split2[0], Domain: split2[1],
				Option: "", Parameter: "", Ip: "10.10.10.10"}
		case "Update":
			comm = l3.Command{Action: 2, Name: split2[0], Domain: split2[1],
				Option: split[2], Parameter: split[3], Ip: ""}
		case "Delete":
			comm = l3.Command{Action: 3, Name: split2[0], Domain: split2[1],
				Option: "", Parameter: "", Ip: ""}
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

		actionRv, _ = dnsc.Action(context.Background(), &comm)

		var rv l3.VectorClock
		rv = l3.VectorClock{Name: comm.Domain, Rv1: 0, Rv2: 0, Rv3: 0}

		newConsistency(comm.Domain, &rv, &dnsIP, &comm)

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

func newConsistency(zfName string, rv *l3.VectorClock, ip *string, com *l3.Command) {
	var flag = false
	if len(consList) != 0 {
		for _, s := range consList {
			if s.zfName == zfName {
				flag = true
			}
		}
	}
	if flag {
		spl := strings.Split(*ip, ".")
		switch spl[3] {
		case "17":
			if rv.Rv1 >= actionRv.Rv1 {
				log.Println("Existe un error en la consistencia")
			}
		case "18":
			if rv.Rv2 >= actionRv.Rv2 {
				log.Println("Existe un error en la consistencia")
			}
		case "19":
			if rv.Rv3 >= actionRv.Rv3 {
				log.Println("Existe un error en la consistencia")
			}
		}

	} else {
		consList = append(consList, &Consistency{zfName: zfName, rv: *rv, ip: *ip, com: *com})
	}
}
