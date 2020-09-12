package util

import (
	"sync"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestWait(t *testing.T) {

	Convey("Given wait group with positive counter", t, func() {
		wg := sync.WaitGroup{}
		wg.Add(1)

		Convey("When wg group is done channel is trigger", func() {
			go func() {
				time.Sleep(1 * time.Second)
				defer wg.Done()
			}()
			err := Wait(&wg, 5 * time.Second)
			Convey("error should be return nil", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("When wg group is run out of time", func() {
			go func() {
				time.Sleep(2 * time.Second)
				defer wg.Done()
			}()
			err := Wait(&wg, 1 * time.Second)
			Convey("error should be existed and message is equals timeout", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "timeout")
			})
		})
	})
}
