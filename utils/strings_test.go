package utils

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUnitGenerateID(t *testing.T) {
	Convey("ID is correct length", t, func() {
		ref := GenerateID()
		So(len(ref), ShouldEqual, 15)
	})

	Convey("Reference Number does not collide", t, func() {
		// generate 10,000 reference numbers and check for any duplicates
		times := 10000 // 10 thousand
		generated := make([]string, times)

		for i := 0; i < times; i++ {
			ref := GenerateID()
			generated[i] = ref
		}

		// check for dups by creating a map of string->int and counting the the entry values whilst
		// iterating through the generated map
		generatedCheck := make(map[string]int)
		for _, reference := range generated {
			_, exists := generatedCheck[reference]
			So(exists, ShouldBeZeroValue)
			generatedCheck[reference] = 1
		}
	})
}

func TestUnitGenerateEtag(t *testing.T) {
	Convey("Generate Etag", t, func() {
		etag, err := GenerateEtag()
		So(len(etag), ShouldEqual, 56)
		So(err, ShouldBeNil)
	})
}
