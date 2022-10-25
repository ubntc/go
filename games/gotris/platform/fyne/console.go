package fyne

import (
	"fmt"
)

func echo(s ...interface{}) { fmt.Println(s...) }

// func format(s ...interface{}) string {
// 	text := fmt.Sprintln(s...)
// 	lines := strings.Split("\n"+text, "\n")
// 	return fmt.Sprintln(strings.Join(lines, "\r\n"))
//}

func (p *Platform) SetRenderingMode(mode string) error {
	return nil
}

func (p *Platform) RenderingModes() (names []string, currentMode int) {
	return []string{"default"}, 0
}

func (p *Platform) RenderingInfo(name string) string {
	return "default mode"
}
