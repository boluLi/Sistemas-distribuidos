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
 )
 
 
 // PRE: verdad = !foundDivisor
 // POST: IsPrime devuelve verdad si n es primo y falso en caso contrario
 func isPrime(n int) (foundDivisor bool) {
	 foundDivisor = false
	 for i := 2; (i < n) && !foundDivisor; i++ {
		 foundDivisor = (n%i == 0)
	 }
	 return !foundDivisor
 }
 
 // PRE: interval.A < interval.B
 // POST: FindPrimes devuelve todos los números primos comprendidos en el
 //
 //	intervalo [interval.A, interval.B]
 func findPrimes(interval com.TPInterval) (primes []int) {
	 for i := interval.Min; i <= interval.Max; i++ {
		 if isPrime(i) {
			 primes = append(primes, i)
		 }
	 }
	 return primes
 }
 
 func receiveRequest(conn net.Conn) {
 
	 var request com.Request
	 decoder := gob.NewDecoder(conn)
	 err := decoder.Decode(&request) //  receive the request
	 log.Println("Request recibida" )
	 com.CheckError(err)
	 if(request.Interval.Max != 0){ // Si no envia mensaje de end Tp interval.max =! 0 
		 sendReply(conn,request)
	 }
 }
 
 func sendReply(conn net.Conn, request com.Request ){
	 var reply com.Reply
 
	 reply.Primes = findPrimes(request.Interval)
	 reply.Id = request.Id
	 encoder := gob.NewEncoder(conn)
	 // Enviar los números primos codificados
	 err := encoder.Encode(reply)
	 com.CheckError(err)
 }
 
 // COMPLETAR EL SERVIDOR  .....
 func main() {
	 endpoint := "localhost:1111"
	 CONN_TYPE := "tcp"
	 listener, err := net.Listen(CONN_TYPE, endpoint)
	 com.CheckError(err)
	 log.SetFlags(log.Lshortfile | log.Lmicroseconds)
	 log.Println("***** Listening for new connection in endpoint ", endpoint)
	 for{
	 conn, err := listener.Accept()
	 log.Println("***** New client ******")
	 receiveRequest(conn)
	 com.CheckError(err)
	 //defer conn.Close()
	 }
 }
 