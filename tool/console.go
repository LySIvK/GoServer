package tool

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
)

type CommandHandler func(args []string) bool

var (
	command                              = make([]byte, 1024)
	reader                               = bufio.NewReader(os.Stdin)
	HandlerMap map[string]CommandHandler = make(map[string]CommandHandler, 20)
)

func StartConsole() {
	go consoleroutine()
}

func StartConsoleWait() {
	consoleroutine()
}

func consoleroutine() {
	for {
		command, _, _ = reader.ReadLine()
		Args := strings.Split(string(command), " ")

		cmdhandler, ok := HandlerMap[Args[0]]
		if ok {
			cmdhandler(Args)
			continue
		}

		switch string(Args[0]) {
		case "cpus":
			fmt.Println(runtime.NumCPU(), " cpus and ", runtime.GOMAXPROCS(0), " in use")

		case "routines":
			fmt.Println("Current number of goroutines: ", runtime.NumGoroutine())

		case "setcpus":
			n, _ := strconv.Atoi(Args[1])
			runtime.GOMAXPROCS(n)
			fmt.Println(runtime.NumCPU(), " cpus and ", runtime.GOMAXPROCS(0), " in use")

		default:
			fmt.Println("Command error, try again.")
		}
	}
}

func HandleFunc(cmd string, mh CommandHandler) {
	if HandlerMap == nil {
		HandlerMap = make(map[string]CommandHandler, 20)
	}

	HandlerMap[cmd] = mh

	return
}
