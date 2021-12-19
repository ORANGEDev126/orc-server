package orc

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"strconv"
)

type Session struct {
	sock       net.Conn
	id         uint64
	playGround *PlayGround
}

const (
	PROTOCOL_LENGTH = 4
	HEADER_LENGTH   = 8
)

var handler map[int]func(*Session, []byte)

func init() {
	handler = make(map[int]func(*Session, []byte))

	handler[int(Protocol_MOVE_JOG_REQ)] = HandleMoveJogReq
}

func NewSession(conn net.Conn) *Session {
	return &Session{
		sock: conn,
		id:   rand.Uint64(),
	}
}

func (session *Session) Start() {
	recv := make([]byte, 8192)
	readLen := 0
	defer func() {
		session.sock.Close()
	}()

	for {
		n, err := session.sock.Read(recv)
		if err != nil {
			session.CloseSession()
			break
		}

		if n == 0 {
			session.CloseSession()
			fmt.Println("read size zero Id:" + strconv.FormatInt(int64(session.id), 10))
			break
		}

		readLen += n

		for readLen > HEADER_LENGTH {
			packetLen := binary.LittleEndian.Uint32(recv)
			if int(packetLen) > readLen {
				break
			}

			protoId := binary.LittleEndian.Uint32(recv[PROTOCOL_LENGTH:])
			handler, ok := handler[int(protoId)]
			if !ok {
				fmt.Println("invalid proto id %d", protoId)
				break
			}

			handler(session, recv)
			copy(recv, recv[packetLen:])
			readLen -= int(packetLen)
		}
	}
}

func (session *Session) GetId() uint64 {
	return session.id
}

func (session *Session) Send(buf []byte) {
	session.sock.Write(buf)
}

func (session *Session) CloseSession() {
	fmt.Println("close session id : %d", session.id)

	playGround := session.playGround

	if playGround != nil {
		UnregisterGlobal(session)
	}

	session.playGround = nil
	session.sock.Close()
}
