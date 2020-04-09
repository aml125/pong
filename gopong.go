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
var totalJugadores int = 0
var totalListos int = 0
var exit bool = false
var conexionj1 *websocket.Conn
var conexionj2 *websocket.Conn

func main() {

	http.HandleFunc("/connect", func(w http.ResponseWriter, r *http.Request) {
		conn, _ := upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity
		if totalJugadores <= 0 {
			conexionj1 = conn
			exit = false
		} else if totalJugadores == 1 {
			conexionj2 = conn
		}
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
				fmt.Printf("%s recieved: %s\n", conn.RemoteAddr(), string(msg))
				totalListos++
				fmt.Println("Jugador listo, total ", totalListos)
				if totalListos >= 2 {
					// Write message back to browser
					fmt.Println("Iniciando el juego")
					if err = conexionj2.WriteMessage(msgType, []byte("nuevojuego 2")); err != nil {
						fmt.Printf("ERROR: Al enviare mensaje al jugador: " + err.Error())
					}

					if err = conexionj1.WriteMessage(msgType, []byte("nuevojuego 1")); err != nil {
						fmt.Printf("ERROR: Al enviare mensaje al jugador: " + err.Error())
					}
				}

				break

			case "sincJugador":
				if spl[1] == "1" {
					if err = conexionj2.WriteMessage(msgType, []byte("sincJugador "+spl[2]+" "+spl[3]+" "+spl[4]+" "+spl[5])); err != nil {
						fmt.Printf("ERROR: Al enviare mensaje al jugador: " + err.Error())
					}
				} else {
					if err = conexionj1.WriteMessage(msgType, []byte("sincJugador "+spl[2]+" "+spl[3]+" "+spl[4]+" "+spl[5])); err != nil {
						fmt.Printf("ERROR: Al enviare mensaje al jugador: " + err.Error())
					}
				}
				break

			case "perdida":
				fmt.Printf("%s recieved: %s\n", conn.RemoteAddr(), string(msg))
				if spl[1] == "1" {
					if err = conexionj2.WriteMessage(msgType, []byte("perdida "+spl[2]+" "+spl[3]+" "+spl[4]+" "+spl[5])); err != nil {
						fmt.Printf("ERROR: Al enviare mensaje al jugador: " + err.Error())
					}
				} else {
					if err = conexionj1.WriteMessage(msgType, []byte("perdida "+spl[2]+" "+spl[3]+" "+spl[4]+" "+spl[5])); err != nil {
						fmt.Printf("ERROR: Al enviare mensaje al jugador: " + err.Error())
					}
				}

			case "devuelta":
				fmt.Printf("%s recieved: %s\n", conn.RemoteAddr(), string(msg))
				if spl[1] == "1" {
					if err = conexionj2.WriteMessage(msgType, []byte("devuelta "+spl[2]+" "+spl[3]+" "+spl[4]+" "+spl[5])); err != nil {
						fmt.Printf("ERROR: Al enviare mensaje al jugador: " + err.Error())
					}
				} else {
					if err = conexionj1.WriteMessage(msgType, []byte("devuelta "+spl[2]+" "+spl[3]+" "+spl[4]+" "+spl[5])); err != nil {
						fmt.Printf("ERROR: Al enviare mensaje al jugador: " + err.Error())
					}
				}
				break
			case "endgame":
				fmt.Println("Fin del juego");
				totalListos=0;
				totalJugadores=0;
				conexionj1 = nil;
				conexionj2 = nil;
				exit = true;
				break;
			default:
				fmt.Println("Error en el servidor. Mensaje incontrolado: " + string(msg))
			}

		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "pong.html")
	})

	http.ListenAndServe(":2020", nil)
}
