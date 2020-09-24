package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/life4/gweb/web"
)

const Flake8QuotesURL = "https://github.com/zheller/flake8-quotes/archive/master.zip"

func Download(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	encoded := base64.StdEncoding.EncodeToString(content)
	return encoded, nil
}

func main() {
	window := web.GetWindow()
	doc := window.Document()
	doc.SetTitle("FlakeHell online")

	// load python
	py := Python{doc: doc, output: doc.Element("py-output")}
	py.PrintIn("Load Python")
	window.Get("languagePluginLoader").Promise().Get()
	py.PrintOut("Python is ready")
	py.pyodide = window.Get("pyodide")

	py.RunAndPrint("'Hello world!'")
	ok := py.InitMicroPip()
	if !ok {
		return
	}

	// skip nighty packages
	skip := map[string]string{
		"flake8-quotes":         "2.1.2",
		"flake8-bugbear":        "19.3",
		"flake8-rst-docstrings": "0.0.12",
		"flake8-eradicate":      "0.3.0",
		"flake8-isort":          "3.0.1",
		"flake8-bandit":         "2.1.1",
	}
	for pname, pversion := range skip {
		cmd := "micropip.PACKAGE_MANAGER.installed_packages['%s'] = '%s'"
		py.Run(fmt.Sprintf(cmd, pname, pversion))
	}

	// install dependencies
	py.Clear()
	py.Install("flake8==3.7.9")
	py.Install("setuptools")
	py.Install("entrypoints")
	py.Install("flake8-builtins==1.5.3")
	py.Install("wemake-python-styleguide==0.14.1")

	// install non-wheel dependencies
	archive, err := Download(Flake8QuotesURL)
	if err != nil {
		log.Fatal(err)
	}
	py.Set("archive", archive)
	scripts := NewScripts()
	script := scripts.ReadExtract()
	py.RunAndPrint(script)

	flakehell := NewFlakeHell(doc, &py)
	flakehell.Run()
	flakehell.Register()

	select {}
}
