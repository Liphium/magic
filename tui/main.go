package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type (
	errMsg error
)

type model struct {
	textInput textinput.Model
	err       error
	width    int
	height   int
}

type sPipe struct {
	mutex *sync.Mutex
	items []string
}

var pn = &sPipe{
	mutex: &sync.Mutex{},
	items: []string{},
}

var ctr = &sPipe{
	mutex: &sync.Mutex{},
	items: []string{},
}

func addItem(pn *sPipe, item string) {
	pn.mutex.Lock()
	defer pn.mutex.Unlock()

	pn.items = append(pn.items, item)
}

func getItem(pn *sPipe) (string, bool){
	pn.mutex.Lock()
	defer pn.mutex.Unlock()
	if len(pn.items) == 0 {
		return "", false
	}
	item := pn.items[0]
	pn.items = pn.items[1:]
	return item, true
}

func runTui(pn *sPipe){
	i := 0
	for  {
		time.Sleep(time.Millisecond * 100)
		addItem(pn, "number of prints: " + fmt.Sprint(i))
		if cmd, newCmd := getItem(ctr); newCmd{
			switch cmd {
			case "cat":
				addItem(pn, "meow :3")
			}
			
		}
		i++
	}
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Enter command"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 0

	return model{
		textInput: ti,
		err:       nil,
	}
}

func (m model) Init() tea.Cmd {
	go runTui(pn)

	return tea.Batch(textinput.Blink)
}
// remove if creates lag or not wanted
func setUpdateTime() tea.Cmd {
	d := time.Millisecond * time.Duration(100) // set update time in ms
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return ""
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			addItem(ctr, m.textInput.Value())
			m.textInput.SetValue("")
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput.Width = m.width
	m.textInput, cmd = m.textInput.Update(msg)

	if itemToPrint, shouldPrint := getItem(pn); shouldPrint{
		return m, tea.Batch(
			tea.Println(itemToPrint),
			cmd,
			setUpdateTime(),
		)
	} else{
		return m, tea.Batch(cmd, setUpdateTime(),)
	}
}

func (m model) View() string {
	return fmt.Sprint(
		m.textInput.View(),
	)
}