package main

import (
	"websocket"
	"time"
	"rand"
	"strconv"
)

type TalkUser struct {
	ws *websocket.Conn
	c chan string
	t chan bool
	salt string
	exit bool
}

type TalkChannel []uint64

var talkUsers = map[uint64] *TalkUser {}
var talkChannels = map[uint64] TalkChannel {}
var currentUserIdx uint64 = 0
var currentChannelIdx uint64 = 0

func NewChannel() uint64 {
	ch := make(TalkChannel, 0, 5000)
	
	// INICIA CODIGO CRITICO
	currentChannelIdx += 1
	
	myIdx := currentChannelIdx
	// TERMINA CODIGO CRITICO
	
	talkChannels[myIdx] = ch
	
	return myIdx
}

func KillChannel(id uint64) {
	ch := talkChannels[id]
	if ch == nil {
		return
	}
	
	talkChannels[id] = nil, false
}

func JoinChannel(id uint64, uid uint64) {
	ch := talkChannels[id]
	if ch == nil {
		return
	}
	
	usr := talkUsers[uid]
	if usr == nil {
		return
	}

/*
	needSort := false
	if len(ch) > 0 {
		if uid < ch[len(ch) - 1] {
			needSort = true
		}
	}
*/
	
	newCh := append(ch, uid)

/*	
	if (needSort) {
		sort.Sort(newCh)
	}
*/

	talkChannels[id] = newCh
}

func LeaveChannel(id uint64, uid uint64) {
	ch := talkChannels[id]
	if ch == nil {
		return
	}
	
	for i := 0; i < len(ch); i++ {
		newCh := append(ch[0 : i], ch[i + 1 :]...)
		
		talkChannels[id] = newCh
		break
	}
}

func SendToChannel(id uint64, message string) {
	ch := talkChannels[id]
	if ch == nil {
		return
	}
	
	for i := 0; i < len(ch); i++ {
		SendToUser(ch[i], message)
	}
}

func SendToUser(id uint64, message string) {
	usr := talkUsers[id]
	if usr == nil {
		return
	}
	
	usr.c <- message
}

func NewUser() string {
	usr := new(TalkUser)
	usr.c = make(chan string)
	usr.t = make(chan bool)
	usr.exit = false
	usr.salt = strconv.Itoa(rand.Int())
	
	go UserSocket(usr)
	go UserTimeout(usr)
	
	// INICIA CODIGO CRITICO
	currentUserIdx += 1
	
	myIdx := currentUserIdx
	// TERMINA CODIGO CRITICO
	
	talkUsers[myIdx] = usr
	
	// Join to the Global Channel
	JoinChannel(1, myIdx)
	
	return strconv.Uitoa64(myIdx) + "!" + usr.salt
}

func KillUser(id uint64) {
	usr := talkUsers[id]
	if usr == nil {
		return
	}
	
	usr.exit = true
	talkUsers[id] = nil, false
}

func LinkUser(id uint64, salt string, ws *websocket.Conn) (*TalkUser) {
	usr := talkUsers[id]
	if usr == nil {
		return nil
	}

	if usr.ws == nil && usr.salt == salt {
		usr.ws = ws
		return usr
	}
	
	return nil
}

func UnlinkUser(id uint64, ws *websocket.Conn) {
	usr := talkUsers[id]
	if usr == nil {
		return
	}

	if usr.ws == ws {
		usr.ws = nil
	}
}

func UserSocket(usr *TalkUser) {
	for usr.exit == false {
		select {
			case <-usr.t:
				continue
			case m := <-usr.c:
				if usr.ws != nil {
					usr.ws.Write([]byte(m))
				}
		}
	}
}

func UserTimeout(usr *TalkUser) {
	for usr.exit == false {
		time.Sleep(1e9)
		usr.t <- false
	}
}
