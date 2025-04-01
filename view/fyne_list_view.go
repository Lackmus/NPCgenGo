// Description: Fyne implementation of the NPC list view.
package view

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/lackmus/npcgengo/controller"
	"github.com/lackmus/npcgengo/model"
	cp "github.com/lackmus/npcgengo/model/npc_components"
	"github.com/lackmus/npcgengo/shared"
)

// rowLayout is a custom layout that positions children in a row.
type rowLayout struct {
	widths []float32
}

// Layout positions each child at the correct X offset. The Y offset is always 0.
func (r *rowLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	xPos := float32(0)
	for i, obj := range objects {
		colWidth := r.widths[i]
		height := obj.MinSize().Height
		obj.Resize(fyne.NewSize(colWidth, height))
		obj.Move(fyne.NewPos(xPos, 0))
		xPos += colWidth
	}
}

// MinSize calculates the total width and max height of this row.
func (r *rowLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	var totalWidth float32
	var maxHeight float32
	for i, obj := range objects {
		w := r.widths[i]
		h := obj.MinSize().Height
		if h > maxHeight {
			maxHeight = h
		}
		totalWidth += w
	}
	return fyne.NewSize(totalWidth, maxHeight)
}

// FyneListView is a view that displays an NPC list and details.
type FyneListView struct {
	controller  *controller.NPCListController
	npcs        []model.NPC
	table       *widget.Table
	detailLabel *fyne.Container
	deleteBtn   *widget.Button
	editBtn     *widget.Button
	createBtn   *widget.Button
	rndmBtn     *widget.Button
	window      fyne.Window
	app         fyne.App
	selectedID  string
}

// NewFyneListView creates and initializes the NPC list view.
func NewFyneListView(controller *controller.NPCListController) shared.NPCListViewer {
	view := &FyneListView{
		npcs:       []model.NPC{},
		app:        app.New(),
		controller: controller,
	}
	controller.RegisterObserver(view)

	// Create window
	view.window = view.app.NewWindow("NPC Manager")
	view.window.Resize(fyne.NewSize(1600, 400))

	// Initialize detail label
	view.detailLabel = container.NewVBox(widget.NewLabel("Select an NPC"))

	// Create header row with fixed column widths.
	// These widths should match the table's column widths.
	headerWidths := []float32{200, 130, 130, 130, 150}
	var headerLabels []fyne.CanvasObject
	// i <6 make bold
	for i := range 5 {
		headerLabels = append(headerLabels, widget.NewLabelWithStyle(cp.CompEnum(i+1).String(), fyne.TextAlignLeading, fyne.TextStyle{Bold: true}))
	}
	header := container.New(&rowLayout{widths: headerWidths}, headerLabels...)

	// Initialize selected row index.
	var selectedRow int = -1

	// Create NPC table.
	view.table = widget.NewTable(
		// Provide the number of rows and columns.
		func() (int, int) { return len(view.npcs), 5 }, // rows, columns
		// Provide the content of the header row.
		func() fyne.CanvasObject {
			return widget.NewLabel("Loading...")
		},
		// Provide the content of each cell.
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			// Make sure we don't panic if table refreshes with no NPCs.
			if id.Row >= len(view.npcs) {
				return
			}

			// Get the NPC at the current row.
			npc := view.npcs[id.Row]
			// Set the text of the cell to the NPC's component value.
			label := cell.(*widget.Label)
			switch id.Col {
			case 0:
				label.SetText(npc.GetComponent(cp.CompName))
			case 1:
				label.SetText(npc.GetComponent(cp.CompType))
			case 2:
				label.SetText(npc.GetComponent(cp.CompSubtype))
			case 3:
				label.SetText(npc.GetComponent(cp.CompSpecies))
			case 4:
				label.SetText(npc.GetComponent(cp.CompFaction))
			}

			// Highlight the selected row (change only text style, no background).
			if id.Row == selectedRow {
				label.TextStyle.Bold = true
			} else {
				label.TextStyle.Bold = false
			}
			label.Refresh() // Apply changes
		},
	)

	// Set fixed column widths for the table (must match header widths).
	view.table.SetColumnWidth(0, headerWidths[0])
	view.table.SetColumnWidth(1, headerWidths[1])
	view.table.SetColumnWidth(2, headerWidths[2])
	view.table.SetColumnWidth(3, headerWidths[3])
	view.table.SetColumnWidth(4, headerWidths[4])

	// Set the table's selection behavior.
	view.table.OnSelected = func(id widget.TableCellID) {
		if id.Row >= 0 && id.Row < len(view.npcs) {
			selectedRow = id.Row // Store the selected row index
			view.selectedID = view.npcs[id.Row].ID
			view.detailLabel.Objects = makeNPCStringFyne(view.npcs[id.Row])
			view.detailLabel.Refresh()
			view.table.Refresh() // Refresh the table to apply highlighting
		}
	}

	// Unselect the row when it is deselected.
	view.table.OnUnselected = func(id widget.TableCellID) {
		selectedRow = -1 // Reset selection when unselected
		view.table.Refresh()
	}

	// Create buttons.
	view.deleteBtn = widget.NewButton("Delete NPC", func() {
		if view.selectedID != "" {
			view.controller.DeleteNPC(view.selectedID)
			view.table.UnselectAll()
			view.selectedID = ""
			view.detailLabel.Objects = []fyne.CanvasObject{}
		}
	})

	view.editBtn = widget.NewButton("Edit NPC", func() {
		if view.selectedID != "" {
			selectedNPC, err := controller.GetNpcByID(view.selectedID)
			if err != nil {
				view.detailLabel.Objects = []fyne.CanvasObject{widget.NewLabel("Error: NPC not found")}
				view.detailLabel.Refresh()
				return
			}
			editCtrl := view.controller.InitEditController()
			editView := NewFyneEditView(editCtrl)
			editCtrl.RegisterObserver(editView)
			editCtrl.LoadNPC(selectedNPC)
			editView.Run()
		}
	})

	view.createBtn = widget.NewButton("Create NPC", func() {
		editCtrl := view.controller.InitEditController()
		editView := NewFyneEditView(editCtrl)
		editCtrl.RegisterObserver(editView)
		editView.Run()
	})

	view.rndmBtn = widget.NewButton("Random NPC", func() {
		controller.CreateRandomNPC()
	})

	// Layout for buttons.
	buttons := container.NewHBox(view.rndmBtn, view.createBtn, view.editBtn, view.deleteBtn)

	// Assemble the left panel with fixed header on top, buttons at bottom, and table in the center.
	leftPanel := container.NewBorder(header, buttons, nil, nil, view.table)

	// Combine left panel and detail label in a horizontal split.
	view.window.SetContent(container.NewHSplit(leftPanel, view.detailLabel))

	return view
}

// Update refreshes the NPC list when notified by the controller.
func (v *FyneListView) Update(npcs []model.NPC) {
	v.npcs = npcs
	v.table.Refresh()
}

// Content returns the main Fyne UI container.
func (v *FyneListView) Content() fyne.CanvasObject {
	return v.window.Content()
}

// Run starts the GUI loop.
func (v *FyneListView) Run() {
	v.window.ShowAndRun()
}

// makeNPCStringFyne returns a slice of widget.Label to be used in a container for display.
func makeNPCStringFyne(npc model.NPC) []fyne.CanvasObject {
	var labels []fyne.CanvasObject
	for i := range cp.CompEnumValues() {
		c := cp.CompEnum(i + 1)
		if comp, ok := npc.Components[c]; ok {
			// Create a label for the component name (bold)
			compNameLabel := widget.NewLabelWithStyle(fmt.Sprintf("%s:", c), fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
			// Create a label for the component value (normal)
			compValueLabel := widget.NewLabel(comp)
			// Add the labels to the list of labels
			labels = append(labels, compNameLabel, compValueLabel)
		}
	}
	return labels
}
