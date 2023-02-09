package main

import (
	"fmt"
	// "google.golang.org/protobuf/types/known/timestamppb"
)

func main() {
	fmt.Println("test")
	test := Person{
		Id:    1234,
		Name:  "TEST",
		Email: "test@test",
		// Time: Times{
		// 	CreatedAt:,
		// },
	}
	fmt.Println(test)
}
