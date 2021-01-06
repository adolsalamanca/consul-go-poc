package internal

import (
	"fmt"
	consul "github.com/hashicorp/consul/api"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	RenewSessionTime      = "5s"
	NotLeaderPoolInterval = 1
	LeaderPoolInterval    = 3
)

var OsExit = os.Exit

type Clienter interface {
	Agent() *consul.Agent
	Session() *consul.Session
	KV() *consul.KV
}

type KeyValuer interface {
	Acquire(p *consul.KVPair, q *consul.WriteOptions) (bool, *consul.WriteMeta, error)
	Release(p *consul.KVPair, q *consul.WriteOptions) (bool, *consul.WriteMeta, error)
}

type KvRetriever func(c Clienter) KeyValuer

func KvFromClient(c Clienter) KeyValuer {
	return c.KV()
}

func ForWaitLeadership(client Clienter, kvRetriever KvRetriever, p *consul.KVPair, isLeader chan bool, h string) {
	for {
		WaitLeadership(client, kvRetriever, p, isLeader, h)
	}
}

func WaitLeadership(client Clienter, kvRetriever KvRetriever, p *consul.KVPair, isLeader chan bool, h string) {
	kv := kvRetriever(client)
	leader, _, err := kv.Acquire(p, nil)
	isLeader <- leader
	if err != nil {
		fmt.Printf("error trying to acquire leadership, %s\n", err)
	}

	if leader {
		fmt.Printf("I'M THE LEADER %s-%d\n\n", h, os.Getpid())
		t := time.NewTimer(LeaderPoolInterval * time.Second)
		<-t.C
		return
	}

	fmt.Printf("Im not the leader %s-%d\n\n", h, os.Getpid())
	t := time.NewTimer(NotLeaderPoolInterval * time.Second)
	<-t.C
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

func ListenShutdown(c Clienter, kvRetriever KvRetriever, sig chan os.Signal, isLeader chan bool, p *consul.KVPair) {
	var leader bool
	var s os.Signal
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT)
	for {
		select {
		case s = <-sig:
			if leader {
				kv := kvRetriever(c)
				success, _, err := kv.Release(p, nil)
				log.Printf("Release leadership from client \n")
				if err != nil || !success {
					log.Printf("could not release session")
					OsExit(1)
				}
			} else {
				log.Printf("Client was not leader, so leadership should not be released \n")
			}

			log.Printf("Shutting down system due to %s \n", s)
			OsExit(0)

		// Add leader values channel processing to avoid locks
		case leader = <-isLeader:
		default:
		}

	}

}
