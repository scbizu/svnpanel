package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"html/template"
	"io"
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
	"github.com/scbizu/svnpanel/gconfig"
	"github.com/scbizu/svnpanel/svn-hook"
	"github.com/scbizu/svnpanel/svn-manager"
)

//Conf ...
type Conf struct {
	Reponame    string
	Author      string
	Aliases     string
	Group       map[string]string
	Directories []map[string]string
	Users       map[string]string
	Anon        string
	Pdb         string
	Adb         string
}

//Template ...
type Template struct {
	templates *template.Template
}

//Render ...
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	server := echo.New()
	gconf := gconfig.NewGconfig()
	SVNPATH := gconf.SVNPATH

	//BasicAuth
	server.Use(middleware.BasicAuth(func(username string, pwd string) bool {
		//solve salt
		salt := gconf.Salt
		rawcode := pwd + "?" + salt + "?"
		hasher := md5.New()
		hasher.Write([]byte(rawcode))
		if username == gconf.Username && hex.EncodeToString(hasher.Sum(nil)) == gconf.Password {
			return true
		}
		return false
	}))
	server.Static("/", "assets")
	t := &Template{
		templates: template.Must(template.ParseGlob("assets/config.html")),
	}
	server.SetRenderer(t)

	server.File("/", "assets/index.html")

	server.GET("/users", func(c echo.Context) error {
		repo := c.QueryParam("repo")
		repopath := svnmanager.InitRepo(repo)
		user, err := repopath.FetchRepoUsers(SVNPATH)
		if err != nil {
			return c.JSON(403, "invaild reponame")
		}

		usersjson, _ := json.Marshal(user)
		return c.JSON(200, string(usersjson))
	})

	server.GET("/groups", func(c echo.Context) error {
		repo := c.QueryParam("repo")
		repopath := svnmanager.InitRepo(repo)
		group, err := repopath.FetchRepoGroup(SVNPATH)
		if err != nil {
			return c.JSON(403, "invaild reponame")
		}
		gjson, _ := json.Marshal(group)
		return c.JSON(200, string(gjson))
	})

	server.GET("/config/:repo", func(c echo.Context) error {
		repo := c.Param("repo")
		hook := svnhook.NewHook(SVNPATH, repo)
		author, err := hook.GetAuthor()
		if err != nil {
			return err
		}

		repopath := svnmanager.InitRepo(repo)
		err = repopath.FetchReplaceDtag(SVNPATH, repo)
		if err != nil {
			return c.JSON(500, []byte("replace dtag error!"))
		}
		groups, err := repopath.FetchRepoGroup(SVNPATH)
		if err != nil {
			return err
		}
		directory, err := repopath.FetchRepoDirectory(SVNPATH)
		if err != nil {
			return err
		}
		// directories := make([]map[string]string, 6)
		var directories []map[string]string

		for k, v := range directory {
			for kk, vv := range v {
				tempDirectories := make(map[string]string)
				tempDirectories["path"] = k
				tempDirectories["user"] = kk
				tempDirectories["auth"] = vv
				directories = append(directories, tempDirectories)

			}
		}

		users, err := repopath.FetchRepoUsers(SVNPATH)
		if err != nil {
			return err
		}
		generals, err := repopath.FetchRepoGeneral(SVNPATH)
		if err != nil {
			return err
		}
		conf := Conf{
			Reponame:    repo,
			Author:      author,
			Group:       groups,
			Directories: directories,
			Users:       users,
			Anon:        generals["anon-access"],
			Pdb:         generals["password-db"],
			Adb:         generals["authz-db"],
		}
		return c.Render(http.StatusOK, "conf", conf)
	})

	server.GET("/repos", func(c echo.Context) error {
		data := svnmanager.LsRepo(SVNPATH)
		// var datamap map[string]string
		var datamaps []map[string]string
		for _, v := range data {
			repo := svnmanager.InitRepo(v)
			infomap, err := repo.FetchRepoInfo(SVNPATH)
			if err != nil {
				return err
			}
			datamaps = append(datamaps, infomap)
		}

		datajson, err := json.Marshal(datamaps)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, string(datajson))
	})

	//edit svnserve.conf file
	server.Put("/edit", func(c echo.Context) error {
		tag := strings.TrimSpace(c.FormValue("tag"))
		action := strings.TrimSpace(c.FormValue("action"))
		reponame := strings.TrimSpace(c.FormValue("reponame"))
		repo := svnmanager.InitRepo(reponame)
		var generals map[string]string
		switch action {
		case "remark":
			var err error
			generals, err = repo.FetchRepoRemarkGeneral(SVNPATH, tag)
			if err != nil {
				return err
			}
		case "rmremark":
			var err error
			generals, err = repo.FetchRepoRmremarkGeneral(SVNPATH, tag)
			if err != nil {
				return err
			}
		}
		conf := new(Conf)
		conf.Adb = generals["authz-db"]
		conf.Pdb = generals["password-db"]
		gjson, err := json.Marshal(conf)
		if err != nil {
			return err
		}
		return c.JSON(200, string(gjson))
	})
	//edit passwd file
	server.Post("/passwd", func(c echo.Context) error {

		olddataUsername := c.FormValue("old_username")
		olddataPwd := c.FormValue("old_pwd")
		newdataUsername := c.FormValue("new_username")
		newdataPwd := c.FormValue("new_pwd")
		reponame := c.FormValue("reponame")
		olddata := make(map[string]string)
		olddata["username"] = olddataUsername
		olddata["pwd"] = olddataPwd
		newdata := make(map[string]string)
		newdata["username"] = newdataUsername
		newdata["pwd"] = newdataPwd
		repo := svnmanager.InitRepo(reponame)
		passwdmap, err := repo.FetchRepoEditPasswd(SVNPATH, olddata, newdata)
		if err != nil {
			return c.JSON(500, err.Error())
		}
		passwdjson, err := json.Marshal(passwdmap)
		if err != nil {
			return c.JSON(500, err.Error())
		}
		return c.JSON(200, string(passwdjson))
	})
	//delete user key-par
	server.Post("/delpasswd", func(c echo.Context) error {
		olddataUsername := c.FormValue("old_username")
		olddataPwd := c.FormValue("old_pwd")
		reponame := c.FormValue("reponame")
		olddata := make(map[string]string)
		olddata["username"] = olddataUsername
		olddata["pwd"] = olddataPwd
		repo := svnmanager.InitRepo(reponame)
		dmap, err := repo.FetchRepoDelPasswd(SVNPATH, olddata)
		if err != nil {
			return c.JSON(500, err.Error())
		}
		djson, err := json.Marshal(dmap)
		if err != nil {
			return c.JSON(500, err.Error())
		}
		return c.JSON(200, string(djson))
	})
	// add a user key-pair
	server.Post("/newpasswd", func(c echo.Context) error {
		newdataUsername := c.FormValue("new_username")
		newdataPwd := c.FormValue("new_pwd")
		reponame := c.FormValue("reponame")
		newdata := make(map[string]string)
		newdata["username"] = newdataUsername
		newdata["pwd"] = newdataPwd
		repo := svnmanager.InitRepo(reponame)
		amap, err := repo.FetchRepoAddPasswd(SVNPATH, newdata)
		if err != nil {
			return c.JSON(500, err.Error())
		}
		ajson, err := json.Marshal(amap)
		if err != nil {
			return c.JSON(500, err.Error())
		}
		return c.JSON(200, string(ajson))
	})
	//edit authz file` groups ` section
	server.Post("/groups", func(c echo.Context) error {
		olddataGroups := c.FormValue("old_groupname")
		olddataUsers := c.FormValue("old_users")
		newdataGroups := c.FormValue("new_groupname")
		newdataUsers := c.FormValue("new_users")
		reponame := c.FormValue("reponame")
		olddata := make(map[string]string)
		olddata["groupname"] = olddataGroups
		olddata["users"] = olddataUsers
		newdata := make(map[string]string)
		newdata["groupname"] = newdataGroups
		newdata["users"] = newdataUsers
		repo := svnmanager.InitRepo(reponame)
		gmap, err := repo.FetchRepoEditedGroup(SVNPATH, olddata, newdata)
		if err != nil {
			return c.JSON(500, err.Error())
		}
		gjson, err := json.Marshal(gmap)
		if err != nil {
			return c.JSON(500, err.Error())
		}
		return c.JSON(200, string(gjson))
	})
	//add a new group
	server.Post("/addgroup", func(c echo.Context) error {
		newdataGroups := c.FormValue("new_groupname")
		newdataUsers := c.FormValue("new_users")
		reponame := c.FormValue("reponame")
		newdata := make(map[string]string)
		newdata["groupname"] = newdataGroups
		newdata["users"] = newdataUsers
		repo := svnmanager.InitRepo(reponame)
		gmap, err := repo.FetchRepoAddGroup(SVNPATH, newdata)
		if err != nil {
			return c.JSON(500, err.Error())
		}
		gjson, err := json.Marshal(gmap)
		if err != nil {
			return c.JSON(500, err.Error())
		}
		return c.JSON(200, gjson)
	})
	//delete a group
	server.Post("/delgroup", func(c echo.Context) error {
		olddataGroups := c.FormValue("old_groupname")
		olddataUsers := c.FormValue("old_users")
		reponame := c.FormValue("reponame")
		olddata := make(map[string]string)
		olddata["groupname"] = olddataGroups
		olddata["users"] = olddataUsers
		repo := svnmanager.InitRepo(reponame)
		gmap, err := repo.FetchRepoDelGroup(SVNPATH, olddata)
		if err != nil {
			return c.JSON(500, err.Error())
		}
		gjson, err := json.Marshal(gmap)
		if err != nil {
			return c.JSON(500, err.Error())
		}
		return c.JSON(200, string(gjson))
	})
	//edit the key-pair of a certain repo directory
	server.Post("/editauth", func(c echo.Context) error {
		olddataUsers := c.FormValue("old_users")
		olddataAuth := c.FormValue("old_auth")
		newdataUsers := c.FormValue("new_users")
		newdataAuth := c.FormValue("new_auth")

		tag := c.FormValue("tag")
		reponame := c.FormValue("reponame")

		olddata := make(map[string]string)
		olddata["users"] = strings.TrimSpace(olddataUsers)
		olddata["auth"] = strings.TrimSpace(olddataAuth)
		newdata := make(map[string]string)
		newdata["users"] = strings.TrimSpace(newdataUsers)
		newdata["auth"] = strings.TrimSpace(newdataAuth)
		repo := svnmanager.InitRepo(reponame)
		amap, err := repo.FetchRepoEditedDirectory(SVNPATH, tag, olddata, newdata)
		if err != nil {
			return c.JSON(500, err.Error())
		}
		ajson, err := json.Marshal(amap)
		if err != nil {
			return c.JSON(500, err.Error())
		}
		return c.JSON(200, string(ajson))
	})

	//delete one auth key-pair under a certain repo directory
	server.Post("/delauth", func(c echo.Context) error {
		olddataUsers := c.FormValue("old_users")
		olddataAuth := c.FormValue("old_auth")
		tag := c.FormValue("tag")
		reponame := c.FormValue("reponame")
		olddata := make(map[string]string)
		olddata["users"] = strings.TrimSpace(olddataUsers)
		olddata["auth"] = strings.TrimSpace(olddataAuth)
		repo := svnmanager.InitRepo(reponame)
		amap, err := repo.FetchRepoDelDirectory(SVNPATH, tag, olddata)
		if err != nil {
			return c.JSON(500, err.Error())
		}
		ajson, err := json.Marshal(amap)
		if err != nil {
			return c.JSON(500, err.Error())
		}
		return c.JSON(200, string(ajson))
	})

	//add one auth key-pair under a certain repo directory
	server.Post("/addauth", func(c echo.Context) error {
		newdataUsers := c.FormValue("new_users")
		newdataAuth := c.FormValue("new_auth")
		tag := c.FormValue("tag")
		reponame := c.FormValue("reponame")

		newdata := make(map[string]string)
		newdata["users"] = strings.TrimSpace(newdataUsers)
		newdata["auth"] = strings.TrimSpace(newdataAuth)
		repo := svnmanager.InitRepo(reponame)
		amap, err := repo.FetchRepoAddDirectory(SVNPATH, tag, newdata)
		if err != nil {
			return c.JSON(500, err.Error())
		}
		ajson, err := json.Marshal(amap)
		if err != nil {
			return c.JSON(500, err.Error())
		}
		return c.JSON(200, string(ajson))
	})

	server.Run(standard.New(":1312"))
}
