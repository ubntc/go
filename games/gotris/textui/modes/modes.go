// static mode management
package modes

import (
	"github.com/ubntc/go/games/gotris/textui/boxart"
	"github.com/ubntc/go/games/gotris/textui/doublewidth"
	"github.com/ubntc/go/games/gotris/textui/fullwidth"
)

// BoxArtMode wraps the core features of arts.BoxArt types
// and allows switching these using the ModeManager.
type BoxArtMode interface {
	GetName() string
	GetDesc() string

	// GetBoxArt returns the concrete BoxArt object.
	GetBoxArt() *boxart.BoxArt

	// TextToBlock converts regular text to BoxArt text using the
	// defined BoxArt characters in the FrameArt implementation.
	TextToBlock(string) string
}

type ModeManager struct {
	art   BoxArtMode
	modes []BoxArtMode
}

func NewModeManager() *ModeManager {
	m := &ModeManager{}
	m.modes = []BoxArtMode{
		doublewidth.New(),
		fullwidth.New(),
	}
	m.art = m.modes[0]
	return m
}

func (m *ModeManager) SetModeByName(name string) {
	m.art = m.GetModeByName(name)
}

func (m *ModeManager) GetModeByName(name string) BoxArtMode {
	return m.modes[m.ModeIndexByName(name)]
}

func (m *ModeManager) ModeNames() (names []string) {
	for _, v := range m.modes {
		names = append(names, v.GetName())
	}
	return
}

func (m *ModeManager) ModeDescs() (descs []string) {
	for _, v := range m.modes {
		descs = append(descs, v.GetDesc())
	}
	return
}

func (m *ModeManager) ModeIndexByName(name string) int {
	for i, v := range m.modes {
		if v.GetName() == name {
			return i
		}
	}
	panic("no modes defined")
}

func (m *ModeManager) Mode() BoxArtMode {
	return m.art
}

func (m *ModeManager) ModeName() string {
	return m.art.GetName()
}

func (m *ModeManager) ModeInfo(name string) string {
	return m.GetModeByName(name).GetDesc()
}

func (m *ModeManager) ModeIndex() int {
	return m.ModeIndexByName(m.ModeName())
}
