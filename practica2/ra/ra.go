/*
* AUTOR: Rafael Tolosana Calasanz
* ASIGNATURA: 30221 Sistemas Distribuidos del Grado en Ingeniería Informática
*			Escuela de Ingeniería y Arquitectura - Universidad de Zaragoza
* FECHA: septiembre de 2021
* FICHERO: ricart-agrawala.go
* DESCRIPCIÓN: Implementación del algoritmo de Ricart-Agrawala Generalizado en Go
*/
package ra

import (
    "ms"
    "sync"
)

type Request struct{
    Clock  []int
    Pid     int
}

type Reply struct{
    Clock  []int
}

type RASharedDB struct {
    pid         int
    vectorPids  []int
    OurSeqNum   int
    HigSeqNum   int
    OutRepCnt   int
    ReqCS       boolean
    RepDefd     bool[]
    ms          *MessageSystem
    done        chan bool
    chrep       chan bool
    Mutex       sync.Mutex // mutex para proteger concurrencia sobre las variables
    // TODO: completar
}


func New(me int, usersFile string) (*RASharedDB) {
    messageTypes := []Message{Request, Reply}
    msgs = ms.New(me, usersFile string, messageTypes)
    ra := RASharedDB{0, []int{}, 0, 0, 0, false, []int{}, &msgs,  make(chan bool),  make(chan bool), &sync.Mutex{}}
    go receiveReply(&ra)
    return &ra
}

//Pre: Verdad
//Post: Realiza  el  PreProtocol  para el  algoritmo de
//      Ricart-Agrawala Generalizado Require mutex
func (ra *RASharedDB) PreProtocol(){
    // TODO completar
    ra.Mutex.Lock()
    // ra.OurSeqNum = ra.HigSeqNum + 1
    ra.vectorPids[ra.pid]++; // Incrementamos el reloj
    ra.OutRepCnt = 0
    ra.ReqCS = true
    ra.RepDefd = make([]bool, ra.N)
    ra.RepDefd[ra.pid] = true
    ra.askPermission()
    ra.Mutex.Unlock()
    <-ra.chrep
    // SECCION CRITICA
}
// enviar a los N - 1 procesos distribuidos una petici´on de acceso a la seccion critica
func (ra *RASharedDB) askPermission(){
    for pidAux := 0; i < n-1; i++ {
        if(ra.pid == pidAux){
            continue;
        }
        ra.Mutex.Lock()
        ra.vectorPids[ra.pid]++; // Incrementamos el reloj
        ra.ms.Send(pid, Request{Clock: ra.vectorPids, Pid: pidAux})
        ra.Mutex.Unlock()
    }
}


func receiveReply(ra *RASharedDB) {
    for {
        msg := <-ra.ms.Inbox
        switch m := msg.Msg.(type) {
        case Request:
            ra.Mutex.Lock()
            ra.handleRequest(m)
            ra.Mutex.Unlock()
        case Reply:
            ra.Mutex.Lock()
            ra.handleReply(m)
            ra.Mutex.Unlock()
        }
    }
}

func (ra *RASharedDB) handleRequest(req Request) {
    refreshClock(ra,req.Clock)
    //ra.HigSeqNum = max(ra.HigSeqNum, req.Clock)
    if ra.RepDefd[req.Pid] {
        ra.ms.Send(req.Pid, Reply{})
    } else {
        ra.OutReqCnt++
        if ra.ReqCS && (ra.vectorPids[ra.pid] < req.Clock[req.pid] || (rra.vectorPids[ra.pid] ==req.Clock[req.pid]
            && ra.pid < req.Pid)) {
            refreshClock(ra,req.Clock)
            ra.ms.Send(req.Pid, Reply{})
          
        } else {
            ra.chrep <- true
        }
    }
}



func (ra *RASharedDB) handleReply(rep Reply) {
    refreshClock(ra,rep.Clock)
    ra.OutRepCnt++
    if ra.OutRepCnt == ra.N-1 {
        ra.chrep <- true
    }
}

//Pre: Verdad
//Post: Realiza  el  PostProtocol  para el  algoritmo de
//      Ricart-Agrawala Generalizado
func (ra *RASharedDB) PostProtocol(){
    ra.Mutex.Lock()
    ra.ReqCS = false
    ra.Mutex.Unlock()

    for pidAux := 0; i < n-1; i++ {
        if(ra.pid == pidAux || !ra.RepDefd[pidAux]){
            continue;
        }
        ra.ms.Send(pidAux, Reply{})
    }
}
func (ra *RASharedDB) refreshClock(vectorPidsAux []int){
    ra.vectorPids[ra.pid]++;
    for pidAux := 0; i < n; i++ {
        ra.vectorPids[pidAux] = max(ra.vectorPids[pidAux],vectorPidsAux[pidAux])
    }
}
func (ra *RASharedDB) maxVectorPid(vectorPids []int){
    max := 0
    for pidAux := 0; i < n; i++ {
        max = max(max, vectorPids[pidAux])
    }
    return max
}

func (ra *RASharedDB,) Stop(){
    ra.ms.Stop()
    ra.done <- true
}
