package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"context"
)

func main() {
	env, err := LoadEnv()
	if err != nil {
		log.Fatalf("Failed to load and parse env variables: %s\n", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Println("Connecting with database")
	client, err := connectDatabase(ctx, env.DatabaseUri())
	if err != nil {
		fmt.Printf("Failed to connect with database : %s\n", err)
		return
	}

	addr := fmt.Sprintf("%s:%d", env.Host(), env.Port())
	log.Printf("Starting server in %s\n", addr)

	server := startServer(addr, client, env.SigningKey())

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigs
	fmt.Printf("\nReceived %s signal\n", sig)

	timeout := 5 * time.Second
	err = shutdownServer(ctx, server, timeout)
	if err != nil {
		fmt.Printf("Failed to shutdown server withing timeout %s : %s\n", timeout, err)
	}

	err = disconnectDatabase(ctx, client)
	if err != nil {
		fmt.Printf("Failed to disconnect database : %s\n", err)
	} else {
		fmt.Println("Disconnected database, exiting")
	}
}
