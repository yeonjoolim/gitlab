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
		"time"
       )

var img_name string

func filew(path string, data []byte) {
	fd, _ := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(0644))
		defer fd.Close()
		_,_ = fd.Write([]byte(data))
		log.Printf("Signed Data save complete!")
}

func handleClient(conn net.Conn){
	var n int
		var err error

		out, err := exec.Command("docker","login","registry.gitlab.example.com").Output()
		if err != nil {
			log.Fatal(err)
				return
		}
	log.Printf(string(out))

		log.Printf("server: conn: waiting")

		buf := make([]byte, 10)
		n, err = conn.Read(buf)	
		if err != nil {
			log.Printf("server: conn: read: %x", err)
				_, err = io.WriteString(conn,"server read error")
				return
		}
	log.Printf("server: conn: read: %s", buf)

		str := string(buf[:n])
		n, _ = strconv.Atoi(str)

		var ret bool

		time.Sleep(3 * time.Second)
		if n == 1{
			ret, err = Push_server(conn)
				if ret == true{
					log.Printf("Success Push")
						return
				}
		}else if n == 2{
			ret, err = Pull_server(conn)
				if ret == true{
					log.Printf("Success Pull")
						return
				}
		}else{
			_, err = io.WriteString(conn,"server read error")
				return
		}

}

func filer(path string) ([]byte, int64) {
	fd, err := os.Open(path)
		if err != nil{
			log.Printf("File read error")
				os.Exit(1)
		}
	fi, _  := fd.Stat()

		defer fd.Close()
		var num = fi.Size()
		var data = make([]byte, num)
		_,_ = fd.Read(data)
		return data, num
}

func Pull_server(conn net.Conn) (bool, error) {
	defer conn.Close()
		buf := make([]byte, 100)
		name := make([]byte, 100)

		n, err := io.WriteString(conn, "OK")
		log.Printf("server: conn: write: OK")	

		n, err = conn.Read(name)	
		if err != nil {
			log.Printf("server: conn: read: %x", err)
				return false, err 
		}
	log.Printf("server: conn: read: %s", name)
		img_name = string(name[:n])

		ret, err := vul_detect()
		if ret == false{
			_, err = io.WriteString(conn, "vulnerability detection alot. Download fail")
				return ret, err
		}else {
			time.Sleep(5 * time.Second)
				_, err = io.WriteString(conn, "OK")
		}

	n, err = conn.Read(buf)	
		if err != nil {
			log.Printf("server: conn: read: %x", err)
				return false, err 
		}

	if string(buf[:n]) == "OK"{
		log.Printf("server: conn: read: %s", buf)
}else{
	log.Printf("server: conn: read: %s", buf)
		return false, nil
}

	message, num := filer(img_name+"-resign.gob")
size := strconv.FormatInt(num,10)

_, err = io.WriteString(conn,size)
	if err != nil {
		log.Fatalf("client: write: %s", err)
	}
log.Printf("client: conn: write: %s", size)

_, err = conn.Write(message)
	if err != nil {
		log.Printf("client: conn: fail send file: %s", img_name+"-resign.gob")
			return false, err
	}
log.Printf("client: conn: send file: %s", img_name+"-resign.gob")

n, err = conn.Read(buf)	
	if string(buf[:n]) == "Send Success"{
	log.Printf("server: conn: read: %s", buf)
	return true, nil
	}else{
		log.Printf("server: conn: read: %s", buf)
			return false, nil
	}
}

func Push_server(conn net.Conn) (bool, error) {
	defer conn.Close()
		buf := make([]byte, 100)
		name := make([]byte, 100)

		n, err := io.WriteString(conn, "OK")
		log.Printf("server: conn: write: OK")	

		time.Sleep(5 * time.Second)

		n, err = conn.Read(name)	
		if err != nil {
			log.Printf("server: conn: read: %x", err)
				return false, err 
		}
	log.Printf("server: conn: read: %s", name)
		img_name = string(name[:n])

		out, err := exec.Command("docker","pull",img_name).Output()
		log.Printf("Running command and waiting for it to finish...")
		if err != nil {
			log.Fatal(err)
				return false, err
		}
	log.Printf(string(out))

		time.Sleep(5 * time.Second)

		n, err = io.WriteString(conn, "OK")
		log.Printf("server: conn: write: OK")	

		n, err = conn.Read(buf)	
		if err != nil {
			log.Printf("server: conn: read: %x", err)
				return false, err 
		}
	log.Printf("server: conn: read: %s", buf)

		str := string(buf[:n])
		size, _ := strconv.Atoi(str)
		log.Printf("server: conn: download %d size",size)

		time.Sleep(5 * time.Second)
		_, err = io.WriteString(conn, "OK")
		log.Printf("server: conn: write: OK")	

		message := make([]byte, size)

		num := size / 1180 + 1
		log.Printf("Ready to receive %d sign data",num)
		message = Read_data(num,conn,size)

		filew(img_name+"-sign.gob",message)

		ret, err := push_verify()
		if ret == false{
			_, err = io.WriteString(conn, "verify fail")
				return ret, err
		}else{
			time.Sleep(5 * time.Second)
				ret, err = vul_detect()
				if ret == false{
					_, err = io.WriteString(conn, "vulnerability detection alot. Upload fail")
						return ret, err
				}else{
					time.Sleep(5 * time.Second)
						ret, err = resign_script()
						if ret == false{
							_, err = io.WriteString(conn, "resign fail")
								return ret, err
						}else{
							time.Sleep(5 * time.Second)
								_, err = io.WriteString(conn,"OK")
								if err != nil {
									log.Printf("server: conn: read: %x", err)
										return false, err
								}else{
									log.Printf("server: conn: write: OK")
										log.Println("server: conn: closed")
								}
						}
				}
		}
	return true, nil
}

func push_verify() (bool, error){
	log.Println("server: verify start")
		log.Printf("Running command and waiting for it to finish...")
		out, err := exec.Command("./layer-verify",img_name).Output()
		if err != nil {
			_, _ = exec.Command("./delete.sh",img_name).Output()
				return false, err
		}
	log.Printf(string(out))
		return true, nil 
}

func vul_detect() (bool, error){
	log.Println("server: Image Vulnerability detection start")
		log.Printf("Running command and waiting for it to finish...")
		out, err := exec.Command("./vul_detect.sh",img_name).Output()
		if err != nil {
			return false, err
		}
	log.Printf(string(out))

		if len(string(out)) > 130 {
			return false, nil
		}
	return true, nil
}

func resign_script() (bool, error){
	log.Println("server: resign start")
		log.Printf("Running command and waiting for it to finish...")
		out, err := exec.Command("./layer-resign",img_name).Output()
		if err != nil {
			_, _ = exec.Command("./delete.sh",img_name).Output()
				return false, err
		}
	log.Printf(string(out))
		return true, nil 
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

func main() {
	fpLog, err := os.OpenFile("access.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
	defer fpLog.Close()
		log.SetOutput(fpLog)

		cert, err := tls.LoadX509KeyPair("./repo/repo.crt", "./repo/repo.key")
		if err != nil {
			log.Fatalf("server: loadkeys: %s", err)
		}
config := tls.Config{Certificates: []tls.Certificate{cert}}
service := "0.0.0.0:50000"
		 config.Rand = rand.Reader

		 listener, err := tls.Listen("tcp4", service, &config)
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
						 log.Printf("hihi\n")
						 state := tlscon.ConnectionState()
						 for _, v := range state.PeerCertificates {
							 log.Print(x509.MarshalPKIXPublicKey(v.PublicKey))
						 }
				 }
			 go handleClient(conn)
		 }
}

