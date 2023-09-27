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
	 //"encoding/gob"
	 "log"
	 "net"
	 "practica1/com"
	 "io"
 )
 
 func worker(canalClient <-chan net.Conn,connWorker net.Conn) {
	for {
        select {
        case connCliente := <-canalClient:
            // Ahora usamos directamente connCliente sin iniciar una nueva conexión
            _, err := io.Copy(connCliente, connWorker)
            com.CheckError(err)
            connCliente.Close()
            connWorker.Close()
        // No necesitamos un caso para workerChanel ya que estamos obteniéndolo directamente
        // dentro del caso de clientChanel. A menos que haya otra lógica específica que quieras implementar.
        }
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
	 workerListener, err := net.Listen(CONN_TYPE, endpointWorker)
	 log.Println("ESCUCHANDO PUERTO para worker"+ endpointWorker)
	 com.CheckError(err)

	 //Canal de cliente
	 canalClient := make(chan net.Conn)
	go func(){
		for {
			conn, err := workerListener.Accept()
			log.Println("***** WORKER CONECTADO ******")
			go worker(canalClient,conn)
			com.CheckError(err)
 		}
	}()

	go func(){
		for {
			conn, err := clientListener.Accept()
			com.CheckError(err)
			canalClient <- conn
 		}
	}()
	for{}
 }
 