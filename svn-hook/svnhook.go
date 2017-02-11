package svnhook

import (
	"io/ioutil"
	"os/exec"
	"strings"
)

//Hook ,,,
type Hook struct {
	svnpath   string
	repo      string
	repospath string
}

//NewHook ...
func NewHook(svnpath string, repo string) *Hook {
	hook := new(Hook)
	hook.svnpath = svnpath
	hook.repo = repo
	hook.repospath = hook.svnpath + "/" + hook.repo
	return hook
}

//GetAuthor hook the svn ,and get its repo author .
func (hook *Hook) GetAuthor() (string, error) {
	cmd := exec.Command("svnlook", "author", hook.repospath)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}
	if err = cmd.Start(); err != nil {
		return "", err
	}
	author, err := ioutil.ReadAll(stdout)
	if err != nil {
		return "", err
	}
	trimauthor := strings.TrimSuffix(string(author), "\n")
	return string(trimauthor), nil
}

//GetChanged  : 打印修改的路径
/**
* 'A':条目添加到版本库
* 'D':条目从版本库删除
* ‘U’:条目内容改变了
* ' U':条目属性改变了
* 'UU':文件内容和属性改变了
**/
func (hook *Hook) GetChanged() (string, error) {
	cmd := exec.Command("svnlook", "changed", hook.repospath)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}
	if err = cmd.Start(); err != nil {
		return "", err
	}
	changed, err := ioutil.ReadAll(stdout)
	if err != nil {
		return "", err
	}
	return string(changed), nil
}

//Getdate get time stamp
func (hook *Hook) Getdate() (string, error) {
	cmd := exec.Command("svnlook", "date", hook.repospath)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}
	if err = cmd.Start(); err != nil {
		return "", err
	}
	date, err := ioutil.ReadAll(stdout)
	if err != nil {
		return "", err
	}
	return string(date), nil
}

//Getdiff print file diff in GNU format
func (hook *Hook) Getdiff() (string, error) {
	cmd := exec.Command("svnlook", "diff", hook.repospath)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}
	if err = cmd.Start(); err != nil {
		return "", err
	}
	diff, err := ioutil.ReadAll(stdout)
	if err != nil {
		return "", err
	}
	return string(diff), nil
}

//Getlog print file log
func (hook *Hook) Getlog() (string, error) {
	cmd := exec.Command("svnlook", "log", hook.repospath)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}
	if err = cmd.Start(); err != nil {
		return "", err
	}
	log, err := ioutil.ReadAll(stdout)
	if err != nil {
		return "", err
	}
	return string(log), nil
}

//Getinfo return a specific map
func (hook *Hook) Getinfo() (map[string]string, error) {
	cmd := exec.Command("svnlook", "info", hook.repospath)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err = cmd.Start(); err != nil {
		return nil, err
	}
	info, err := ioutil.ReadAll(stdout)
	if err != nil {
		return nil, err
	}
	splitinfo := strings.SplitN(string(info), "\n", -1)
	infomap := make(map[string]string)
	infomap["author"] = splitinfo[0]
	infomap["timestamp"] = splitinfo[1]
	infomap["logsize"] = splitinfo[2]
	infomap["log"] = splitinfo[3]
	return infomap, nil
}
