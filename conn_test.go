package go_neu_ipgw

import (
	"log"
	"testing"
)

// TODO: Replace this with your own username and password
var testUser = IPGWUser{"20188888", "888888"}

func TestConnectPC(t *testing.T) {
	err := Connect(&testUser, TargetURLPC, UserAgentPCDefault)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func TestConnectMobile(t *testing.T) {
	err := Connect(&testUser, TargetURLMobile, UserAgentMobileDefault)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func TestDisconnect(t *testing.T) {
	err := Disconnect(&testUser)
	if err != nil {
		log.Fatalln(err.Error())
	}
}
