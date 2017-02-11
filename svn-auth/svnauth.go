package svnauth

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

//SVNAuth the auth structure
type SVNAuth struct {
	//SVNPATH  the repo path
	SVNPATH string `json:"svnpath"`
	Aliases string `json:"aliases"`
	Groups  string `json:"groups"`
}

//NewSVNAuth init a svn structure
func NewSVNAuth(path string) *SVNAuth {
	auth := &SVNAuth{}
	auth.SVNPATH = path
	return auth
}

//Readfill will read the authz to the buffer
func (auth *SVNAuth) readfile() ([]byte, error) {
	authz := filepath.Join(auth.SVNPATH, "conf", "authz")
	filehandle, err := os.Open(authz)
	defer filehandle.Close()
	if err != nil {
		return nil, errors.New(`Please check your ./conf/authz file,if you change it's name,rename it as "authz".`)
	}

	content, err := ioutil.ReadAll(filehandle)
	if err != nil {
		return nil, err
	}
	return content, nil
}

//file2json parse authz content to readable json .
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

//edit
//Attention: keep the groupname as a primary key .
func (auth *SVNAuth) changeGroup(tagname string, olddata map[string]string, newdata map[string]string, srcContent []byte) error {
	authzPath := filepath.Join(auth.SVNPATH, "conf", "authz")
	var lines []string
	var dstContent bytes.Buffer
	lines = strings.Split(string(srcContent), "\n")
	// O(n) loop
	for k, v := range lines {
		t := strings.TrimPrefix(v, "[")
		t = strings.TrimSuffix(t, "]")
		if strings.TrimSpace(t) == tagname {
			dstContent.WriteString(lines[k] + "\n")
			for i := k + 1; i < len(lines); i++ {
				if strings.HasPrefix(lines[i], "[") {
					break
				}
				tsplit := strings.SplitN(lines[i], "=", 2)
				//find the username ;  not  change username
				if strings.TrimSpace(tsplit[0]) == olddata["groupname"] {
					// replace rudely
					lines[i] = newdata["groupname"] + "=" + newdata["users"]
				}
			}
		} else {
			dstContent.WriteString(v + "\n")
		}
	}
	err := ioutil.WriteFile(authzPath, dstContent.Bytes(), 0666)
	if err != nil {
		return err
	}
	return nil
}

//delete
//Attention: keep the groupname as a primary key .
func (auth *SVNAuth) delGroup(tagname string, olddata map[string]string, srcContent []byte) error {
	authzPath := filepath.Join(auth.SVNPATH, "conf", "authz")
	var lines []string
	var dstContent bytes.Buffer
	lines = strings.Split(string(srcContent), "\n")
	// O(n) loop
	for k, v := range lines {
		t := strings.TrimPrefix(v, "[")
		t = strings.TrimSuffix(t, "]")
		if strings.TrimSpace(t) == tagname {
			dstContent.WriteString(lines[k] + "\n")
			for i := k + 1; i < len(lines); i++ {
				if strings.HasPrefix(lines[i], "[") {
					break
				}
				tsplit := strings.SplitN(lines[i], "=", 2)
				//find the username ;  not  change username
				if strings.TrimSpace(tsplit[0]) == olddata["groupname"] {
					// not really delete ,avoid the risk that lose a user key-pair permanently
					lines[i] = "#" + olddata["groupname"] + "=" + olddata["users"]
				}
			}
		} else {
			dstContent.WriteString(v + "\n")
		}
	}
	err := ioutil.WriteFile(authzPath, dstContent.Bytes(), 0666)
	if err != nil {
		return err
	}
	return nil
}

//delete
//Attention: keep the groupname as a primary key .
func (auth *SVNAuth) addGroup(tagname string, newdata map[string]string, srcContent []byte) error {
	authzPath := filepath.Join(auth.SVNPATH, "conf", "authz")
	var lines []string
	var dstContent bytes.Buffer
	lines = strings.Split(string(srcContent), "\n")
	// O(n) loop
	for k, v := range lines {
		t := strings.TrimPrefix(v, "[")
		t = strings.TrimSuffix(t, "]")
		if strings.TrimSpace(t) == tagname {
			dstContent.WriteString(lines[k] + "\n")
			for i := k + 1; i < len(lines); i++ {
				if strings.HasPrefix(lines[i], "[") {
					break
				}
				tsplit := strings.SplitN(lines[i], "=", 2)
				//find the username ;  not  change username
				if strings.TrimSpace(tsplit[0]) == newdata["groupname"] {
					//return a existed error
					return errors.New("group existed")
				}
			}
			dstContent.WriteString(newdata["groupname"] + "=" + newdata["users"] + "\n")
		} else {
			dstContent.WriteString(v + "\n")
		}
	}
	err := ioutil.WriteFile(authzPath, dstContent.Bytes(), 0666)
	if err != nil {
		return err
	}
	return nil
}

func (auth *SVNAuth) changeDirectory(tagname string, olddata map[string]string, newdata map[string]string, srcContent []byte) error {
	authzPath := filepath.Join(auth.SVNPATH, "conf", "authz")
	var lines []string
	var dstContent bytes.Buffer
	lines = strings.Split(string(srcContent), "\n")
	// O(n) loop
	for k, v := range lines {
		t := strings.TrimPrefix(v, "[")
		t = strings.TrimSuffix(t, "]")
		if strings.TrimSpace(t) == tagname {
			dstContent.WriteString(lines[k] + "\n")
			for i := k + 1; i < len(lines); i++ {
				if strings.HasPrefix(lines[i], "[") {
					break
				}
				tsplit := strings.SplitN(lines[i], "=", 2)
				//find the username ;  not  change username
				if strings.TrimSpace(tsplit[0]) == olddata["users"] {
					// replace rudely
					lines[i] = newdata["users"] + "=" + newdata["auth"]
				}
			}
		} else {
			dstContent.WriteString(v + "\n")
		}
	}
	err := ioutil.WriteFile(authzPath, dstContent.Bytes(), 0666)
	if err != nil {
		return err
	}
	return nil
}

//replaceDtagPrefix replace initial reponame or some bad reponames
func (auth SVNAuth) replaceDtagPrefix(newprefix string, srcContent []byte) error {
	authzPath := filepath.Join(auth.SVNPATH, "conf", "authz")
	var lines []string
	var dstContent bytes.Buffer
	lines = strings.Split(string(srcContent), "\n")
	// O(n) loop
	for _, v := range lines {
		if strings.Contains(v, "[") && strings.Contains(v, "]") && !strings.HasPrefix(v, "#") {
			t := strings.TrimPrefix(v, "[")
			t = strings.TrimSuffix(t, "]")
			if strings.Contains(t, ":") {
				tsplit := strings.SplitN(t, ":", 2)
				tstring := "[" + newprefix + ":" + tsplit[1] + "]"
				dstContent.WriteString(tstring + "\n")
			} else {
				dstContent.WriteString(v + "\n")
			}
		} else {
			dstContent.WriteString(v + "\n")
		}
	}
	err := ioutil.WriteFile(authzPath, dstContent.Bytes(), 0666)
	if err != nil {
		return err
	}
	return nil
}

//Attention: keep the users as a primary key .
func (auth *SVNAuth) addDirectory(tagname string, newdata map[string]string, srcContent []byte) error {
	authzPath := filepath.Join(auth.SVNPATH, "conf", "authz")
	var lines []string
	var dstContent bytes.Buffer
	tagExisted := false
	lines = strings.Split(string(srcContent), "\n")
	// O(n) loop
	for k, v := range lines {
		t := strings.TrimPrefix(v, "[")
		t = strings.TrimSuffix(t, "]")
		if strings.TrimSpace(t) == tagname {
			tagExisted = true
			dstContent.WriteString(lines[k] + "\n")
			for i := k + 1; i < len(lines); i++ {
				if strings.HasPrefix(lines[i], "[") {
					break
				}
				tsplit := strings.SplitN(lines[i], "=", 2)
				//find the username ;  not  change username
				if strings.TrimSpace(tsplit[0]) == newdata["users"] {
					//return a existed error
					return errors.New("Directory existed")
				}
			}
			dstContent.WriteString(newdata["users"] + "=" + newdata["auth"] + "\n")
		} else {
			if strings.TrimSpace(v) != "" {
				dstContent.WriteString(v + "\n")
			}
		}
	}
	//tag was not existed yet
	if !tagExisted {
		dstContent.WriteString("[" + tagname + "]" + "\n")
		dstContent.WriteString(newdata["users"] + "=" + newdata["auth"] + "\n")
	}
	err := ioutil.WriteFile(authzPath, dstContent.Bytes(), 0666)
	if err != nil {
		return err
	}
	return nil
}

//Attention: keep the users as a primary key .
func (auth *SVNAuth) delDirectory(tagname string, olddata map[string]string, srcContent []byte) error {
	authzPath := filepath.Join(auth.SVNPATH, "conf", "authz")
	var lines []string
	var dstContent bytes.Buffer
	lines = strings.Split(string(srcContent), "\n")
	// O(n) loop
	for k, v := range lines {
		t := strings.TrimPrefix(v, "[")
		t = strings.TrimSuffix(t, "]")
		if strings.TrimSpace(t) == tagname {
			dstContent.WriteString(lines[k] + "\n")
			for i := k + 1; i < len(lines); i++ {
				if strings.HasPrefix(lines[i], "[") {
					break
				}
				tsplit := strings.SplitN(lines[i], "=", 2)
				//find the username ;  not  change username
				if strings.TrimSpace(tsplit[0]) == olddata["users"] {
					// not really delete ,avoid the risk that lose a user key-pair permanently
					lines[i] = "#" + olddata["users"] + "=" + olddata["auth"]
				}
			}
		} else {
			dstContent.WriteString(v + "\n")
		}
		err := ioutil.WriteFile(authzPath, dstContent.Bytes(), 0666)
		if err != nil {
			return err
		}
	}
	return nil
}

/////////////////////////parse function /////////////

//parse groups json
func parseGroups(rawjson []byte) ([]byte, error) {
	authz := make(map[string][]string)
	err := json.Unmarshal(rawjson, &authz)
	if err != nil {
		return []byte(""), err
	}
	usergroups := make(map[string]string)
	groups := authz["groups"]
	for _, v := range groups {
		trimv := strings.TrimSpace(string(v))
		splitv := strings.SplitN(trimv, "=", 2)
		usergroups[strings.TrimSpace(splitv[0])] = strings.TrimSpace(splitv[1])
	}
	groupsjson, err := json.Marshal(usergroups)
	if err != nil {
		return []byte(""), err
	}
	return groupsjson, nil
}

//pasreDirectory parse other section except [aliases] [groups]
func parseDirectory(rawjson []byte) ([]byte, error) {
	authz := make(map[string][]string)
	err := json.Unmarshal(rawjson, &authz)
	if err != nil {
		return []byte(""), err
	}
	var userauth map[string]string
	directory := make(map[string]map[string]string)
	for k, v := range authz {
		userauth = make(map[string]string)
		if k != "aliases" && k != "groups" {
			for _, authrow := range v {
				trim := strings.TrimSpace(string(authrow))
				split := strings.SplitN(trim, "=", 2)
				userauth[split[0]] = split[1]
			}
			directory[k] = userauth
		}
	}
	directoryjson, err := json.Marshal(directory)
	if err != nil {
		return []byte(""), err
	}
	return directoryjson, nil
}

//Unsupport aliases
// func parseAliases(rawjson []byte) ([]byte, error) {
// 	authz := make(map[string]string)
// 	err := json.Unmarshal(rawjson, &authz)
// 	if err != nil {
// 		return []byte(""), err
// 	}
// 	useraliases := make(map[string]string)
// 	aliases := authz["aliases"]
// 	for _, v := range aliases {
// 		trimv := strings.TrimSpace(string(v))
// 		splitv := strings.SplitN(trimv, "=", 2)
// 		useraliases[splitv[0]] = splitv[1]
// 	}
// 	aliasesjson, err := json.Marshal(useraliases)
// 	if err != nil {
// 		return []byte(""), err
// 	}
// 	return aliasesjson, nil
// }

//ExportGroups ...
func (auth *SVNAuth) ExportGroups() ([]byte, error) {
	content, err := auth.readfile()
	if err != nil {
		return []byte(""), err
	}
	rawjson, err := file2json(content)
	if err != nil {
		return []byte(""), err
	}
	groups, err := parseGroups(rawjson)
	if err != nil {
		return []byte(""), err
	}
	return groups, nil
}

//ExportDirectory export directory json data
func (auth *SVNAuth) ExportDirectory() ([]byte, error) {
	content, err := auth.readfile()
	if err != nil {
		return []byte(""), err
	}
	rawjson, err := file2json(content)
	if err != nil {
		return []byte(""), err
	}
	directory, err := parseDirectory(rawjson)
	if err != nil {
		return []byte(""), err
	}
	return directory, nil
}

//ExportEditedGroup export edited group data
func (auth *SVNAuth) ExportEditedGroup(tag string, olddata map[string]string, newdata map[string]string) ([]byte, error) {
	content, err := auth.readfile()
	if err != nil {
		return []byte(""), err
	}
	err = auth.changeGroup(tag, olddata, newdata, content)
	if err != nil {
		return []byte(""), err
	}
	content, err = auth.readfile()
	if err != nil {
		return []byte(""), err
	}
	rawjson, err := file2json(content)
	if err != nil {
		return []byte(""), err
	}
	groups, err := parseGroups(rawjson)
	if err != nil {
		return []byte(""), err
	}
	return groups, nil
}

//ExportDelGroup export group info after  a delete .
func (auth *SVNAuth) ExportDelGroup(tag string, olddata map[string]string) ([]byte, error) {
	content, err := auth.readfile()
	if err != nil {
		return []byte(""), err
	}
	err = auth.delGroup(tag, olddata, content)
	if err != nil {
		return []byte(""), err
	}
	content, err = auth.readfile()
	if err != nil {
		return []byte(""), err
	}
	rawjson, err := file2json(content)
	if err != nil {
		return []byte(""), err
	}
	groups, err := parseGroups(rawjson)
	if err != nil {
		return []byte(""), err
	}
	return groups, nil
}

//ExportAddGroup export group info after  an add operation .
func (auth *SVNAuth) ExportAddGroup(tag string, newdata map[string]string) ([]byte, error) {
	content, err := auth.readfile()
	if err != nil {
		return []byte(""), err
	}
	err = auth.addGroup(tag, newdata, content)
	if err != nil {
		return []byte(""), err
	}
	content, err = auth.readfile()
	if err != nil {
		return []byte(""), err
	}
	rawjson, err := file2json(content)
	if err != nil {
		return []byte(""), err
	}
	groups, err := parseGroups(rawjson)
	if err != nil {
		return []byte(""), err
	}
	return groups, nil
}

//ExportEditedDirectory export directory info after edited
func (auth *SVNAuth) ExportEditedDirectory(tag string, olddata map[string]string, newdata map[string]string) ([]byte, error) {
	content, err := auth.readfile()
	if err != nil {
		return []byte(""), err
	}
	err = auth.changeDirectory(tag, olddata, newdata, content)
	if err != nil {
		return []byte(""), err
	}
	content, err = auth.readfile()
	if err != nil {
		return []byte(""), err
	}
	rawjson, err := file2json(content)
	if err != nil {
		return []byte(""), err
	}
	directory, err := parseDirectory(rawjson)
	if err != nil {
		return []byte(""), err
	}
	return directory, nil
}

//ExportAddDirectory export group info after  an add operation .
func (auth *SVNAuth) ExportAddDirectory(tag string, newdata map[string]string) ([]byte, error) {
	content, err := auth.readfile()
	if err != nil {
		return []byte(""), err
	}
	err = auth.addDirectory(tag, newdata, content)
	if err != nil {
		return []byte(""), err
	}
	content, err = auth.readfile()
	if err != nil {
		return []byte(""), err
	}
	rawjson, err := file2json(content)
	if err != nil {
		return []byte(""), err
	}
	directory, err := parseDirectory(rawjson)
	if err != nil {
		return []byte(""), err
	}
	return directory, nil
}

//ExportDelDirectory export auth key-pair after  a delete .
func (auth *SVNAuth) ExportDelDirectory(tag string, olddata map[string]string) ([]byte, error) {
	content, err := auth.readfile()
	if err != nil {
		return []byte(""), err
	}
	err = auth.delDirectory(tag, olddata, content)
	if err != nil {
		return []byte(""), err
	}
	content, err = auth.readfile()
	if err != nil {
		return []byte(""), err
	}
	rawjson, err := file2json(content)
	if err != nil {
		return []byte(""), err
	}
	directory, err := parseDirectory(rawjson)
	if err != nil {
		return []byte(""), err
	}
	return directory, nil
}

//WrapReplaceDtag wrap replaceDtagPrefix function
func (auth *SVNAuth) WrapReplaceDtag(newtag string) error {
	content, err := auth.readfile()
	if err != nil {
		return err
	}
	err = auth.replaceDtagPrefix(newtag, content)
	if err != nil {
		return err
	}
	return nil
}
