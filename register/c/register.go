package main

import (
	"C"
	"log"
	"strings"

	"golang.org/x/sys/unix"

	"github.com/go-service-register/register"
)
import (
	"os"
	"runtime"
	"strconv"
)

func init() {
	runtime.GOMAXPROCS(1)
	runtime.LockOSThread()
}

//export Register
func Register(etcdAddrsStr string, serviceName, serviceValue string, dialTimeout, ttl int) int {
	// TODO: using c++ client
	// as hard to convert []string from C++ to Go, it's using multiple string emulating array
	etcdAddrs := strings.Split(etcdAddrsStr, "|")
	// TODO: will it be recycled after return and is it safe for goroutine?
	node := register.New(etcdAddrs, serviceName, serviceValue, dialTimeout, ttl)
	if node == nil {
		// failure
		return 0
	}
	err := node.Register()
	if err != nil {
		// failure
		return 0
	}
	// success
	return 1
}

// Please use RegisterWithAffinity before other core binding

//export RegisterWithAffinity
func RegisterWithAffinity(etcdAddrsStr string, serviceName, serviceValue string, dialTimeout, ttl, core int) int {
	// TODO: using c++ client
	// as hard to convert []string from C++ to Go, it's using multiple string emulating array
	etcdAddrs := strings.Split(etcdAddrsStr, "|")
	// TODO: will it be recycled after return and is it safe for goroutine?
	node := register.New(etcdAddrs, serviceName, serviceValue, dialTimeout, ttl)
	if node == nil {
		// failure
		return 0
	}
	err := node.Register()
	if err != nil {
		// failure
		return 0
	}
	// set all of thread
	cpuSet := &unix.CPUSet{}
	cpuSet.Zero()
	cpuSet.Set(core)
	pid := os.Getpid()
	f, err := os.Open("/proc/" + strconv.Itoa(pid) + "/task")
	if err != nil {
		log.Println(err.Error())
		return 0
	}
	tasks, err := f.Readdirnames(0)
	if err != nil {
		log.Println(err.Error())
		return 0
	}
	for _, tid := range tasks {
		nTid, err := strconv.Atoi(tid)
		if err != nil {
			log.Println(err.Error())
			return 0
		}
		log.Println("nTid = ", nTid)
		err = unix.SchedSetaffinity(nTid, cpuSet)
		if err != nil {
			log.Println(err.Error())
			return 0
		}
	}
	log.Println("All Done")
	// success
	return 1
}

func main() {
}
