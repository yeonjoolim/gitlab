package main

import (
		"strconv"
		"io"
		"crypto/tls"
		"log"
		"fmt"
		"net"
		"crypto/x509"
		"crypto/rand"
		"os"
		"os/exec"
		"time"
                "./mysql"
       )

var img_name string

func filew(path string, data []byte) {
	fd, _ := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(0644))
		defer fd.Close()
		_,_ = fd.Write([]byte(data))
		mysql.Insert("Signed Data save complete!")
}

func handleClient(conn net.Conn){
	var n int
		var err error

		out, err := exec.Command("docker","login","registry.gitlab.example.com").Output()
		if err != nil {
			mysql.Insert("Docker login fail. Check your Docker registry status")
				return
		}
		mysql.Insert("Docker login success")
		log.Printf(string(out))
		buf := make([]byte, 10)
		n, err = conn.Read(buf)	
		if err != nil {
			mysql.Insert("server: conn: read:" +err.Error())
			_, err = io.WriteString(conn,"server read error")
			return
		}

		str := string(buf[:n])
		n, _ = strconv.Atoi(str)
		mysql.Insert("server: conn: read: "+str)

		var ret bool

		if n == 1{
			ret, err = Push_server(conn)
				if ret == true{
					mysql.Insert("Complete Push")
						return
				}
		}else if n == 2{
			ret, err = Pull_server(conn)
				if ret == true{
					mysql.Insert("Complete Pull")
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
			mysql.Insert("Client Signiture File read error")
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
		mysql.Insert("server: conn: write: OK")	

		n, err = conn.Read(name)	
		if err != nil {
			mysql.Insert("server: conn: read: "+ err.Error())
				return false, err 
		}
	
	img_name = string(name[:n])
	mysql.Insert("server: conn: read: "+ img_name)

		ret, err := vul_detect()
		if ret == false{
			_, err = io.WriteString(conn, "vulnerability detection alot. Download fail")
			mysql.Insert("vulnerability detection alot. Download fail")
				return ret, err
		}else {
			time.Sleep(5 * time.Second)
				_, err = io.WriteString(conn, "OK")
			mysql.Insert("pass vulnerability detection")
		}

	n, err = conn.Read(buf)	
		if err != nil {
			mysql.Insert("server: conn: read: "+ err.Error())
				return false, err 
		}

 	result := string(buf[:n])
	if result == "OK"{
		mysql.Insert("server: conn: read: " + result)
	}else{
		mysql.Insert("server: conn: read: "+result)
		return false, nil
	}

	message, num := filer(img_name+"-resign.gob")
size := strconv.FormatInt(num,10)

_, err = io.WriteString(conn,size)
	if err != nil {
		mysql.Insert("client: write: "+ err.Error())
	}

_, err = conn.Write(message)
	if err != nil {
		mysql.Insert("client: conn: fail send file: "+ img_name+"-resign.gob")
			return false, err
	}
	mysql.Insert("client: conn: send file: "+ img_name+"-resign.gob")

n, err = conn.Read(buf)	
	if string(buf[:n]) == "Send Success"{
		return true, nil
	}else{
		return false, nil
	}
}

func Push_server(conn net.Conn) (bool, error) {
		defer conn.Close()
		buf := make([]byte, 100)
		name := make([]byte, 100)

		n, err := io.WriteString(conn, "OK")
		mysql.Insert("server: conn: write: OK")	

		time.Sleep(5 * time.Second)

		n, err = conn.Read(name)	
		if err != nil {
			mysql.Insert("server: conn: read: "+ err.Error())
			return false, err 
		}
		img_name = string(name[:n])
		mysql.Insert("server: conn: read: "+ img_name)

		out, err := exec.Command("docker","pull",img_name).Output()
		log.Printf("Running command and waiting for it to finish...")
		log.Printf("Uploading Image to GitLab server...")
		if err != nil {
			log.Fatal(err)
				return false, err
		}
		log.Printf(string(out))


		n, err = io.WriteString(conn, "OK")
		mysql.Insert("server: conn: write: Push Image checking.....")	

		n, err = conn.Read(buf)	
		if err != nil {
			mysql.Insert("server: conn: read: "+ err.Error())
				return false, err 
		}

		str := string(buf[:n])
		size, _ := strconv.Atoi(str)
		mysql.Insert("server: conn: download "+str+" size")

		_, err = io.WriteString(conn, "OK")

		message := make([]byte, size)

		num := size / 1180 + 1
		message = Read_data(num,conn,size)
		mysql.Insert("Save signature file complete")

		filew(img_name+"-sign.gob",message)

		ret, err := push_verify()
		if ret == false{
			_, err = io.WriteString(conn, "verify fail")
			mysql.Insert("verify fail")
				return ret, err
		}else{
			mysql.Insert("verify success")
			time.Sleep(5 * time.Second)
				ret, err = vul_detect()
				if ret == false{
					_, err = io.WriteString(conn, "vulnerability detection alot. Upload fail")
					mysql.Insert("vulnerability detection alot. Upload fail")
						return ret, err
				}else{
					time.Sleep(5 * time.Second)
						mysql.Insert("pass vulnerability detection score")
						ret, err = resign_script()
						if ret == false{
							_, err = io.WriteString(conn, "resign fail")
							mysql.Insert("resign fail")
								return ret, err
						}else{
							time.Sleep(5 * time.Second)
								_, err = io.WriteString(conn,"OK")
								if err != nil {
									mysql.Insert("server: conn: read: "+ err.Error())
										return false, err
								}else{
									mysql.Insert("server: conn: write: OK")	
									mysql.Insert("server: conn: closed")	
								}
						}
				}
		}
	return true, nil
}

func push_verify() (bool, error){
	mysql.Insert("server: verify start")
	mysql.Insert("Running command and waiting for it to finish...")
		out, err := exec.Command("./layer-verify",img_name).Output()
		if err != nil {
			mysql.Insert("Docker layers verify fail")
			return false, err
		}
		log.Printf(string(out))
		return true, nil 
}

func vul_detect() (bool, error){
		mysql.Insert("server: Image Vulnerability detection start")
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
		mysql.Insert("server: resign start")
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
			mysql.Insert("server: loadkeys: " +err.Error())
		}
config := tls.Config{Certificates: []tls.Certificate{cert}}
service := "0.0.0.0:50000"
		 config.Rand = rand.Reader

		 listener, err := tls.Listen("tcp4", service, &config)
		 if err != nil {
			mysql.Insert("server: listen: " +err.Error())
		 }
		 mysql.Insert("server: listening")
		 for {
			 conn, err := listener.Accept()
			 if err != nil {
					 mysql.Insert("server: accept: %s"+ err.Error())
					break
				 }
			 defer conn.Close()
				 s := fmt.Sprintf("server: accepted from %s", conn.RemoteAddr())
				mysql.Insert(s) 
				tlscon, ok := conn.(*tls.Conn)
				 if ok {
					 mysql.Insert("Certifiaction correct!\n")
						 state := tlscon.ConnectionState()
						 for _, v := range state.PeerCertificates {
							 tr, _ := x509.MarshalPKIXPublicKey(v.PublicKey)
							 str := string(tr)
							 mysql.Insert(str)
						 }
				 }
			 go handleClient(conn)
		 }
}

