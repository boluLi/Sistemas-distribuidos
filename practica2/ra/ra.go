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
    op      string
}

type Reply struct{
    Clock  []int
}

const (
    Reader  = "reader"
    Writer = "writer"
)
type ExclusionMatrix map[string]map[string]bool


type RASharedDB struct {
    pid         int             // identificador de proceso
    vectorPids  []int           // reloj vectorial
    OurSeqNum   int             // numero de secuencia de la peticion
    HigSeqNum   int             // numero de secuencia mas alto recibido
    OutRepCnt   int             // numero de respuestas recibidas
    ReqCS       boolean         // true si el proceso quiere entrar en la seccion critica
    RepDefd     bool[]          // true si la respuesta a un proceso ha sido pospuesuta
    exclude ExclusionMatrix     // matriz de exclusion mutua
    ms          *MessageSystem  // sistema de mensajes
    done        chan bool       // canal de comunicacion para parar el proceso
    chrep       chan bool       // canal de comunicacion para esperar todas las respuestas recibidas
    Mutex       sync.Mutex      // mutex para proteger las variables compartidas
    op          string          // operacion a realizar
}


func New(me int, usersFile string) (*RASharedDB) {
    messageTypes := []Message{Request, Reply}
    msgs = ms.New(me, usersFile string, messageTypes)
    // Definicion de las reglas de exclusion mutua
    var exclusionRules = ExclusionMatrix{
        Read: {
            Read:  false,
            Write: true,
        },
        Write: {
            Read:  true,
            Write: true,
        },
    }
    ra := RASharedDB{0, []int{}, 0, 0, 0, false, []int{}, exclusionRules, &msgs,  make(chan bool),  make(chan bool), &sync.Mutex{}}
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
    // propia respuesta ya recibida
    ra.RepDefd[ra.pid] = true
    ra.askPermission()
    ra.Mutex.Unlock()
    // Esperamos a recibir todas las respuestas
    <-ra.chrep
    // SECCION CRITICA
}
// enviar a los N - 1 procesos distribuidos una petici´on de acceso a la seccion critica
func (ra *RASharedDB) askPermission(){
    for pidAux := 0; i < n-1; i++ {
        // No enviamos la peticion a nosotros mismos
        if(ra.pid == pidAux){
            continue;
        }
        ra.Mutex.Lock()
        ra.vectorPids[ra.pid]++; // Incrementamos el reloj
        // Enviar la peticion a todos los procesos
        ra.ms.Send(pid, Request{Clock: ra.vectorPids, Pid: pidAux, op: ra.op})
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
    // prioridad a la peticion mas antigua
    var noPriority bool = ra.vectorPids[ra.pid] < req.Clock[req.pid] || (rra.vectorPids[ra.pid] ==req.Clock[req.pid])
        && (ra.pid < req.Pid)

    refreshClock(ra,req.Clock)
    //ra.HigSeqNum = max(ra.HigSeqNum, req.Clock)
    if !ra.ReqCS || !ra.exclude[req.op] {
        ra.ms.Send(req.Pid, Reply{})
    } 
    else {
        // si no tenemos prioridad, respondemos
         if noPriority {
            refreshClock(ra,req.Clock)
            ra.ms.Send(req.Pid, Reply{})
          
        } else { // si tenemos prioridad, posponemos la respuesta
            ra.RepDefd[req.Pid] = true
        }
    }
}



func (ra *RASharedDB) handleReply(rep Reply) {
    refreshClock(ra,rep.Clock)
    // incrementamos número de respuestas recibidas
    ra.OutRepCnt++
    // si hemos recibido todas las respuestas, entramos en la seccion critica
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
    // Una vez hemos salido de la seccion critica, enviamos las respuestas pospuestas
    for pidAux := 0; i < n-1; i++ {
        // No enviamos la peticion a nosotros mismos ni a los que no han solicitado entrar en la seccion critica
        if(ra.pid == pidAux || !ra.RepDefd[pidAux]){
            continue;
        }
        ra.ms.Send(pidAux, Reply{})
    }
}
func (ra *RASharedDB) refreshClock(vectorPidsAux []int){
    // Incremento del reloj propio
    ra.vectorPids[ra.pid]++;
    // Para cada reloj de cada proceso, cogemos el maximo entre el local y el recibido
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