//go:generate fyne bundle -o data.go Icon.png

package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/Knetic/govaluate"
)

type calc struct {
	precision           int
	equation            string
	isHistoryWindowOpen bool

	app           fyne.App
	output        *widget.Label
	process       *widget.Label
	historyText   *widget.TextGrid
	buttons       map[string]*widget.Button
	window        fyne.Window
	historyWindow fyne.Window
}

func (c *calc) display(newText string) {
	c.equation = newText
	c.output.SetText(newText)
}

func (c *calc) showProcess(newText string) {
	c.process.SetText(newText)
}

func (c *calc) character(char rune) {
	c.display(c.equation + string(char))
}

func (c *calc) digit(d int) {
	c.character(rune(d) + '0')
}

func (c *calc) clear() {
	c.display("")
	c.showProcess("")
}

func (c *calc) backspace() {
	if len(c.equation) == 0 {
		return
	} else if c.equation == "error" {
		c.clear()
		return
	}

	c.display(c.equation[:len(c.equation)-1])
}

func (c *calc) evaluate() {
	if strings.Contains(c.output.Text, "error") {
		c.display("error")
		return
	}

	_, err := strconv.ParseFloat(c.output.Text, 64)
	if err == nil {
		log.Println("Invalid equation", c.output.Text)
		c.display("error")
		return
	}

	expression, err := govaluate.NewEvaluableExpression(c.output.Text)
	if err != nil {
		log.Println("Error in calculation", err)
		c.display("error")
		return
	}

	result, err := expression.Evaluate(nil)
	if err != nil {
		log.Println("Error in calculation", err)
		c.display("error")
		return
	}

	value, ok := result.(float64)
	if !ok {
		log.Println("Invalid input:", c.output.Text)
		c.display("error")
		return
	}

	newText := strconv.FormatFloat(value, 'f', c.precision, 64)
	c.showProcess(fmt.Sprintf("%s=", c.equation))
	c.writeHistory(fmt.Sprintf("%s=%s", c.equation, newText))
	c.display(newText)
}

func (c *calc) addButton(text string, action func()) *widget.Button {
	button := widget.NewButton(text, action)
	c.buttons[text] = button
	return button
}

func (c *calc) digitButton(number int) *widget.Button {
	str := strconv.Itoa(number)
	return c.addButton(str, func() {
		// Clear the calculator if the output is "error" or
		//if the process is not empty and the equation doesn't contain any of the operators: +, -, *, /
		if c.output.Text == "error" ||
			(c.process.Text != "" && !strings.ContainsAny(c.equation, "+-*/")) {
			c.clear()
		}
		c.digit(number)
	})
}

func (c *calc) charButton(char rune) *widget.Button {
	return c.addButton(string(char), func() {
		c.character(char)
	})
}

func (c *calc) onTypedRune(r rune) {
	if r == 'c' {
		r = 'C' // The button is using a capital C.
	}

	if button, ok := c.buttons[string(r)]; ok {
		button.OnTapped()
	}
}

func (c *calc) onTypedKey(ev *fyne.KeyEvent) {
	if ev.Name == fyne.KeyReturn || ev.Name == fyne.KeyEnter {
		c.evaluate()
	} else if ev.Name == fyne.KeyBackspace {
		c.backspace()
	}
}

func (c *calc) onPasteShortcut(shortcut fyne.Shortcut) {
	content := shortcut.(*fyne.ShortcutPaste).Clipboard.Content()
	if content == "" {
		return
	}

	if _, err := strconv.ParseFloat(content, 64); err != nil {
		return
	}

	if !strings.ContainsAny(c.equation, "+-*/") {
		c.display(content)
		c.msgBubble("paste from clipboard")
	}

	if strings.LastIndexAny(c.equation, "+-*/") == len(c.equation)-1 {
		c.display(c.equation + content)
		c.msgBubble("paste from clipboard")
	}
}

func (c *calc) onCopyShortcut(shortcut fyne.Shortcut) {
	if c.equation == "" {
		return
	}
	shortcut.(*fyne.ShortcutCopy).Clipboard.SetContent(c.equation)
	c.msgBubble("copy to clipboard")
}

func (c *calc) msgBubble(text string) {
	tipText := canvas.NewText(text, nil)

	popUp := widget.NewPopUp(container.NewWithoutLayout(tipText), c.window.Canvas())
	popUp.ShowAtPosition(fyne.CurrentApp().Driver().AbsolutePositionForObject(c.output))

	time.AfterFunc(1*time.Second, func() {
		popUp.Hide()
	})
}

func (c *calc) menu() *fyne.MainMenu {
	settingsMenu := fyne.NewMenu("Settings",
		c.showHistoryMenu(),
		c.clearHistoryMenu(),
		c.setPrecisionMenu(),
	)
	return fyne.NewMainMenu(settingsMenu)
}

func (c *calc) showHistoryMenu() *fyne.MenuItem {
	pref := c.app.Preferences()
	var historyMenu *fyne.MenuItem

	newHistoryWindow := func() {
		historyMenu.Checked = true
		c.isHistoryWindowOpen = true
		c.historyWindow = c.app.NewWindow("History")
		c.historyWindow.SetOnClosed(func() {
			historyMenu.Checked = false
			c.isHistoryWindowOpen = false
			pref.SetBool(ShowHistoryPref, false)
		})
		scrollContainer := container.NewScroll(c.historyText)
		scrollContainer.Resize(fyne.NewSize(300, 450))
		scrollContainer.Position()
		c.historyWindow.SetContent(scrollContainer)
		c.historyWindow.Resize(fyne.NewSize(300, 450))
		c.historyWindow.Show()
	}

	historyMenu = fyne.NewMenuItem("Show History", func() {
		if c.isHistoryWindowOpen {
			c.historyWindow.Close()
			c.isHistoryWindowOpen = false
			historyMenu.Checked = false
			pref.SetBool(ShowHistoryPref, false)
		} else {
			pref.SetBool(ShowHistoryPref, true)
			newHistoryWindow()
		}
	})

	if c.isHistoryWindowOpen {
		log.Println("History window already open")
		newHistoryWindow()
	}
	return historyMenu
}

func (c *calc) clearHistoryMenu() *fyne.MenuItem {
	return fyne.NewMenuItem("Clear History", func() {
		c.clearHistory()
	})
}

func (c *calc) setPrecisionMenu() *fyne.MenuItem {
	pref := c.app.Preferences()
	var options []*fyne.MenuItem
	for i := 0; i <= 15; i++ {
		option := fyne.NewMenuItem(strconv.Itoa(i), func() {
			pref.SetInt(PrecisionPref, i)
			options[c.precision].Checked = false
			options[i].Checked = true
			c.precision = i
			log.Println("Precision set to", c.precision)
		})
		if i == c.precision {
			option.Checked = true
		}
		options = append(options, option)
	}
	configureSubMenu := fyne.NewMenu("", options...)

	precisionMenu := fyne.NewMenuItem("Precision", nil)
	precisionMenu.ChildMenu = configureSubMenu
	return precisionMenu
}

func (c *calc) writeHistory(line string) {
	c.historyText.SetText(line + "\n" + c.historyText.Text())
	c.app.Preferences().SetString(HistoryTextPref, c.historyText.Text())
}

func (c *calc) clearHistory() {
	c.historyText.SetText("")
	c.app.Preferences().SetString(HistoryTextPref, "")
}

func (c *calc) loadPreferences() {
	p := c.app.Preferences()
	c.precision = p.Int(PrecisionPref)
	c.isHistoryWindowOpen = p.Bool(ShowHistoryPref)
	text := c.app.Preferences().String(HistoryTextPref)
	if text != "" {
		c.historyText.SetText(text)
	}
}

func (c *calc) loadUI() {
	c.output = &widget.Label{
		Alignment: fyne.TextAlignTrailing,
		TextStyle: fyne.TextStyle{
			Monospace: true,
			Bold:      true,
		},
		Truncation: fyne.TextTruncateEllipsis,
	}

	c.process = &widget.Label{
		Alignment: fyne.TextAlignTrailing,
		TextStyle: fyne.TextStyle{
			Monospace: true,
			Bold:      true,
		},
		Truncation: fyne.TextTruncateEllipsis,
	}

	c.historyText.ShowLineNumbers = true

	equals := c.addButton("=", c.evaluate)
	equals.Importance = widget.HighImportance

	c.window = c.app.NewWindow("Calc")
	c.window.SetContent(container.NewGridWithColumns(1,
		container.NewGridWithColumns(1,
			c.process, c.output),
		container.NewGridWithColumns(4,
			c.addButton("C", c.clear),
			c.charButton('('),
			c.charButton(')'),
			c.charButton('/')),
		container.NewGridWithColumns(4,
			c.digitButton(7),
			c.digitButton(8),
			c.digitButton(9),
			c.charButton('*')),
		container.NewGridWithColumns(4,
			c.digitButton(4),
			c.digitButton(5),
			c.digitButton(6),
			c.charButton('-')),
		container.NewGridWithColumns(4,
			c.digitButton(1),
			c.digitButton(2),
			c.digitButton(3),
			c.charButton('+')),
		container.NewGridWithColumns(2,
			container.NewGridWithColumns(2,
				c.digitButton(0),
				c.charButton('.')),
			equals)),
	)

	canvas := c.window.Canvas()
	canvas.SetOnTypedRune(c.onTypedRune)
	canvas.SetOnTypedKey(c.onTypedKey)
	canvas.AddShortcut(&fyne.ShortcutCopy{}, c.onCopyShortcut)
	canvas.AddShortcut(&fyne.ShortcutPaste{}, c.onPasteShortcut)

	c.window.SetMaster()
	c.window.Resize(fyne.NewSize(300, 450))
	c.window.CenterOnScreen()
	c.window.SetMainMenu(c.menu())
	c.window.Show()
}

func newCalculator(app fyne.App) *calc {
	return &calc{
		app:         app,
		historyText: widget.NewTextGrid(),
		buttons:     make(map[string]*widget.Button, 19),
	}
}
