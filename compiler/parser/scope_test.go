package parser

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_scope(t *testing.T) {
	Convey("Given that global scope exists", t, func() {
		InitializeScope()
		var c1 = NewScope(Method)
		var c1c1 = NewScope(Statement)
		var c2 = NewScope(Statement)

		c1.PushChild(c1c1)
		GlobalScope.PushChild(c1, c2)

		Convey("Children of children return the correct parent", func() {
			So(c1c1.getParentScope(), ShouldEqual, c1)
		})
		Convey("Children of the global scope should return the global scope", func() {
			So(c2.getParentScope(), ShouldEqual, GlobalScope)
		})
		Convey("Once getParentScope() is ren on the global scope, it should return nil", func() {
			So(GlobalScope.getParentScope(), ShouldBeNil)
		})
	})
}
