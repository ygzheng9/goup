package main

// 20210327 go mod init tidy
// 论语记录 只用到一张表，两个页面，挺 cool

import (
	"bytes"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"

	"pickup/lunyu"
)

var Db *sql.DB

const (
	PORT       = ":8088"
	VIEWS_ROOT = "./views/"
)

type Cite struct {
	ID      int64
	Seq     string
	Content string
	Tags    string
	Refs    string
}

func getAll(cond string) ([]*Cite, error) {
	cmd := "select id, SEQ, CTENT, TAGS, REFS from cite "
	if len(cond) > 0 {
		cmd = cmd + cond
	}
	log.Printf("cmd: %s\n", cmd)

	stmt, err := Db.Prepare(cmd)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}

	var results []*Cite
	for rows.Next() {
		c := Cite{}

		err = rows.Scan(&c.ID, &c.Seq, &c.Content, &c.Tags, &c.Refs)
		if err != nil {
			return nil, err
		}

		results = append(results, &c)
	}

	//log.Printf("Total get: %d\n", len(results))
	return results, nil
}

func findBySeq(seq string) ([]*Cite, error) {
	cmd := "select id, SEQ, CTENT, TAGS, REFS from cite where seq = ?"

	stmt, err := Db.Prepare(cmd)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(seq)
	if err != nil {
		return nil, err
	}

	var results []*Cite
	for rows.Next() {
		c := Cite{}

		err = rows.Scan(&c.ID, &c.Seq, &c.Content, &c.Tags, &c.Refs)
		if err != nil {
			return nil, err
		}

		results = append(results, &c)
	}

	//log.Printf("Total get: %d\n", len(results))
	return results, nil
}

func findByID(id int) (*Cite, error) {
	stmt, err := Db.Prepare("select id, SEQ, CTENT, TAGS, REFS from cite where ID = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	c := Cite{}

	err = stmt.QueryRow(id).Scan(&c.ID, &c.Seq, &c.Content, &c.Tags, &c.Refs)
	if err != nil {
		return nil, err
	}

	//log.Printf("Get by Id: %d\n", c.ID)
	return &c, nil
}

// 根据 id，取得 cite，以及该 cite 的所有 refs
func getBundle(id int) ([]*Cite, error) {
	c, err := findByID(id)
	if err != nil {
		return nil, err
	}

	refs := strings.Split(c.Refs, " ")
	cond := bytes.Buffer{}
	for _, v := range refs {
		fmt.Fprintf(&cond, ", %q ", v)
	}

	s := " where seq in (\"\"" + cond.String() + ")"

	cites, err := getAll(s)
	if err != nil {
		return nil, err
	}
	cites = append(cites, c)
	return cites, nil
}

func (c *Cite) save() (err error) {
	stmt, err := Db.Prepare("insert into cite (SEQ, CTENT, TAGS, REFS) values (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(c.Seq, c.Content, c.Tags, c.Refs)
	if err != nil {
		return err
	}

	c.ID, err = result.LastInsertId()
	if err != nil {
		return err
	}

	log.Printf("Insert successfully. %+v\n", c)
	return nil
}

func (c *Cite) update() (err error) {
	stmt, err := Db.Prepare("update cite set SEQ = ?, CTENT = ?, TAGS = ?, REFS = ? where id =	?")
	if err != nil {
		return nil
	}
	defer stmt.Close()

	_, err = stmt.Exec(c.Seq, c.Content, c.Tags, c.Refs, c.ID)
	if err != nil {
		return nil
	}

	log.Printf("Update successfully. %+v\n", c)
	return nil
}

// 当前 cite 在保存前，对于 refs 需要两边同时增加引用；
// 然后，当前版本和历史版本作比对，对于 去掉的 refs，需要更新 对应的 refs
func (c *Cite) updateRelatedRefs() {
	c.addRefs()

	removed := c.getRemovedRefs()

	//log.Printf("%s removed: %s", c.Seq, removed)

	c.removeRefs(removed)
}

// 修改时，去掉了哪些 refs
// 新增/修改时，增加了哪些 refs
func (c *Cite) getRemovedRefs() []string {
	currRefs := strings.Split(c.Refs, " ")

	oldRefs := []string{}
	if c.ID > 0 {
		oldValue, err := findByID(int(c.ID))
		if err != nil {
			log.Println(err)
			return nil
		}
		oldRefs = strings.Split(oldValue.Refs, " ")
	}

	return lunyu.Minus(oldRefs, currRefs)
}

// 文本替换，去掉连续的两个空格，以及开头的空格
func (c *Cite) removeRefs(refs []string) {
	for _, v := range refs {
		cites, err := findBySeq(v)
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, r := range cites {
			s := strings.Replace(r.Refs, c.Seq, " ", -1)
			s = strings.Replace(s, "  ", " ", -1)
			s = strings.TrimSpace(s)
			r.Refs = s
			r.update()
		}
	}
}

// 拼接字符串
func (c *Cite) addRefs() {
	refs := strings.Split(c.Refs, " ")
	for _, v := range refs {
		cites, err := findBySeq(v)
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, a := range cites {
			if strings.Contains(a.Refs, c.Seq) {
				continue
			} else {
				if len(a.Refs) == 0 {
					a.Refs = c.Seq
				} else {
					a.Refs = a.Refs + " " + c.Seq
				}
				a.update()
			}
		}
	}
}

func (c *Cite) deleteMe() (err error) {
	stmt, err := Db.Prepare("delete from cite WHERE id = ?")
	if err != nil {
		return nil
	}
	defer stmt.Close()

	_, err = stmt.Exec(c.ID)
	if err != nil {
		return nil
	}

	//log.Printf("Delete successfully. %+v\n", c)
	return nil
}

func init() {
	var err error
	Db, err = sql.Open("mysql", "root:mysql@/world?parseTime=true")
	if err != nil {
		panic(err)
	}

	// Test the connection to the database
	err = Db.Ping()
	if err != nil {
		panic(err.Error())
	}
}

func staticDirHandler(mux *http.ServeMux, prefix string, staticDir string) {
	mux.HandleFunc(prefix, func(w http.ResponseWriter, r *http.Request) {
		file := staticDir + r.URL.Path[len(prefix)-1:]

		http.ServeFile(w, r, file)
	})
}

func blankCite(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(VIEWS_ROOT + "cite_form.html")
	if err != nil {
		log.Fatal(err)
	}

	c := Cite{}

	t.Execute(w, c)
}

func createOrUpdateCite(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	log.Printf("%+v\n", r.Form)

	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		log.Println(err)
		return
	}

	var c Cite
	c.ID = int64(id)
	c.Seq = r.FormValue("seq")
	c.Content = r.FormValue("content")
	c.Tags = r.FormValue("tags")
	c.Refs = r.FormValue("refs")

	c.updateRelatedRefs()

	if c.ID == 0 {
		if len(c.Content) > 0 {
			c.save()
		}
	} else {
		c.update()
	}

	http.Redirect(w, r, "/list", http.StatusMovedPermanently)
}

func split(s string) []string {
	return strings.Split(s, " ")
}

func listCite(w http.ResponseWriter, r *http.Request) {
	all, err := getAll("")
	if err != nil {
		log.Println(err)
	}

	// tags 是 空格 分割的多个值
	// 过滤条件是：同时包含多个值
	tags := strings.TrimSpace(r.FormValue("tags"))
	cites := make([]*Cite, 0)

	if len(tags) > 0 {
		tagList := strings.Split(tags, " ")
		for _, v := range all {
			ok := false

			for _, t := range tagList {
				ok = strings.Contains(v.Tags, t)
				if !ok {
					break
				}
			}

			if ok {
				cites = append(cites, v)
			}
		}
	} else {
		cites = all
	}

	params := make(map[string]interface{})
	params["cites"] = cites
	params["tags"] = tags

	renderList(w, params)
}

func renderList(w http.ResponseWriter, params map[string]interface{}) {
	funcMap := template.FuncMap{"split": split}

	// New 返回的是一个 template group
	t := template.New("anyName").Funcs(funcMap)

	_, err := t.ParseFiles(VIEWS_ROOT + "cite_list.html")
	if err != nil {
		log.Fatal(err)
	}

	// 由于 group 中有多个 template，所以需要指定 具体的名字；
	// 具体的名字，就是第一个文件的 basename
	t.ExecuteTemplate(w, "cite_list.html", params)
}

func findRefs(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		log.Println(err)
		return
	}

	cites, err := getBundle(id)
	if err != nil {
		log.Println(err)
		return
	}

	params := make(map[string]interface{})
	params["cites"] = cites
	params["tags"] = ""

	renderList(w, params)
}

func editCite(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		log.Println(err)
		return
	}

	c, err := findByID(id)
	if err != nil {
		log.Println(err)
		return
	}

	t, err := template.ParseFiles(VIEWS_ROOT + "cite_form.html")
	if err != nil {
		log.Fatal(err)
	}

	t.Execute(w, c)
}

func main() {
	startServer()
}

func startServer() {
	mux := http.NewServeMux()
	staticDirHandler(mux, "/public/", "./public")

	mux.HandleFunc("/new", blankCite)
	mux.HandleFunc("/createCite", createOrUpdateCite)
	mux.HandleFunc("/list", listCite)
	mux.HandleFunc("/edit", editCite)
	mux.HandleFunc("/findRefs", findRefs)

	log.Printf("Listening at: %s.\n", PORT)

	err := http.ListenAndServe(PORT, mux)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}
}
