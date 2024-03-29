// Package main launches the calculator app
//
//go:generate fyne bundle -o data.go Icon.png
package main

import "fyne.io/fyne/v2/app"

func main() {
	a := app.New()
	a.SetIcon(resourceIconPng)

	c := newCalculator()
	c.loadUI(a)
	a.Run()
}
