package gconfig

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

//GlobalConf ...
type GlobalConf struct {
	//SVNPATH ...
	SVNPATH string `json:"svnroot"`
	//Username ...
	Username string `json:"username"`
	//Password ...
	Password string `json:"password"`
	//Salt ...
	Salt string `json:"salt"`
}

//NewGconfig init a configuration (read json from svnadmin.json)
func NewGconfig() *GlobalConf {
	file, err := os.Open("../svnadmin.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	//content is a json file
	conf := new(GlobalConf)
	scontent := string(content)
	err = json.Unmarshal([]byte(scontent), &conf)
	if err != nil {
		panic(err)
	}
	return conf
}
