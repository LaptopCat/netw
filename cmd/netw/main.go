package main

import (
	"flag"

	"log"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	primaryStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff4151"))
	boldStyle    = lipgloss.NewStyle().Bold(true).Render
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#bf616a")).Render
)

var p *tea.Program

func main() {
	flag.Parse()

	var app tea.Model
	switch flag.Arg(0) {
	case "":
		m := probing{}
		m.v4spinner = spinner.New(spinner.WithSpinner(spinner.Meter), spinner.WithStyle(primaryStyle))
		m.v6spinner = spinner.New(spinner.WithSpinner(spinner.Meter), spinner.WithStyle(primaryStyle))

		app = &m
	}

	p = tea.NewProgram(app)
	if _, err := p.Run(); err != nil {
		log.Fatalln("could not run program:", err)
	}
}
