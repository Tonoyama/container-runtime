package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
)

var repo string = "nginx"
var tag string = "latest"

func mainImage(){
	getToken(repo)
	getManifest(repo, tag, "access-token")
}

// docker hubからトークンを入手する関数
func getToken(repo string) {
	url := "https://auth.docker.io/token?service=registry.docker.io&scope=repository:library/" + repo +":pull"

	resp,_ := http.Get(url)
	defer resp.Body.Close()

	byteArray,_ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(byteArray))
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Autorization", "Bearer access-token")

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
        panic(err)
    }
}

// マニフェストを取得する関数
func getManifest(repo, tag, token string) {
	url := "https://registry.hub.docker.com/v2/" + repo + "/manifests/" + tag
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer " + token)

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}