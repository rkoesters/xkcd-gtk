package main

import (
	"github.com/gotk3/gotk3/gtk"
	"log"
)

const shortcutsWindowUI = `<?xml version="1.0" encoding="UTF-8"?>
<interface>
  <object class="GtkShortcutsWindow" id="shortcuts-window">
    <property name="title">Keyboard Shortcuts</property>
    <child>
      <object class="GtkShortcutsSection">
        <property name="section-name">shortcuts</property>
        <property name="visible">1</property>
        <child>
          <object class="GtkShortcutsGroup">
            <property name="title">Comic Navigation</property>
            <property name="visible">1</property>
            <child>
              <object class="GtkShortcutsShortcut">
                <property name="title">Go to the first comic</property>
                <property name="accelerator">&lt;ctrl&gt;Home</property>
                <property name="visible">1</property>
              </object>
            </child>
            <child>
              <object class="GtkShortcutsShortcut">
                <property name="title">Go to the previous comic</property>
                <property name="accelerator">&lt;ctrl&gt;Left</property>
                <property name="visible">1</property>
              </object>
            </child>
            <child>
              <object class="GtkShortcutsShortcut">
                <property name="title">Go to the next comic</property>
                <property name="accelerator">&lt;ctrl&gt;Right</property>
                <property name="visible">1</property>
              </object>
            </child>
            <child>
              <object class="GtkShortcutsShortcut">
                <property name="title">Go to the newest comic</property>
                <property name="accelerator">&lt;ctrl&gt;End</property>
                <property name="visible">1</property>
              </object>
            </child>
            <child>
              <object class="GtkShortcutsShortcut">
                <property name="title">Go to a random comic</property>
                <property name="accelerator">&lt;ctrl&gt;r</property>
                <property name="visible">1</property>
              </object>
            </child>
            <child>
              <object class="GtkShortcutsShortcut">
                <property name="title">Search comics</property>
                <property name="accelerator">&lt;ctrl&gt;f</property>
                <property name="visible">1</property>
              </object>
            </child>
          </object>
        </child>
        <child>
          <object class="GtkShortcutsGroup">
            <property name="title">Comic Viewing</property>
            <property name="visible">1</property>
            <child>
              <object class="GtkShortcutsShortcut">
                <property name="title">Show comic properties</property>
                <property name="accelerator">&lt;ctrl&gt;p</property>
                <property name="visible">1</property>
              </object>
            </child>
          </object>
        </child>
        <child>
          <object class="GtkShortcutsGroup">
            <property name="title">Application Actions</property>
            <property name="visible">1</property>
            <child>
              <object class="GtkShortcutsShortcut">
                <property name="title">Open new window</property>
                <property name="accelerator">&lt;ctrl&gt;n</property>
                <property name="visible">1</property>
              </object>
            </child>
            <child>
              <object class="GtkShortcutsShortcut">
                <property name="title">Toggle dark mode</property>
                <property name="accelerator">&lt;ctrl&gt;d</property>
                <property name="visible">1</property>
              </object>
            </child>
            <child>
              <object class="GtkShortcutsShortcut">
                <property name="title">Open shortcuts window</property>
                <property name="accelerator">&lt;ctrl&gt;question</property>
                <property name="visible">1</property>
              </object>
            </child>
            <child>
              <object class="GtkShortcutsShortcut">
                <property name="title">Quit application</property>
                <property name="accelerator">&lt;ctrl&gt;q</property>
                <property name="visible">1</property>
              </object>
            </child>
          </object>
        </child>
      </object>
    </child>
  </object>
</interface>
`

// NewShortcutsWindow creates a gtk.ShortcutsWindow populated with our
// application's keyboard shortcuts.
func NewShortcutsWindow() (*gtk.ShortcutsWindow, error) {
	builder, err := gtk.BuilderNew()
	if err != nil {
		return nil, err
	}

	err = builder.AddFromString(shortcutsWindowUI)
	if err != nil {
		return nil, err
	}

	obj, err := builder.GetObject("shortcuts-window")
	if err != nil {
		return nil, err
	}

	return obj.(*gtk.ShortcutsWindow), nil
}

var shortcutsWindow *gtk.ShortcutsWindow

// ShowShortcuts shows a shortcuts window to the user.
func (app *Application) ShowShortcuts() {
	var err error
	if shortcutsWindow == nil {
		shortcutsWindow, err = NewShortcutsWindow()
		if err != nil {
			log.Print(err)
			return
		}

		// We want to keep the shortcuts window around in case
		// we want to show it again.
		shortcutsWindow.HideOnDelete()
		shortcutsWindow.Connect("hide", func() {
			app.application.RemoveWindow(&shortcutsWindow.Window)
		})
	}

	app.application.AddWindow(&shortcutsWindow.Window)
	shortcutsWindow.Present()
}
