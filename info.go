package go_neu_ipgw

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type IPGWInfo struct {
	BytesUsed     uint64
	SecondsOnline uint64
	BalanceLeft   float32
	UserIP        string
}

func randKey() string {
	return strconv.Itoa(rand.Intn(100001))
}

func GetInfo() (*IPGWInfo, error) {
	key := randKey()
	addr := fmt.Sprintf("http://ipgw.neu.edu.cn/include/auth_action.php?k=%s", key)
	resp, err := http.PostForm(addr, url.Values{
		"action": {"get_online_info"},
		"key":    {key},
	})

	if err != nil {
		return nil, err
	}

	body, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(body)
	defer resp.Body.Close()

	infoParts := strings.SplitN(bodyString, ",", 6)
	if len(infoParts) != 6 {
		return nil, errors.New(bodyString)
	}

	if balance, err := strconv.ParseFloat(infoParts[2], 32); err != nil {
		return nil, err
	} else if bytesUsed, err := strconv.ParseUint(infoParts[0], 10, 64); err != nil {
		return nil, err
	} else if secondsOnline, err := strconv.ParseUint(infoParts[1], 10, 64); err != nil {
		return nil, err
	} else {
		return &IPGWInfo{
			bytesUsed,
			secondsOnline,
			float32(balance),
			infoParts[5],
		}, nil
	}
}
