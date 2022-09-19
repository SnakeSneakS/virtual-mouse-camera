package cmd

import (
	"log"
	"net"

	"github.com/snakesneaks/virtual-mouse-camera/mouse-controller/controller"
)

const BUFF_SIZE = 2048

func Run(addr string, port int) {

	udpLn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.ParseIP(addr),
		Port: port,
	})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("listening udp on %s:%d\n", addr, port)

	buf := make([]byte, BUFF_SIZE)
	handler := controller.NewHandler()
	for {
		n, addr, err := udpLn.ReadFromUDP(buf)
		if err != nil {
			log.Fatalln(err)
		}

		if addr.IP.Equal(net.ParseIP("127.0.0.1")) {
			//log.Println("n: ", n)
			if err := handler.HandleHandLandmarksBytes(buf[:n]); err != nil {
				log.Println(err)
				log.Printf("[addr:%s] received %s\n", addr, string(buf[:n]))
			}
		} else {
			log.Printf("not handling socket from unknown address: %s", addr)
		}
	}
}
