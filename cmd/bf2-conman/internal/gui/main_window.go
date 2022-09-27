package gui

import (
	_ "embed"
	"strconv"

	"github.com/cetteup/conman/cmd/bf2-conman/internal/actions"
	"github.com/cetteup/conman/pkg/game"
	"github.com/cetteup/conman/pkg/handler"
	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
	"github.com/lxn/win"
)

const (
	windowWidth  = 300
	windowHeight = 425
)

type DropDownItem struct { // Used in the ComboBox dropdown
	Key  int
	Name string
}

func CreateMainWindow(h *handler.Handler, profiles []game.Profile, defaultProfileKey string) (*walk.MainWindow, error) {
	icon, err := walk.NewIconFromResourceIdWithSize(2, walk.Size{Width: 256, Height: 256})
	if err != nil {
		return nil, err
	}

	screenWidth := win.GetSystemMetrics(win.SM_CXSCREEN)
	screenHeight := win.GetSystemMetrics(win.SM_CYSCREEN)

	profileOptions, selectedProfile, err := computeProfileSelectOptions(profiles, defaultProfileKey)
	if err != nil {
		return nil, err
	}

	var mw *walk.MainWindow
	var profileSelection *walk.ComboBox
	var passwordGB *walk.GroupBox

	if err := (declarative.MainWindow{
		AssignTo: &mw,
		Title:    "BF2 conman",
		Name:     "BF2 conman",
		Bounds: declarative.Rectangle{
			X:      int((screenWidth - windowWidth) / 2),
			Y:      int((screenHeight - windowHeight) / 2),
			Width:  windowWidth,
			Height: windowHeight,
		},
		Layout:  declarative.VBox{},
		Icon:    icon,
		ToolBar: declarative.ToolBar{},
		Children: []declarative.Widget{
			declarative.Label{
				Text:       "Select profile",
				TextColor:  walk.Color(win.GetSysColor(win.COLOR_CAPTIONTEXT)),
				Background: declarative.SolidColorBrush{Color: walk.Color(win.GetSysColor(win.COLOR_BTNFACE))},
			},
			declarative.ComboBox{
				AssignTo:      &profileSelection,
				Value:         profileOptions[selectedProfile].Key,
				Model:         profileOptions,
				DisplayMember: "Name",
				BindingMember: "Key",
				Name:          "Select profile",
				ToolTipText:   "Select profile",
				OnCurrentIndexChanged: func() {
					// Password actions cannot be used with singleplayer profiles, since those don't have passwords
					if profiles[profileSelection.CurrentIndex()].Type == game.ProfileTypeMultiplayer {
						passwordGB.SetEnabled(true)
					} else {
						passwordGB.SetEnabled(false)
					}
				},
			},
			declarative.GroupBox{
				Title:  "Profile actions",
				Name:   "Profile actions",
				Layout: declarative.VBox{},
				Children: []declarative.Widget{
					declarative.PushButton{
						Text: "Set as default profile",
						OnClicked: func() {
							err := actions.SetDefaultProfile(h, profiles[profileSelection.CurrentIndex()].Key)
							if err != nil {
								walk.MsgBox(mw, "Error", "Failed to set default profile", walk.MsgBoxIconError)
							} else {
								walk.MsgBox(mw, "Success", "Updated default profile", walk.MsgBoxIconInformation)
							}
						},
					},
					declarative.PushButton{
						Text: "Purge server history",
						OnClicked: func() {
							err := actions.PurgeServerHistory(h, profiles[profileSelection.CurrentIndex()].Key)
							if err != nil {
								walk.MsgBox(mw, "Error", "Failed to purge server history", walk.MsgBoxIconError)
							} else {
								walk.MsgBox(mw, "Success", "Purged server history", walk.MsgBoxIconInformation)
							}
						},
					},
					declarative.PushButton{
						Text: "Purge old demo bookmarks",
						OnClicked: func() {
							err := actions.PurgeOldDemoBookmarks(h, profiles[profileSelection.CurrentIndex()].Key)
							if err != nil {
								walk.MsgBox(mw, "Error", "Failed to purge demo bookmarks", walk.MsgBoxIconError)
							} else {
								walk.MsgBox(mw, "Success", "Purged all demo bookmarks older than one week", walk.MsgBoxIconInformation)
							}
						},
					},
					declarative.PushButton{
						Text: "Disable help voice overs",
						OnClicked: func() {
							err := actions.MarkAllVoiceOverHelpAsPlayed(h, profiles[profileSelection.CurrentIndex()].Key)
							if err != nil {
								walk.MsgBox(mw, "Error", "Failed to disable help voice overs", walk.MsgBoxIconError)
							} else {
								walk.MsgBox(mw, "Success", "Disabled help voice overs", walk.MsgBoxIconInformation)
							}
						},
					},
					declarative.GroupBox{
						AssignTo: &passwordGB,
						Title:    "Password (multiplayer profiles only)",
						Name:     "Password",
						Layout:   declarative.HBox{},
						MaxSize:  declarative.Size{Width: windowWidth, Height: 60},
						Children: []declarative.Widget{
							declarative.Composite{
								Layout: declarative.HBox{},
								Children: []declarative.Widget{
									declarative.PushButton{
										Text: "Copy",
										OnClicked: func() {
											profile := profiles[profileSelection.CurrentIndex()]
											password, err := actions.GetProfilePassword(h, profile.Key)
											if err != nil {
												walk.MsgBox(mw, "Error", "Failed get profile password\n\nIf the password was encrypted on another Windows installation or with another user, the password cannot be decrypted by design.\n\nYou can provide the password via the password \"Edit\" functionality or by logging in in game.", walk.MsgBoxIconError)
												return
											}
											if err := walk.Clipboard().SetText(password); err != nil {
												walk.MsgBox(mw, "Error", "Failed to copy password to clipboard", walk.MsgBoxIconError)
												return
											}
											walk.MsgBox(mw, "Success", "Password copied to clipboard", walk.MsgBoxIconInformation)
										},
									},
									declarative.PushButton{
										Text: "Edit",
										OnClicked: func() {
											profile := profiles[profileSelection.CurrentIndex()]
											password, err := actions.GetProfilePassword(h, profile.Key)
											if err != nil {
												walk.MsgBox(mw, "Error", "Failed get profile password\n\nIf the password was encrypted on another Windows installation or with another user, the password cannot be decrypted by design.\n\nYou can still update the password, but ensure you know the current password before continuing.", walk.MsgBoxIconError)
											}
											if cmd, err := RunPasswordEditDialog(mw, h, profile, password); err != nil {
												walk.MsgBox(mw, "Error", "Failed to set password", walk.MsgBoxIconError)
											} else if cmd == walk.DlgCmdOK {
												walk.MsgBox(mw, "Success", "Password updated", walk.MsgBoxIconInformation)
											}
										},
									},
								},
							},
						},
					},
				},
			},
			declarative.GroupBox{
				Title:  "Global actions",
				Name:   "Global actions",
				Layout: declarative.VBox{},
				Children: []declarative.Widget{
					declarative.PushButton{
						Text: "Purge shader cache",
						OnClicked: func() {
							err := actions.PurgeShareCache(h)
							if err != nil {
								walk.MsgBox(mw, "Error", "Failed to purge shader cache", walk.MsgBoxIconError)
							} else {
								walk.MsgBox(mw, "Success", "Purged shader cache", walk.MsgBoxIconInformation)
							}
						},
					},
					declarative.PushButton{
						Text: "Purge logo cache",
						OnClicked: func() {
							err := actions.PurgeLogoCache(h)
							if err != nil {
								walk.MsgBox(mw, "Error", "Failed to purge logo cache", walk.MsgBoxIconError)
							} else {
								walk.MsgBox(mw, "Success", "Purged logo cache", walk.MsgBoxIconInformation)
							}
						},
					},
				},
			},
			declarative.Label{
				Text:       "BF2 conman v0.1.5",
				Alignment:  declarative.AlignHCenterVCenter,
				TextColor:  walk.Color(win.GetSysColor(win.COLOR_GRAYTEXT)),
				Background: declarative.SolidColorBrush{Color: walk.Color(win.GetSysColor(win.COLOR_BTNFACE))},
			},
		},
	}).Create(); err != nil {
		return nil, err
	}

	// Disable minimize/maximize buttons and fix size
	win.SetWindowLong(mw.Handle(), win.GWL_STYLE, win.GetWindowLong(mw.Handle(), win.GWL_STYLE) & ^win.WS_MINIMIZEBOX & ^win.WS_MAXIMIZEBOX & ^win.WS_SIZEBOX)

	return mw, nil
}

func computeProfileSelectOptions(profiles []game.Profile, defaultProfileKey string) ([]DropDownItem, int, error) {
	defaultOption := 0
	options := make([]DropDownItem, 0, len(profiles))
	for i, profile := range profiles {
		key, err := strconv.Atoi(profile.Key)
		if err != nil {
			return nil, 0, err
		}

		if profile.Key == defaultProfileKey {
			defaultOption = i
		}

		options = append(options, DropDownItem{
			Key:  key,
			Name: profile.Name,
		})
	}

	return options, defaultOption, nil
}
