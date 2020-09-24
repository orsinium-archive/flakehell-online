package main

import "github.com/life4/gweb/web"

type FlakeHell struct {
	script string
	input  web.HTMLElement
	btn    web.HTMLElement
	py     *Python
}

func NewFlakeHell(doc web.Document, py *Python) FlakeHell {
	scripts := NewScripts()
	script := scripts.ReadFlakeHell()
	return FlakeHell{
		script: script,
		input:  doc.Element("py-code"),
		btn:    doc.Element("py-lint"),
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
	fh.py.Clear()
	fh.py.Set("text", fh.input.Text())
	fh.py.RunAndPrint(fh.script)

	fh.py.Clear()
	fh.py.RunAndPrint("code")
	fh.py.RunAndPrint("'\\n'.join(app.formatter._out)")
}
