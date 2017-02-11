package svnconf

import "testing"

func TestExportGeneral(t *testing.T) {
	conf := NewSVNconf("/home/svn/repo")
	general, err := conf.ExportGeneral()
	if err != nil {
		t.Error(err)
	}
	t.Log(string(general))
}

func TestExporRemarkGeneral(t *testing.T) {
	conf := NewSVNconf("/home/svn/repo")
	general, err := conf.ExportRemarkGeneral("authz-db")
	if err != nil {
		t.Error(err)
	}
	t.Log(string(general))
}

func TestExporRmremarkGeneral(t *testing.T) {
	conf := NewSVNconf("/home/svn/repo")
	general, err := conf.ExportRmremarkGeneral("authz-db")
	if err != nil {
		t.Error(err)
	}
	t.Log(string(general))
}
