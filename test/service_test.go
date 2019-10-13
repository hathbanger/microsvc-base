package service_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/go-kit/kit/log"
	fakes "github.com/hathbanger/microsvc-base/test/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/hathbanger/microsvc-base/pkg/microsvc"
)

var _ = Describe("Service", func() {
	var (
		request  *http.Request
		recorder *httptest.ResponseRecorder
		service  = &fakes.FakeService{}
		server   = microsvc.MakeRoutes(service, log.NewNopLogger(), nil)
	)

	BeforeEach(func() {
		recorder = httptest.NewRecorder()
	})

	AfterEach(func() {})

	// test.txt
})
