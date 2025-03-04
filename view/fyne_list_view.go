package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/lackmus/npcgengo/controller"
	"github.com/lackmus/npcgengo/model"
	"github.com/lackmus/npcgengo/shared"
)

// FyneListView is a view that displays an NPC.
type FyneListView struct {
	controller  *controller.NPCListController
	npcs        []model.NPC
	list        *widget.List
	detailLabel *widget.Label
	deleteBtn   *widget.Button
	editBtn     *widget.Button
	createBtn   *widget.Button
	rndmBtn     *widget.Button
	content     fyne.CanvasObject
	window      fyne.Window
	app         fyne.App
	selectedID  string
}

func NewFyneListView(controller *controller.NPCListController) shared.NPCListViewer {
	view := &FyneListView{
		npcs:       []model.NPC{},
		app:        app.New(),
		controller: controller,
	}
	controller.RegisterObserver(view)

	// Create window
	view.window = view.app.NewWindow("NPC Manager")
	// size of the window
	view.window.Resize(fyne.NewSize(1500, 400))

	// Create NPC list on the left
	view.list = widget.NewList(
		func() int { return len(view.npcs) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(view.npcs[id].ShortString())
		},
	)

	// Detail label on the right
	view.detailLabel = widget.NewLabel("Select an NPC")

	// Handle list selection
	view.list.OnSelected = func(id widget.ListItemID) {
		if id >= 0 && id < len(view.npcs) {
			view.selectedID = view.npcs[id].ID
			view.detailLabel.SetText(view.npcs[id].String())
		}
	}

	// "Delete NPC" button
	view.deleteBtn = widget.NewButton("Delete NPC", func() {
		if view.selectedID != "" {
			view.controller.DeleteNPC(view.selectedID)
			view.list.UnselectAll()
			view.selectedID = ""
			view.detailLabel.SetText("Select an NPC") // Clear details
		}
	})

	// "Edit NPC" button
	view.editBtn = widget.NewButton("Edit NPC", func() {
		if view.selectedID != "" {
			selectedNPC, err := controller.GetNpcByID(view.selectedID)
			if err != nil {
				view.detailLabel.SetText("Error: " + err.Error())
				return
			}
			editCtrl := view.controller.InitEditController()
			editView := NewFyneEditView(editCtrl)
			editCtrl.RegisterObserver(editView)
			editCtrl.LoadNPC(selectedNPC)
			editView.Run()
		}
	})

	// "Create NPC" button
	view.createBtn = widget.NewButton("Create NPC", func() {
		editCtrl := view.controller.InitEditController()
		editView := NewFyneEditView(editCtrl)
		editCtrl.RegisterObserver(editView)
		editView.Run()
	})

	// "Random NPC" button
	view.rndmBtn = widget.NewButton("Random NPC", func() {
		controller.CreateRandomNPC()
	})

	// Layout
	buttons := container.NewHBox(view.rndmBtn, view.createBtn, view.editBtn, view.deleteBtn)
	leftPanel := container.NewBorder(nil, buttons, nil, nil, view.list)
	view.window.SetContent(container.NewHSplit(leftPanel, view.detailLabel))
	return view
}

// Update refreshes the NPC list when the controller notifies observers.
func (v *FyneListView) Update(npcs []model.NPC) {
	v.npcs = npcs
	v.list.Refresh()
}

// Content returns the main Fyne UI container.
func (v *FyneListView) Content() fyne.CanvasObject {
	return v.content
}

// Run starts the GUI loop.
func (v *FyneListView) Run() {
	v.window.ShowAndRun()
}
