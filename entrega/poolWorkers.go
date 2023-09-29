
 package main

 import (
	 "encoding/gob"
	 "log"
	 "net"
	 "practica1/com"
	 //"encoding/binary"
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
 
 func receiveRequest(connChanel <-chan net.Conn) {
	for conn := range connChanel {
		var request com.Request
		decoder := gob.NewDecoder(conn)
		err := decoder.Decode(&request) //  receive the request
		com.CheckError(err)
		if(request.Interval.Max != 0){ // Si no envia mensaje de end Tp interval.max =! 0 
			sendReply(conn,request)
		}else{
		   log.Println("***** FIN ******")
		   conn.Close()
		}
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
 
	 endpoint := "localhost:30000"
	 CONN_TYPE := "tcp"
	 listener, err := net.Listen(CONN_TYPE, endpoint)
	 com.CheckError(err)
	 //Creamos el pool de workers con sus canales
	 workers := 10
	 var canalJob chan net.Conn
	 canalJob = make(chan net.Conn)
	for i := 1; i <= workers; i++ {
		go receiveRequest(canalJob) // we're finished with this client
	}
	for {
		conn, err := listener.Accept()
		com.CheckError(err)
		canalJob <- conn
 	}
		
	
	 
 }
 