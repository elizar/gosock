// CRAP CHAT
// ----------------------------------
// Testing out tcp socket with  GO
// @version 0.0.1
// @author Penzur Desu
package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func main() {
	socks := make(map[float64]net.Conn)

	sock, err := net.Listen("tcp", ":31337")
	failOnError(err, "Can't bind to address and port")

	for {
		conn, err := sock.Accept()
		failOnError(err, "socket eror")

		go func(con net.Conn) {
			// create random identifier
			id := rand.New(rand.NewSource(time.Now().UnixNano())).Float64()
			// add new con to socks
			socks[id] = con

			writeToSocks(
				socks,
				id,
				fmt.Sprintf("client-%s connected\n", strconv.FormatFloat(id, 'f', 6, 64)),
			)

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
						fmt.Sprintf("client [%f] disconnected\n", id),
					) // print out dc notice
					return
				}

				matched, _ = regexp.MatchString(".{2,}", line)
				// if message is not empty
				if matched == true {
					sid := strings.Replace(
						strconv.FormatFloat(id, 'f', 6, 64),
						"0.",
						"",
						-1,
					)

					writeToSocks(socks, id, fmt.Sprintf("client-%s> %s", sid, line))
				}
			}
		}(conn)
	}

}

func writeToSocks(socks map[float64]net.Conn, id float64, msg string) {
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
		fmt.Printf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}
