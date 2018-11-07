package chain

import (
	"errors"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestC_Handles(t *testing.T) {
	convey.Convey("Chain", t, func() {
		var (
			val     = 0
			errTest = errors.New("test")
		)

		f1 := func() error { val++; return nil }
		f2 := func(i int) { val += i }
		f3 := func() error { return errTest }
		f4 := func() error { val++; return nil }

		convey.Convey("Error", func() {
			err := New().Handles(f1, func() error { f2(2); return nil }, f3, f4).Run()
			convey.So(err, convey.ShouldBeError, errTest)
			convey.So(val, convey.ShouldEqual, 3)
		})

		convey.Convey("No Error", func() {
			err := New().Handles(f1, f4).Run()
			convey.So(err, convey.ShouldBeNil)
			convey.So(val, convey.ShouldEqual, 2)
		})
	})
}

