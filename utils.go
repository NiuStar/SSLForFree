package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"crypto/tls"
	"strings"
	"net/url"
)

func GETS(uri string,PHPSESSID string) string {
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	//向服务端发送get请求
	request, err := http.NewRequest("GET", uri, nil)

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Cookie", "PHPSESSID=" + PHPSESSID)
	response, err := client.Do(request)

	if err != nil {
		panic(err)
	}

	if response.StatusCode == 200 {

		body, _ := ioutil.ReadAll(response.Body)
		response.Body.Close()
		return string(body)

	}
	response.Body.Close()
	return GETS(uri,PHPSESSID)
}


func POSTS(uri string, datas map[string]string,PHPSESSID string) string {
	postValues := url.Values{}
	for key, value := range datas {
		postValues.Add(key, value)
	}
	body := ioutil.NopCloser(strings.NewReader(postValues.Encode())) //把form数据编下码
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	request, err := http.NewRequest("POST", uri, body)

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Cookie", "PHPSESSID=" + PHPSESSID)
	resp, err := client.Do(request)

	if err != nil {
		// handle error
		panic(err)
		return "{}"
	}

	if resp.StatusCode != 200 {
		fmt.Println("重新开始请求:",uri)
		return POSTS(uri,datas,PHPSESSID)

	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
		// handle error
	}

	fmt.Println("resp. cookies: ",resp.Header)
	resp.Body.Close()
	return string(data)

}
