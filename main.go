package main

import (
	"flag"
	"fmt"
	"github.com/adolsalamanca/consul-go-poc/api"
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

type Clienter interface {
	Agent() *consul.Agent
	Session() *consul.Session
	KV() *consul.KV
}

func main() {
	isLeader := make(chan bool, 500)
	var client Clienter
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
	go ListenShutdown(sig, isLeader, client, p)
	go WaitForLeadership(isLeader, client, p, h)
	go RenewSession(client, sessionID, doneChan)
	go StartApi(*port, h, pid)

	<-doneChan
}

func WaitForLeadership(isLeader chan bool, client Clienter, p *consul.KVPair, h string) {
	for {
		leader, _, err := client.KV().Acquire(p, nil)
		isLeader <- leader
		if err != nil {
			fmt.Printf("error trying to acquire leadership, %s\n", err)
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

func RenewSession(client Clienter, sessionID string, doneChan chan struct{}) {
	// RenewPeriodic is used to periodically invoke Session until a doneChan is closed.
	// This is used in a long running goroutine to ensure a session stays valid.
	err := client.Session().RenewPeriodic(
		RenewSessionTime,
		sessionID,
		nil,
		doneChan,
	)

	if err != nil {
		log.Fatalf("could not renew consul session, %s", err)
	}

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

func ListenShutdown(sig chan os.Signal, isLeader chan bool, c Clienter, p *consul.KVPair) {
	var leader bool
	var s os.Signal
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT)

	for {
		select {
		case s = <-sig:
			if leader {
				success, _, err := c.KV().Release(p, nil)
				log.Printf("About to release leadership from this client \n")
				if err != nil || !success {
					log.Printf("could not release session")
					os.Exit(1)
				}
			} else {
				log.Printf("This client was not leader, so nothing occurs \n")
			}

			log.Printf("Shutting down system due to %s \n", s)
			os.Exit(0)

		case leader = <-isLeader:
		default:

		}

	}

}
