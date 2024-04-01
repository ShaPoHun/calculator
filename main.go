// Package main launches the calculator app
//
//go:generate fyne bundle -o data.go Icon.png
package main

import "fyne.io/fyne/v2/app"

func main() {
	a := app.NewWithID("io.github.shapohun.calculator")
	a.SetIcon(resourceIconPng)

	c := newCalculator(a)
	c.loadPreferences()
	c.loadUI()
	a.Run()
}
