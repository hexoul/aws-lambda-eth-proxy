package predefined

import (
	"fmt"
)

func foo() {
	fmt.Println("foo")
}

var predefinedPaths = map[string]interface{}{
	"foo": foo,
}
