package main

import "github.com/life4/gweb/web"

func main() {
	window := web.GetWindow()
	doc := window.Document()
	doc.SetTitle("FlakeHell online")
	// body := doc.Body()

	// output := doc.Element("py-output")

}
