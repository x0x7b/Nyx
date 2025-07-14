package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	fyne2 "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var (
	peers     = make(map[net.Conn]string)
	peersMu   sync.Mutex
	peersList []string
	nickname  = ""
	chatView  *tview.TextView
	peersView *tview.TextView
)

func main() {

	app := tview.NewApplication()
	chatView = tview.NewTextView()
	chatView.SetBorder(true)
	chatView.SetDynamicColors(true)
	chatView.SetTitle("Chat")
	chatView.SetChangedFunc(func() {
		app.Draw()
	})

	input := tview.NewInputField()
	input.SetLabel("> ")
	input.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			text := input.GetText()
			if strings.TrimSpace(text) != "" {
				broadcast(text)
				fmt.Fprintf(chatView, "[ You ]: %s\n", text)
			}
			input.SetText("")
		}
	})

	peersView = tview.NewTextView()
	peersView.SetBorder(true)
	peersView.SetDynamicColors(true)
	peersView.SetTitle("Peers")
	peersView.SetChangedFunc(func() {
		app.Draw()
	})
	part1 := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(chatView, 0, 5, false).
		AddItem(peersView, 0, 1, true)

	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(part1, 0, 1, false).
		AddItem(input, 3, 0, true)

	if len(os.Args) < 4 {
		chatView.Write([]byte("Using: go run main.go <port> <peer_ip:port> <nickname>"))
		os.Exit(1)
	}
	nickname = os.Args[3]
	port := os.Args[1]
	chatView.Write([]byte(fmt.Sprintf("Starting peer at: %s\n", port)))
	go listen(port)

	if len(os.Args) >= 3 {
		peerAddr := os.Args[2]
		if peerAddr != "-" {
			go connectToPeer(peerAddr)
		}

	}
	if err := app.SetRoot(layout, true).Run(); err != nil {
		panic(err)
	}
	go ui()
}

func listen(port string) {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal("Failed to start listening: ", err)
	}
	chatView.Write([]byte(fmt.Sprintf("Listening at %s\n", port)))

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal("Failed to accept: ", err)
			continue
		}

		chatView.Write([]byte(fmt.Sprintf("New connection: %v\n", conn.RemoteAddr())))
		fmt.Fprintf(conn, "%s", "NICKNAME|"+nickname+"\n")
		go handleConn(conn)

	}
}

func connectToPeer(addr string) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal("Failed to connect", addr, err)
	}
	fmt.Fprintf(conn, "%s", "NICKNAME|"+nickname+"\n")
	go handleConn(conn)

}

func addPeer(conn net.Conn, peerName string) {
	peersMu.Lock()
	peers[conn] = peerName
	peersList = append(peersList, peerName)
	updatePeerList()
	peersMu.Unlock()
}

func removePeer(conn net.Conn) {
	peersMu.Lock()
	delete(peers, conn)
	peersMu.Unlock()
	conn.Close()
	fmt.Printf("Peer disconected: %v\n", conn.RemoteAddr())
}

func handleConn(conn net.Conn) {
	var peerName string = conn.RemoteAddr().String()
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		msg := scanner.Text()
		if strings.HasPrefix(msg, "NICKNAME|") {
			peerName = strings.TrimPrefix(msg, "NICKNAME|")
			addPeer(conn, peerName)
			chatView.Write([]byte(fmt.Sprintf("[ %s connected]\n", peerName)))
			continue

		} else {
			chatView.Write([]byte(fmt.Sprintf("[ %s ]: %s\n", peerName, msg)))
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

func updatePeerList() {
	peersView.SetText(strings.Join(peersList, "\n"))
}

func ui() {
	myApp := app.New()
	myWindow := myApp.NewWindow("P2P chat")

	messageLog := widget.NewMultiLineEntry()
	messageLog.SetPlaceHolder("Chat will appears here..")
	messageLog.Disable()

	input := widget.NewEntry()
	input.SetPlaceHolder("enter message..")

	sendButton := widget.NewButton("send", func() {
		text := input.Text
		if strings.TrimSpace(text) != "" {
			broadcast(text)
			messageLog.SetText(messageLog.Text + "\n[You]: " + text)
			input.SetText("")
		}
	})

	inputArea := container.NewBorder(nil, nil, nil, sendButton, input)
	content := container.NewBorder(nil, inputArea, nil, nil, messageLog)

	myWindow.Resize(fyne2.NewSize(500, 400))
	myWindow.SetContent(content)
	myWindow.ShowAndRun()
}
