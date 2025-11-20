package slices

import (
	"fmt"
	"time"
)

func Start() {

	slice := []string{"sad", "asd", "dsad"}
	sliceP := &slice
	(*sliceP)[0] = "SUKA"
	fmt.Println(slice[0])
	takeValueSliceAndSetValueByIndex(slice, 0, "PIDORAS")
	fmt.Println(slice[0])
	takePointerSliceAndSetValueByIndex(&slice, 0, "BLYAT")
	fmt.Println(slice[0])

	time.Sleep(time.Second * 5)
}

func takePointerSliceAndSetValueByIndex(slice *[]string, index int, value string) {
	(*slice)[index] = value
	// next line will not compile, because of "pointer type doesn't support indexing"
	// slice[index] = value
}

func takeValueSliceAndSetValueByIndex(slice []string, index int, value string) {
	slice[index] = value
}
