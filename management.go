package main

import (
	"strconv"
	"net"
	"strings"
)

func RunCommand(c *net.TCPConn, frame *Frame) {
	switch frame.command {
		case "newuser":
			respond(c, "ok", NewUser())
		case "killuser":
			id, err := strconv.Atoui64(frame.data)
			if (err != nil) {
				respond(c, "error", "Invalid User: " + frame.data + ", " + err.String())
				return
			}
			KillUser(id)
			respond(c, "ok", "user killed")
		case "whisper":
			idx := strings.Index(frame.data, "!")
			if idx < 0 {
				respond(c, "error", "User separator not found (!)")
			}
			
			uidstr := frame.data[0:idx]
			msg := frame.data[idx+1:]
			
			id, err := strconv.Atoui64(uidstr)
			if (err != nil) {
				respond(c, "error", "Invalid User: " + uidstr + ", " + err.String())
				return
			}
			
			SendToUser(id, msg)
			respond(c, "ok", "message sent")
		case "join":
			idx := strings.Index(frame.data, "!")
			if idx < 0 {
				respond(c, "error", "Channel/User separator not found (!)")
			}
		
			cidstr := frame.data[0:idx]
			uidstr := frame.data[idx+1:]
		
			cid, err := strconv.Atoui64(cidstr)
			if (err != nil) {
				respond(c, "error", "Invalid Channel: " + cidstr + ", " + err.String())
				return
			}

			uid, err := strconv.Atoui64(uidstr)
			if (err != nil) {
				respond(c, "error", "Invalid User: " + uidstr + ", " + err.String())
				return
			}
		
			JoinChannel(cid, uid)
			respond(c, "ok", "user joined")
		case "leave":
			idx := strings.Index(frame.data, "!")
			if idx < 0 {
				respond(c, "error", "Channel/User separator not found (!)")
			}
	
			cidstr := frame.data[0:idx]
			uidstr := frame.data[idx+1:]
	
			cid, err := strconv.Atoui64(cidstr)
			if (err != nil) {
				respond(c, "error", "Invalid Channel: " + cidstr + ", " + err.String())
				return
			}

			uid, err := strconv.Atoui64(uidstr)
			if (err != nil) {
				respond(c, "error", "Invalid User: " + uidstr + ", " + err.String())
				return
			}
	
			LeaveChannel(cid, uid)
			respond(c, "ok", "user leaved")
		case "newchannel":
			respond(c, "ok", strconv.Uitoa64(NewChannel()))
		case "killchannel":
			id, err := strconv.Atoui64(frame.data)
			if (err != nil) {
				respond(c, "error", "Invalid Channel: " + frame.data + ", " + err.String())
				return
			}
			KillChannel(id)
			respond(c, "ok", "channel killed")
		case "shout":
			idx := strings.Index(frame.data, "!")
			if idx < 0 {
				respond(c, "error", "Channel separator not found (!)")
			}
		
			cidstr := frame.data[0:idx]
			msg := frame.data[idx+1:]
		
			id, err := strconv.Atoui64(cidstr)
			if (err != nil) {
				respond(c, "error", "Invalid Channel: " + cidstr + ", " + err.String())
				return
			}
		
			SendToChannel(id, msg)
			respond(c, "ok", "message sent")
		default:
			respond(c, "error", "Invalid command: " + frame.command)
	}
}
