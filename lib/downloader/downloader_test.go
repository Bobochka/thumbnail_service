package downloader

import (
	"testing"

	"github.com/Bobochka/thumbnail_service/lib"
	"github.com/go-errors/errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/h2non/gock.v1"
)

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}

func Test(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Downloader Suite")
}

var _ = Describe("Downloader", func() {
	AfterEach(func() {
		Expect(gock.IsDone()).To(BeTrue())
		gock.Off()
	})

	var subject *Http
	var allowedTypes []string

	JustBeforeEach(func() {
		subject = New(allowedTypes)
	})

	Describe("Download", func() {
		var host, path string
		var data []byte
		var err error

		BeforeEach(func() {
			host = "http://foo.bar"
			path = "/baz"
		})

		JustBeforeEach(func() {
			data, err = subject.Download(host + path)
		})

		Context("When request failure", func() {
			BeforeEach(func() {
				gock.New(host).
					Get(path).
					ReplyError(errors.New("err"))
			})

			It("Responds without data", func() {
				Expect(data).To(BeEmpty())
			})

			It("Responds with error code 404", func() {
				typedErr, ok := err.(lib.Error)
				Expect(ok).To(BeTrue())
				Expect(typedErr.Code()).To(Equal(404))
			})
		})

		Context("When non 2xx status code", func() {
			BeforeEach(func() {
				gock.New(host).
					Get(path).
					Reply(500)
			})

			It("Responds without data", func() {
				Expect(data).To(BeEmpty())
			})

			It("Responds with error code 404", func() {
				typedErr, ok := err.(lib.Error)
				Expect(ok).To(BeTrue())
				Expect(typedErr.Code()).To(Equal(404))
			})
		})

		Context("When can't read body", func() {
			BeforeEach(func() {
				gock.New(host).
					Get(path).
					Body(errReader(0))
			})

			It("Responds without data", func() {
				Expect(data).To(BeEmpty())
			})

			It("Responds with error code 404", func() {
				typedErr, ok := err.(lib.Error)
				Expect(ok).To(BeTrue())
				Expect(typedErr.Code()).To(Equal(404))
			})
		})

		Context("When unsupported content type", func() {
			BeforeEach(func() {
				allowedTypes = []string{"image/png"}
			})

			BeforeEach(func() {
				gock.New(host).
					Get(path).
					Reply(200).
					BodyString("something")
			})

			It("Responds without data", func() {
				Expect(data).To(BeEmpty())
			})

			It("Responds with error code 400", func() {
				typedErr, ok := err.(lib.Error)
				Expect(ok).To(BeTrue())
				Expect(typedErr.Code()).To(Equal(400))
			})
		})

		Context("When supported content type", func() {
			BeforeEach(func() {
				allowedTypes = []string{"text/plain; charset=utf-8"}
			})

			BeforeEach(func() {
				gock.New(host).
					Get(path).
					Reply(200).
					BodyString("something")
			})

			It("Responds with data", func() {
				Expect(data).To(Equal([]byte(`something`)))
			})

			It("Responds without error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

})
