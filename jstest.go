package main

import (
	"fmt"
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/render"
	"github.com/kennygrant/sanitize"
	"io/ioutil"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"net/http"
)

// Our configuration struct.
type Config struct {
	Path string
}

// A helper function to get the path from the configuration settings.
func (config *Config) GetPath() string {
	return config.Path
}

// A helper function to set the path in the configuration settings.
func (config *Config) SetPath(path string) {
	config.Path = path
}

// The main struct for our app.
type App struct {
	AppName  string
	AppTitle string
}

// A struct to contain our app's code.
type AppCode struct {
	HTMLcode string
	CSScode  string
	JScode   string
}

// Create a handler for the configuration settings.
// TODO: Store the configuration in the database and provide a form to edit the configuration.
func Cfg() martini.Handler {
	return func() {
		config := Config{}
		config.SetPath("/home/mike/go/src/github.com/prinsmike/jstest")
	}
}

// Create a handler for the database connection.
func DB() martini.Handler {
	session, e := mgo.Dial("mongodb://localhost")
	check(e)

	return func(c martini.Context) {
		s := session.Clone()
		c.Map(s.DB("jstest"))
		defer s.Close()
		c.Next()
	}
}

// Get a list of all apps from the database.
func GetAll(db *mgo.Database) []App {
	var appslist []App
	db.C("apps").Find(nil).All(&appslist)
	return appslist
}

// Get the app info from the database.
func GetAppInfo(db *mgo.Database, appname string) App {
	var app App
	db.C("apps").Find(bson.M{"Name": appname}).One(&app)
	return app
}

// Insert a new app into the database.
func InsertApp(db *mgo.Database, appname, title string) {
	app := App{appname, title}
	db.C("apps").Insert(app)
}

// Create a new app.
func CreateApp(db *mgo.Database, cfg *Config, appname, title string) {
	InsertApp(db, appname, title)

	// Create app files.

}

// GetCode fetches a given app's code from its files.
func GetCode(appname string, cfg *Config) AppCode {
	html_code, e := ioutil.ReadFile(fmt.Sprintf("%s/templates/%s/content.tmpl", cfg.GetPath(), appname))
	check(e)
	css_code, e := ioutil.ReadFile(fmt.Sprintf("%s/assets/css/%s/style.css", cfg.GetPath(), appname))
	check(e)
	js_code, e := ioutil.ReadFile(fmt.Sprintf("%s/assets/js/%s/script.js", cfg.GetPath(), appname))
	check(e)

	return AppCode{string(html_code), string(css_code), string(js_code)}
}

// Save the code for an app.
func SaveCode(appname, htmlcode, csscode, jscode string, cfg *Config) {
	if htmlcode != "" {
		StringToFile(htmlcode, fmt.Sprintf("%s/templates/%s/content.tmpl", cfg.GetPath(), appname))
	}
	if csscode != "" {
		StringToFile(csscode, fmt.Sprintf("%s/assets/css/%s/style.css", cfg.GetPath(), appname))
	}
	if jscode != "" {
		StringToFile(jscode, fmt.Sprintf("%s/assets/js/%s/script.js", cfg.GetPath(), appname))
	}
}

// Check if files exist and create them if not.
func CheckFiles(appname string, cfg *Config) {

}

// StringToFile saves a string to a given file.
func StringToFile(str, file string) {
	b := []byte(str)
	e := ioutil.WriteFile(file, b, 0644)
	check(e)
}

// Simple error checking.
func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Simple debugging function.
func debug(str string) {
	fmt.Printf("Debug: %s\n", str)
}

func main() {

	m := martini.Classic()
	m.Use(render.Renderer())
	m.Use(martini.Static("assets"))
	m.Use(DB())

	// Get the list of apps.
	m.Get("/", func(r render.Render, db *mgo.Database) {
		r.HTML(200, "list", nil)
	})

	// Get an app.
	m.Get("/:app", func(r render.Render, db *mgo.Database, cfg *Config, params martini.Params) {
		appinfo := GetAppInfo(db, params["app"])
		appedit := GetCode(params["app"], cfg.Path)
		r.HTML(200, "app", appedit)
	})

	// Post a modified app.
	m.Post("/:app", func(r render.Render, db *mgo.Database, cfg *Config) {

		// Insert or update the app's information in the database.
		// Fields: Title and Name
		title := req.FormValue("Title")
		appname := sanitize.Name(req.FormValue("Title"))

		// Write the app files.
		SaveCode(appname, req.FormValue("HTMLcode"), req.FormValue("CSScode"), req.FormValue("JScode"), cfg)

		appedit := GetCode(appname, cfg)
		r.HTML(200, "app", appedit)
	})

	m.Run()
}
