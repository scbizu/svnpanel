package svnauth

import "testing"

func Test_NewSVNAuth(t *testing.T) {}

func Test_ExportGroups(t *testing.T) {
	auth := NewSVNAuth("/home/svn/repo")
	data, err := auth.ExportGroups()
	if err != nil {
		t.Error(err)
	}
	t.Log(string(data))
}

func TestExportDirectory(t *testing.T) {
	auth := NewSVNAuth("/home/svn/repo")
	data, err := auth.ExportDirectory()
	if err != nil {
		t.Error(err)
	}
	t.Log(string(data))
}

func TestExportEditedGroup(t *testing.T) {
	auth := NewSVNAuth("/home/svn/repo")
	olddata := make(map[string]string)
	olddata["groupname"] = "harry_and_sally"
	olddata["users"] = "harry,sally"
	newdata := make(map[string]string)
	newdata["groupname"] = "hahaha"
	newdata["users"] = "scnace,scbizu"
	data, err := auth.ExportEditedGroup("groups", olddata, newdata)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(data))
}

func TestExportDelGroup(t *testing.T) {
	auth := NewSVNAuth("/home/svn/repo")
	olddata := make(map[string]string)
	olddata["groupname"] = "nace"
	olddata["users"] = "scnace"
	data, err := auth.ExportDelGroup("groups", olddata)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(data))
}

func TestExportAddGroup(t *testing.T) {
	auth := NewSVNAuth("/home/svn/repo")
	newdata := make(map[string]string)
	newdata["groupname"] = "nace"
	newdata["users"] = "scnace,scbizu"
	data, err := auth.ExportAddGroup("groups", newdata)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(data))
}

func TestExportEditedDirectory(t *testing.T) {
	auth := NewSVNAuth("/home/svn/repo")
	olddata := make(map[string]string)
	olddata["users"] = "@harry_and_sally"
	olddata["auth"] = "rw"
	newdata := make(map[string]string)
	newdata["users"] = "@nace"
	newdata["auth"] = "rw"
	data, err := auth.ExportEditedDirectory("repository:/baz/fuz", olddata, newdata)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(data))
}

func TestExportAddDirectory(t *testing.T) {
	auth := NewSVNAuth("/home/svn/repo")
	newdata := make(map[string]string)
	newdata["users"] = "scbizu"
	newdata["auth"] = "r"
	data, err := auth.ExportAddDirectory("/foo", newdata)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(data))
}

func TestExportDelDirectory(t *testing.T) {
	auth := NewSVNAuth("/home/svn/repo")
	olddata := make(map[string]string)
	olddata["users"] = "scbizu"
	olddata["auth"] = "r"
	data, err := auth.ExportDelDirectory("/foo/bar", olddata)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(data))
}

func TestWrapReplaceDtag(t *testing.T) {
	auth := NewSVNAuth("/home/svn/repo")
	err := auth.WrapReplaceDtag("repo")
	if err != nil {
		t.Error(err)
	}
	t.Log("replaced Dtag.")
}
