package main

import "5g-v2x-user-service/internal/container"

func main() {
	if err := container.NewContainer().Run().Error; err != nil {
		panic(err)
	}

}
