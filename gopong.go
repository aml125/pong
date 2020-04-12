package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type partida struct {
	j1 *websocket.Conn
	j2 *websocket.Conn
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
		totalJugadores++
		fmt.Println("Jugador conectado. Total: ", totalJugadores)

		for exit == false {
			// Read message from browser
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
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
				} else {
					nuevaPartida.j2 = conn
					partidaActual.j1 = nuevaPartida.j1
					partidaActual.j2 = nuevaPartida.j2
					totalListos = -1
					nueva = true
				}
				totalListos++
				fmt.Println("Jugador listo, total ", totalListos)
				if nueva {
					// Write message back to browser
					fmt.Println("Iniciando el juego")
					if err = partidaActual.j2.WriteMessage(msgType, []byte("nuevojuego 2")); err != nil {
						fmt.Printf("ERROR: Al enviare mensaje al jugador: " + err.Error())
						return
					}

					if err = partidaActual.j1.WriteMessage(msgType, []byte("nuevojuego 1")); err != nil {
						fmt.Printf("ERROR: Al enviare mensaje al jugador: " + err.Error())
						return
					}
				}

				break

			case "sincJugador":
				if nueva == false {
					partidaActual.j1 = nuevaPartida.j1
					partidaActual.j2 = nuevaPartida.j2
					nuevaPartida.j1 = nil
					nuevaPartida.j2 = nil
					nueva = true
				}
				if spl[1] == "1" {
					if err = partidaActual.j2.WriteMessage(msgType, []byte("sincJugador "+spl[2]+" "+spl[3]+" "+spl[4]+" "+spl[5])); err != nil {
						fmt.Printf("ERROR: Al enviare mensaje al jugador: " + err.Error())
						return
					}
				} else {
					if err = partidaActual.j1.WriteMessage(msgType, []byte("sincJugador "+spl[2]+" "+spl[3]+" "+spl[4]+" "+spl[5])); err != nil {
						fmt.Printf("ERROR: Al enviare mensaje al jugador: " + err.Error())
						return
					}
				}
				break

			case "perdida":
				fmt.Printf("%s recieved: %s\n", conn.RemoteAddr(), string(msg))
				// if spl[1] == "1" {
				// 	if err = partidaActual.j2.WriteMessage(msgType, []byte("perdida "+spl[2]+" "+spl[3]+" "+spl[4]+" "+spl[5])); err != nil {
				// 		fmt.Printf("ERROR: Al enviare mensaje al jugador: " + err.Error())
				// 	}
				// } else {
				// 	if err = partidaActual.j1.WriteMessage(msgType, []byte("perdida "+spl[2]+" "+spl[3]+" "+spl[4]+" "+spl[5])); err != nil {
				// 		fmt.Printf("ERROR: Al enviare mensaje al jugador: " + err.Error())
				// 	}
				// }
				if spl[1] == "1" {
					if err = partidaActual.j2.WriteMessage(msgType, []byte("perdida")); err != nil {
						fmt.Printf("ERROR: Al enviare mensaje al jugador: " + err.Error())
						return
					}
				} else {
					if err = partidaActual.j1.WriteMessage(msgType, []byte("perdida")); err != nil {
						fmt.Printf("ERROR: Al enviare mensaje al jugador: " + err.Error())
						return
					}
				}

			case "devuelta":
				fmt.Printf("%s recieved: %s\n", conn.RemoteAddr(), string(msg))
				if spl[1] == "1" {
					if err = partidaActual.j2.WriteMessage(msgType, []byte("devuelta "+spl[2]+" "+spl[3]+" "+spl[4]+" "+spl[5])); err != nil {
						fmt.Printf("ERROR: Al enviare mensaje al jugador: " + err.Error())
						return
					}
				} else {
					if err = partidaActual.j1.WriteMessage(msgType, []byte("devuelta "+spl[2]+" "+spl[3]+" "+spl[4]+" "+spl[5])); err != nil {
						fmt.Printf("ERROR: Al enviare mensaje al jugador: " + err.Error())
						return
					}
				}
				break
			case "endgame":
				fmt.Println("Fin del juego")
				exit = true
				break
			default:
				fmt.Println("Error en el servidor. Mensaje incontrolado: " + string(msg))
			}

		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "pong.html")
	})

	http.HandleFunc("/juego", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "pong.html")
	})

	http.ListenAndServe(":2020", nil)
}
