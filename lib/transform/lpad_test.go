package transform

import (
	"testing"

	"image"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func Test(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "LPad Suite")
}

var _ = Describe("findSp", func() {
	var rect image.Rectangle
	var subject image.Point

	var newWidth, newHeight int

	BeforeEach(func() {
		// actual thumbnail is 200x100 px
		rect = image.Rect(0, 0, 200, 100)
	})

	JustBeforeEach(func() {
		// new frame is newWidth x newHeight
		lpad := LPad{Height: newHeight, Width: newWidth}
		subject = lpad.findSp(rect)
	})

	Context("When padding everywhere", func() {
		BeforeEach(func() {
			newHeight = 200
			newWidth = 300
		})

		It("Is diagonally shifted", func() {
			Expect(subject).To(Equal(image.Pt(-50, -50)))
		})
	})

	Context("When vertical padding", func() {
		BeforeEach(func() {
			newHeight = 200
			newWidth = 200
		})

		It("Is vertically shifted", func() {
			Expect(subject).To(Equal(image.Pt(0, -50)))
		})
	})

	Context("Horizontal padding", func() {
		BeforeEach(func() {
			newHeight = 100
			newWidth = 300
		})

		It("Is horizontally shifted", func() {
			Expect(subject).To(Equal(image.Pt(-50, 0)))
		})
	})

	Context("No padding", func() {
		BeforeEach(func() {
			newHeight = 100
			newWidth = 50
		})

		It("Is horizontally shifted", func() {
			Expect(subject).To(Equal(image.Pt(0, 0)))
		})
	})
})

var _ = Describe("isScaledDownsize", func() {
	var subject bool
	var newWidth, newHeight int

	JustBeforeEach(func() {
		lpad := LPad{Height: newHeight, Width: newWidth}
		subject = lpad.isScaledDownsize(250, 167)
	})

	Context("When same aspect ratio", func() {
		BeforeEach(func() {
			newWidth = 200
			newHeight = 134
		})

		It("Is truthy", func() {
			Expect(subject).To(BeTrue())
		})
	})

	Context("When different aspect ratio", func() {
		BeforeEach(func() {
			newWidth = 200
			newHeight = 200
		})

		It("Is falsey", func() {
			Expect(subject).To(BeFalse())
		})
	})
})
