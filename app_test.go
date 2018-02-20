package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"fmt"

	"encoding/json"

	"gopkg.in/h2non/gock.v1"

	"io/ioutil"

	"github.com/Bobochka/thumbnail_service/lib/service"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

func Test(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Suite")
}

var _ = Describe("App", func() {
	var app *App

	BeforeSuite(func() {
		cfg, err := ReadConfig()

		Expect(err).NotTo(HaveOccurred())

		app = &App{service: service.New(cfg)}
	})

	Describe("/thumbnail", func() {
		DescribeTable("Invalid Params",
			func(url, width, height, desc string) {
				query := fmt.Sprintf("?url=%s&width=%s&height=%s", url, width, height)

				rr, err := Request(app, query)
				Expect(err).NotTo(HaveOccurred())

				resp := struct{ Error string }{}

				err = json.Unmarshal(rr.Body.Bytes(), &resp)
				Expect(err).NotTo(HaveOccurred())

				Expect(resp.Error).To(Equal(desc))
				Expect(rr.Code).To(Equal(400))
			},
			Entry("url invalid", "malformed.com", "", "", "url malformed.com is not valid"),
			Entry("width NaN", "http://google.com", "width", "", "width width is not valid: should be positive integer"),
			Entry("width < 0", "http://google.com", "-42", "", "width -42 is not valid: should be positive integer"),
			Entry("width = 0", "http://google.com", "0", "", "width 0 is not valid: should be positive integer"),
			Entry("height NaN", "http://google.com", "42", "height", "height height is not valid: should be positive integer"),
			Entry("height < 0", "http://google.com", "42", "-42", "height -42 is not valid: should be positive integer"),
			Entry("height = 0", "http://google.com", "42", "0", "height 0 is not valid: should be positive integer"),
			Entry("too big", "http://google.com", "42000", "42000", "requested size of 42000 x 42000 is too big"),
		)

		Context("Presumably Valid params", func() {
			var rr *httptest.ResponseRecorder

			BeforeEach(func() {
				gock.EnableNetworking() // in order to access s3
			})

			BeforeEach(func() {
				gock.Off()
			})

			JustBeforeEach(func() {
				query := "?url=http://foo.com/sample.jpg&width=200&height=200"

				var err error
				rr, err = Request(app, query)
				Expect(err).NotTo(HaveOccurred())
			})

			Context("When valid indeed params", func() {
				BeforeEach(func() {
					gock.New("http://foo.com").
						Get("/sample.jpg").
						Reply(200).
						File("./testdata/sample.jpg")
				})

				AfterEach(func() {
					Expect(gock.IsDone()).To(BeTrue())
				})

				It("Renders thumbnail", func() {
					data, err := ioutil.ReadFile("./testdata/result.jpg")

					Expect(err).NotTo(HaveOccurred())
					Expect(rr.Code).To(Equal(200))
					Expect(rr.Body.Bytes()).To(Equal(data))
				})
			})

			Context("When actually not an image", func() {
				BeforeEach(func() {
					gock.New("http://foo.com").
						Get("/sample.jpg").
						Reply(200).
						JSON("trap")
				})

				AfterEach(func() {
					Expect(gock.IsDone()).To(BeTrue())
				})

				It("Responds with error", func() {
					resp := struct{ Error string }{}

					err := json.Unmarshal(rr.Body.Bytes(), &resp)
					Expect(err).NotTo(HaveOccurred())

					Expect(resp.Error).To(Equal("Content type is not supported, supported formats: jpeg, gif, png"))
					Expect(rr.Code).To(Equal(400))
				})
			})
		})
	})
})

func Request(app *App, query string) (*httptest.ResponseRecorder, error) {
	req, err := http.NewRequest("GET", "/thumbnail"+query, nil)

	if err != nil {
		return nil, err
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.thumbnail)
	handler.ServeHTTP(rr, req)

	return rr, nil
}
