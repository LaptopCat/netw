package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/laptopcat/netw/lib/netw"
)

type probe struct {
	result netw.Probe
	err    error
	t      string

	done bool
}

type probing struct {
	v4spinner spinner.Model
	v6spinner spinner.Model
	v4probe   probe
	v6probe   probe

	done bool
}

func (m *probing) Init() tea.Cmd {
	return tea.Batch(
		func() tea.Msg {
			res, err := netw.ProbeV6()
			return probe{res, err, "v6", false}
		},
		func() tea.Msg {
			res, err := netw.ProbeV4()
			return probe{res, err, "v4", false}
		},
		m.v4spinner.Tick,
		m.v6spinner.Tick,
	)
}

func (m *probing) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		default:
			return m, nil
		}
	case spinner.TickMsg:
		var cmd tea.Cmd
		switch msg.ID {
		case m.v4spinner.ID():
			if !m.v4probe.done {
				m.v4spinner, cmd = m.v4spinner.Update(msg)
			}
		case m.v6spinner.ID():
			if !m.v6probe.done {
				m.v6spinner, cmd = m.v6spinner.Update(msg)
			}
		}

		return m, cmd
	case probe:
		switch msg.t {
		case "v4":
			m.v4probe = msg
		case "v6":
			m.v6probe = msg
		}
	}

	return m, nil
}

func probeview(name string, p *probe, sp spinner.Model) (s string) {
	if p.t == "" {
		s += fmt.Sprintf("\n  %s Probing %s...\n", sp.View(), name)
	} else {
		p.done = true
		if p.err != nil {
			s += fmt.Sprintf("\n  %s probe failed: %s\n", name, errorStyle(p.err.Error()))
		} else if !p.result.Worked {
			s += errorStyle("\n   %s probe failed\n")
		} else {
			s += fmt.Sprintf("\n  %s: %s (%s)", boldStyle(name), p.result.IPA, p.result.IPP)
			s += fmt.Sprintf("\n        AS%d (%s)\n", p.result.ASN, p.result.ASS)
		}
	}
	return
}

func (m *probing) View() (s string) {
	s = boldStyle("You are connecting from:")

	s += probeview("IPv4", &m.v4probe, m.v4spinner)
	s += probeview("IPv6", &m.v6probe, m.v6spinner)

	if m.done {
		go p.Quit()
		return
	} else if m.v4probe.done && m.v6probe.done {
		m.done = true
	}

	return
}
