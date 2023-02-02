// main.go
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/insomnius/inapp-queue/queue"
)

var (
	server   *http.Server
	osSignal chan os.Signal
)

func main() {
	osSignal = make(chan os.Signal, 1)
	signal.Notify(osSignal, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	emailQueue := queue.NewEmailQueue()

	appEngine := gin.Default()
	appEngine.POST("/users", func(c *gin.Context) {
		startTime := time.Now()
		// Enqueue into go channel
		emailQueue.Enqueue("Send email to the user")

		// Return the response
		c.JSON(http.StatusCreated, gin.H{
			"data": gin.H{
				"username": "user1",
				"email":    "user1@gmail.com",
			},
			"message": fmt.Sprintf("Success create user in %v", time.Since(startTime)),
			"status":  http.StatusCreated,
		})
	})
	server = &http.Server{
		Addr:    ":8080",
		Handler: appEngine,
	}

	// Start the server concurrently
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Unexpected server error because of: %v\n", err)
		}
	}()

	// Start the email queue concurrently, with 10 worker
	for i := 0; i < 10; i++ {
		go emailQueue.Work()
	}

	// Catch the exit signal
	<-osSignal

	fmt.Println("Terminating server")
	server.Shutdown(context.Background())

	fmt.Println("Terminating email queue")
	// Wait untuk there is no active job in the queue
	for emailQueue.Size() > 0 {
		time.Sleep(time.Millisecond * 500)
	}

	fmt.Println("Complete terminating application")
}
