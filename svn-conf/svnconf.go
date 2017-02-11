package svnconf

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

//SVNconf ...
type SVNconf struct {
	//SVNPATH ...
	SVNPATH string `json:"svnpath"`
	anonAC  string
	authAC  string
	pdb     string
	adb     string
	gdb     string
}

//NewSVNconf ...
func NewSVNconf(path string) *SVNconf {
	conf := &SVNconf{}
	conf.SVNPATH = path
	return conf
}

//Readfill will read the authz to the buffer
func (conf *SVNconf) readfile() ([]byte, error) {
	confs := filepath.Join(conf.SVNPATH, "conf", "svnserve.conf")
	filehandle, err := os.Open(confs)
	defer filehandle.Close()
	if err != nil {
		return nil, errors.New(`Please check your ./conf/svnserve.conf file,if you change it's name,rename it as "svnserve.conf".`)
	}
	content, err := ioutil.ReadAll(filehandle)
	if err != nil {
		return nil, err
	}
	return content, nil
}

//remarktag把tag重新改到注释状态
func (conf *SVNconf) remarktag(srcContent []byte, tag string) error {
	confs := filepath.Join(conf.SVNPATH, "conf", "svnserve.conf")
	var lines []string
	var dstContent bytes.Buffer
	lines = strings.Split(string(srcContent), "\n")
	for k, v := range lines {
		if strings.HasPrefix(v, tag) {
			v = strings.Replace(lines[k], tag, "#"+tag, -1)
		}
		dstContent.WriteString(v + "\n")
	}
	//write to file
	err := ioutil.WriteFile(confs, dstContent.Bytes(), 0666)
	if err != nil {
		return err
	}
	return nil
}

//Rmremarktag 删除注释状态
func (conf *SVNconf) rmremarktag(srcContent []byte, tag string) error {
	confs := filepath.Join(conf.SVNPATH, "conf", "svnserve.conf")
	var lines []string
	var dstContent bytes.Buffer
	lines = strings.Split(string(srcContent), "\n")
	for k, v := range lines {
		if strings.HasPrefix(strings.TrimSpace(v), "#"+tag) {
			v = strings.Replace(lines[k], "#"+tag, tag, -1)
		}
		dstContent.WriteString(v + "\n")
	}
	//write to file
	err := ioutil.WriteFile(confs, dstContent.Bytes(), 0644)
	if err != nil {
		return err
	}
	return nil
}

//file2json parse authz content to readable json .
//[general]=>["anon-access = read","auth-access = write",...]
func file2json(filecontent []byte) ([]byte, error) {
	var lines []string
	lines = strings.Split(string(filecontent), "\n")
	jsonm := make(map[string][]string)
	for k, v := range lines {
		if strings.HasPrefix(v, "[") {
			key := strings.TrimPrefix(v, "[")
			key = strings.TrimSuffix(key, "]")
			var value []string
			for i := k + 1; i < len(lines); i++ {
				if strings.HasPrefix(lines[i], "[") {
					break
				}
				if lines[i] != "" && !strings.HasPrefix(lines[i], "#") {

					value = append(value, lines[i])

				}
			}
			jsonm[key] = value
		}
	}
	jsondata, err := json.Marshal(jsonm)
	if err != nil {
		return []byte(""), err
	}
	return jsondata, nil
}

func parseGeneral(rawjson []byte) ([]byte, error) {
	conf := make(map[string][]string)
	err := json.Unmarshal(rawjson, &conf)
	if err != nil {
		return []byte(""), err
	}
	generals := make(map[string]string)
	general := conf["general"]
	for _, v := range general {
		trimv := strings.TrimSpace(string(v))
		splitv := strings.SplitN(trimv, "=", 2)
		generals[strings.TrimSpace(splitv[0])] = strings.TrimSpace(splitv[1])
	}
	generalsjson, err := json.Marshal(generals)
	if err != nil {
		return []byte(""), err
	}
	return generalsjson, nil
}

//ExportGeneral export the content .
func (conf *SVNconf) ExportGeneral() ([]byte, error) {
	content, err := conf.readfile()
	if err != nil {
		return []byte(""), err
	}
	rawjson, err := file2json(content)
	if err != nil {
		return []byte(""), err
	}
	general, err := parseGeneral(rawjson)
	if err != nil {
		return []byte(""), err
	}
	return general, nil
}

//ExportRemarkGeneral export the content after edited ..
func (conf *SVNconf) ExportRemarkGeneral(tag string) ([]byte, error) {
	content, err := conf.readfile()
	if err != nil {
		return []byte(""), err
	}
	err = conf.remarktag(content, tag)
	if err != nil {
		return []byte(""), err
	}
	content, err = conf.readfile()
	if err != nil {
		return []byte(""), err
	}
	rawjson, err := file2json(content)
	if err != nil {
		return []byte(""), err
	}
	general, err := parseGeneral(rawjson)
	if err != nil {
		return []byte(""), err
	}
	return general, nil
}

//ExportRmremarkGeneral export the content after edited ..
func (conf *SVNconf) ExportRmremarkGeneral(tag string) ([]byte, error) {
	content, err := conf.readfile()
	if err != nil {
		return []byte(""), err
	}
	err = conf.rmremarktag(content, tag)
	if err != nil {
		return []byte(""), err
	}
	content, err = conf.readfile()
	if err != nil {
		return []byte(""), err
	}
	rawjson, err := file2json(content)
	if err != nil {
		return []byte(""), err
	}
	general, err := parseGeneral(rawjson)
	if err != nil {
		return []byte(""), err
	}
	return general, nil
}
