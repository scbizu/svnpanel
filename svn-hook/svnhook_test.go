package svnhook

import "testing"

func TestGetAuthor(t *testing.T) {
	hook := NewHook("/home/svn", "repo")
	author, err := hook.GetAuthor()
	if err != nil {
		t.Error(err)
	}
	t.Log(author)
}

func TestGetChanged(t *testing.T) {
	hook := NewHook("/home/svn", "repo")
	changed, err := hook.GetChanged()
	if err != nil {
		t.Error(err)
	}
	t.Log(changed)
}

func TestGetdate(t *testing.T) {
	hook := NewHook("/home/svn", "repo")
	date, err := hook.Getdate()
	if err != nil {
		t.Error(err)
	}
	t.Log(date)
}

func TestGetinfo(t *testing.T) {
	hook := NewHook("/home/svn", "repo")
	info, err := hook.Getinfo()
	if err != nil {
		t.Error(err)
	}
	for k, v := range info {
		t.Log(k, ":", v)
	}
}
