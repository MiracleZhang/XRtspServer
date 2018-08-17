package rtsp

import (
	"net"
	"github.com/NiuStar/log/fmt"
	"strconv"
	"github.com/MiracleZhang/XRtspServer/RtspClientManager"
	"github.com/MiracleZhang/XRtspServer/media"
	"github.com/MiracleZhang/XRtspServer/util"
)

const (
	PLAYER_COMMAND_DESCRIBE = 1
	PLAYER_COMMAND_SETUP_1 = 2
	PLAYER_COMMAND_SETUP_2 = 3
	PLAYER_COMMAND_OPTIONS = 4
	PLAYER_COMMAND_PLAY = 5
	PLAYER_COMMAND_DATA = 6
	PLAYER_COMMAND_TEARDOWN = 7
)

type RtspPlayer struct {
	tcpCon net.Conn
	cSeq int
	Signals    chan bool
	Outgoing   chan *RtspClientManager.Response
	url string
	pushurl string
	next int
	session int
}

var G_Player *RtspPlayer

func NewRtspPlayer(address string) (*RtspPlayer, error) {

	player := &RtspPlayer{cSeq: 1, Signals: make(chan bool, 1), Outgoing: make(chan *RtspClientManager.Response, 1)}

	tcpCon, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("ERROR: connect to (", address, ") failed -", err)
		return nil, err
	}

	player.tcpCon = tcpCon
	player.next = PLAYER_COMMAND_DESCRIBE


	//test
	player.url = "rtsp://127.0.0.1:8554/test.sdp"
	player.pushurl = "rtsp://127.0.0.1:8556/test.sdp"

	player.session = 1

	G_Player = player

	return player, nil
}

func (p *RtspPlayer)sendDescribe() {
	cmd := "DESCRIBE"
	header := "application/sdp"
	str := cmd + " " + p.url + " RTSP/1.0\r\n"
	str += "Accept: " + header + "\r\n"
	str += "CSeq: " + strconv.Itoa(p.cSeq) + "\r\n"
	p.tcpCon.Write([]byte(str))
	p.cSeq++

	resp, err := RtspClientManager.ReadResponse(p.tcpCon)
	if err != nil {
		//panic(err)
		fmt.Println("err")

		return
	}
	fmt.Println("no err", resp)

	sdpName := util.GetSdpName(p.url)
	_, exits := media.NewMediaSession(sdpName, resp.Body)
	if exits != nil {
		fmt.Println(exits)
	}

	RtspClientManager.NewClientManager(p.pushurl)

	p.Outgoing <- resp
}

func (p *RtspPlayer)sendSetup1() {
	cmd := "SETUP"
	// todo: tcp/udp 区分
	header := "RTP/AVP/TCP"
	str := cmd + " " + p.url + " RTSP/1.0\r\n"
	//todo: 0 = tracktype * 2， 1 = tracktype * 2 + 1； tracktype = 0
	str += "Transport: " + header + ";unicast;interleaved=" + "0" + "-" + "1" + "\r\n"
	str += "CSeq: " + strconv.Itoa(p.cSeq) + "\r\n"
	p.tcpCon.Write([]byte(str))
	p.cSeq++

	resp, err := RtspClientManager.ReadResponse(p.tcpCon)
	if err != nil {
		//panic(err)
		fmt.Println("err")

		return
	}
	fmt.Println("no err", resp)

	p.Outgoing <- resp
}

func (p *RtspPlayer)sendSetup2() {
	cmd := "SETUP"
	// todo: tcp/udp 区分
	header := "RTP/AVP/TCP"
	str := cmd + " " + p.url + " RTSP/1.0\r\n"
	//todo: 0 = tracktype * 2， 1 = tracktype * 2 + 1； tracktype = 0
	str += "Transport: " + header + ";unicast;interleaved=" + "2" + "-" + "3" + "\r\n"
	str += "CSeq: " + strconv.Itoa(p.cSeq) + "\r\n"
	p.tcpCon.Write([]byte(str))
	p.cSeq++

	resp, err := RtspClientManager.ReadResponse(p.tcpCon)
	if err != nil {
		//panic(err)
		fmt.Println("err")

		return
	}
	fmt.Println("no err", resp)

	p.Outgoing <- resp
}

func (p *RtspPlayer)sendOptions() {
	cmd := "OPTIONS"
	str := cmd + " " + p.url + " RTSP/1.0\r\n"
	str += "Session: " + strconv.Itoa(p.session) + "\r\n"
	str += "CSeq: " + strconv.Itoa(p.cSeq) + "\r\n"
	p.tcpCon.Write([]byte(str))
	p.cSeq++

	resp, err := RtspClientManager.ReadResponse(p.tcpCon)
	if err != nil {
		//panic(err)
		fmt.Println("err")

		return
	}
	fmt.Println("no err", resp)

	p.Outgoing <- resp
}

func (p *RtspPlayer)sendPlay() {
	cmd := "PLAY"
	header := "npt=0.00-"
	str := cmd + " " + p.url + " RTSP/1.0\r\n"
	str += "Range: " + header + "\r\n"
	str += "Session: " + strconv.Itoa(p.session) + "\r\n"
	str += "CSeq: " + strconv.Itoa(p.cSeq) + "\r\n"
	p.tcpCon.Write([]byte(str))
	p.cSeq++

	resp, err := RtspClientManager.ReadResponse(p.tcpCon)
	if err != nil {
		//panic(err)
		fmt.Println("err")

		return
	}
	fmt.Println("no err", resp)

	p.Outgoing <- resp
}

// todo: 退出
func (p *RtspPlayer)sendTeardown() {
	cmd := "TEARDOWN"
	str := cmd + " " + p.url + " RTSP/1.0\r\n"
	str += "Session: " + strconv.Itoa(p.session) + "\r\n"
	str += "CSeq: " + strconv.Itoa(p.cSeq) + "\r\n"
	p.tcpCon.Write([]byte(str))

	p.Signals <- true
}

func (p *RtspPlayer)readData() {
	data, _ := RtspClientManager.ReadSocket(p.tcpCon)

	if data != nil {
		req := &RtspClientManager.Response{
			Header: make(map[string][]string),
		}
		req.Body = string(data)

		// todo: test send to outer

		manager := RtspClientManager.GetCurrManager(p.pushurl)
		if manager != nil {
			manager.Write([]byte(req.Body))
		}

		//fmt.Println(data)
		p.Outgoing <- req
	}
}

func (p *RtspPlayer)Run() {
	fmt.Start()
	defer fmt.Over()
	fmt.Println("------ rtsp player connection: handling ------\n", p.tcpCon.RemoteAddr().String())


	p.sendDescribe()

	for {
		fmt.Start()
		select {
		case <-p.Signals:
			fmt.Println("Exit signals by rtsp")
			//if !conn.pushClient && conn.manager != nil {
			//	conn.manager.RemoveClient(conn.conn)
			//} else {
			//	RtspClientManager.RemoveManager(conn.url)
			//}
			//RtspClientManager.GetCurrManager(conn.url).RemoveClient(conn.conn)
			fmt.Println("------ Session[%s] : closed ------\n", p.tcpCon.RemoteAddr())
			return
		case resp := <-p.Outgoing:

			fmt.Println("------ rtsp player connection: get resp ------ \n%s\n", p.tcpCon.RemoteAddr().String(), resp)

			if resp.Status == "OK" {
				p.next++
			} else {
				//fmt.Println("------ rtsp player error: get resp ------ \n%s\n", p.tcpCon.RemoteAddr().String(), resp)
			}

			switch p.next {
			case PLAYER_COMMAND_DESCRIBE:
				p.sendDescribe()
			case PLAYER_COMMAND_SETUP_1:
				p.sendSetup1()
			case PLAYER_COMMAND_SETUP_2:
				p.sendSetup2()
			case PLAYER_COMMAND_OPTIONS:
				p.sendOptions()
			case PLAYER_COMMAND_PLAY:
				p.sendPlay()
			case PLAYER_COMMAND_DATA:
				p.readData()
			}

			//if len(req.URL) != 0 && len(conn.url) == 0 {
			//	conn.url = req.URL
			//}
			//
			//resp := conn.handleRequestAndReturnResponse(req)
			//if resp != nil {
			//	//time.Sleep(1 * time.Second)
			//	_, err := conn.conn.Write([]byte(resp.String()))
			//	if err != nil && !conn.pushClient && conn.manager != nil {
			//		fmt.Println("有人断开链接了1")
			//		conn.manager.RemoveClient(conn.conn)
			//		conn.conn.Close()
			//		return
			//	} else if err != nil {
			//		conn.conn.Close()
			//		return
			//	}
			//	fmt.Println("------ rtsp client connection: get request ------ \n%s\n", conn.conn.RemoteAddr().String(), req)
			//	fmt.Println("------ Session : set response ------ \n%s\n", conn.conn.RemoteAddr().String(), resp)
			//}
			////处理RTSP请求
			//if req.Method != RtspClientManager.DATA {
			//
			//	if conn.pushClient {
			//		if req.Method != RtspClientManager.PLAY && req.Method != RtspClientManager.RECORD {
			//			fmt.Println("Player ")
			//			client.ReadRequest()
			//		}
			//	} else {
			//		if req.Method != RtspClientManager.RECORD {
			//			fmt.Println("Player ")
			//			client.ReadRequest()
			//		}
			//	}
			//
			//}
		}
	}
}

