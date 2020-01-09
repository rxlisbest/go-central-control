package main

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/beanstalkd/go-beanstalk"
	"net"
	"time"
)

func main() {
	// simple tcp server
	// 1.listen ip+port
	listener, err := net.Listen("tcp", "0.0.0.0:2345")
	if err != nil {
		fmt.Printf("listen fail, err: %v\n", err)
		return
	}

	// 2.accept client request
	// 3.create goroutine for each request
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("accept fail, err: %v\n", err)
			continue
		}

		//create goroutine for each connect
		go process(conn)
	}
}

func process(conn net.Conn) {
	defer conn.Close()
	addr := beego.AppConfig.String("beanstalkdaddr") + ":" + beego.AppConfig.String("beanstalkdport")
	b, err := beanstalk.Dial("tcp", addr)
	tubeSet := beanstalk.NewTubeSet(b, "test")
	if err != nil {
		return
	}
	timeout, err := beego.AppConfig.Int("beanstalkdreservetimeout")

	for {
		id, body, err := tubeSet.Reserve(time.Duration(int(time.Duration(timeout) * time.Second)))
		if err != nil {
			return
		}

		b.Delete(id)
		//var buf [128]byte
		//n, err := conn.Read(buf[:])
		//if err != nil {
		//	fmt.Printf("read from connect failed, err: %v\n", err)
		//	break
		//}
		str := body
		//fmt.Printf("receive from client, data: %v\n", str)
		conn.Write(str)
	}
}
