package main

import (
	"fmt"

	"github.com/life4/gweb/web"
)

type Python struct {
	pyodide web.Value
	doc     web.Document
	output  web.HTMLElement
}

func (py Python) print(text string, cls string) {
	el := py.doc.CreateElement("div")
	el.Attribute("class").Set("alert alert-" + cls)
	el.SetText(text)
	py.output.Node().AppendChild(el.Node())
}

func (py Python) PrintIn(text string) {
	py.print(text, "secondary")
}

func (py Python) PrintOut(text string) {
	py.print(text, "success")
}

func (py Python) PrintErr(text string) {
	py.print(text, "danger")
}

func (py Python) Run(cmd string) string {
	return py.pyodide.Call("runPython", cmd).String()
}

func (py Python) RunAndPrint(cmd string) {
	py.PrintIn(cmd)
	result := py.Run(cmd)
	py.PrintOut(result)
}

func (py Python) Install(pkg string) bool {
	cmd := fmt.Sprintf("micropip.install('%s')", pkg)
	py.PrintIn(cmd)
	_, fail := py.pyodide.Call("runPython", cmd).Promise().Get()
	if fail.Truthy() {
		py.PrintErr(fail.String())
		return false
	}
	py.PrintOut("True")
	return true
}

func (py Python) Clear() {
	py.output.SetText("")
}

func (py Python) InitMicroPip() bool {
	py.PrintIn("import micropip")
	_, fail := py.pyodide.Call("loadPackage", "micropip").Promise().Get()
	if fail.Truthy() {
		py.PrintErr(fail.String())
		return false
	}
	py.Run("import micropip")
	py.PrintOut("True")
	return true
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

	py.RunAndPrint("print('Hello world!')")
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

	select {}
}
