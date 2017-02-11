package svnmanager

import (
	"encoding/json"
	"io/ioutil"
	"os/exec"
	"strings"

	"github.com/scbizu/svnpanel/svn-auth"
	"github.com/scbizu/svnpanel/svn-conf"
	"github.com/scbizu/svnpanel/svn-hook"
	"github.com/scbizu/svnpanel/svn-passwd"
)

//Repo ....
type Repo struct {
	//Reponame ...
	Reponame string `json:"name"`
	//RepoAuthor
	RepoAuthor string `json:"author"`
	//RepoStatus
	RepoStatus string `json:"status"`
}

//InitRepo ...
func InitRepo(repo string) *Repo {
	Newrepo := new(Repo)
	Newrepo.Reponame = repo
	return Newrepo
}

//LsRepo return a slice of repos name
func LsRepo(svnpath string) []string {
	cmd := exec.Command("ls", svnpath)
	stdout, err := cmd.StdoutPipe()
	defer stdout.Close()
	if err != nil {
		panic(err)
	}
	cmd.Start()
	data, err := ioutil.ReadAll(stdout)
	if err != nil {
		panic(err)
	}
	trimdata := strings.TrimSuffix(string(data), "\n")
	// fmt.Printf("%q", string(data))
	splitrepos := strings.SplitN(string(trimdata), "\n", -1)
	return splitrepos
}

//FetchRepoInfo return repo json format info
func (repo *Repo) FetchRepoInfo(svnpath string) (map[string]string, error) {
	hook := svnhook.NewHook(svnpath, repo.Reponame)
	author, err := hook.GetAuthor()
	if err != nil {
		return nil, err
	}
	info, err := hook.Getinfo()
	if err != nil {
		return nil, err
	}
	info["owner"] = author
	info["reponame"] = repo.Reponame
	return info, nil
}

//FetchRepoGroup ...
func (repo *Repo) FetchRepoGroup(svnpath string) (map[string]string, error) {
	auth := svnauth.NewSVNAuth(svnpath + "/" + repo.Reponame)
	groups, err := auth.ExportGroups()
	if err != nil {
		return nil, err
	}
	var groupsmap map[string]string
	err = json.Unmarshal(groups, &groupsmap)
	if err != nil {
		return nil, err
	}
	return groupsmap, nil
}

//FetchRepoEditedGroup  fetch edited group info .
func (repo *Repo) FetchRepoEditedGroup(svnpath string, olddata map[string]string, newdata map[string]string) (map[string]string, error) {
	auth := svnauth.NewSVNAuth(svnpath + "/" + repo.Reponame)
	groups, err := auth.ExportEditedGroup("groups", olddata, newdata)
	if err != nil {
		return nil, err
	}
	var groupsmap map[string]string
	err = json.Unmarshal(groups, &groupsmap)
	if err != nil {
		return nil, err
	}
	return groupsmap, nil
}

//FetchRepoDelGroup  fetch  group info  after deleted.
func (repo *Repo) FetchRepoDelGroup(svnpath string, olddata map[string]string) (map[string]string, error) {
	auth := svnauth.NewSVNAuth(svnpath + "/" + repo.Reponame)
	groups, err := auth.ExportDelGroup("groups", olddata)
	if err != nil {
		return nil, err
	}
	var groupsmap map[string]string
	err = json.Unmarshal(groups, &groupsmap)
	if err != nil {
		return nil, err
	}
	return groupsmap, nil
}

//FetchRepoAddGroup  fetch  group info  after inserted.
func (repo *Repo) FetchRepoAddGroup(svnpath string, newdata map[string]string) (map[string]string, error) {
	auth := svnauth.NewSVNAuth(svnpath + "/" + repo.Reponame)
	groups, err := auth.ExportAddGroup("groups", newdata)
	if err != nil {
		return nil, err
	}
	var groupsmap map[string]string
	err = json.Unmarshal(groups, &groupsmap)
	if err != nil {
		return nil, err
	}
	return groupsmap, nil
}

//FetchRepoDirectory ..
func (repo *Repo) FetchRepoDirectory(svnpath string) (map[string]map[string]string, error) {
	auth := svnauth.NewSVNAuth(svnpath + "/" + repo.Reponame)
	directory, err := auth.ExportDirectory()
	if err != nil {
		return nil, err
	}
	var directorymap map[string]map[string]string
	err = json.Unmarshal(directory, &directorymap)
	if err != nil {
		return nil, err
	}
	return directorymap, nil
}

//FetchRepoEditedDirectory fetch auth key-pair after edited.
func (repo *Repo) FetchRepoEditedDirectory(svnpath string, tag string, olddata map[string]string, newdata map[string]string) (map[string]map[string]string, error) {
	auth := svnauth.NewSVNAuth(svnpath + "/" + repo.Reponame)
	directory, err := auth.ExportEditedDirectory(tag, olddata, newdata)
	if err != nil {
		return nil, err
	}
	var directorymap map[string]map[string]string
	err = json.Unmarshal(directory, &directorymap)
	if err != nil {
		return nil, err
	}
	return directorymap, nil
}

//FetchRepoDelDirectory  fetch  auth key-pair after deleted.
func (repo *Repo) FetchRepoDelDirectory(svnpath string, tag string, olddata map[string]string) (map[string]map[string]string, error) {
	auth := svnauth.NewSVNAuth(svnpath + "/" + repo.Reponame)
	directory, err := auth.ExportDelDirectory(tag, olddata)
	if err != nil {
		return nil, err
	}
	var directorymap map[string]map[string]string
	err = json.Unmarshal(directory, &directorymap)
	if err != nil {
		return nil, err
	}
	return directorymap, nil
}

//FetchRepoAddDirectory  fetch auth key-par  after inserted.
func (repo *Repo) FetchRepoAddDirectory(svnpath string, tag string, newdata map[string]string) (map[string]map[string]string, error) {
	auth := svnauth.NewSVNAuth(svnpath + "/" + repo.Reponame)
	directory, err := auth.ExportAddDirectory(tag, newdata)
	if err != nil {
		return nil, err
	}
	var directorymap map[string]map[string]string
	err = json.Unmarshal(directory, &directorymap)
	if err != nil {
		return nil, err
	}
	return directorymap, nil
}

//FetchRepoUsers ...
func (repo *Repo) FetchRepoUsers(svnpath string) (map[string]string, error) {
	user := svnpasswd.NewUser(svnpath + "/" + repo.Reponame)
	users, err := user.ExportUser()
	if err != nil {
		return nil, err
	}
	var usermap map[string]string
	err = json.Unmarshal(users, &usermap)
	if err != nil {
		return nil, err
	}
	return usermap, nil
}

//FetchRepoGeneral ...
func (repo *Repo) FetchRepoGeneral(svnpath string) (map[string]string, error) {
	conf := svnconf.NewSVNconf(svnpath + "/" + repo.Reponame)
	general, err := conf.ExportGeneral()
	if err != nil {
		return nil, err
	}
	var generalmap map[string]string
	err = json.Unmarshal(general, &generalmap)
	if err != nil {
		return nil, err
	}
	return generalmap, nil
}

//FetchRepoRemarkGeneral remark section .
func (repo *Repo) FetchRepoRemarkGeneral(svnpath string, tag string) (map[string]string, error) {
	conf := svnconf.NewSVNconf(svnpath + "/" + repo.Reponame)
	general, err := conf.ExportRemarkGeneral(tag)
	if err != nil {
		return nil, err
	}
	var generalmap map[string]string
	err = json.Unmarshal(general, &generalmap)
	if err != nil {
		return nil, err
	}
	return generalmap, nil
}

//FetchRepoRmremarkGeneral rm remark section .
func (repo *Repo) FetchRepoRmremarkGeneral(svnpath string, tag string) (map[string]string, error) {
	conf := svnconf.NewSVNconf(svnpath + "/" + repo.Reponame)
	general, err := conf.ExportRmremarkGeneral(tag)
	if err != nil {
		return nil, err
	}
	var generalmap map[string]string
	err = json.Unmarshal(general, &generalmap)
	if err != nil {
		return nil, err
	}
	return generalmap, nil
}

//FetchRepoEditPasswd fetch edited passwd info
func (repo *Repo) FetchRepoEditPasswd(svnpath string, olddata map[string]string, newdata map[string]string) (map[string]string, error) {
	user := svnpasswd.NewUser(svnpath + "/" + repo.Reponame)
	editedUser, err := user.ExportEditedUser("users", olddata, newdata)
	if err != nil {
		return nil, err
	}
	var usermap map[string]string
	err = json.Unmarshal(editedUser, &usermap)
	if err != nil {
		return nil, err
	}
	return usermap, nil
}

//FetchRepoDelPasswd fetch user key-pair after a delete operation .
func (repo *Repo) FetchRepoDelPasswd(svnpath string, olddata map[string]string) (map[string]string, error) {
	user := svnpasswd.NewUser(svnpath + "/" + repo.Reponame)
	delUser, err := user.ExportAfterDelUser("users", olddata)
	if err != nil {
		return nil, err
	}
	var usermap map[string]string
	err = json.Unmarshal(delUser, &usermap)
	if err != nil {
		return nil, err
	}
	return usermap, nil
}

//FetchRepoAddPasswd fetch all key-pair after a adding operation
func (repo *Repo) FetchRepoAddPasswd(svnpath string, newdata map[string]string) (map[string]string, error) {
	user := svnpasswd.NewUser(svnpath + "/" + repo.Reponame)
	addUser, err := user.ExportAddUser("users", newdata)
	if err != nil {
		return nil, err
	}
	var usermap map[string]string
	err = json.Unmarshal(addUser, &usermap)
	if err != nil {
		return nil, err
	}
	return usermap, nil
}

//FetchReplaceDtag ...
func (repo Repo) FetchReplaceDtag(svnpath string, newtag string) error {
	auth := svnauth.NewSVNAuth(svnpath + "/" + repo.Reponame)
	err := auth.WrapReplaceDtag(newtag)
	if err != nil {
		return err
	}
	return nil
}
