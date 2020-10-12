package main

import (
    "crypto/tls"
    "crypto/x509"
    "io"
    "log"
    "os"
    "strconv"
    "net"
)

func main() {
    cert, err := tls.LoadX509KeyPair("./dev2/dev2.crt", "./dev2/dev2.key")
    if err != nil {
        log.Fatalf("server: loadkeys: %s", err)
    }
    config := tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}
    conn, err := tls.Dial("tcp", "192.168.1.49:5555", &config)
    if err != nil {
        log.Fatalf("client: dial: %s", err)
    }
    defer conn.Close()
    log.Println("client: connected to: ", conn.RemoteAddr())

    state := conn.ConnectionState()
    for _, v := range state.PeerCertificates {
        x509.MarshalPKIXPublicKey(v.PublicKey)
    }
    log.Println("client: handshake: ", state.HandshakeComplete)
    log.Println("client: mutual: ", state.NegotiatedProtocolIsMutual)
    
    buf := make([]byte, 40)

    var flag string = "2"
    _, err = io.WriteString(conn, flag)
    if err != nil {
	    log.Fatalf("client: write: %s", err)
    }
    log.Printf("client: conn: write: %s", flag)

    n, err := conn.Read(buf)
    if err != nil {
            log.Printf("server: conn: read: %x", err)
            _, err = io.WriteString(conn,"server read error")
            os.Exit(1)
    }

    var name string = os.Args[1]
    _, err = io.WriteString(conn,name)
    log.Printf("client: conn: write: %s", name)

    n, err = conn.Read(buf)
    if err != nil {
            log.Printf("server: conn: read: %s", err)
            os.Exit(1)
    }

    if string(buf[:n]) == "OK"{
    	log.Printf("server: conn: read: %s", buf)
    }else{ 
    	log.Printf("server: conn: read: %s", buf)
    	os.Exit(1)
    }

    _, err = io.WriteString(conn, "OK")

    n, err = conn.Read(buf)
    if err != nil {
            log.Printf("server: conn: read: %s", err)
            os.Exit(1)
    }

    str := string(buf[:n])
    size, _ := strconv.Atoi(str)
    log.Printf("client: conn: download %d size",size)

    _, err = io.WriteString(conn, "OK")
    if err != nil {
            _, _ = io.WriteString(conn, "download error")
            os.Exit(1)
    }
    log.Printf("client: conn: write: %s", "OK")
 
    message := make([]byte, size)

    num := size / 1180 + 1
    log.Printf("Ready to receive %d sign data",num)
    message = Read_data(num,conn,size)

    filew(name+"-resign.gob",message)

    _, err = io.WriteString(conn, "Receive Success")
    if err != nil {
            _, _ = io.WriteString(conn, "Receive fail")
            os.Exit(1)
    }
    log.Printf("client: conn: write: %s", "Receive Success")

    log.Print("client: exiting")
}

func Read_data (num int, conn net.Conn, size int) []byte{
    m := make([][]byte, num)
    message := make([]byte,size)
    var n int
    for i:=0;i<num;i++{
        m[i] = make([]byte, 1180)
        n, _ = conn.Read(m[i])
        log.Printf("server: conn: read  %d size",n)
        if i > 0{
            message = append(message, m[i][:n]...)
        }else {
            message = append(message[:0:0], m[i][:n]...)
        }
    }
    return message
}


func filew(path string, data []byte) {
    fd, _ := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(0644))
    defer fd.Close()
    _,_ = fd.Write([]byte(data))
    log.Printf("Signed Data save complete!")
}

