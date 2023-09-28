/*
* AUTOR: Rafael Tolosana Calasanz y Unai Arronategui
* ASIGNATURA: 30221 Sistemas Distribuidos del Grado en Ingeniería Informática
*			Escuela de Ingeniería y Arquitectura - Universidad de Zaragoza
* FECHA: septiembre de 2022
* FICHERO: server-draft.go
* DESCRIPCIÓN: contiene la funcionalidad esencial para realizar los servidores
*				correspondientes a la práctica 1
 */
 package main

 import (
	 "encoding/gob"
	 "log"
	 "net"
	 "practica1/com"
	 //"io"
 )

 

 
 
 func worker(canalClient <-chan net.Conn,endpointWorker string) {
	for connCliente := range canalClient { 
		connWorker, err := net.Dial("tcp", endpointWorker)
		com.CheckError(err)
		var request com.Request
		var reply com.Reply
		
		decoder := gob.NewDecoder(connCliente)
		encoder := gob.NewEncoder(connWorker)
		err = decoder.Decode(&request)
		com.CheckError(err)
		err = encoder.Encode(request)
		com.CheckError(err)
		
		log.Println("***** Recibir request y enviar al coworker ******")
		decoder = gob.NewDecoder(connWorker)
		encoder = gob.NewEncoder(connCliente)
		err = decoder.Decode(&reply)
		com.CheckError(err)
		err = encoder.Encode(reply)
		com.CheckError(err)
		log.Println("***** Recibir replu y enviar al cliente ******")
		
		

		//Envia bytes al cliente
		/*
		log.Println("***** Recibir request y enviar al coworker ******")
		io.Copy(connWorker, connCliente)
		log.Println("***** A ******")
		io.Copy(connCliente, connWorker)
		log.Println("***** Enviado reply al cliente ******")
		*/
		connCliente.Close()    
	}
 }
 // COMPLETAR EL SERVIDOR  .....
 func main() {
	//Cliente
	 endpointCliente := "localhost:30000"
	 CONN_TYPE := "tcp"
	 clientListener, err := net.Listen(CONN_TYPE, endpointCliente)
	 log.Println("ESCUCHANDO PUERTO para clientes"+endpointCliente)
	 com.CheckError(err)
	//Workers
	endpointWorker := "localhost:1111"
	 //Canal de cliente
	canalClient := make(chan net.Conn)
	go worker(canalClient,endpointWorker) // we're finished with this client
	for{
		conn, err := clientListener.Accept()
		com.CheckError(err)
		canalClient <- conn
 	}
 }
 