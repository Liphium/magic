package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/urfave/cli/v3"
)

const refreshDelay = 100

var cmdError = ""

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
	width     int
	height    int
}

type sPipe struct {
	mutex *sync.Mutex
	items []string
}

var console = &sPipe{
	mutex: &sync.Mutex{},
	items: []string{},
}

var ctr = &sPipe{
	mutex: &sync.Mutex{},
	items: []string{},
}

func (pn *sPipe) addItem(item string) {
	pn.mutex.Lock()
	defer pn.mutex.Unlock()

	pn.items = append(pn.items, item)
}

func (pn *sPipe) getItem() (string, bool) {
	pn.mutex.Lock()
	defer pn.mutex.Unlock()
	if len(pn.items) == 0 {
		return "", false
	}
	item := pn.items[0]
	pn.items = pn.items[1:]
	return item, true
}

func runTui() {
	i := 0
	for {
		time.Sleep(time.Millisecond * refreshDelay)
		console.addItem("number of prints: " + fmt.Sprint(i))
		if cmd, newCmd := ctr.getItem(); newCmd {
			var testPath string = ""
			var scriptPath string = ""
			var commands = &cli.Command{
				HideHelp:        true,
				OnUsageError:    func(ctx context.Context, cmd *cli.Command, err error, isSubcommand bool) error { return nil },
				CommandNotFound: func(ctx context.Context, c *cli.Command, s string) {cmdError = "command doesnt exist"},
				Commands: []*cli.Command{
					{
						Name:    "run",
						Usage:   "",
						Aliases: []string{"r"},
						Arguments: []cli.Argument{
							&cli.StringArg{
								Name:        "path",
								Destination: &scriptPath,
							},
						},
						Action: func(ctx context.Context, cmd *cli.Command) error {
							if scriptPath != "" {
								go runCommand(scriptPath, console)
							} else {
								cmdError = "usage: run [path]"
							}
							return nil
						},
					},
					{
						Name:    "test",
						Usage:   "",
						Aliases: []string{"t"},
						Arguments: []cli.Argument{
							&cli.StringArg{
								Name:        "path",
								Destination: &testPath,
							},
						},
						Action: func(ctx context.Context, cmd *cli.Command) error {
							if testPath != "" {
								go testCommand(testPath, console)
							} else {
								cmdError = "usage: test [path]"
							}
							return nil
						},
					},
				},
			}
			if err := commands.Run(context.Background(), append([]string{""}, strings.Split(strings.Trim(cmd, " "), " ")...)); err != nil {
			}
		}
		i++
	}
}

func initialModel() model {
	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 0

	return model{
		textInput: ti,
		err:       nil,
	}
}

func (m model) Init() tea.Cmd {
	go runTui()

	return tea.Batch(textinput.Blink)
}

// remove if creates lag or not wanted
func setUpdateTime() tea.Cmd {
	d := time.Millisecond * time.Duration(refreshDelay)
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
			if value := m.textInput.Value(); value != "" {
				m.textInput.SetValue("")
				ctr.addItem(value)
			}
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	if m.textInput.Value() == "" && cmdError != "" {
		m.textInput.Placeholder = cmdError
	} else {
		m.textInput.Placeholder = "Enter command"
	}
	if m.textInput.Value() != "" {
		cmdError = ""
	}

	m.textInput.Width = m.width
	m.textInput, cmd = m.textInput.Update(msg)

	if itemToPrint, shouldPrint := console.getItem(); shouldPrint {
		return m, tea.Batch(
			tea.Println(itemToPrint),
			cmd,
			setUpdateTime(),
		)
	} else {
		return m, tea.Batch(cmd, setUpdateTime())
	}
}

func (m model) View() string {
	return fmt.Sprint(
		m.textInput.View(),
	)
}
