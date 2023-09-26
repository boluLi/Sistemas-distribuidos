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

// COMPLETAR EL SERVIDOR  .....
func main() {

	endpoint := "localhost:30000"
	CONN_TYPE := "tcp"
	listener, err := net.Listen(CONN_TYPE, endpoint)
	com.CheckError(err)
	log.SetFlags(log.Lshortfile | log.Lmicroseconds)
	var reply com.Reply
	myInterval := com.TPInterval{
		Min: 10,
		Max: 20,
	}
	for{
		log.Println("***** Listening for new connection in endpoint ", endpoint)
		conn, err := listener.Accept()
		com.CheckError(err)
		reply.Primes = findPrimes(myInterval)
		reply.Id = 1
		encoder := gob.NewEncoder(conn)
		// Enviar los números primos codificados
		err = encoder.Encode(reply)
		defer conn.Close()
	}
		com.CheckError(err)
}
