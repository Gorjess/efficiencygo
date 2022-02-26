package kconn

import (
	"bufio"
	"github.com/Gorjess/kitten/lib/codec"
	"github.com/Gorjess/kitten/lib/plainstr"
	"io"
	"log"
	"net"
)

type TMsgCB func(msg codec.IContent)

// NetConn wrapper for net.Conn
// which impl custom read and write
type NetConn struct {
	socketFD net.Conn
	writeCh  chan string
	cc       codec.Codec
	msgCB    TMsgCB
}

// New just create and return a ptr of NetConn
func New(connection net.Conn) *NetConn {
	return &NetConn{
		cc:       codec.New(),
		socketFD: connection,
		writeCh:  make(chan string, 1),
	}
}

func (nc *NetConn) SetMsgCB(cb TMsgCB) {
	nc.msgCB = cb
}

func (nc *NetConn) Run() {
	// write to peer
	go func() {
		er := nc.doWrite2()
		if er != nil {
			log.Println("write failed:", er.Error())
		}
	}()
	// read from peer
	er := nc.Read()
	log.Println("Aborted:", er.Error(), "|", nc.socketFD.RemoteAddr())
}

// Establish create a NetConn ptr and
// starts reading and writing right away
func Establish(connection net.Conn) {
	nc := New(connection)
	nc.Run()
}

// Read reads from NetConn.socketFD,
// it will never return until an error occurs
func (nc *NetConn) Read() error {
	var (
		buf    = make([]byte, 4096)
		reader = bufio.NewReader(nc.socketFD)
		bsz    byte
		sz     int
		er     error

		decoded *plainstr.PlainStr
	)

	for {
		// read header size
		bsz, er = reader.ReadByte()
		if er != nil {
			break
		}
		sz = int(bsz)
		// read body
		_, er = io.ReadFull(reader, buf[:sz])
		if er != nil {
			break
		}
		// decode bytes
		decoded = plainstr.NewEmpty()
		er = nc.cc.Unmarshal(buf[:sz], decoded)
		if er != nil {
			break
		}
		nc.msgCB(decoded)
		// reset buf
		buf = buf[:0]
	}
	return er
}

// Write push s to channel
func (nc *NetConn) Write(s string) {
	nc.writeCh <- s
}

// doWrite consumes strings in NetConn.writeCh
func (nc *NetConn) doWrite() error {
	var (
		er error
		bs []byte
		ps *plainstr.PlainStr
	)
	for s := range nc.writeCh {
		ps = plainstr.New(s)
		bs, er = nc.cc.Marshal(ps)
		if er != nil {
			break
		}
		_, er = nc.socketFD.Write(bs)
		if er != nil {
			break
		} else {
			log.Println("write ", s)
		}
	}
	return er
}

// doWrite2 consumes strings in NetConn.writeCh
func (nc *NetConn) doWrite2() error {
	var (
		er error
		bs []byte
		ps *plainstr.PlainStr
	)
	for {
		select {
		case s := <-nc.writeCh:
			ps = plainstr.New(s)
			bs, er = nc.cc.Marshal(ps)
			if er != nil {
				break
			}
			_, er = nc.socketFD.Write(bs)
			if er != nil {
				break
			} else {
				log.Println("write ", s)
			}
		}
	}
}
