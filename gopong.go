package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type partida struct {
	muj1   *sync.Mutex
	muj2   *sync.Mutex
	j1     *websocket.Conn
	j2     *websocket.Conn
	pingj1 int
	pingj2 int
}

var totalJugadores int = 0
var totalListos int = 0
var nuevaPartida partida

func main() {

	http.HandleFunc("/connect", func(w http.ResponseWriter, r *http.Request) {
		conn, _ := upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity
		var exit bool = false
		var partidaActual partida
		var nueva bool = false
		var conectado bool = false
		// var numjugador = 0
		var myMutex sync.Mutex
		totalJugadores++
		fmt.Println("Jugador conectado. Total: ", totalJugadores)

		for exit == false {
			// Read message from browser
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				if conectado && !nueva {
					totalListos--
					nuevaPartida.j1 = nil
					nuevaPartida.j2 = nil
					nuevaPartida.muj1 = nil
					nuevaPartida.muj2 = nil
				}
				return
			}

			// Print the message to the console
			// fmt.Printf("%s recieved: %s\n", conn.RemoteAddr(), string(msg))
			spl := strings.Split(string(msg), " ")
			tipoMensaje := spl[0]
			switch string(tipoMensaje) {
			case "connect":
				if conectado == true {
					continue
				}
				conectado = true
				fmt.Printf("%s recieved: %s\n", conn.RemoteAddr(), string(msg))
				if totalListos == 0 {
					nuevaPartida.j1 = conn
					nuevaPartida.muj1 = &myMutex
					// numjugador = 1
				} else {
					nuevaPartida.j2 = conn
					nuevaPartida.muj2 = &myMutex
					partidaActual.j1 = nuevaPartida.j1
					partidaActual.j2 = nuevaPartida.j2
					partidaActual.muj1 = nuevaPartida.muj1
					partidaActual.muj2 = nuevaPartida.muj2
					totalListos = -1
					nueva = true
					// numjugador = 2
				}
				totalListos++
				fmt.Println("Jugador listo, total ", totalListos)
				if nueva {
					// Write message back to browser
					fmt.Println("Iniciando el juego")
					partidaActual.muj2.Lock()
					if err = partidaActual.j2.WriteMessage(msgType, []byte("nuevojuego 2")); err != nil {
						fmt.Printf("ERROR: Al enviare mensaje al jugador: " + err.Error())
						partidaActual.muj2.Unlock()
						return
					}
					partidaActual.muj2.Unlock()

					partidaActual.muj1.Lock()
					if err = partidaActual.j1.WriteMessage(msgType, []byte("nuevojuego 1")); err != nil {
						fmt.Printf("ERROR: Al enviare mensaje al jugador: " + err.Error())
						return
					}
					partidaActual.muj1.Unlock()
				}

				break

			case "sincJugador":
				if nueva == false {
					partidaActual.j1 = nuevaPartida.j1
					partidaActual.j2 = nuevaPartida.j2
					nuevaPartida.j1 = nil
					nuevaPartida.j2 = nil

					partidaActual.muj1 = nuevaPartida.muj1
					partidaActual.muj2 = nuevaPartida.muj2
					nuevaPartida.muj1 = nil
					nuevaPartida.muj2 = nil
					nueva = true
				}
				// if spl[1] == "1" {
				// 	if err = partidaActual.j2.WriteMessage(msgType, []byte("sincJugador "+spl[2]+" "+spl[3]+" "+spl[4]+" "+spl[5]+" "+strconv.Itoa(partidaActual.pingj1))); err != nil {
				// 		fmt.Printf("ERROR: Al enviare mensaje al jugador: " + err.Error())
				// 		return
				// 	}
				// } else {
				// 	if err = partidaActual.j1.WriteMessage(msgType, []byte("sincJugador "+spl[2]+" "+spl[3]+" "+spl[4]+" "+spl[5]+" "+strconv.Itoa(partidaActual.pingj2))); err != nil {
				// 		fmt.Printf("ERROR: Al enviare mensaje al jugador: " + err.Error())
				// 		return
				// 	}
				// }
				if spl[1] == "1" {
					partidaActual.muj2.Lock()
					if err = partidaActual.j2.WriteMessage(msgType, []byte("sincJugador "+spl[2]+" "+spl[3]+" "+"0"+" "+"0"+" "+strconv.Itoa(partidaActual.pingj1))); err != nil {
						partidaActual.muj2.Unlock()
						fmt.Printf("ERROR: Al enviare mensaje al jugador: " + err.Error())
						return
					}
					partidaActual.muj2.Unlock()
				} else {
					partidaActual.muj1.Lock()
					if err = partidaActual.j1.WriteMessage(msgType, []byte("sincJugador "+spl[2]+" "+spl[3]+" "+"0"+" "+"0"+" "+strconv.Itoa(partidaActual.pingj2))); err != nil {
						partidaActual.muj1.Unlock()
						fmt.Printf("ERROR: Al enviare mensaje al jugador: " + err.Error())
						return
					}
					partidaActual.muj1.Unlock()
				}
				break

			case "perdida":
				fmt.Printf("%s recieved: %s\n", conn.RemoteAddr(), string(msg))
				if spl[1] == "1" {
					partidaActual.muj2.Lock()
					if err = partidaActual.j2.WriteMessage(msgType, []byte("perdida "+strconv.Itoa(partidaActual.pingj1))); err != nil {
						partidaActual.muj2.Unlock()
						fmt.Printf("ERROR: Al enviare mensaje al jugador: " + err.Error())
						return
					}
					partidaActual.muj2.Unlock()
				} else {
					partidaActual.muj1.Lock()
					if err = partidaActual.j1.WriteMessage(msgType, []byte("perdida "+strconv.Itoa(partidaActual.pingj2))); err != nil {
						partidaActual.muj1.Unlock()
						fmt.Printf("ERROR: Al enviare mensaje al jugador: " + err.Error())
						return
					}
					partidaActual.muj1.Unlock()
				}

			case "devuelta":
				fmt.Printf("%s recieved: %s\n", conn.RemoteAddr(), string(msg))
				if spl[1] == "1" {
					partidaActual.muj2.Lock()
					if err = partidaActual.j2.WriteMessage(msgType, []byte("devuelta "+spl[2]+" "+spl[3]+" "+spl[4]+" "+spl[5]+" "+strconv.Itoa(partidaActual.pingj1))); err != nil {
						partidaActual.muj2.Unlock()
						fmt.Printf("ERROR: Al enviare mensaje al jugador: " + err.Error())
						return
					}
					partidaActual.muj2.Unlock()
				} else {
					partidaActual.muj1.Lock()
					if err = partidaActual.j1.WriteMessage(msgType, []byte("devuelta "+spl[2]+" "+spl[3]+" "+spl[4]+" "+spl[5]+" "+strconv.Itoa(partidaActual.pingj2))); err != nil {
						partidaActual.muj1.Unlock()
						fmt.Printf("ERROR: Al enviare mensaje al jugador: " + err.Error())
						return
					}
					partidaActual.muj1.Unlock()
				}
				break
			case "endgame":
				fmt.Println("Fin del juego")
				exit = true
				break
			case "ping":
				// if numjugador == 1 {
				// 	partidaActual.muj1.Lock()
				// } else if numjugador == 2 {
				// 	partidaActual.muj2.Lock()
				// }
				myMutex.Lock()
				if err = conn.WriteMessage(msgType, []byte("ping")); err != nil {
					fmt.Printf("ERROR: Al enviare mensaje al jugador: " + err.Error())
					return
				}
				// if numjugador == 1 {
				// 	partidaActual.muj1.Unlock()
				// } else if numjugador == 2 {
				// 	partidaActual.muj2.Unlock()
				// }
				myMutex.Unlock()

				break
			case "setPing":
				if spl[1] == "1" {
					partidaActual.pingj1, _ = strconv.Atoi(spl[2])
				} else if spl[1] == "2" {
					partidaActual.pingj2, _ = strconv.Atoi(spl[2])
				}
				break
			default:
				fmt.Println("Error en el servidor. Mensaje incontrolado: " + string(msg))
			}

		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "redirect.html")
	})

	http.HandleFunc("/juego", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "pong.html")
	})

	http.ListenAndServe(":2020", nil)
}
