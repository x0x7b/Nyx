# Nyx

A simple peer-to-peer chat application built with Go and Fyne GUI.

## Features

- Peer-to-peer TCP connections without a central server  
- Nickname exchange for friendly identification  
- Broadcast messages to all connected peers  
- Direct messages to individual peers  
- Dark-themed GUI for comfortable use

## Usage

Run the app with:  
```bash
go run main.go <port> <peer_ip:port or '-'> <nickname>
```

- `<port>` — port to listen on  
- `<peer_ip:port>` — peer to connect to, or `-` to listen only  
- `<nickname>` — your username in chat  

Example:  
```bash
go run main.go 33333 - 0x7b
```

## Requirements

- Go 1.18 or higher  
- Fyne v2  

Install Fyne with:  
```bash
go get fyne.io/fyne/v2
```

## License

MIT License
