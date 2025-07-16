# Nyx 🚀

A **simple** and **secure** peer-to-peer chat application built with Go and [Fyne](https://fyne.io/) GUI framework.

---


[![Go Version](https://img.shields.io/badge/go-1.18+-blue.svg)](https://golang.org/dl/)
[![Fyne Version](https://img.shields.io/badge/fyne-v2.x-brightgreen)](https://github.com/fyne-io/fyne)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![GitHub stars](https://img.shields.io/github/stars/x0x7b/Nyx?style=social)](https://github.com/x0x7b/Nyx/stargazers)


## ✨ Features

- 🔗 Peer-to-peer TCP connections — **no central server required**  
- 👤 Nickname exchange for friendly identification  
- 📢 Broadcast messages to **all** connected peers  
- ✉️ Direct messages to **individual** peers  
- 🌙 Dark-themed GUI for comfortable night use  

---

## 🚀 Usage

Run the app with:

```bash
go run main.go <port> <peer_ip:port or '-'> <nickname>
```

| Argument        | Description                     |
|-----------------|--------------------------------|
| `<port>`        | Port to listen on               |
| `<peer_ip:port>`| Peer address to connect or `-` to listen only |
| `<nickname>`    | Your username in chat           |

Example:

```bash
go run main.go 33333 - 0x7b
```

---

## 🛠️ Requirements

- Go **1.18+**  
- Fyne **v2**  

Install Fyne:

```bash
go get fyne.io/fyne/v2
```

---

## 📜 License

[MIT License](LICENSE)

---

<details>
  <summary>🤔 Want to contribute?</summary>

  Feel free to fork, open issues or send pull requests!  
  Make sure to follow the Go and Fyne coding standards.  

</details>

---

> *“Simplicity is the ultimate sophistication.”* — Leonardo da Vinci

---

Happy chatting! 💬
