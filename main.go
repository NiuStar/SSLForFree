package main

import (

	"nqc.cn/utils"

	"encoding/json"
	"fmt"
	"strings"
	"time"
)


var config map[string]interface{}

func main() {

	data := utils.ReadFile("config.json")
	err := json.Unmarshal([]byte(data), &config)
	if err != nil {
		panic(err)
	}

	data1 := utils.PostData("http://125.70.230.211:10010/GetCert","config=" + data)

	var result map[string]interface{}
	err = json.Unmarshal([]byte(data1), &result)
	if err != nil {
		panic(err)
	}

	fmt.Println("data1:",data1)
	fmt.Println("result:",result)
	list := result["list"].([]interface{})
	for index,value := range list {

		fmt.Println("list : ", strings.Split(value.(string),"?verification="))

		fmt.Println("uri :",config["checkPath"].(string) + "/" + strings.Split(value.(string),"?verification=")[1])

		data2 := GETS(value.(string),result["PHPSESSID"].(string))

		fmt.Println(index,"  : 验证文件： ",value.(string) + " : " ,data2)

		utils.WriteToFile(config["checkPath"].(string) + "/" + strings.Split(value.(string),"?verification=")[1],data2)


	}

	var times time.Duration = (time.Duration)(config["SleepTime"].(float64))
	time.Sleep(time.Second * times)

	data1 = utils.PostData("http://125.70.230.211:10010/CheckDomain","config=" + data)
	fmt.Println("data2:",data1)
	err = json.Unmarshal([]byte(data1), &result)
	if err != nil {
		panic(err)
	}
	Cert := result["Cert"].(map[string]interface{})

	crt := Cert["certificate"].(string) + "\n\n" + Cert["certificate_bundle"].(string)

	utils.WriteToFile(config["savePath"].(string) + "/server.crt",crt)
	utils.WriteToFile(config["savePath"].(string) + "/server.key",Cert["certificate_private"].(string))
}
