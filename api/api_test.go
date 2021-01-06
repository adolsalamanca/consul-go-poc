package api_test

import (
	"fmt"
	"github.com/adolsalamanca/consul-go-poc/api"
	"github.com/labstack/echo/v4"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
)

const (
	pid      = 123
	hostname = "hostname"
)

var _ = Describe("Api handlers test", func() {

	It("main working as expected", func() {
		w := httptest.NewRecorder()
		e := NewEchoContext(w)

		err := api.MainHandler(e, hostname, pid)

		Expect(err).ToNot(HaveOccurred())
		Expect(w.Code).To(BeEquivalentTo(http.StatusOK))
		Expect(w.Body.String()).To(BeEquivalentTo(fmt.Sprintf("API working from %s, process %d", hostname, pid)))

	})

	It("health working as expected", func() {
		w := httptest.NewRecorder()
		e := NewEchoContext(w)

		err := api.HealthHandler(e)

		Expect(err).ToNot(HaveOccurred())
		Expect(w.Code).To(BeEquivalentTo(http.StatusOK))
		Expect(w.Body.String()).To(BeEquivalentTo("Health check working"))

	})

})

func NewEchoContext(w http.ResponseWriter) echo.Context {
	e := echo.New()
	r, err := http.NewRequest("GET", "http://www.fakeurl.com", nil)
	Expect(err).ToNot(HaveOccurred())

	return e.NewContext(r, w)
}
