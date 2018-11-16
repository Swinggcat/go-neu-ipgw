package go_neu_ipgw

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

const TargetURLPC = "http://ipgw.neu.edu.cn/srun_portal_pc.php?ac_id=1&"
const TargetURLMobile = "http://ipgw.neu.edu.cn/srun_portal_phone.php?ac_id=1&"
const UserAgentPCDefault = "Mozilla/5.0 (X11; Linux x86_64; rv:63.0) Gecko/20100101 Firefox/63.0"
const UserAgentMobileDefault = "Mozilla/5.0 (iPhone; CPU iPhone OS 10_3_1 like Mac OS X) AppleWebKit/603.1.30 (KHTML, like Gecko) Version/10.0 Mobile/14E304 Safari/602.1"

type IPGWUser struct {
	Username string
	Password string
}

var connStatusRegexPC, _ = regexp.Compile(`<input type="hidden" name="url" value="" >[\w\W]+?<p>(.+?)</p>`)
var connStatusRegexMobile, _ = regexp.Compile(`class="weui_toptips weui_warn js_tooltips">(.+?)</div>`)
var connStatusRegexPCSpecial, _ = regexp.Compile(`style="font-weight:bold;color:orange;">[\n\r ]+(.+?) `)

func getConnStatusText(response []byte) (string, error) {
	respString := string(response)
	if matches := connStatusRegexPCSpecial.FindStringSubmatch(respString); len(matches) == 2 {
		return matches[1], nil
	} else if matches := connStatusRegexPC.FindStringSubmatch(respString); len(matches) == 2 {
		return matches[1], nil
	} else if matches := connStatusRegexMobile.FindStringSubmatch(respString); len(matches) == 2 {
		return matches[1], nil
	} else if strings.Contains(respString, "注册、自服务以及忘记密码") {
		return "网络已连接", nil
	} else {
		fmt.Println(respString)
		return "", errors.New("服务器应答无效")
	}
}

func Connect(user *IPGWUser, targetURL string, userAgent string) error {
	body := []byte(url.Values{
		"ac_id":    {"1"},
		"action":   {"login"},
		"username": {user.Username},
		"password": {user.Password},
		"save_me":  {"0"}}.Encode())

	req, _ := http.NewRequest(http.MethodPost, targetURL, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", userAgent)

	if response, e := http.DefaultClient.Do(req); e == nil {
		respBody, _ := ioutil.ReadAll(response.Body)
		defer response.Body.Close()

		if statusText, e := getConnStatusText(respBody); e == nil {
			if statusText == "网络已连接" {
				return nil
			} else {
				return errors.New(statusText)
			}
		} else {
			return e
		}
	} else {
		return e
	}
}

func Disconnect(user *IPGWUser) error {
	resp, err := http.PostForm("http://ipgw.neu.edu.cn/include/auth_action.php", url.Values{
		"action":   {"logout"},
		"ajax":     {"1"},
		"username": {user.Username},
		"password": {user.Password},
	})
	if err != nil {
		return err
	}

	body, _ := ioutil.ReadAll(resp.Body)
	bodyStr := string(body)

	if bodyStr == "网络已断开" {
		return nil
	} else {
		return errors.New(bodyStr)
	}
}
