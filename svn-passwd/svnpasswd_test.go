package svnpasswd

import "testing"

func TestExportUser(t *testing.T) {
	user := NewUser("/home/svn/repo")
	res, err := user.ExportUser()
	if err != nil {
		t.Error(err)
	}
	t.Log(string(res))
}

func TestExportEditedUser(t *testing.T) {
	user := NewUser("/home/svn/repo")
	olddata := make(map[string]string)
	olddata["username"] = "scnace"
	olddata["pwd"] = "scnace"
	newdata := make(map[string]string)
	newdata["username"] = "scnace"
	newdata["pwd"] = "123456"
	editUser, err := user.ExportEditedUser("users", olddata, newdata)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(editUser))
}

func TestExportAfterDelUser(t *testing.T) {
	user := NewUser("/home/svn/repo")
	olddata := make(map[string]string)
	olddata["username"] = "scnace"
	olddata["pwd"] = "scnace"
	afterDel, err := user.ExportAfterDelUser("users", olddata)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(afterDel))
}

func TestExportAddUser(t *testing.T) {
	user := NewUser("/home/svn/repo")
	newdata := make(map[string]string)
	newdata["username"] = "heiheihei"
	newdata["pwd"] = "123456"
	afterAdd, err := user.ExportAddUser("users", newdata)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(afterAdd))
}
