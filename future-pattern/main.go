package main

import (
	"context"
	"fmt"

	"github.com/pvoliveira/go-study/future-pattern/future"
)

func main() {
	ctx := context.Background()
	f := future.SlowFunction(ctx)

	res, err := f.Result()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(res)
}
