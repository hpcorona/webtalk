package main

import (
	"fmt"
	"http"
	"websocket"
	"time"
	"net"
	"os"
	"strconv"
	"strings"
	"io/ioutil"
	"json"
)

const (
	second int64 = 1000000000
)

func indexHandler(ws *websocket.Conn) {
	var ibuff [50]byte
	
	i, err := ws.Read(ibuff[0:50])
	if err != nil {
		fmt.Printf("WebSocket error: " + err.String())
		return
	}
	
	msg := string(ibuff[0:i])	
	idx := strings.Index(msg, "!")
	if idx < 0 {
		fmt.Printf("Salt not found")
		return
	}
	
	idstr := msg[0:idx]
	salt := msg[idx+1:]
	
	uid, err := strconv.Atoui64(idstr)
	if err != nil {
		fmt.Printf("User ID invalid")
		return
	}
	
	defer UnlinkUser(uid, ws)
	
	usr := LinkUser(uid, salt, ws)
	if usr == nil {
		fmt.Printf("Cannot link with User ID")
		return
	}
	
	for usr.exit == false {
		time.Sleep(1e8)
	}
}

type Config struct {
	Public string
	Private string
}

var config *Config

func loadConfig() {
	content, err := ioutil.ReadFile("webtalk.ini")
	if err != nil {
		fmt.Printf("No configuration webtalk.ini found...\n")
		fmt.Printf("Using Public Address:  0.0.0.0:12345\n")
		fmt.Printf("Using Private Address: 0.0.0.0:12344\n")
		
		config = &Config { Public: "0.0.0.0:12345", Private: "0.0.0.0:12344" }
		return
	}
	
	config = new(Config)
	
	err = json.Unmarshal(content, config)
	if err != nil {
		fmt.Printf("Invalid config file: " + err.String())
		os.Exit(1)
	}
	
	fmt.Printf("Loaded configuration from webtalk.ini...\n")
	fmt.Printf("Using Public Address:  %s\n", config.Public)
	fmt.Printf("Using Private Address: %s\n", config.Private)
}

func main() {
	var err os.Error
	
	loadConfig()
	
	fmt.Printf("Creating global channel 1...\n")
	NewChannel()
	
	fmt.Printf("Server webtalk started...\n")
	defer fmt.Printf("Server webtalk terminated\n")
	
	fmt.Printf("Running management interop...\n")
	_, err = Start(config.Private, management)
	if err != nil {
		panic("[error] Could not initialize management server: " + err.String())
	}

	fmt.Printf("Running webtalk service...\n")
	http.Handle("/webtalk", websocket.Handler(indexHandler))
	err = http.ListenAndServe(config.Public, nil)
	if err != nil {
		panic("[error] While trying to serve: " + err.String())
	}
}

func respond(c *net.TCPConn, command, data string) {
	c.Write([]byte(NewFrameString(command, data)))
}

func management(c *net.TCPConn) {
	defer c.Close()
	
	var buffer [2048]byte
	
	for {
		i, err := c.Read(buffer[0:5])
		if err != nil {
			respond(c, "error", "Invalid Frame: " + err.String())
			break
		}
		if i < 5 {
			respond(c, "error", "Insufficient Frame Size: " + strconv.Itoa(i))
			break
		}
		
		frameSize, err := strconv.Atoi(string(buffer[0:5]))
		if err != nil {
			respond(c, "error", "Invalid Frame Size: " + err.String())
			break
		}
		
		if frameSize == 0 {
			respond(c, "ok", "goodbye")
			break
		}
		
		var toRead int = frameSize
		
		for toRead > 0 {
			i, err = c.Read(buffer[frameSize - toRead:frameSize])
			if err != nil {
				respond(c, "error", "Error reading Frame: " + err.String())
				break
			}
			
			toRead -= i
		}
		
		frame, err := NewFrame(string(buffer[0:frameSize]))
		if err != nil {
			respond(c, "error", "Invalid Frame: " + err.String())
			break
		}
		
		RunCommand(c, frame)
	}
}
