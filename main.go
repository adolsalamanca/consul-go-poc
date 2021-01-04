package main

import (
	"flag"
	"fmt"
	"github.com/adolsalamanca/consul-kv/api"
	consul "github.com/hashicorp/consul/api"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	LocalAddress     = "0.0.0.0"
	ServiceAddress   = "host.docker.internal"
	ConsulAddress    = "http://127.0.0.1:8500"
	RenewSessionTime = "5s"
)

func main() {
	var leader = false
	isLeader := &leader

	port := flag.Int("port", 3001, "port of http server")
	flag.Parse()

	fmt.Printf("the port is %d\n", *port)
	client, err := consul.NewClient(&consul.Config{
		Address: ConsulAddress,
		Scheme:  "http",
	})
	if err != nil {
		log.Fatalf("could not create new consul client, %s", err)
	}

	h, err := os.Hostname()
	if err != nil {
		log.Fatalf("could not retrieve hostname, %s", err)
	}

	err = client.Agent().ServiceRegister(&consul.AgentServiceRegistration{
		Address: fmt.Sprintf("%s:%d", ServiceAddress, *port),
		ID:      fmt.Sprintf("%s-%d", h, os.Getpid()),
		Name:    fmt.Sprintf("%s-%d", h, os.Getpid()),
		Tags:    []string{"api"},
		Check: &consul.AgentServiceCheck{
			HTTP:     fmt.Sprintf("http://%s:%d/_health", ServiceAddress, *port),
			Interval: "10s",
		},
	})
	if err != nil {
		log.Fatalf("could not register service, %s", err)
	}

	sessionID, _, err := client.Session().Create(&consul.SessionEntry{
		Name:     "service/api/leader",
		Behavior: consul.SessionBehaviorDelete,
		TTL:      "10s",
	}, nil)
	if err != nil {
		log.Fatalf("could not create consul session, %s", err)
	}

	p := &consul.KVPair{
		Key:     "service/api/leader",
		Value:   []byte(sessionID),
		Session: sessionID,
	}
	fmt.Printf("Session to be acquired: %s\n", sessionID)

	doneChan := make(chan struct{})
	defer close(doneChan)

	go WaitForLeadership(isLeader, client, p, h)
	go RenewSession(client, sessionID, doneChan)

	go ListenShutdown(isLeader, client, p)
	go StartApi(*port)

	<-doneChan
}

func WaitForLeadership(isLeader *bool, client *consul.Client, p *consul.KVPair, h string) {
	for {
		leader, _, err := client.KV().Acquire(p, nil)
		isLeader = &leader
		if err != nil {
			fmt.Errorf("error trying to acquire leadership, %s\n", err)
		}

		if !leader {
			fmt.Printf("Im not the leader %s-%d\n\n", h, os.Getpid())
			t := time.NewTimer(1 * time.Second)
			<-t.C
			fmt.Printf("\n\n About to check leadership again... %s-%d\n\n", h, os.Getpid())
		} else {
			fmt.Printf("I'M THE LEADER %s-%d\n\n", h, os.Getpid())
			t := time.NewTimer(3 * time.Second)
			<-t.C
			fmt.Printf("\n\n Will check leadership again... %s-%d\n\n", h, os.Getpid())
		}
	}
}

func RenewSession(client *consul.Client, sessionID string, doneChan chan struct{}) {
	// RenewPeriodic is used to periodically invoke Session until a doneChan is closed.
	// This is used in a long running goroutine to ensure a session stays valid.
	client.Session().RenewPeriodic(
		RenewSessionTime,
		sessionID,
		nil,
		doneChan,
	)
}

func StartApi(port int) {
	e := echo.New()
	e.Use(middleware.Logger())
	e.GET("/api", api.MainHandler)
	e.GET("/_health", api.HealthHandler)
	e.Start(fmt.Sprintf("%s:%d", LocalAddress, port))
}

func ListenShutdown(isLeader *bool, c *consul.Client, p *consul.KVPair) {
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT)
	sig := <-s

	if *isLeader {
		success, _, err := c.KV().Release(p, nil)
		if err != nil || !success {
			log.Printf("could not release session")
			os.Exit(1)
		}
	}

	log.Printf("Shutting down system due to %s", sig)
	os.Exit(0)
}
