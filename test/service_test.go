package service_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/go-kit/kit/log"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/hathbanger/microsvc-base/pkg/microsvc"
	fakes "github.com/hathbanger/microsvc-base/test/microsvcfakes"
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

	Describe("POST /api/v1/foo", func() {
		Context("with an invalid request/json", func() {
			BeforeEach(func() {
				payload, _ := json.Marshal(
					`{"invalid": "json"}`,
				)
				request, _ = http.NewRequest(
					http.MethodPost,
					"/api/v1/foo",
					bytes.NewReader(payload),
				)
			})

			It("returns http status code 400", func() {
				server.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(400))
			})
		})
	}) // end of /api/v1/foo
	//
	// replace me
})
