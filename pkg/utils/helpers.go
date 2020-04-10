package utils

import "fmt"

func logg(x interface{}) {
	fmt.Printf("%+v\n", x)
}
