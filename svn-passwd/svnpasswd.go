package svnpasswd

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

//User ...
type User struct {
	//Username ...
	Username string `json:"username"`
	//Password ...
	Password string `json:"password"`
	//SVNPATH ...
	SVNPATH string `json:"svnpath"`
}

//NewUser ...
func NewUser(path string) *User {
	user := new(User)
	user.SVNPATH = path
	return user
}

//Readfill will read the passwd to the buffer
func (user *User) readfile() ([]byte, error) {
	authz := filepath.Join(user.SVNPATH, "conf", "passwd")
	filehandle, err := os.Open(authz)
	defer filehandle.Close()
	if err != nil {
		return nil, errors.New(`Please check your ./conf/passwd file,if you change it's name,rename it as "passwd".`)
	}

	content, err := ioutil.ReadAll(filehandle)
	if err != nil {
		return nil, err
	}
	return content, nil
}

//file2json parse passwd content to readable json .
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

//Attention: keep the user as a primary key .
func (user *User) changeUser(tagname string, olddata map[string]string, newdata map[string]string, srcContent []byte) error {
	passwdPath := filepath.Join(user.SVNPATH, "conf", "passwd")
	var lines []string
	var dstContent bytes.Buffer
	lines = strings.Split(string(srcContent), "\n")
	// O(n) loop
	for k, v := range lines {
		v = strings.TrimPrefix(v, "[")
		v = strings.TrimSuffix(v, "]")
		if strings.TrimSpace(v) == tagname {
			dstContent.WriteString(lines[k] + "\n")
			for i := k + 1; i < len(lines); i++ {
				if strings.HasPrefix(lines[i], "[") {
					break
				}
				tsplit := strings.SplitN(lines[i], "=", 2)
				//find the username ;  not  change username
				if strings.TrimSpace(tsplit[0]) == olddata["username"] {
					// replace rudely
					lines[i] = newdata["username"] + "=" + newdata["pwd"]
				}
			}
		} else {
			dstContent.WriteString(v + "\n")
		}
		err := ioutil.WriteFile(passwdPath, dstContent.Bytes(), 0666)
		if err != nil {
			return err
		}
	}
	return nil
}

//Attention: keep the user as a primary key .
func (user *User) delUser(tagname string, olddata map[string]string, srcContent []byte) error {
	passwdPath := filepath.Join(user.SVNPATH, "conf", "passwd")
	var lines []string
	var dstContent bytes.Buffer
	lines = strings.Split(string(srcContent), "\n")
	// O(n) loop
	for k, v := range lines {
		v = strings.TrimPrefix(v, "[")
		v = strings.TrimSuffix(v, "]")
		if strings.TrimSpace(v) == tagname {
			dstContent.WriteString(lines[k] + "\n")
			for i := k + 1; i < len(lines); i++ {
				if strings.HasPrefix(lines[i], "[") {
					break
				}
				tsplit := strings.SplitN(lines[i], "=", 2)
				//find the username ;  not  change username
				if strings.TrimSpace(tsplit[0]) == olddata["username"] {
					// not really delete ,avoid the risk that lose a user key-pair permanently
					lines[i] = "#" + olddata["username"] + "=" + olddata["pwd"]
				}
			}
		} else {
			dstContent.WriteString(v + "\n")
		}
		err := ioutil.WriteFile(passwdPath, dstContent.Bytes(), 0666)
		if err != nil {
			return err
		}
	}
	return nil
}

func (user *User) addUser(tagname string, newdata map[string]string, srcContent []byte) error {
	passwdPath := filepath.Join(user.SVNPATH, "conf", "passwd")
	var lines []string
	var dstContent bytes.Buffer
	lines = strings.Split(string(srcContent), "\n")
	// O(n) loop
	for k, v := range lines {
		v = strings.TrimPrefix(v, "[")
		v = strings.TrimSuffix(v, "]")
		if strings.TrimSpace(v) == tagname {
			dstContent.WriteString(lines[k] + "\n")
			for i := k + 1; i < len(lines); i++ {
				if strings.HasPrefix(lines[i], "[") {
					break
				}
				tsplit := strings.SplitN(lines[i], "=", 2)
				//find the username ;  not  change username
				if strings.TrimSpace(tsplit[0]) == newdata["username"] {
					//return a existed error
					return errors.New("user existed")
				}
			}
			dstContent.WriteString(newdata["username"] + "=" + newdata["pwd"] + "\n")
		} else {
			dstContent.WriteString(v + "\n")
		}
		err := ioutil.WriteFile(passwdPath, dstContent.Bytes(), 0666)
		if err != nil {
			return err
		}
	}
	return nil
}

//parse users json
func parseUsers(rawjson []byte) ([]byte, error) {
	passwd := make(map[string][]string)
	err := json.Unmarshal(rawjson, &passwd)
	if err != nil {
		return []byte(""), err
	}
	userspasswd := make(map[string]string)
	users := passwd["users"]
	for _, v := range users {
		trimv := strings.TrimSpace(string(v))
		splitv := strings.SplitN(trimv, "=", 2)
		userspasswd[strings.TrimSpace(splitv[0])] = strings.TrimSpace(splitv[1])
	}
	passwdjson, err := json.Marshal(userspasswd)
	if err != nil {
		return []byte(""), err
	}
	return passwdjson, nil
}

//ExportUser ...
func (user *User) ExportUser() ([]byte, error) {
	content, err := user.readfile()
	if err != nil {
		return []byte(""), err
	}
	rawjson, err := file2json(content)
	if err != nil {
		return []byte(""), err
	}
	users, err := parseUsers(rawjson)
	if err != nil {
		return []byte(""), err
	}
	return users, nil
}

//ExportEditedUser export edited user data
func (user *User) ExportEditedUser(tag string, olddata map[string]string, newdata map[string]string) ([]byte, error) {
	content, err := user.readfile()
	if err != nil {
		return []byte(""), err
	}
	err = user.changeUser(tag, olddata, newdata, content)
	if err != nil {
		return []byte(""), err
	}
	content, err = user.readfile()
	if err != nil {
		return []byte(""), err
	}
	rawjson, err := file2json(content)
	if err != nil {
		return []byte(""), err
	}
	users, err := parseUsers(rawjson)
	if err != nil {
		return []byte(""), err
	}
	return users, nil
}

//ExportAfterDelUser export all the user key-pair after a "delete " operation
func (user *User) ExportAfterDelUser(tag string, olddata map[string]string) ([]byte, error) {
	content, err := user.readfile()
	if err != nil {
		return []byte(""), err
	}
	err = user.delUser(tag, olddata, content)
	if err != nil {
		return []byte(""), err
	}
	content, err = user.readfile()
	if err != nil {
		return []byte(""), err
	}
	rawjson, err := file2json(content)
	if err != nil {
		return []byte(""), err
	}
	users, err := parseUsers(rawjson)
	if err != nil {
		return []byte(""), err
	}
	return users, nil
}

//ExportAddUser export user key-pair after add operation
func (user *User) ExportAddUser(tag string, newdata map[string]string) ([]byte, error) {
	content, err := user.readfile()
	if err != nil {
		return []byte(""), err
	}
	err = user.addUser(tag, newdata, content)
	if err != nil {
		return []byte(""), err
	}
	content, err = user.readfile()
	if err != nil {
		return []byte(""), err
	}
	rawjson, err := file2json(content)
	if err != nil {
		return []byte(""), err
	}
	users, err := parseUsers(rawjson)
	if err != nil {
		return []byte(""), err
	}
	return users, nil
}
