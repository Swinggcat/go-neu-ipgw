package go_neu_ipgw

import (
	"fmt"
	"log"
	"testing"
)

func TestGetInfo(t *testing.T) {
	info, err := GetInfo()
	if err != nil {
		log.Fatalln(err.Error())
	}

	fmt.Println(info)
}
