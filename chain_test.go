package chain

import (
	"errors"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestC_Handles(t *testing.T) {
	Convey("Chain", t, func() {
		var (
			val     = 0
			errTest = errors.New("test")
		)

		last := func(ctx *Context, err error) {
			val = 0
			t.Log("Last error: ", err)
		}

		f1 := func(ctx *Context) error { val++; return nil }
		f2 := func(i int) { val += i }
		f3 := func(ctx *Context) error { return errTest }
		f4 := func(ctx *Context) error { val++; return nil }

		Convey("Error", func() {
			err := New().Handles(f1, func(ctx *Context) error { f2(2); return nil }, f3, f4).Run()
			So(err, ShouldBeError, errTest)
			So(val, ShouldEqual, 3)
		})

		Convey("No Error", func() {
			err := New().Handles(f1, f4).Run()
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 2)
		})

		Convey("Add Last", func() {
			r := New().Last(last)
			err := r.Handles(
				f1, func(ctx *Context) error {
					f2(10)
					return nil
				},
				f3).Run()
			So(val, ShouldEqual, 0)
			So(err, ShouldEqual, errTest)
		})

		Convey("Run", func() {
			err := Run(f1, f4)
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 2)
		})

		Convey("Each", func() {
			New().EachBefore(func(ctx *Context) {
				t.Log("each before")
			}).EachAfter(func(ctx *Context, err error) {
				t.Logf("each after %v", err)
			}).Handles(f1, f3).Run()
		})
	})
}

func TestWithContext(t *testing.T) {
	Convey("with context", t, func() {

		Convey("1", func() {
			Run(
				func(ctx *Context) error {
					ctx.Set("k", 1)
					return nil
				},
				func(ctx *Context) error {
					t.Logf("find k %d", ctx.MustGet("k").(int))
					ctx.Set("k", 2)
					return nil
				},
				func(ctx *Context) error {
					t.Logf("find k %d", ctx.MustGet("k").(int))
					return nil
				},
			)
		})
	})
}
