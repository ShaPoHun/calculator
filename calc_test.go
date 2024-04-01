package main

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"

	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	calc := newCalculator(test.NewApp())
	calc.loadUI()

	test.Tap(calc.buttons["1"])
	test.Tap(calc.buttons["+"])
	test.Tap(calc.buttons["1"])
	test.Tap(calc.buttons["="])

	assert.Equal(t, "2", calc.output.Text)
}

func TestSubtract(t *testing.T) {
	calc := newCalculator(test.NewApp())
	calc.loadUI()

	test.Tap(calc.buttons["2"])
	test.Tap(calc.buttons["-"])
	test.Tap(calc.buttons["1"])
	test.Tap(calc.buttons["="])

	assert.Equal(t, "1", calc.output.Text)
}

func TestDivide(t *testing.T) {
	calc := newCalculator(test.NewApp())
	calc.loadUI()

	test.Tap(calc.buttons["4"])
	test.Tap(calc.buttons["/"])
	test.Tap(calc.buttons["2"])
	test.Tap(calc.buttons["="])

	assert.Equal(t, "2", calc.output.Text)
}

func TestDividePrecision(t *testing.T) {
	calc := newCalculator(test.NewApp())
	calc.app.Preferences().SetFloat(PrecisionPref, 2)
	calc.loadPreferences()
	calc.loadUI()

	test.Tap(calc.buttons["3"])
	test.Tap(calc.buttons["."])
	test.Tap(calc.buttons["5"])
	test.Tap(calc.buttons["/"])
	test.Tap(calc.buttons["2"])
	test.Tap(calc.buttons["="])

	assert.Equal(t, "1.75", calc.output.Text)
}

func TestMultiply(t *testing.T) {
	calc := newCalculator(test.NewApp())
	calc.loadUI()

	test.Tap(calc.buttons["5"])
	test.Tap(calc.buttons["*"])
	test.Tap(calc.buttons["2"])
	test.Tap(calc.buttons["="])

	assert.Equal(t, "10", calc.output.Text)
}

func TestParenthesis(t *testing.T) {
	calc := newCalculator(test.NewApp())
	calc.loadUI()

	test.Tap(calc.buttons["2"])
	test.Tap(calc.buttons["*"])
	test.Tap(calc.buttons["("])
	test.Tap(calc.buttons["3"])
	test.Tap(calc.buttons["+"])
	test.Tap(calc.buttons["4"])
	test.Tap(calc.buttons[")"])
	test.Tap(calc.buttons["="])

	assert.Equal(t, "14", calc.output.Text)
}

func TestDot(t *testing.T) {
	calc := newCalculator(test.NewApp())
	calc.loadUI()

	test.Tap(calc.buttons["2"])
	test.Tap(calc.buttons["."])
	test.Tap(calc.buttons["2"])
	test.Tap(calc.buttons["+"])
	test.Tap(calc.buttons["7"])
	test.Tap(calc.buttons["."])
	test.Tap(calc.buttons["8"])
	test.Tap(calc.buttons["="])

	assert.Equal(t, "10", calc.output.Text)
}

func TestClear(t *testing.T) {
	calc := newCalculator(test.NewApp())
	calc.loadUI()

	test.Tap(calc.buttons["1"])
	test.Tap(calc.buttons["2"])
	test.Tap(calc.buttons["C"])

	assert.Equal(t, "", calc.output.Text)
}

func TestContinueAfterResult(t *testing.T) {
	calc := newCalculator(test.NewApp())
	calc.loadUI()

	test.Tap(calc.buttons["6"])
	test.Tap(calc.buttons["+"])
	test.Tap(calc.buttons["4"])
	test.Tap(calc.buttons["="])
	test.Tap(calc.buttons["-"])
	test.Tap(calc.buttons["2"])
	test.Tap(calc.buttons["="])

	assert.Equal(t, "8", calc.output.Text)
}

func TestKeyboard(t *testing.T) {
	calc := newCalculator(test.NewApp())
	calc.loadUI()

	test.TypeOnCanvas(calc.window.Canvas(), "1+1")
	assert.Equal(t, "1+1", calc.output.Text)

	test.TypeOnCanvas(calc.window.Canvas(), "=")
	assert.Equal(t, "2", calc.output.Text)

	test.TypeOnCanvas(calc.window.Canvas(), "c")
	assert.Equal(t, "", calc.output.Text)
}

func TestKeyboard_Buttons(t *testing.T) {
	calc := newCalculator(test.NewApp())
	calc.loadUI()

	test.TypeOnCanvas(calc.window.Canvas(), "1+1")
	calc.onTypedKey(&fyne.KeyEvent{Name: fyne.KeyReturn})
	assert.Equal(t, "2", calc.output.Text)

	test.TypeOnCanvas(calc.window.Canvas(), "c")

	test.TypeOnCanvas(calc.window.Canvas(), "1+1")
	calc.onTypedKey(&fyne.KeyEvent{Name: fyne.KeyEnter})
	assert.Equal(t, "2", calc.output.Text)
}

func TestKeyboard_Backspace(t *testing.T) {
	calc := newCalculator(test.NewApp())
	calc.loadUI()

	calc.onTypedKey(&fyne.KeyEvent{Name: fyne.KeyBackspace})
	assert.Equal(t, "", calc.output.Text)

	test.TypeOnCanvas(calc.window.Canvas(), "1/2")
	calc.onTypedKey(&fyne.KeyEvent{Name: fyne.KeyBackspace})
	assert.Equal(t, "1/", calc.output.Text)

	calc.onTypedKey(&fyne.KeyEvent{Name: fyne.KeyEnter})
	assert.Equal(t, "error", calc.output.Text)

	calc.onTypedKey(&fyne.KeyEvent{Name: fyne.KeyBackspace})
	assert.Equal(t, "", calc.output.Text)
}

func TestError(t *testing.T) {
	calc := newCalculator(test.NewApp())
	calc.loadUI()

	test.TypeOnCanvas(calc.window.Canvas(), "1//1=")
	assert.Equal(t, "error", calc.output.Text)

	test.TypeOnCanvas(calc.window.Canvas(), "c")

	test.TypeOnCanvas(calc.window.Canvas(), "()9=")
	assert.Equal(t, "error", calc.output.Text)

	test.TypeOnCanvas(calc.window.Canvas(), "=")
	assert.Equal(t, "error", calc.output.Text)

	test.TypeOnCanvas(calc.window.Canvas(), "55=")
	assert.Equal(t, "error", calc.output.Text)
}

func TestShortcuts(t *testing.T) {
	app := test.NewApp()
	calc := newCalculator(app)
	calc.loadUI()
	clipboard := app.Driver().AllWindows()[0].Clipboard()

	test.TypeOnCanvas(calc.window.Canvas(), "720 + 80")
	calc.onCopyShortcut(&fyne.ShortcutCopy{Clipboard: clipboard})
	assert.Equal(t, clipboard.Content(), calc.output.Text)

	test.TypeOnCanvas(calc.window.Canvas(), "+")
	clipboard.SetContent("50")
	calc.onPasteShortcut(&fyne.ShortcutPaste{Clipboard: clipboard})
	test.TypeOnCanvas(calc.window.Canvas(), "=")
	assert.Equal(t, "850", calc.output.Text)

	clipboard.SetContent("not a valid number")
	calc.onPasteShortcut(&fyne.ShortcutPaste{Clipboard: clipboard})
	assert.Equal(t, "850", calc.output.Text)
}

func TestCalculationProcess(t *testing.T) {
	calc := newCalculator(test.NewApp())
	calc.loadUI()

	test.Tap(calc.buttons["9"])
	test.Tap(calc.buttons["."])
	test.Tap(calc.buttons["2"])
	test.Tap(calc.buttons["*"])
	test.Tap(calc.buttons["3"])
	test.Tap(calc.buttons["="])

	assert.Equal(t, "9.2*3", calc.process.Text)

	test.Tap(calc.buttons["C"])
	assert.Equal(t, "", calc.process.Text)
}

func TestHistory(t *testing.T) {
	calc := newCalculator(test.NewApp())
	calc.loadUI()

	test.Tap(calc.buttons["1"])
	test.Tap(calc.buttons["+"])
	test.Tap(calc.buttons["2"])
	test.Tap(calc.buttons["="])
	assert.Contains(t, calc.historyText.Text(), "1+2 = 3")

	test.Tap(calc.buttons["C"])
	test.Tap(calc.buttons["1"])
	test.Tap(calc.buttons["*"])
	test.Tap(calc.buttons["2"])
	test.Tap(calc.buttons["="])
	assert.Contains(t, calc.historyText.Text(), "1*2 = 2")

	for _, item := range calc.menu().Items {
		if item.Label == "Settings" {
			for _, settingItem := range item.Items {
				if settingItem.Label == "Show History" {
					settingItem.Action()
					assert.Equal(t, calc.isHistoryWindowOpen, true)
				} else if settingItem.Label == "Clear History" {
					settingItem.Action()
					assert.Equal(t, "", calc.historyText.Text())
				}
			}
			break
		}
	}
}

func TestPrecision(t *testing.T) {
	calc := newCalculator(test.NewApp())
	calc.loadUI()

	for _, item := range calc.menu().Items {
		if item.Label == "Settings" {
			for _, settingItem := range item.Items {
				if settingItem.Label == "Precision" {
					for i, m := range settingItem.ChildMenu.Items {
						m.Action()
						assert.Equal(t, calc.precision, i)
					}
				}
			}
			break
		}
	}
}
