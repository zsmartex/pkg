package main

import "github.com/zsmartex/pkg/services"

func main() {
	logger := services.NewLoggerService("Finex")

	logger.Info("This is an info message")
}
