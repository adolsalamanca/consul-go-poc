package internal_test

import (
	"github.com/adolsalamanca/consul-go-poc/internal"
	"github.com/adolsalamanca/consul-go-poc/internal/mocks"
	"github.com/golang/mock/gomock"
	consul "github.com/hashicorp/consul/api"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
	"syscall"
	"time"
)

var kv *mocks.MockKeyValuer

var _ = Describe("Main test suite", func() {

	const (
		hostname = "hostname"
	)
	var (
		ctrl     *gomock.Controller
		client   *mocks.MockClienter
		p        *consul.KVPair
		isLeader chan bool
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		client = mocks.NewMockClienter(ctrl)
		kv = mocks.NewMockKeyValuer(ctrl)
		p = &consul.KVPair{}
		isLeader = make(chan bool)
	})

	AfterEach(func() {
		defer ctrl.Finish()
	})

	Context("waitForLeadership", func() {

		It("should spend more than LeaderPoolInterval seconds if the kv Acquire returned as leader", func() {
			kv.EXPECT().Acquire(gomock.Any(), gomock.Any()).Return(true, nil, nil)
			tBef := time.Now()
			go ProcessLeaderValues(isLeader)

			internal.WaitLeadership(client, InjectKvMock, p, isLeader, hostname)
			spentTime := time.Since(tBef)

			Expect(spentTime).To(BeNumerically(">=", internal.LeaderPoolInterval*time.Second))
		})

		It("should spend more than NotLeaderPoolInterval seconds if the kv Acquire returned as not leader", func() {
			kv.EXPECT().Acquire(gomock.Any(), gomock.Any()).Return(false, nil, nil)
			tBef := time.Now()
			go ProcessLeaderValues(isLeader)

			internal.WaitLeadership(client, InjectKvMock, p, isLeader, hostname)
			spentTime := time.Since(tBef)

			Expect(spentTime).To(BeNumerically("<", internal.LeaderPoolInterval*time.Second))
			Expect(spentTime).To(BeNumerically(">=", internal.NotLeaderPoolInterval*time.Second))
		})

	})

	/*Context("listenForShutdown", func() {
		var (
			fakeExit func(code int)
			exitCode int
		)

		BeforeEach(func() {
			fakeExit = func(code int) {
				exitCode = code
				os.Exit(code)
			}
			internal.OsExit = fakeExit
		})

		It("should try to release leadership if it was leader", func() {
			sigChan := make(chan os.Signal)
			isLeader := make(chan bool)
			kv.EXPECT().Release(p, nil).Times(1).Return(true, nil, nil)

			go SendSignalWithDelay(sigChan, time.Millisecond*250)
			go SendLeaderValues(isLeader, true)

			internal.ListenShutdown(client, InjectKvMock, sigChan, isLeader, p)

			Expect(exitCode).To(BeEquivalentTo(0))
		})

		It("should exit with exit(1) if there was an error releasing kv", func() {
			sigChan := make(chan os.Signal)
			isLeader := make(chan bool)
			kv.EXPECT().Release(p, nil).Times(1).Return(false, nil, errors.New("fake err"))

			go SendSignalWithDelay(sigChan, time.Millisecond*250)
			go SendLeaderValues(isLeader, true)

			internal.ListenShutdown(client, InjectKvMock, sigChan, isLeader, p)

			Expect(exitCode).To(BeEquivalentTo(1))
		})

		It("should not release leadership if it was not leader", func() {
			sigChan := make(chan os.Signal)
			isLeader := make(chan bool)

			go SendSignalWithDelay(sigChan, time.Millisecond*50)
			go SendLeaderValues(isLeader, false)

			internal.ListenShutdown(client, InjectKvMock, sigChan, isLeader, p)

			Expect(exitCode).To(BeEquivalentTo(0))
		})

	})*/

})

func SendSignalWithDelay(sigChan chan os.Signal, delay time.Duration) {
	t := time.NewTimer(delay)
	<-t.C

	sigChan <- syscall.SIGINT
}

func SendLeaderValues(leader chan bool, v bool) {
	leader <- v
}

func ProcessLeaderValues(leader chan bool) {
	select {
	case _ = <-leader:
	}
}

func InjectKvMock(client internal.Clienter) internal.KeyValuer {
	return kv
}
