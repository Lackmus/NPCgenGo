package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/lackmus/npcgengo/controller"
	"github.com/lackmus/npcgengo/model"
	cp "github.com/lackmus/npcgengo/model/npc_components"
	"github.com/lackmus/npcgengo/shared"
)

type FyneEditView struct {
	editCtrl    *controller.NPCEditController
	window      fyne.Window
	nameEntry   *widget.Entry
	typeSelect  *widget.Select
	subtypeSel  *widget.Select
	speciesSel  *widget.Select
	factionSel  *widget.Select
	traitSel    *widget.Select
	statsEntry  *widget.Entry
	itemsEntry  *widget.Entry
	description *widget.Entry
	saveBtn     *widget.Button
	cancelBtn   *widget.Button
	rndmNameBtn *widget.Button
	statsBtn    *widget.Button
	itemsBtn    *widget.Button
	descBtn     *widget.Button
}

func NewFyneEditView(editCtrl *controller.NPCEditController) shared.NPCEditViewer {
	view := &FyneEditView{
		editCtrl:    editCtrl,
		window:      fyne.CurrentApp().NewWindow("Edit NPC"),
		nameEntry:   widget.NewEntry(),
		statsEntry:  widget.NewMultiLineEntry(),
		itemsEntry:  widget.NewMultiLineEntry(),
		description: widget.NewMultiLineEntry(),
	}

	view.typeSelect = widget.NewSelect(editCtrl.GetFieldOptions(cp.CompType), func(selected string) {
		editCtrl.SaveField(cp.CompType, selected)

		newOptions := editCtrl.GetFieldOptions(cp.CompSubtype)
		view.subtypeSel.Options = newOptions

		if len(newOptions) > 0 {
			view.subtypeSel.SetSelected(newOptions[0])
		}
		view.subtypeSel.Refresh()
	})

	view.subtypeSel = widget.NewSelect([]string{}, func(selected string) {
		editCtrl.SaveField(cp.CompSubtype, selected)
	})

	view.speciesSel = widget.NewSelect(editCtrl.GetFieldOptions(cp.CompSpecies), func(selected string) {
		editCtrl.SaveField(cp.CompSpecies, selected)
	})

	view.factionSel = widget.NewSelect(editCtrl.GetFieldOptions(cp.CompFaction), func(selected string) {
		editCtrl.SaveField(cp.CompFaction, selected)
	})

	view.traitSel = widget.NewSelect(editCtrl.GetFieldOptions(cp.CompTrait), func(selected string) {
		editCtrl.SaveField(cp.CompTrait, selected)
	})

	view.saveBtn = widget.NewButton("Save", func() {
		editCtrl.SaveNPC()
		view.window.Close()
	})

	view.cancelBtn = widget.NewButton("Cancel", func() {
		view.window.Close()
	})

	view.rndmNameBtn = widget.NewButton("Random Name", func() {
		if view.speciesSel.Selected != "" {
			view.nameEntry.SetText(editCtrl.RandomizeField(cp.CompName))
		}
	})

	view.statsBtn = widget.NewButton("Random Stats", func() {
		if view.subtypeSel.Selected != "" {
			view.statsEntry.SetText(editCtrl.RandomizeField(cp.CompStats))
		}
	})

	view.itemsBtn = widget.NewButton("Random Items", func() {
		if view.subtypeSel.Selected != "" {
			view.itemsEntry.SetText(editCtrl.RandomizeField(cp.CompItems))
		}
	})

	view.descBtn = widget.NewButton("Random Description", func() {
		if view.subtypeSel.Selected != "" {
			view.description.SetText(editCtrl.RandomizeField(cp.CompDescription))
		}
	})
	formItems := []fyne.CanvasObject{}
	selections := []fyne.CanvasObject{
		view.nameEntry,
		view.typeSelect,
		view.subtypeSel,
		view.speciesSel,
		view.factionSel,
		view.traitSel,
		view.statsEntry,
		view.itemsEntry,
		view.description,
	}
	for i := range 9 {
		label := widget.NewLabelWithStyle(cp.CompEnum(i+1).String(), fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
		formItems = append(formItems, label, selections[i])
	}

	formItems = append(formItems, container.NewHBox(view.saveBtn, view.cancelBtn, view.rndmNameBtn, view.statsBtn, view.itemsBtn, view.descBtn))

	form := container.NewVBox(formItems...)

	view.window.SetContent(form)
	view.window.Show()
	return view
}

func (v *FyneEditView) UpdateNPC(npc model.NPC) {
	v.nameEntry.SetText(npc.GetComponent(cp.CompName))
	v.typeSelect.SetSelected(npc.GetComponent(cp.CompType))
	v.subtypeSel.SetSelected(npc.GetComponent(cp.CompSubtype))
	v.speciesSel.SetSelected(npc.GetComponent(cp.CompSpecies))
	v.factionSel.SetSelected(npc.GetComponent(cp.CompFaction))
	v.traitSel.SetSelected(npc.GetComponent(cp.CompTrait))
	v.statsEntry.SetText(npc.GetComponent(cp.CompStats))
	v.itemsEntry.SetText(npc.GetComponent(cp.CompItems))
	v.description.SetText(npc.GetComponent(cp.CompDescription))
}

// UpdateField updates a field in the view.
func (v *FyneEditView) UpdateField(field cp.CompEnum, value interface{}) {
	switch field {
	case cp.CompType:
		v.typeSelect.SetSelected(value.(string))
	}
}

func (v *FyneEditView) OnNPCEditError(err error) {
	widget.NewLabel(err.Error())
}

func (v *FyneEditView) Run() {
	v.window.Show()
}
