package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/life4/gweb/web"
)

type FlakeHell struct {
	script string
	input  web.HTMLElement
	btn    web.HTMLElement
	doc    web.Document
	py     *Python
}

type Violation struct {
	Code        string
	Description string
	Context     string
	Line        int
	Column      int
	Plugin      string
}

func NewFlakeHell(doc web.Document, py *Python) FlakeHell {
	scripts := NewScripts()
	script := scripts.ReadFlakeHell()
	return FlakeHell{
		script: script,
		input:  doc.Element("py-code"),
		btn:    doc.Element("py-lint"),
		doc:    doc,
		py:     py,
	}

}

func (fh *FlakeHell) Register() {
	fh.btn.Set("disabled", false)
	fh.btn.EventTarget().Listen(web.EventTypeClick, func(event web.Event) {
		fh.btn.Set("disabled", true)
		fh.Run()
		fh.btn.Set("disabled", false)
	})
}

func (fh *FlakeHell) Run() {
	// fh.py.Clear()
	fh.py.Set("text", fh.input.Text())
	fh.py.RunAndPrint(fh.script)

	fh.py.Clear()
	fh.py.RunAndPrint("code")

	cmd := "'\\n'.join(app.formatter._out)"
	fh.py.PrintIn(cmd)
	result := fh.py.Run(cmd)
	fh.py.PrintOut(result)

	violations := make([]Violation, 0)
	for _, line := range strings.Split(result, "\n") {
		v := Violation{}
		err := json.Unmarshal([]byte(line), &v)
		if err != nil {
			fh.py.PrintErr(err.Error())
			return
		}
		violations = append(violations, v)
	}

	fh.py.Clear()
	fh.table(violations)
}

func (fh *FlakeHell) table(violations []Violation) {
	table := fh.doc.CreateElement("table")
	table.Attribute("class").Set("table table-sm")

	thead := fh.doc.CreateElement("thead")
	table.Node().AppendChild(thead.Node())
	tr := fh.doc.CreateElement("tr")
	thead.Node().AppendChild(tr.Node())

	cols := []string{"plugin", "code", "descr", "pos", "context"}
	for _, name := range cols {
		th := fh.doc.CreateElement("th")
		th.SetText(name)
		tr.Node().AppendChild(th.Node())
	}

	tbody := fh.doc.CreateElement("tbody")
	table.Node().AppendChild(tbody.Node())

	for _, vl := range violations {
		tr := fh.doc.CreateElement("tr")

		td := fh.doc.CreateElement("td")
		td.SetText(vl.Plugin)
		tr.Node().AppendChild(td.Node())

		td = fh.doc.CreateElement("td")
		td.SetText(vl.Code)
		tr.Node().AppendChild(td.Node())

		td = fh.doc.CreateElement("td")
		td.SetText(vl.Description)
		tr.Node().AppendChild(td.Node())

		td = fh.doc.CreateElement("td")
		td.SetText(fmt.Sprintf("%d:%d", vl.Line, vl.Column))
		tr.Node().AppendChild(td.Node())

		td = fh.doc.CreateElement("td")
		td.SetText(vl.Context)
		tr.Node().AppendChild(td.Node())

		tbody.Node().AppendChild(tr.Node())
	}

	fh.py.output.Node().AppendChild(table.Node())
}
