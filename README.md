# go-shutdown
A lightweight package to assist in graceful shutdown

## Installation
go get github.com/hemaan/go-shutdown

## Utility

- <b>shutdown.Requested() bool</b> : returns a boolean indicating if interrupt was received.
- <b>shutdown.SubscribeFunctionForInterrupt(callback func())</b> : lets a function subscribe for a shutdown interrupt
- <b>shutdown.SubscribeChannelForInterrupt(channen chan<- struct{})</b> : lets a channel subscribe for a shutdown interrupt
- <b>shutdown.Add(num int)</b> : lets you add a method to the shutdown wait group. sync.WaitGroup.Add() underhood
- <b>shutdown.Done()</b> : lets you remove a method to the shutdown wait group. sync.WaitGroup.Done() underhood
- <b>shutdown.WaitForProcesses()</b> : waits for processes currently in the shutdown wait group. sync.WaitGroup.Wait() underhood. 
- <b>shutdown.RemainingProcesses()</b> : returns remaining processes currently in the shutdown wait group
