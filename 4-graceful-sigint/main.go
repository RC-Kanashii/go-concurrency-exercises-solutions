//////////////////////////////////////////////////////////////////////
//
// Given is a mock process which runs indefinitely and blocks the
// program. Right now the only way to stop the program is to send a
// SIGINT (Ctrl-C). Killing a process like that is not graceful, so we
// want to try to gracefully stop the process first.
//
// Change the program to do the following:
//   1. On SIGINT try to gracefully stop the process using
//          `proc.Stop()`
//   2. If SIGINT is called again, just kill the program (last resort)
//

package main

import (
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Create a process
	proc := MockProcess{}

	sigCh := make(chan os.Signal)        // 传递信号的通道
	shutdownCh := make(chan bool)        // 传递结束程序信息的通道
	signal.Notify(sigCh, syscall.SIGINT) // 监听SIGINT信号

	// Run the process (blocking)
	go proc.Run()

	go func() {
		<-sigCh
		go proc.Stop()
		signal.Reset()
		<-sigCh
		shutdownCh <- true
	}()

	<-shutdownCh
}
