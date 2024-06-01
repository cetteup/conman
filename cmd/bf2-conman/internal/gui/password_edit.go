package gui

import (
	"fmt"

	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
	"github.com/lxn/win"

	"github.com/cetteup/conman/cmd/bf2-conman/internal/actions"
	"github.com/cetteup/conman/pkg/game"
	"github.com/cetteup/conman/pkg/handler"
)

func RunPasswordEditDialog(owner walk.Form, h *handler.Handler, profile game.Profile, currentPassword string) (int, error) {
	var dialog *walk.Dialog
	var passwordTE *walk.TextEdit
	var savePB, cancelPB *walk.PushButton

	return declarative.Dialog{
		AssignTo:      &dialog,
		Title:         fmt.Sprintf("Edit password for %s", profile.Name),
		Name:          fmt.Sprintf("Edit password for %s", profile.Name),
		MinSize:       declarative.Size{Width: windowWidth - 20},
		FixedSize:     true,
		Layout:        declarative.VBox{},
		Icon:          owner.Icon(),
		DefaultButton: &savePB,
		CancelButton:  &cancelPB,
		Children: []declarative.Widget{
			declarative.Label{
				Text:       "Enter/update password",
				TextColor:  walk.Color(win.GetSysColor(win.COLOR_CAPTIONTEXT)),
				Background: declarative.SolidColorBrush{Color: walk.Color(win.GetSysColor(win.COLOR_BTNFACE))},
			},
			declarative.TextEdit{
				AssignTo: &passwordTE,
				Name:     "Password",
				Text:     currentPassword,
			},
			declarative.Composite{
				Layout: declarative.HBox{},
				Children: []declarative.Widget{
					declarative.HSpacer{},
					declarative.PushButton{
						AssignTo: &savePB,
						Text:     "Save",
						OnClicked: func() {
							if err := actions.SetProfilePassword(h, profile.Key, passwordTE.Text()); err != nil {
								return
							}
							dialog.Accept()
						},
					},
					declarative.PushButton{
						AssignTo: &cancelPB,
						Text:     "Cancel",
						OnClicked: func() {
							dialog.Cancel()
						},
					},
				},
			},
		},
	}.Run(owner)
}
