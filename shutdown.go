package shutdown

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var shutdownWG sync.WaitGroup
var requested bool = false
var intr = make(chan os.Signal, 10)
var subscribedFunctions []func()
var subscribedChannels []chan<- struct{}
var processes int // keeps track of number processes in wait group
var processesMux sync.Mutex
var channelMux sync.Mutex
var functionMux sync.Mutex

// will run before everything
func init() {
	go listenForInterrupt()
}

func listenForInterrupt() {
	signal.Notify(intr, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-intr

	var wg sync.WaitGroup = sync.WaitGroup{}
	wg.Add(2)
	go notifyAllChannels(&wg)
	go notifyAllFunctions(&wg)
	requested = true
	wg.Wait()
}

func notifyAllFunctions(wg *sync.WaitGroup) {
	defer wg.Done()
	functionMux.Lock()
	defer functionMux.Unlock()

	for _, callback := range subscribedFunctions {
		go callback()
	}
}

func notifyAllChannels(wg *sync.WaitGroup) {
	defer wg.Done()
	channelMux.Lock()
	defer channelMux.Unlock()

	for _, channel := range subscribedChannels {
		go sendSafelyToChannel(channel)
	}
}

// sending on closed channel generated panic.
// this function prevents main function from panicking and helps continue our main loop
func sendSafelyToChannel(c chan<- struct{}) {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()

	select {
	case c <- struct{}{}:
	default:
		return
	}

}

// Lets a channel listen for a shutdown event
//
// If a channel is blocking or closed, simply skips it
func SubscribeChannelForInterrupt(channel chan<- struct{}) {
	channelMux.Lock()
	defer channelMux.Unlock()
	subscribedChannels = append(subscribedChannels, channel)
}

// Subscribes a callback function to listen for interrupt event.
//
// The callback function is called asynchronous.
func SubscribeFunctionForInterrupt(callback func()) {
	functionMux.Lock()
	defer functionMux.Unlock()
	subscribedFunctions = append(subscribedFunctions, callback)
}

func Requested() bool {
	return requested
}

// Returns number of processes that are currently shutting down
func RemainingProcesses() int {
	processesMux.Lock()
	defer processesMux.Unlock()

	return processes
}

// Registers a process for graceful shutdown
//
// wg.Add() underhood
//
func Add(num int) {
	processesMux.Lock()
	defer processesMux.Unlock()
	shutdownWG.Add(num)
	processes++
}

// Decerements registered process by one.
//
// wg.Done() underhood
//
func Done() {
	processesMux.Lock()
	defer processesMux.Unlock()
	shutdownWG.Done()
	processes--
}

// Waits for all the processes registered in shutdown wait group
//
// wg.Wait() underhood
func WaitForProcesses() {
	shutdownWG.Wait()
}

