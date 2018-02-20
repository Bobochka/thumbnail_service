package service

import (
	"errors"
	"testing"

	"github.com/Bobochka/thumbnail_service/lib"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func Test(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Service Suite")
}

const fprint = "transformation_id"

var ErrOups = errors.New("oups")

var _ = Describe("Service", func() {
	var subject *Service
	var url string
	var t *MockTransformation
	var store *MockStore
	var downloader *MockDownloader
	var locker *MockLocker

	var storeGetCalls []*gomock.Call
	var lockerNewMutexCalls []*gomock.Call
	var mtxLockCalls []*gomock.Call

	var mockCtrl *gomock.Controller

	AfterEach(func() {
		mockCtrl.Finish()
	})

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		store = NewMockStore(mockCtrl)
		t = NewMockTransformation(mockCtrl)
		downloader = NewMockDownloader(mockCtrl)
		locker = NewMockLocker(mockCtrl)
	})

	JustBeforeEach(func() {
		gomock.InOrder(storeGetCalls...)
		gomock.InOrder(lockerNewMutexCalls...)
		gomock.InOrder(mtxLockCalls...)

		subject = New(&Config{
			Store:      store,
			Downloader: downloader,
			Locker:     locker,
		})
	})

	AfterEach(func() {
		storeGetCalls = []*gomock.Call{}
		lockerNewMutexCalls = []*gomock.Call{}
		mtxLockCalls = []*gomock.Call{}
	})

	Describe("Perform", func() {
		var result []byte
		var err error
		var data []byte
		var resData []byte
		var mtx *lib.MockMutex

		BeforeEach(func() {
			mtx = lib.NewMockMutex(mockCtrl)

			t.EXPECT().Fingerprint(gomock.Any()).Return(fprint).AnyTimes()

			data = []byte("image of flower")
			resData = []byte("thumbed image of flower")
		})

		JustBeforeEach(func() {
			result, err = subject.Perform(url, t)
		})

		// shared examples
		ItBehavesAsPerformed := func() {
			It("Returns transformed data", func() {
				Expect(result).To(Equal(resData))
			})

			It("Does not return error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		}

		ItBehavesAsNotPerformed := func() {
			It("Does not return transformed data", func() {
				Expect(result).To(BeEmpty())
			})

			It("Returns error", func() {
				Expect(err).To(MatchError(ErrOups.Error()))
			})
		}

		Context("When downloader can't download from url", func() {
			BeforeEach(func() {
				downloader.EXPECT().Download(gomock.Any()).Return([]byte{}, ErrOups)
			})

			ItBehavesAsNotPerformed()
		})

		Context("When url is downloadable", func() {
			BeforeEach(func() {
				downloader.EXPECT().Download(gomock.Any()).Return(data, nil)
			})

			Context("When data already in store", func() {
				BeforeEach(func() {
					storeGetCalls = append(storeGetCalls, store.EXPECT().Get(fprint).Return(resData))
				})

				ItBehavesAsPerformed()
			})

			Context("When data not in store", func() {
				BeforeEach(func() {
					storeGetCalls = append(storeGetCalls, store.EXPECT().Get(fprint).Return([]byte{}))
					lockerNewMutexCalls = append(lockerNewMutexCalls, locker.EXPECT().NewMutex(fprint).Return(mtx))
				})

				WhenMutexAquired := func() {
					Context("When mutex acquired", func() {
						BeforeEach(func() {
							mtxLockCalls = append(mtxLockCalls, mtx.EXPECT().Lock().Return(nil))
							mtx.EXPECT().Unlock().Return(true)
						})

						Context("When data already in store", func() {
							BeforeEach(func() {
								storeGetCalls = append(storeGetCalls, store.EXPECT().Get(fprint).Return(resData))
							})

							ItBehavesAsPerformed()
						})

						Context("When data still not in store", func() {
							BeforeEach(func() {
								storeGetCalls = append(storeGetCalls, store.EXPECT().Get(fprint).Return([]byte{}))
							})

							Context("When can't perform transformation", func() {
								BeforeEach(func() {
									t.EXPECT().Perform(data).Return([]byte{}, ErrOups)
								})

								ItBehavesAsNotPerformed()
							})

							Context("When transformation is performed", func() {
								BeforeEach(func() {
									t.EXPECT().Perform(data).Return(resData, nil)
								})

								Context("When transformed value is stored successfully", func() {
									BeforeEach(func() {
										storeGetCalls = append(storeGetCalls, store.EXPECT().Set(fprint, resData).Return(nil))
									})

									ItBehavesAsPerformed()
								})

								Context("When transformed value is not stored", func() {
									BeforeEach(func() {
										storeGetCalls = append(storeGetCalls, store.EXPECT().Set(fprint, resData).Return(ErrOups))
									})

									ItBehavesAsPerformed()
								})
							})

							Context("When transformation is not performed", func() {
								BeforeEach(func() {
									t.EXPECT().Perform(data).Return([]byte{}, ErrOups)
								})

								ItBehavesAsNotPerformed()
							})
						})
					})
				}

				WhenMutexAquired()

				Context("When mutex not acquired", func() {
					BeforeEach(func() {
						mtxLockCalls = append(mtxLockCalls, mtx.EXPECT().Lock().Return(ErrOups))
					})

					Describe("Store Polling", func() {
						Context("When data is immediately in store", func() {
							BeforeEach(func() {
								storeGetCalls = append(storeGetCalls, store.EXPECT().Get(fprint).Return(resData))
							})

							ItBehavesAsPerformed()
						})

						Context("When data is in store after another poll", func() {
							BeforeEach(func() {
								storeGetCalls = append(storeGetCalls, store.EXPECT().Get(fprint).Return([]byte{}))
								storeGetCalls = append(storeGetCalls, store.EXPECT().Get(fprint).Return(resData))
							})

							ItBehavesAsPerformed()
						})

						Context("When data is not in store after all polls", func() {
							BeforeEach(func() {
								storeGetCalls = append(storeGetCalls, store.EXPECT().Get(fprint).Return([]byte{}))
								storeGetCalls = append(storeGetCalls, store.EXPECT().Get(fprint).Return([]byte{}))
								storeGetCalls = append(storeGetCalls, store.EXPECT().Get(fprint).Return([]byte{}))

								lockerNewMutexCalls = append(lockerNewMutexCalls, locker.EXPECT().NewMutex(fprint).Return(mtx))
							})

							WhenMutexAquired()

							Context("When mutex not acquired", func() {
								BeforeEach(func() {
									storeGetCalls = append(storeGetCalls, store.EXPECT().Get(fprint).Return([]byte{}))
									storeGetCalls = append(storeGetCalls, store.EXPECT().Get(fprint).Return([]byte{}))
									storeGetCalls = append(storeGetCalls, store.EXPECT().Get(fprint).Return([]byte{}))

									mtxLockCalls = append(mtxLockCalls, mtx.EXPECT().Lock().Return(ErrOups))

									t.EXPECT().Perform(data).Return(resData, nil)

									storeGetCalls = append(storeGetCalls, store.EXPECT().Set(fprint, resData).Return(nil))
								})

								ItBehavesAsPerformed()
							})
						})
					})
				})
			})
		})
	})
})
