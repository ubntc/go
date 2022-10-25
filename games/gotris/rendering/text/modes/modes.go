// static mode management
package modes

import (
	"github.com/ubntc/go/games/gotris/rendering/text/arts"
	"github.com/ubntc/go/games/gotris/rendering/text/doublewidth"
	"github.com/ubntc/go/games/gotris/rendering/text/fullwidth"
)

type ModeManager struct {
	art   arts.FrameArt
	modes []arts.FrameArt
}

func NewModeManager() *ModeManager {
	m := &ModeManager{}
	m.art = fullwidth.New()
	m.modes = []arts.FrameArt{
		m.art,
		doublewidth.New(),
	}
	return m
}

func (m *ModeManager) SetModeByName(name string) {
	m.art = m.GetModeByName(name)
}

func (m *ModeManager) GetModeByName(name string) arts.FrameArt {
	return m.modes[m.ModeIndexByName(name)]
}

func (m *ModeManager) ModeNames() (names []string) {
	for _, v := range m.modes {
		names = append(names, v.Art().Name)
	}
	return
}

func (m *ModeManager) ModeIndexByName(name string) int {
	for i, v := range m.modes {
		if v.Art().Name == name {
			return i
		}
	}
	panic("no modes defined")
}

func (m *ModeManager) Mode() arts.FrameArt {
	return m.art
}

func (m *ModeManager) ModeName() string {
	return m.art.Art().Name
}

func (m *ModeManager) ModeInfo(name string) string {
	return m.GetModeByName(name).Art().Desc
}

func (m *ModeManager) ModeIndex() int {
	return m.ModeIndexByName(m.ModeName())
}
