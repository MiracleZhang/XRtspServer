package main

import (
	"github.com/MiracleZhang/XRtspServer/stream_server"
	"github.com/MiracleZhang/XRtspServer/rtsp"
)

var rtspServer *rtsp.RtspServer

func main() {
	s := stream_server.NewStreamServer(":8554")

	var err error
	rtspServer, err = s.NewRtspServer()
	if err != nil {
		return
	}

	s.Run(rtspServer)
}
