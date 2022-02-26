package klient

import (
	"fmt"
	"github.com/Gorjess/kitten/lib/codec"
	"github.com/Gorjess/kitten/lib/kconn"
	"log"
	"net"
	"time"
)

type Klient struct {
	conn    *kconn.NetConn
	counter int
}

func New() *Klient {
	k := &Klient{}
	return k
}

func execWithDelay(delay time.Duration, cb func() error) error {
	var (
		er error
		c  = time.NewTicker(delay * time.Second)
	)
	for {
		select {
		case <-c.C:
			er = cb()
			if er != nil {
				goto END
			}
		}
	}
END:
	return er
}

// Communicate send msg to echo-server every n seconds
func (k *Klient) Communicate(n time.Duration) error {
	return execWithDelay(n, func() error {
		//rand.Seed(time.Now().UnixNano())
		// write to echo-server
		//k.conn.Write(fmt.Sprintf("Hello_%d", rand.Int()%100))
		k.conn.Write(fmt.Sprintf("Hello_%d", k.counter))
		k.counter++
		return nil
	})
}

func (k *Klient) Run() error {
	con, er := net.Dial("tcp", "127.0.0.1:11008")
	if er != nil {
		return er
	}
	k.conn = kconn.New(con)

	// to watch if any error occurs while communicating
	quit := make(chan error)

	// write msg to echo-server
	go func() {
		er := k.Communicate(5)
		if er != nil {
			log.Println("comm failed:", er.Error())
		}
	}()

	select {
	case er = <-quit:
		break
	default:
		k.conn.SetMsgCB(func(msg codec.IContent) {
			log.Println("recv:", msg)
		})
		k.conn.Run()
	}

	return er
}
