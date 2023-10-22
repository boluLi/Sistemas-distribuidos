package main

import (
	"fmt"
	"ms"
	"os"
	"ra"
	"strconv"
	"sync"
)

func writer(msgs *ms.MessageSystem, radb *ra.RASharedDB, myFile string, text string, me int, wg *sync.WaitGroup) {
	for {
		radb.PreProtocol()
		gf.EscribirFichero(myFile, text)
		for i := 1; i <= ra.N; i++ {
			if i != me {
				msgs.Send(i, mr.Update{text})
			}
		}
		radb.PostProtocol()
	}
}

func main() {
	meString := os.Args[1]
	fmt.Println("Proceso escritor con pid " + meString)
	me, _ := strconv.Atoi(meString)
	text := "a"
	myFile := "fichero_" + meString + ".txt"
	usersFile := "./ms/users.txt"
	gf.CrearFichero(myFile)

	// Declaración de canales
	reqch := make(chan ra.Request) // canal para la redirección de requests al ra
	repch := make(chan ra.Reply)   // canal para la redirección de replies al ra

	// Inicialización del ms
	messageTypes := []ms.Message{ra.Request{}, ra.Reply{}, mr.Update{}, mr.Barrier{}}
	msgs := ms.New(me, usersFile, messageTypes)
	go mr.ReceiveMessage(msgs, myFile, reqch, repch)

	// Inicialización del ra
	radb := ra.New(msgs, me, "write", reqch, repch)

	// Lanzamiento del proceso escritor
	var wg sync.WaitGroup
	wg.Add(1)
	go writer(msgs, radb, myFile, text, me, &wg)
	wg.Wait()
}
