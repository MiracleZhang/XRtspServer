package main

import (
	"github.com/MiracleZhang/XRtspServer/stream_server"
	"github.com/MiracleZhang/XRtspServer/rtsp"
	"net/http"
	"os"
	"log"
)

var rtspServer *rtsp.RtspServer

func httpServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		p, e:= rtsp.NewRtspPlayer("127.0.0.1:8554")
		if e != nil {
			log.Print("error := ", e)
			os.Exit(1)
		}

		go p.Run()
	})
	http.ListenAndServe(":8555", nil)
}

func main() {
	s := stream_server.NewStreamServer(":8554")

	var err error
	rtspServer, err = s.NewRtspServer()
	if err != nil {
		return
	}

	go httpServer()

	s.Run(rtspServer)
}
