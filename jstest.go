package main

import (
	"fmt"
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/render"
	"io/ioutil"
	"net/http"
	///"labix.org/v2/mgo"
)

type App struct {
	Title    string
	HTMLcode string
	CSScode  string
	JScode   string
}

// func DB() martini.Handler {
// 	session, err := mgo.Dial
// }

func debug(str string) {
	fmt.Printf("Debug: %s\n", str)
}

func getCode() App {
	html_code, e := ioutil.ReadFile("/home/mike/go/src/github.com/prinsmike/jstest/templates/content.tmpl")
	check(e)
	css_code, e := ioutil.ReadFile("/home/mike/go/src/github.com/prinsmike/jstest/assets/css/style.css")
	check(e)
	js_code, e := ioutil.ReadFile("/home/mike/go/src/github.com/prinsmike/jstest/assets/js/script.js")
	check(e)

	return App{"JSTest Application", string(html_code), string(css_code), string(js_code)}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func StringToFile(str, file string) {
	b := []byte(str)
	e := ioutil.WriteFile(file, b, 0644)
	check(e)
}

func main() {
	m := martini.Classic()
	m.Use(render.Renderer())
	m.Use(martini.Static("assets"))

	m.Get("/", func(r render.Render) {
		r.HTML(200, "app", nil)
	})

	m.Get("/edit", func(r render.Render) {
		appedit := getCode()
		r.HTML(200, "edit", appedit)
	})

	m.Post("/edit", func(req *http.Request, r render.Render) {

		debug(req.FormValue("HTMLcode"))
		StringToFile(req.FormValue("HTMLcode"), "/home/mike/go/src/github.com/prinsmike/jstest/templates/content.tmpl")

		debug(req.FormValue("CSScode"))
		StringToFile(req.FormValue("CSScode"), "/home/mike/go/src/github.com/prinsmike/jstest/assets/css/style.css")

		debug(req.FormValue("JScode"))
		StringToFile(req.FormValue("JScode"), "/home/mike/go/src/github.com/prinsmike/jstest/assets/js/script.js")

		apppost := getCode()
		r.HTML(200, "edit", apppost)
	})

	m.Run()
}
