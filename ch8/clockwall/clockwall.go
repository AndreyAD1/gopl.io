package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

const columnWidth = 20
var indent = strings.Repeat(" ", columnWidth)


func main() {
	if len(os.Args) < 2 {
		log.Fatal(
			"ERROR. Enter at least one server name and address. ",
			"Example: ./clockwall NewYork=localhost:8010 London=localhost:8030",
		)
	}
	servers := make(map[string]string)
	for _, argument := range os.Args[1:] {
		splittedArg := strings.Split(argument, "=")
		if len(splittedArg) != 2 {
			log.Fatalf("Invalid argument: %s", argument)
		}
		serverName, serverAddress := splittedArg[0], splittedArg[1]
		if serverName == "" || serverAddress == "" {
			log.Fatalf("Invalid argument: %s", argument)
		}
		servers[serverName] = serverAddress
	}

	var tableHeaders string
	for serverName := range servers {
		tableHeaders += fmt.Sprintf("%s%s", indent, serverName)
	}
	fmt.Println(tableHeaders)
	i := 0
	for serverName, serverAddress := range servers {
		i++
		serverIsLast := i == len(servers)
		go printTimeFromServer(serverName, serverAddress, i, serverIsLast)
	}
	for {}
}


func printTimeFromServer(serverName, serverAddress string, serverNumber int, serverIsLast bool) {
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	for {
		output := make([]byte, 40)
		_, err := conn.Read(output)
		if err != nil {
			log.Fatalf(
				"Connection error with server %s: %s", 
				serverName, 
				serverAddress,
			)
		}
		recordIndent := strings.Repeat(indent, serverNumber)
		outputString := fmt.Sprintf("%s%s", recordIndent, string(output))
		if serverIsLast {
			outputString += "\n"
		}
		fmt.Print(outputString)
	}
}
