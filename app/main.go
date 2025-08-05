package main

import (
	"encoding/binary"
	"fmt"
	"net"
)

func main() {
	fmt.Println("Logs from your program will appear here!")

	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:2053")
	if err != nil {
		fmt.Println("Failed to resolve UDP address:", err)
		return
	}

	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("Failed to bind to address:", err)
		return
	}
	defer udpConn.Close()

	buf := make([]byte, 512)

	for {
		size, source, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error receiving data:", err)
			break
		}

		receivedData := string(buf[:size])
		fmt.Printf("Received %d bytes from %s: %s\n", size, source, receivedData)

		// Create an empty response
		response := make([]byte, 12)

		id := binary.BigEndian.Uint16(buf[0:2])
		qdcount := binary.BigEndian.Uint16(buf[4:6]) // QDCOUNT from request
		reqFlags := binary.BigEndian.Uint16(buf[2:4])

		opcode := (reqFlags >> 11) & 0xF
		rd := (reqFlags >> 8) & 1

		var respFlags uint16 = 0
		respFlags |= 1 << 15      // QR = 1 (response)
		respFlags |= opcode << 11 // Copy OPCODE from request
		if rd == 1 {
			respFlags |= 1 << 8 // Copy RD if set
		}
		// RA, AA, TC, Z are all left as 0

		// Set RCODE
		if opcode == 0 {
			respFlags |= 0 // No error
		} else {
			respFlags |= 4 // Not implemented
		}

		binary.BigEndian.PutUint16(response[0:2], id)        //random id
		binary.BigEndian.PutUint16(response[2:4], respFlags) //Flags - telling that its a response
		binary.BigEndian.PutUint16(response[4:6], qdcount)   //QDCOUNT-  zero question count
		binary.BigEndian.PutUint16(response[6:8], 1)         // ANCOUNT -> zero answer count
		binary.BigEndian.PutUint16(response[8:10], 0)        // NSCOUNT = Authority Section Count (boss of the domain)
		binary.BigEndian.PutUint16(response[10:12], 0)       // ARCOUNT = Additional Section Count

		i := 12
		for {
			length := int(buf[i])
			i++
			if length == 0 {
				break
			}
			i += length
		}

		qtype := binary.BigEndian.Uint16(buf[i : i+2])    // QTYPE: e.g., 1 = A record
		qclass := binary.BigEndian.Uint16(buf[i+2 : i+4]) // QCLASS: e.g., 1 = IN (Internet)

		if qtype != 1 || qclass != 1 {
			return
		}
		// Append the original question section to the response
		response = append(response, buf[12:i+4]...) // domain name + qtype + qclass

		// Answer section begins here

		// Name: pointer to offset 12 (0xC00C)
		response = append(response, 0xC0, 0x0C) // compression pointer to question section

		// Type: A (0x0001)
		response = append(response, buf[i:i+2]...) // same QTYPE

		// Class: IN (0x0001)
		response = append(response, buf[i+2:i+4]...) // same QCLASS

		// TTL (Time to Live): 300 seconds (0x0000012c)
		response = append(response, 0x00, 0x00, 0x01, 0x2c)

		// RDLENGTH: 4 bytes (IPv4 address)
		response = append(response, 0x00, 0x04)

		// RDATA: 127.0.0.1 (IPv4 loopback)
		response = append(response, 127, 0, 0, 1)

		_, err = udpConn.WriteToUDP(response, source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}
