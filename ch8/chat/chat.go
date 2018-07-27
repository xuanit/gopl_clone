// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 254.
//!+

// Chat is a server that lets clients chat with each other.
package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net"
	"time"
)

//!+broadcaster
type client struct {
	ch   chan<- string // an outgoing message channel
	name string
}

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string) // all incoming client messages
)

func broadcaster() {
	clients := make(map[client]bool) // all connected clients
	for {
		select {
		case msg := <-messages:
			// Broadcast incoming message to all
			// clients' outgoing message channels.
			for cli := range clients {
				cli.ch <- msg
			}

		case cli := <-entering:
			clients[cli] = true
			var buf bytes.Buffer
			fmt.Fprintf(&buf, "There are")
			for existingCli := range clients {
				fmt.Fprintf(&buf, " %s,", existingCli.name)
			}
			cli.ch <- buf.String()

		case cli := <-leaving:
			delete(clients, cli)
			close(cli.ch)
		}
	}
}

//!-broadcaster

type bye bool

//!+handleConn
func handleConn(conn net.Conn) {
	ch := make(chan string) // outgoing client messages
	go clientWriter(conn, ch)

	who := conn.RemoteAddr().String()
	ch <- "You are " + who
	messages <- who + " has arrived"
	entering <- client{name: who, ch: ch}
	active := make(chan bye)
	go func() {
		input := bufio.NewScanner(conn)
		for input.Scan() {
			messages <- who + ": " + input.Text()
			active <- bye(false)
		}
		active <- bye(true)
	}()

	var IdleTime = 30 * time.Second
	ticker := time.NewTicker(IdleTime)
	var lastActive = time.Now()
loop:
	for {
		select {
		case bye := <-active:
			if bye {
				break loop
			} else {
				lastActive = time.Now()
			}
		case <-ticker.C:
			if time.Since(lastActive) > IdleTime {
				break loop
			}
		}
	}
	// NOTE: ignoring potential errors from input.Err()

	leaving <- client{ch: ch, name: who}
	messages <- who + " has left"
	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg) // NOTE: ignoring network errors
	}
}

//!-handleConn

//!+main
func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

//!-main
