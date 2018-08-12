package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"
)

func main() {
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	conn, err := dialer.Dial("tcp", "localhost:18888")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	reader := bufio.NewReader(conn)

	request, _ := http.NewRequest("GET", "http://localhost:18888/chunked", nil)
	err = request.Write(conn)
	if err != nil {
		panic(err)
	}
	resp, err := http.ReadResponse(reader, request)
	if err != nil {
		panic(err)
	}
	if resp.TransferEncoding[0] != "chunked" {
		panic("wrong transfer encoding")
	}

	for {
		sizeStr, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		}

		size, err := strconv.ParseInt(string(sizeStr[:len(sizeStr)-2]), 16, 64)
		if size == 0 {
			break
		}
		if err != nil {
			panic(err)
		}
		line := make([]byte, int(size))
		reader.Read(line)
		reader.Discard(2)
		log.Println(" ", string(line))
	}

	// log.Println("Status:", resp.Status)
	// log.Println("Headers:", resp.Header)

	// counter := 10
	// for {
	// 	data, err := reader.ReadBytes('\n')
	// 	if err == io.EOF {
	// 		break
	// 	}
	// 	fmt.Println("<-", string(bytes.TrimSpace(data)))
	// 	fmt.Fprintf(conn, "%d\n", counter)
	// 	fmt.Println("->", counter)
	// 	counter--
	// }

	// cert, err := ioutil.ReadFile("../zoo/ca.crt")
	// if err != nil {
	// 	panic(err)
	// }
	// certPool := x509.NewCertPool()
	// certPool.AppendCertsFromPEM(cert)
	// tlsConfig := &tls.Config{
	// 	RootCAs: certPool,
	// }
	// tlsConfig.BuildNameToCertificate()

	// client := &http.Client{
	// 	Transport: &http.Transport{
	// 		TLSClientConfig: tlsConfig,
	// 	},
	// }

	// resp, err := http.Get("http://localhost:18888/chunked")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer resp.Body.Close()
	// reader := bufio.NewReader(resp.Body)
	// for {
	// 	line, err := reader.ReadBytes('\n')
	// 	if err == io.EOF {
	// 		break
	// 	}
	// 	log.Println(string(bytes.TrimSpace(line)))
	// }
}
