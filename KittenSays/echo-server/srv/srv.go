package srv

import (
	"github.com/Gorjess/kitten/lib/codec"
	"github.com/Gorjess/kitten/lib/kconn"
	"log"
	"net"
)

// TEchoServer A typical echo echo-server
type TEchoServer struct {
	killSig  chan bool
	protocol string
	addr     string
}

// New create a typical echo echo-server which supports kill signal
func New() (*TEchoServer, error) {
	tes := &TEchoServer{
		killSig:  make(chan bool, 1),
		protocol: "tcp",
		addr:     "0.0.0.0:11008",
	}
	return tes, nil
}

// Run .
func (t *TEchoServer) Run() {
	var (
		er error
		cn net.Conn
		ln net.Listener
		kc *kconn.NetConn
	)

	// listen
	ln, er = net.Listen(t.protocol, t.addr)
	if er != nil {
		goto END
	}
	log.Println("Typical echo echo-server started...on", ln.Addr())

	for {
		select {
		case <-t.killSig:
			goto END
		default:
			cn, er = ln.Accept()
			if er != nil {
				goto END
			}
			log.Println(cn.RemoteAddr(), " connected")

			kc = kconn.New(cn)
			kc.SetMsgCB(func(msg codec.IContent) {
				log.Println("recv:", msg)
				// echo
				kc.Write(msg.String())
			})
			kc.Run()
		}
	}

END:
	log.Println("Typical echo echo-server stopped...", er)
}
