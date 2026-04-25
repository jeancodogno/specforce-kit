package tui

import (
	"fmt"
	"github.com/jeancodogno/specforce-kit/src/internal/core"
)

type terminalUI struct {
	spinner *NeonSpinner
}

func NewUI() core.UI {
	return &terminalUI{}
}

func (u *terminalUI) Log(message string) {
	fmt.Println(BodyStyle.Render(message))
}

func (u *terminalUI) Warn(message string) {
	fmt.Print(RenderBadge("warning", message))
}

func (u *terminalUI) Error(message string) {
	fmt.Print(RenderBadge("error", message))
}

func (u *terminalUI) Success(message string) {
	fmt.Print(RenderBadge("success", message))
}

func (u *terminalUI) SubTask(message string) {
	LogSubTask(message)
}

func (u *terminalUI) StartSpinner(message string) {
	if u.spinner == nil {
		s := NewNeonSpinner(message)
		u.spinner = &s
	}
	fmt.Print(u.spinner.View())
}

func (u *terminalUI) StopSpinner() {
	// For now, our spinner implementation is simple and just prints once in View().
	// A more complex spinner would need to be updated.
	u.spinner = nil
}

func (u *terminalUI) Confirm(question string) bool {
	return Confirm(question)
}
