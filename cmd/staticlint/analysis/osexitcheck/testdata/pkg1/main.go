package main

import (
	"os"
)

func main() {

}

func osExitCheckFunc() {
	// формулируем ожидания: анализатор должен находить ошибку,
	// описанную в комментарии want
	os.Exit(1) // want "os.Exit in main file"

}
