//** This file was code generated by got. ***

package testOut

import (
	"bytes"
	"fmt"
	"strconv"
)

func givesString() string {
	return "Me"
}

func givesInt() int {
	return -5
}

func givesUint() uint {
	return 6
}

func TestVars(buf *bytes.Buffer) {

	buf.WriteString("Evaluates to a string")

	buf.WriteString(`Here is a number: `)
	buf.WriteString(strconv.Itoa(givesInt()))
	buf.WriteString(`
And another: `)
	buf.WriteString(strconv.FormatUint(uint64(givesUint()), 10))
	buf.WriteString(`
And a float: `)
	buf.WriteString(strconv.FormatFloat(float64(45/6), 'g', -1, 64))
	buf.WriteString(`
Stringer: `)
	buf.WriteString(fmt.Sprintf("%v", "something"))
	buf.WriteString(`
`)

}
