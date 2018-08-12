package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
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

	request, _ := http.NewRequest("GET", "http://localhost:18888/upgrade", nil)
	request.Header.Set("Connection", "Upgrade")
	request.Header.Set("Upgrade", "MyProtocol")
	err = request.Write(conn)
	if err != nil {
		panic(err)
	}
	resp, err := http.ReadResponse(reader, request)
	if err != nil {
		panic(err)
	}
	log.Println("Status:", resp.Status)
	log.Println("Headers:", resp.Header)

	counter := 10
	for {
		data, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		fmt.Println("<-", string(bytes.TrimSpace(data)))
		fmt.Fprintf(conn, "%d\n", counter)
		fmt.Println("->", counter)
		counter--
	}

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

	// resp, err := client.Get("https://localhost:18443")
	// if err != nil {
	// 	panic(err)
	// }
	// defer resp.Body.Close()
	// dump, err := httputil.DumpResponse(resp, true)
	// if err != nil {
	// 	panic(err)
	// }
	// log.Println(string(dump))
}
