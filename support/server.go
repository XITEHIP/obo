package support

import (
	"fmt"
	"net"
	"time"
	"os"
	"os/signal"
	"sync"
)

const (
	WorkerCount = 2
)

type Task struct {
	Id int32
	Message string
}

var wg sync.WaitGroup
var taskChannel = make(chan Task)
var signChannel = make(chan os.Signal, 1)
var exitChanel  = make(chan int)

var rmsg chan string

func TcpServer(msg chan string) {
	rmsg = msg
	go installSign()
	tcpAddr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:2202")
	if err != nil {
		panic("解析ip地址失败: " + err.Error())
	}
	fmt.Println("Listening 127.0.0.1:2202 ....")
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		panic("监听TCP失败: " + err.Error())
	}
	fmt.Println("Listen success on 127.0.0.1:2202 with tcp4")
	defer func() {
		fmt.Println("Close listenning ....")
		listener.Close()
		fmt.Println("Shutdown")
	}()

	connChannel := make(chan net.Conn)

	go accept(listener, connChannel)
	go handleConn(connChannel)
	go taskDispatch()

	for {
		select {
		case <- signChannel:
			fmt.Println("Get shutdown sign")
			go notifyGoroutingExit()
			goto EXIT
		}
	}

EXIT:
	fmt.Println("Waiting gorouting exit ....")
	wg.Wait()
}

func accept(listener * net.TCPListener, connChannel chan net.Conn) {
	for {
		connection, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println("Accept 失败: " + err.Error())
		} else {
			connChannel <- connection
		}
	}
}

func handleConn(connChannel chan net.Conn) {
	fmt.Println("Wating connection ....")
	for {
		select {
		case conn := <- connChannel:
			remoteAddr := conn.RemoteAddr()
			fmt.Println("Client " + remoteAddr.String() + " connected")
			readConn(&conn)
		}
	}

}

func readConn(conn *net.Conn) {
	for {
		(*conn).SetReadDeadline(time.Now().Add(5 * time.Second))
		buf := make([]byte, 1024)
		_, err := (*conn).Read(buf)
		if err != nil {
			fmt.Println("Read connection error: " + err.Error())
			if  err.Error() == "EOF" {
				(*conn).Close()
				fmt.Println("Close connection " + (*conn).RemoteAddr().String())
				break
			}
		}
		if buf != nil {
			fmt.Printf("Read message from connect: %s", string(buf))
			writeConn(conn, buf)
			var task Task
			task.Id = 1
			task.Message = string(buf)
			taskChannel <- task
		}
	}
}

func writeConn(conn *net.Conn, msg []byte)  {
	_, err := (*conn).Write(msg)
	if err != nil {
		fmt.Println("Write connection error: " + err.Error())
		if  err.Error() == "EOF" {
			(*conn).Close();
			fmt.Println("Close connection " + (*conn).RemoteAddr().String())
		}
	}
}


func taskDispatch() {
	fmt.Println("Init task moniter ....")
	for i := 0; i < WorkerCount; i ++ {
		go loop()
	}
	fmt.Println("Init task moniter DONE!")
}

func loop() {
	ticker := time.NewTicker(10 * time.Second)
	wg.Add(1)
	defer func() {
		defer wg.Done()
		defer ticker.Stop()
	}()
	for {
		fmt.Println("Wating task ....")
		select {
		case task := <- taskChannel:
			fmt.Println("Task comming: " + task.Message)
			rmsg <- task.Message
		case <- exitChanel:
			fmt.Println("Woker get exit sign")
			goto STOP
			//default:
		}
		// Epoll, 去读任务数据, 不需要处理超时的情况
		//select {
		//  case <- ticker.C:
		//      fmt.Println(time.Now().String() +  " No task after 10 second")
		//      break
		//}
	}
STOP:
	//TODO: Clear undo task
}

func installSign() {
	signal.Notify(signChannel, os.Interrupt, os.Kill)
}

func notifyGoroutingExit() {
	for i := 0; i < WorkerCount; i ++ {
		exitChanel <- 1
	}
}