package main

import (
	"fmt"
	"log"
	"yaoj-go/judger"
)

func TestTradJudger() {
	log.SetPrefix("\033[35mDEBUG \033[0m")

	var judger judger.Judger = &judger.TradJudger{
		Executable: "test",
		InpFile:    "test.in",
		OutFile:    "test.out",
		ErrFile:    "test.err",
		Fileio:     false,
		Policy:     "builtin:cstdio",
	}
	result, err := judger.Perform()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("result: %v\n", result)
}
