package main

import (
	"flag"
	"fmt"
	"github.com/adolsalamanca/consul-go-poc/api"
	"github.com/adolsalamanca/consul-go-poc/internal"
	consul "github.com/hashicorp/consul/api"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"os"
	"time"
)

const (
	LocalAddress   = "0.0.0.0"
	ServiceAddress = "host.docker.internal"
	ConsulAddress  = "http://127.0.0.1:8500"
)

func main() {
	isLeader := make(chan bool)
	var client internal.Clienter
	var err error

	port := flag.Int("port", 3001, "port of http server")
	flag.Parse()

	client, err = consul.NewClient(&consul.Config{
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
	pid := os.Getpid()

	fmt.Printf("Started client, id: %s-%d\n", h, pid)
	err = client.Agent().ServiceRegister(&consul.AgentServiceRegistration{
		Address: fmt.Sprintf("%s:%d", ServiceAddress, *port),
		ID:      fmt.Sprintf("%s-%d", h, pid),
		Name:    fmt.Sprintf("%s-%d", h, pid),
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
		Name:      "service/api/leader",
		Behavior:  consul.SessionBehaviorDelete,
		TTL:       "10s",
		LockDelay: 5 * time.Second,
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

	sig := make(chan os.Signal, 1)
	go StartApi(*port, h, pid)
	go internal.ListenShutdown(client, internal.KvFromClient, sig, isLeader, p)
	go internal.ForWaitLeadership(client, internal.KvFromClient, p, isLeader, h)
	go internal.RenewSession(client, sessionID, doneChan)

	<-doneChan
}

func StartApi(port int, hostname string, pid int) {
	e := echo.New()
	e.Use(middleware.Logger())
	e.GET("/api", func(c echo.Context) error {
		return api.MainHandler(c, hostname, pid)
	})
	e.GET("/_health", api.HealthHandler)
	err := e.Start(fmt.Sprintf("%s:%d", LocalAddress, port))
	if err != nil {
		log.Fatalf("could not start API, %s", err)
	}
}
