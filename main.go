package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"

	fyne2 "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"image/color"
)

var (
	peers      = make(map[net.Conn]string)
	peersMu    sync.Mutex
	peersList  []string
	nickname   = ""
	messageLog = widget.NewMultiLineEntry()
	logChannel = make(chan string, 100)
	myApp      = app.New()
	myWindow   = myApp.NewWindow("P2P chat")
	radioGroup *widget.RadioGroup
)

type customTheme struct {
	fyne.Theme
}

func (c *customTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameForeground: // текст
		return color.RGBA{255, 255, 255, 255} // білий
	case theme.ColorNameBackground: // фон Entry
		return color.RGBA{30, 30, 30, 255} // темно-сірий
	case theme.ColorNamePlaceHolder: // колір плейсхолдера
		return color.RGBA{150, 150, 150, 255} // світло-сірий
	default:
		return c.Theme.Color(name, variant)
	}
}

func main() {
	myApp.Settings().SetTheme(&customTheme{Theme: theme.DarkTheme()})

	peersList = []string{"All"}
	radioGroup = widget.NewRadioGroup(peersList, func(selected string) {
		logChannel <- fmt.Sprintf("[!] target set to %s\n", selected)
	})
	radioGroup.Horizontal = false
	radioGroup.Selected = "All"
	messageLog.SetPlaceHolder("Chat will appears here..")

	peersView := widget.NewMultiLineEntry()
	peersView.Disable()

	go func() {
		for msg := range logChannel {
			messageLog.SetText(messageLog.Text + msg)
		}
	}()

	input := widget.NewEntry()
	input.SetPlaceHolder("enter message..")

	sendButton := widget.NewButton("send", func() {
		text := input.Text
		if strings.TrimSpace(text) != "" {
			selected := radioGroup.Selected
			if selected == "All" {
				broadcast(text)
				logChannel <- fmt.Sprintf("[You > All]: %s\n", text)
			} else {
				sendToPeer(selected, text)
				logChannel <- fmt.Sprintf("[You > %s]: %s\n", selected, text)
			}
			input.SetText("")
		}
	})
	rightSide := container.NewVBox(
		widget.NewLabel("send to:"),
		radioGroup,
	)
	inputArea := container.NewBorder(nil, nil, nil, sendButton, input)
	content := container.NewBorder(nil, inputArea, nil, rightSide, messageLog)

	myWindow.Resize(fyne2.NewSize(500, 400))
	myWindow.SetContent(content)

	if len(os.Args) < 4 {
		logChannel <- "Using: go run main.go <port> <peer_ip:port> <nickname>\n"
		os.Exit(1)
	}
	nickname = os.Args[3]
	port := os.Args[1]
	logChannel <- fmt.Sprintf("Starting peer at: %s\n", port)
	go listen(port)

	if len(os.Args) >= 3 {
		peerAddr := os.Args[2]
		if peerAddr != "-" {
			go connectToPeer(peerAddr)
		}

	}

	myWindow.ShowAndRun()

}

func listen(port string) {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		dialog.ShowInformation("Error!", fmt.Sprintf("Failed to start listening: %v", err), myWindow)
		return
	}
	logChannel <- fmt.Sprintf("Listening at %s\n", port)
	ip := getOutboundIP().String()
	myWindow.SetTitle(fmt.Sprintf("Nyx: listening at %s:%s", ip, port))

	for {
		conn, err := ln.Accept()
		if err != nil {
			dialog.ShowInformation("Error!", fmt.Sprintf("Failed to accept: %v", err), myWindow)
			continue
		}

		logChannel <- fmt.Sprintf("New connection: %v\n", conn.RemoteAddr())
		fmt.Fprintf(conn, "%s", "NICKNAME|"+nickname+"\n")
		go handleConn(conn)

	}
}

func connectToPeer(addr string) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		dialog.ShowInformation("Error!", fmt.Sprintf("Failed to connect to %v: %v", addr, err), myWindow)
		return
	}
	fmt.Fprintf(conn, "%s", "NICKNAME|"+nickname+"\n")
	go handleConn(conn)

}

func addPeer(conn net.Conn, peerName string) {
	peersMu.Lock()
	peers[conn] = peerName
	peersMu.Unlock()
	updateRadioUI()
}

func removePeer(conn net.Conn) {
	peersMu.Lock()
	delete(peers, conn)
	peersMu.Unlock()
	conn.Close()
	fmt.Printf("Peer disconected: %v\n", conn.RemoteAddr())
	updateRadioUI()
}

func handleConn(conn net.Conn) {
	var peerName string = conn.RemoteAddr().String()
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		msg := scanner.Text()
		if strings.HasPrefix(msg, "NICKNAME|") {
			peerName = strings.TrimPrefix(msg, "NICKNAME|")
			addPeer(conn, peerName)
			logChannel <- fmt.Sprintf("[ %s connected]\n", peerName)
			continue

		} else {
			logChannel <- fmt.Sprintf("[ %s ]: %s\n", peerName, msg)
		}

	}
	if err := scanner.Err(); err != nil {
		log.Println("error reading from peer:", err)
	}
	removePeer(conn)
}

func broadcast(msg string) {
	peersMu.Lock()
	defer peersMu.Unlock()
	for conn := range peers {
		_, err := fmt.Fprintln(conn, msg)
		if err != nil {
			log.Println("Error sending: ", err)
			conn.Close()
			delete(peers, conn)
		}

	}
}

func updateRadioUI() {
	peersMu.Lock()
	defer peersMu.Unlock()
	names := []string{"All"}
	for _, name := range peers {
		names = append(names, name)
	}
	radioGroup.Options = names
	radioGroup.Refresh()

}

func sendToPeer(peer string, text string) {
	peersMu.Lock()
	defer peersMu.Unlock()
	for conn, name := range peers {
		if name == peer {
			_, err := fmt.Fprintln(conn, text)
			if err != nil {
				log.Println("Error sending: ", err)
				conn.Close()
				delete(peers, conn)
			}
		}
	}
}

func getOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return net.IPv4(127, 0, 0, 1)
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP
}
