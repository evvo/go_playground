package main

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"log"
)

import "sort"

// IDs to access the tree view columns by
const (
	COLUMN_NAME = iota
	COLUMN_LINK
)

type Row struct {
	name string
	link string
}

type Listing []Row

func (v Listing) Len() int           { return len(v) }
func (v Listing) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v Listing) Less(i, j int) bool { return v[i].name < v[j].name }

// Add a column to the tree view (during the initialization of the tree view)
func createColumn(title string, id int) *gtk.TreeViewColumn {
	cellRenderer, err := gtk.CellRendererTextNew()
	if err != nil {
		log.Fatal("Unable to create text cell renderer:", err)
	}

	column, err := gtk.TreeViewColumnNewWithAttribute(title, cellRenderer, "text", id)
	if err != nil {
		log.Fatal("Unable to create cell column:", err)
	}

	return column
}

// Creates a tree view and the list store that holds its data
func setupTreeView() (*gtk.TreeView, *gtk.ListStore) {
	treeView, err := gtk.TreeViewNew()
	if err != nil {
		log.Fatal("Unable to create tree view:", err)
	}

	treeView.AppendColumn(createColumn("Name", COLUMN_NAME))
	treeView.AppendColumn(createColumn("Link", COLUMN_LINK))

	// Creating a list store. This is what holds the data that will be shown on our tree view.
	listStore, err := gtk.ListStoreNew(glib.TYPE_STRING, glib.TYPE_STRING)
	if err != nil {
		log.Fatal("Unable to create list store:", err)
	}
	treeView.SetModel(listStore)

	return treeView, listStore
}

// Append a row to the list store for the tree view
func addRow(listStore *gtk.ListStore, row Row) *gtk.TreeIter {
	// Get an iterator for a new row at the end of the list store
	iter := listStore.Append()

	// Set the contents of the list store row that the iterator represents
	err := listStore.Set(iter,
		[]int{COLUMN_NAME, COLUMN_LINK},
		[]interface{}{row.name, row.link})

	if err != nil {
		log.Fatal("Unable to add row:", err)
	}

	return iter
}

// Create and initialize the window
func setupWindow(title string) *gtk.Window {
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}

	win.SetTitle(title)
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})
	win.SetPosition(gtk.WIN_POS_CENTER)
	win.SetDefaultSize(600, 300)
	return win
}

func main() {
	gtk.Init(nil)

	win := setupWindow("Go Feature Timeline")

	treeView, listStore := setupTreeView()
	treeView.Connect("button-press-event", func() {
		log.Print("Button Was Clicked")
	})
	win.Add(treeView)

	// Add some rows to the list store
	list := Listing{
		{"test", "http://sometestwebsite.com"},
		{"test2", "http://anothertestwebsite.com"},
	}

	sort.Sort(sort.Reverse(Listing(list)))

	for _, row := range list {
		addRow(listStore, row)
	}

	win.ShowAll()
	gtk.Main()
}
