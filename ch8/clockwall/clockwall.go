package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"text/tabwriter"
)

type ServerConnection struct {
	ServerName string
	Connection net.Conn
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal(
			"ERROR. Enter at least one server name and address. ",
			"Example: ./clockwall NewYork=localhost:8010 London=localhost:8030",
		)
	}
	var serverConnections []ServerConnection
	for _, argument := range os.Args[1:] {
		splittedArg := strings.Split(argument, "=")
		if len(splittedArg) != 2 {
			log.Fatalf("Invalid argument: %s", argument)
		}
		serverName, serverAddress := splittedArg[0], splittedArg[1]
		if serverName == "" || serverAddress == "" {
			log.Fatalf("Invalid argument: %s", argument)
		}
		connection, err := net.Dial("tcp", serverAddress)
		if err != nil {
			log.Fatal(err)
		}
		defer connection.Close()
		serverConnection := ServerConnection{serverName, connection}
		serverConnections = append(serverConnections, serverConnection)
	}
	tableWriter := tabwriter.NewWriter(
		os.Stdout,
		20,
		0,
		0,
		' ',
		0,
	)
	var tableHeaders string
	for _, serverConnection := range serverConnections {
		tableHeaders += fmt.Sprintf("%s\t", serverConnection.ServerName)
	}
	fmt.Fprintln(tableWriter, tableHeaders)
	tableWriter.Flush()
	for {
		var reportRow string
		for _, serverConnection := range serverConnections {
			output := make([]byte, 9)
			_, err := serverConnection.Connection.Read(output)
			if err != nil {
				log.Fatalf(
					"Connection error with server '%s': %v",
					serverConnection.ServerName,
					err,
				)
			}
			trimmedOutput := strings.TrimSpace(string(output))
			reportRow += fmt.Sprintf("%s\t", trimmedOutput)
		}
		fmt.Fprintln(tableWriter, reportRow)
		tableWriter.Flush()
	}
}
