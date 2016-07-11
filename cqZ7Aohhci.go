package main

import (
	&#34;os&#34;

	&#34;github.com/mattn/go-gtk/gtk&#34;
)

var last int

func main() {
	gtk.Init(&amp;os.Args)

	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	window.SetTitle(&#34;Notes&#34;)
	window.SetIconName(&#34;gtk-about&#34;)
	window.Connect(&#34;destroy&#34;, func() {
		gtk.MainQuit()
	})

	notebook := gtk.NewNotebook()
	notebook.AppendPage(gtk.NewVBox(false, 1),
		gtk.NewImageFromStock(gtk.STOCK_ADD, gtk.ICON_SIZE_MENU))

	notebook.Connect(&#34;button-release-event&#34;, func() bool {
		n := notebook.GetCurrentPage()
		if n == last {
			addPage(notebook)
		}
		return false
	})

	window.Add(notebook)
	window.SetSizeRequest(500, 300)
	window.ShowAll()

	gtk.Main()
}

func addPage(notebook *gtk.Notebook) {
	dialog := gtk.NewDialog()
	dialog.SetTitle(&#34;Title?&#34;)
	dVbox := dialog.GetVBox()

	input := gtk.NewEntry()
	input.SetEditable(true)
	vbox := gtk.NewVBox(false, 1)
	input.Connect(&#34;activate&#34;, func() {
		s := input.GetText()
		if s != &#34;&#34; {
			last = notebook.InsertPage(vbox, gtk.NewLabel(s), last) &#43; 1
			notebook.ShowAll()
		}
		notebook.PrevPage()
		dialog.Destroy()
	})

	dVbox.Add(input)
	button := gtk.NewButtonWithLabel(&#34;OK&#34;)
	button.Connect(&#34;clicked&#34;, func() {
		input.Emit(&#34;activate&#34;)
	})
	dVbox.Add(button)
	dialog.SetModal(true)
	dialog.ShowAll()

	hbox := gtk.NewHBox(false, 1)
	swin := gtk.NewScrolledWindow(nil, nil)
	swin.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	textview := gtk.NewTextView()
	swin.Add(textview)

	butBold := gtk.NewToolButtonFromStock(gtk.STOCK_BOLD)
	butIta := gtk.NewToolButtonFromStock(gtk.STOCK_ITALIC)

	butFont := gtk.NewFontButton()
	butFont.Connect(&#34;font-set&#34;, func() {
		textview.ModifyFontEasy(butFont.GetFontName())
	})

	butClose := gtk.NewToolButtonFromStock(gtk.STOCK_DELETE)
	butClose.Connect(&#34;clicked&#34;, func() {
		n := notebook.GetCurrentPage()
		notebook.RemovePage(notebook, n)
		last--
		notebook.PrevPage()
	})

	hbox.PackStart(butBold, false, false, 0)
	hbox.PackStart(butIta, false, false, 0)
	hbox.PackStart(butFont, false, false, 0)
	hbox.PackEnd(butClose, false, false, 0)
	vbox.PackStart(hbox, false, false, 0)
	vbox.Add(swin)
	notebook.ShowAll()
}
