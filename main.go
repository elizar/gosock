// CRAP CHAT
// ----------------------------------
// Testing out tcp socket with  GO
// @version 0.0.1
// @author Penzur Desu
package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func main() {
	socks := make(map[string]net.Conn)

	sock, err := net.Listen("tcp", ":31337")
	failOnError(err, "Can't bind to address and port")
	log.Println("[SERVER] - Now accepting connection on port 31337")

	for {
		conn, err := sock.Accept()
		failOnError(err, "socket eror")

		go func(con net.Conn) {
			// create random identifier
			rnd := rand.New(rand.NewSource(time.Now().UnixNano())).Float64()
			id := strings.Replace(
				strconv.FormatFloat(rnd, 'f', 6, 64),
				"0.",
				"",
				-1,
			)

			// add new con to socks
			socks[id] = con
			connection_msg := fmt.Sprintf("client-%s connected\n", id)
			writeToSocks(
				socks,
				id,
				connection_msg,
			)

			log.Printf("%s", connection_msg)

			// Loop
			for {
				line, err := bufio.NewReader(socks[id]).ReadString('\n')
				line = strings.Trim(line, " ")

				// Enable quitting thru: .dc, .quit, or .disconnect command
				matched, _ := regexp.MatchString(`^\.+(dc|quit|disconnect)`, line)
				if matched == true {
					err = struct {
						error
					}{}
				}

				// If client disconnects
				if err != nil {
					socks[id].Close() // Close the connection
					delete(socks, id) // and remove socket from the  map

					writeToSocks(
						socks,
						id,
						fmt.Sprintf("client-%s disconnected\n", id),
					) // print out dc notice
					return
				}

				matched, _ = regexp.MatchString(".{2,}", line)
				// if message is not empty
				if matched == true {
					writeToSocks(socks, id, fmt.Sprintf("client-%s> %s", id, line))
				}
			}
		}(conn)
	}

}

func writeToSocks(socks map[string]net.Conn, id string, msg string) {
	for k, s := range socks {
		// only write to other sockets
		if k != id {
			hh, mm, ss := time.Now().Clock()
			s.Write([]byte(fmt.Sprintf("%02d:%02d:%02d %s", hh, mm, ss, msg)))
		}
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Printf("%s: %s\n", msg, err)
		panic(fmt.Sprintf("%s: %s\n", msg, err))
	}
}
