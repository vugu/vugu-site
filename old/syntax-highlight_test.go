package main

import (
	"log"
	"testing"
)

func TestShowCode(t *testing.T) {

	log.Print(showCode("go", "var a = b\n"))

}
