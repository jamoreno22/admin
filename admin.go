package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
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

	var dnsIP string
	dnsIP = runDNSIsAvailable(bc, command)
	var conn1 *grpc.ClientConn

	conn1, err1 := grpc.Dial(dnsIP, grpc.WithInsecure())
	if err1 != nil {
		log.Fatalf("did not connect: %s", err)
	}

	dnsc := l3.NewDNSClient(conn1)

	defer conn.Close()

	for {
		fmt.Println("Ingrese comando")
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			command = scanner.Text()
		}
		split := strings.Split(command, " ")
		split2 := strings.Split(split[1], ".")
		log.Println(len(split))
		log.Println(len(split2))
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

		actionRv, _ = dnsc.Action(context.Background(), &comm)

		//newConsistency(comm.Domain, &dnsIP, &comm)

		if len(consList) != 0 {
			for _, s := range consList {
				if s.zfName == comm.Domain {
					localConsistency := s
					spl := strings.Split(dnsIP, ".")
					switch spl[3] {
					case "17":
						if localConsistency.rv.Rv1 >= actionRv.Rv1 {
							log.Println("Existe un error en la consistencia")
						} else {
							s.ip = dnsIP
							s.rv.Rv1 = actionRv.Rv1
							s.rv.Rv2 = actionRv.Rv2
							s.rv.Rv3 = actionRv.Rv3
							s.com.Domain = comm.Domain
							s.com.Ip = comm.Ip
							s.com.Name = comm.Name
							s.com.Option = comm.Option
							s.com.Parameter = comm.Parameter
						}
					case "18":
						if localConsistency.rv.Rv2 >= actionRv.Rv2 {
							log.Println("Existe un error en la consistencia")
						} else {
							s.ip = dnsIP
							s.ip = dnsIP
							s.rv.Rv1 = actionRv.Rv1
							s.rv.Rv2 = actionRv.Rv2
							s.rv.Rv3 = actionRv.Rv3
							s.com.Domain = comm.Domain
							s.com.Ip = comm.Ip
							s.com.Name = comm.Name
							s.com.Option = comm.Option
							s.com.Parameter = comm.Parameter
						}
					case "19":
						if localConsistency.rv.Rv3 >= actionRv.Rv3 {
							log.Println("Existe un error en la consistencia")
						} else {
							s.ip = dnsIP
							s.rv.Rv1 = actionRv.Rv1
							s.rv.Rv2 = actionRv.Rv2
							s.rv.Rv3 = actionRv.Rv3
							s.com.Domain = comm.Domain
							s.com.Ip = comm.Ip
							s.com.Name = comm.Name
							s.com.Option = comm.Option
							s.com.Parameter = comm.Parameter
						}
					}
				}
			}
		} else {
			consList = append(consList, &Consistency{zfName: comm.Domain, rv: l3.VectorClock{Name: comm.Domain, Rv1: 0, Rv2: 0, Rv3: 0},
				ip: dnsIP, com: l3.Command{Action: comm.Action, Name: comm.Name, Domain: comm.Domain, Option: comm.Option, Parameter: comm.Parameter, Ip: comm.Ip}})
		}

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
