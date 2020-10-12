package main

import (
	"strconv"
	"io"
    	"crypto/tls"
    	"log"
    	"net"
    	"crypto/x509"
   	"crypto/rand"
	"os"
	"os/exec"
)

func handleClient(conn net.Conn){
	var n int
	var err error
	
	log.Printf("server: conn: waiting")

    	buf := make([]byte, 10)
	n, err = conn.Read(buf)	
	if err != nil {
			log.Printf("server: conn: read: %x", err)
  			_, err = io.WriteString(conn,"server read error")
			return
	}
	
	str := string(buf[:n])
	count, _ := strconv.Atoi(str)
	log.Printf("server: conn: read: %s", str)

    	buf2 := make([]byte, 80)
	for i:=0; i<count; i++ {
		n, err = conn.Read(buf2)	
		tmp := string(buf2[:n])
		ret, _ := Pull_node(tmp)
		if ret == true{
			log.Printf("Success Pull %s",tmp)
			_, _ = io.WriteString(conn, "OK")
 			log.Printf("server: conn: write: OK")	
		}else{
			_, err = io.WriteString(conn,"server read error")
        		return
		}
	}
	
}

func Pull_node(img_name string) (bool, error) {
	log.Printf("Will download %s image",img_name)
	out, err := exec.Command("./client-pull.sh",img_name).Output()
    	log.Printf("Running command and waiting for it to finish...")
    	if err != nil {
        	log.Fatal(err)
        	return false, err
    	}
	log.Printf(string(out))
	if len(string(out)) < 100 {
		return false, nil
	}
	return true, nil
}

func main() {
    fpLog, err := os.OpenFile("work-node2.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        panic(err)
    }
    defer fpLog.Close()
    log.SetOutput(fpLog)
    
    cert, err := tls.LoadX509KeyPair("../pull/dev1/dev1.crt", "../pull/dev1/dev1.key")
    if err != nil {
        log.Fatalf("server: loadkeys: %s", err)
    }
    config := tls.Config{Certificates: []tls.Certificate{cert}}
    service := "0.0.0.0:10000"
	config.Rand = rand.Reader

    listener, err := tls.Listen("tcp", service, &config)
    if err != nil {
        log.Fatalf("server: listen: %s", err)
    }
    log.Print("server: listening")
	for {
        conn, err := listener.Accept()
        if err != nil {
            log.Printf("server: accept: %s", err)
            break
        }
        defer conn.Close()
        log.Printf("server: accepted from %s", conn.RemoteAddr())
        tlscon, ok := conn.(*tls.Conn)
        if ok {
			log.Printf("Certifiaction correct!\n")
            state := tlscon.ConnectionState()
            for _, v := range state.PeerCertificates {
                log.Print(x509.MarshalPKIXPublicKey(v.PublicKey))
            }
        }
        go handleClient(conn)
	}
}

