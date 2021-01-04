package main

import (
	"flag"
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
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
	client, err := api.NewClient(&api.Config{
		Address: ConsulAddress,
		Scheme:  "http",
	})

	if err != nil {
		panic(err)
	}

	h, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	err = client.Agent().ServiceRegister(&api.AgentServiceRegistration{
		Address: fmt.Sprintf("%s:%d", ServiceAddress, *port),
		ID:      fmt.Sprintf("%s-%d", h, os.Getpid()),
		Name:    fmt.Sprintf("%s-%d", h, os.Getpid()),
		Tags:    []string{"monitoring"},
		Check: &api.AgentServiceCheck{
			HTTP:     fmt.Sprintf("http://%s:%d/_health", ServiceAddress, *port),
			Interval: "10s",
		},
	})
	if err != nil {
		panic(err)
	}

	sessionID, _, err := client.Session().Create(&api.SessionEntry{
		Name:     "service/monitoring/leader",
		Behavior: api.SessionBehaviorDelete,
		TTL:      "10s",
	}, nil)
	if err != nil {
		panic(err)
	}

	p := &api.KVPair{
		Key:     "service/monitoring/leader",
		Value:   []byte(sessionID),
		Session: sessionID,
	}

	fmt.Printf("Session to be acquired: %s\n", sessionID)

	doneChan := make(chan struct{})

	go func(isLeader *bool) {
		for {
			leader, _, err = client.KV().Acquire(p, nil)
			isLeader = &leader
			if err != nil {
				fmt.Errorf("error trying to acquire leadership, %s\n", err)
			}

			if !leader {
				fmt.Printf("Im not the leader %s-%d\n\n", h, os.Getpid())
				t := time.NewTimer(1 * time.Second)
				<-t.C
				fmt.Printf("About to check leadership again... %s-%d\n\n", h, os.Getpid())
			} else {
				fmt.Printf("I'M THE LEADER %s-%d\n\n", h, os.Getpid())
				t := time.NewTimer(3 * time.Second)
				<-t.C
				fmt.Printf("Will check leadership again... \n\n%s-%d\n\n", h, os.Getpid())
			}
		}
	}(isLeader)

	go RenewSession(client, sessionID, doneChan)()

	defer close(doneChan)
	go ListenShutdown(isLeader, client, p)
	go StartApi(*port)

	<-doneChan

}

func RenewSession(client *api.Client, sessionID string, doneChan chan struct{}) func() {
	return func() {
		// RenewPeriodic is used to periodically invoke Session until a doneChan is closed.
		// This is used in a long running goroutine to ensure a session stays valid.
		client.Session().RenewPeriodic(
			RenewSessionTime,
			sessionID,
			nil,
			doneChan,
		)
	}
}

func StartApi(port int) {
	e := echo.New()
	e.Use(middleware.Logger())
	e.GET("/api", apiHandler)
	e.GET("/_health", healthHandler)
	e.Start(fmt.Sprintf("%s:%d", LocalAddress, port))
}

func ListenShutdown(isLeader *bool, c *api.Client, p *api.KVPair) {
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

func apiHandler(c echo.Context) error {
	h, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	return c.String(http.StatusOK, fmt.Sprintf("api service working from %s, process %d", h, os.Getpid()))
}

func healthHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Health check working")
}
