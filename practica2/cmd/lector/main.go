package main

import (
	"fmt"
	"ms"
	"os"
	"ra"
	"strconv"
	"sync"
)

func reader(radb *ra.RASharedDB, myFile string, wg *sync.WaitGroup) {
	for {
		radb.PreProtocol()
		_ = gf.LeerFichero(myFile)
		radb.PostProtocol()
	}
}

func main() {
	meString := os.Args[1]
	fmt.Println("Proceso lector con pid " + meString)
	me, _ := strconv.Atoi(os.Args[1])
	myFile := "fichero_" + meString + ".txt"
	usersFile := "./ms/users.txt"
	gf.CrearFichero(myFile)

	// Declaración de canales
	reqch := make(chan ra.Request) // canal para la redirección de requests al ra
	repch := make(chan ra.Reply)   // canal para la redirección de replies al ra

	// Inicialización del ms
	messageTypes := []ms.Message{ra.Request{}, ra.Reply{}, mr.Update{}, mr.Barrier{}}
	msgs := ms.New(me, usersFile, messageTypes)
	go mr.ReceiveMessage(&msgs, myFile, reqch, repch)

	// Inicialización del ra
	radb := ra.New(&msgs, me, "read", reqch, repch)

	// Lanzamiento del proceso lector
	var wg sync.WaitGroup
	wg.Add(1)
	go reader(radb, myFile, &wg)
	wg.Wait()
}
