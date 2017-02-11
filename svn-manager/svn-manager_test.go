package svnmanager

import "testing"

func TestLsrepo(t *testing.T) {
	data := LsRepo("/home/svn/")
	// t.Log(len(data))
	t.Log(data[0])
}

func TestFetchRepoinfo(t *testing.T) {
	repo := InitRepo("repo")
	repoinfo, err := repo.FetchRepoInfo("/home/svn")
	// repoinfo, err := json.Marshal(repo)
	if err != nil {
		t.Error(err)
	}
	for k, v := range repoinfo {
		t.Log(k, ":", v)
	}
}

func TestFetchGroup(t *testing.T) {
	repo := InitRepo("repo")
	groups, err := repo.FetchRepoGroup("/home/svn")
	if err != nil {
		t.Error(err)
	}
	for k, v := range groups {
		t.Log(k, ":", v)
	}
}

func TestFetchEditedGroup(t *testing.T) {
	repo := InitRepo("repo")
	olddata := make(map[string]string)
	olddata["groupname"] = "harry_and_sally"
	olddata["users"] = "harry,sally"
	newdata := make(map[string]string)
	newdata["groupname"] = "hahaha"
	newdata["users"] = "scnace,scbizu"
	groups, err := repo.FetchRepoEditedGroup("/home/svn", olddata, newdata)
	if err != nil {
		t.Error(err)
	}
	for k, v := range groups {
		t.Log(k, ":", v)
	}
}

func TestFetchAddGroup(t *testing.T) {
	repo := InitRepo("repo")
	newdata := make(map[string]string)
	newdata["groupname"] = "hahaha"
	newdata["users"] = "scnace,scbizu"
	groups, err := repo.FetchRepoAddGroup("/home/svn", newdata)
	if err != nil {
		t.Error(err)
	}
	for k, v := range groups {
		t.Log(k, ":", v)
	}
}

func TestFetchDelGroup(t *testing.T) {
	repo := InitRepo("repo")
	olddata := make(map[string]string)
	olddata["groupname"] = "hahaha"
	olddata["users"] = "scnace,scbizu"
	groups, err := repo.FetchRepoDelGroup("/home/svn", olddata)
	if err != nil {
		t.Error(err)
	}
	for k, v := range groups {
		t.Log(k, ":", v)
	}
}

func TestFetchRepoDirectory(t *testing.T) {
	repo := InitRepo("repo")
	directory, err := repo.FetchRepoDirectory("/home/svn")
	if err != nil {
		t.Error(err)
	}
	for k, v := range directory {
		for kk, vv := range v {
			t.Log(k, kk, vv)
		}

	}
}

func TestFetchRepoUsers(t *testing.T) {
	repo := InitRepo("repo")
	users, err := repo.FetchRepoUsers("/home/svn")
	if err != nil {
		t.Error(err)
	}
	for k, v := range users {
		t.Log(k, v)
	}
}

func TestFetchRepoGeneral(t *testing.T) {
	repo := InitRepo("repo")
	generals, err := repo.FetchRepoGeneral("/home/svn")
	if err != nil {
		t.Error(err)
	}
	for k, v := range generals {
		t.Log(k, v)
	}
}

func TestFetchRepoRemarkGeneral(t *testing.T) {
	repo := InitRepo("repo")
	generals, err := repo.FetchRepoRemarkGeneral("/home/svn", "authz-db")
	if err != nil {
		t.Error(err)
	}
	for k, v := range generals {
		t.Log(k, v)
	}
}

func TestFetchRepoRmremarkGeneral(t *testing.T) {
	repo := InitRepo("repo")
	generals, err := repo.FetchRepoRmremarkGeneral("/home/svn", "authz-db")
	if err != nil {
		t.Error(err)
	}
	for k, v := range generals {
		t.Log(k, v)
	}
}

func TestFetchRepoEditPasswd(t *testing.T) {
	repo := InitRepo("repo")
	olddata := make(map[string]string)
	olddata["username"] = "scnace"
	olddata["pwd"] = "scnace"
	newdata := make(map[string]string)
	newdata["username"] = "scnace"
	newdata["pwd"] = "123456"
	editUser, err := repo.FetchRepoEditPasswd("/home/svn", olddata, newdata)
	if err != nil {
		t.Error(err)
	}
	for k, v := range editUser {
		t.Log(k, v)
	}
}

func TestFetchRepoDelPasswd(t *testing.T) {
	repo := InitRepo("repo")
	olddata := make(map[string]string)
	olddata["username"] = "scnace"
	olddata["pwd"] = "scnace"
	delUser, err := repo.FetchRepoDelPasswd("/home/svn", olddata)
	if err != nil {
		t.Error(err)
	}
	for k, v := range delUser {
		t.Log(k, v)
	}
}

func TestFetchRepoAddPasswd(t *testing.T) {
	repo := InitRepo("repo")
	newdata := make(map[string]string)
	newdata["username"] = "scnace"
	newdata["pwd"] = "scnace"
	addUser, err := repo.FetchRepoAddPasswd("/home/svn", newdata)
	if err != nil {
		t.Error(err)
	}
	for k, v := range addUser {
		t.Log(k, v)
	}
}

func TestFetchRepoEditedDirectory(t *testing.T) {
	repo := InitRepo("repo")
	olddata := make(map[string]string)
	olddata["users"] = "@nace"
	olddata["auth"] = "rw"
	newdata := make(map[string]string)
	newdata["users"] = "scnace"
	newdata["auth"] = "rw"
	editDirectory, err := repo.FetchRepoEditedDirectory("/home/svn", "repository:/baz/fuz", olddata, newdata)
	if err != nil {
		t.Error(err)
	}
	for k, v := range editDirectory {
		t.Log(k, v)
	}
}

func TestFetchRepoAddDirectory(t *testing.T) {
	repo := InitRepo("repo")
	newdata := make(map[string]string)
	newdata["users"] = "scbizu"
	newdata["auth"] = "r"
	addDirectory, err := repo.FetchRepoAddDirectory("/home/svn", "/foo/bar", newdata)
	if err != nil {
		t.Error(err)
	}
	for k, v := range addDirectory {
		t.Log(k, v)
	}
}

func TestFetchRepoDelDirectory(t *testing.T) {
	repo := InitRepo("repo")
	olddata := make(map[string]string)
	olddata["users"] = "scbizu"
	olddata["auth"] = "r"
	delDirectory, err := repo.FetchRepoDelDirectory("/home/svn", "/foo/bar", olddata)
	if err != nil {
		t.Error(err)
	}
	for k, v := range delDirectory {
		t.Log(k, v)
	}
}

func TestFetchReplaceDtag(t *testing.T) {
	repo := InitRepo("repo")
	err := repo.FetchReplaceDtag("/home/svn", "repo")
	if err != nil {
		t.Error(err)
	}
	t.Log("dtag repalced.")
}
