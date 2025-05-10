package ui

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type State int

const (
	StateWelcome State = iota
	StateSettings
	StateGame
	StateGameOver
)

type Model struct {
	state            State
	words            []string
	input            string
	currentIdx       int
	startTime        time.Time
	timeLeft         time.Duration
	timeLimit        time.Duration
	finished         bool
	mistakes         int
	lastMistakePos   int
	incorrectLetters map[int]map[int]bool // wordIdx -> letterIdx -> bool
	finalElapsed     float64
	prevInputs       map[int]string
	animFrame        int    // for welcome screen animation
	barLetters       []rune // for Pacman animation
}

func NewModel(words []string) Model {
	return Model{
		state:            StateWelcome,
		words:            words,
		timeLimit:        30 * time.Second,
		timeLeft:         30 * time.Second,
		lastMistakePos:   -1,
		incorrectLetters: make(map[int]map[int]bool),
		prevInputs:       make(map[int]string),
		animFrame:        0,
		barLetters:       nil,
	}
}

func (m Model) Init() tea.Cmd {
	if m.state == StateWelcome {
		return welcomeAnimTick()
	}
	return nil
}

func welcomeAnimTick() tea.Cmd {
	return tea.Tick(60*time.Millisecond, func(time.Time) tea.Msg {
		return welcomeAnimMsg{}
	})
}

type welcomeAnimMsg struct{}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case welcomeAnimMsg:
		if m.state == StateWelcome {
			barWidth := 66
			blockSpacing := 6
			letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
			headPos := m.animFrame % ((barWidth - 2) * 2)
			if headPos >= barWidth-2 {
				headPos = (barWidth-2)*2 - headPos
			}
			// If at rightmost edge, generate new random letters
			if headPos == barWidth-2 || m.barLetters == nil {
				newLetters := make([]rune, barWidth)
				r := rand.New(rand.NewSource(time.Now().UnixNano()))
				for i := 0; i < barWidth; i += blockSpacing {
					newLetters[i] = rune(letters[r.Intn(len(letters))])
				}
				m.barLetters = newLetters
			}
			m.animFrame++
			return m, welcomeAnimTick()
		}
		return m, nil
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
		switch m.state {
		case StateWelcome:
			if msg.String() == "enter" || msg.String() == " " {
				m.state = StateGame
				m.startTime = time.Now()
				m.timeLeft = m.timeLimit
				m.input = ""
				m.currentIdx = 0
				m.mistakes = 0
				m.lastMistakePos = -1
				m.incorrectLetters = make(map[int]map[int]bool)
				m.finalElapsed = 0
				m.prevInputs = make(map[int]string)
				return m, tick()
			}
			if msg.String() == "s" {
				m.state = StateSettings
				return m, nil
			}
		case StateSettings:
			if msg.String() == "esc" {
				m.state = StateWelcome
				m.animFrame = 0
				m.barLetters = nil
				return m, welcomeAnimTick()
			}
			if msg.String() == "up" {
				m.timeLimit += 5 * time.Second
				if m.timeLimit > 300*time.Second {
					m.timeLimit = 300 * time.Second
				}
			}
			if msg.String() == "down" {
				m.timeLimit -= 5 * time.Second
				if m.timeLimit < 10*time.Second {
					m.timeLimit = 10 * time.Second
				}
			}
			return m, nil
		case StateGame:
			if m.finished {
				return m, nil
			}
			if msg.String() == "esc" {
				m.state = StateWelcome
				m.animFrame = 0
				m.barLetters = nil
				return m, welcomeAnimTick()
			}
			if msg.String() == "backspace" {
				if len(m.input) > 0 {
					m.input = m.input[:len(m.input)-1]
					m.lastMistakePos = -1
					return m, nil
				} else if m.currentIdx > 0 {
					m.prevInputs[m.currentIdx] = ""
					delete(m.incorrectLetters, m.currentIdx)
					m.currentIdx--
					m.input = m.prevInputs[m.currentIdx]
					m.lastMistakePos = -1
					return m, nil
				}
			}
			if msg.String() == "left" && m.currentIdx > 0 {
				m.currentIdx--
				m.input = m.prevInputs[m.currentIdx]
				m.lastMistakePos = -1
				return m, nil
			}
			if msg.String() == "right" && m.currentIdx < len(m.words)-1 {
				m.currentIdx++
				m.input = m.prevInputs[m.currentIdx]
				m.lastMistakePos = -1
				return m, nil
			}
			if msg.String() == " " {
				word := m.words[m.currentIdx]
				m.prevInputs[m.currentIdx] = m.input
				if m.input != word {
					if m.incorrectLetters[m.currentIdx] == nil {
						m.incorrectLetters[m.currentIdx] = make(map[int]bool)
					}
					for i := 0; i < len(m.input) && i < len(word); i++ {
						if m.input[i] != word[i] {
							m.incorrectLetters[m.currentIdx][i] = true
							m.mistakes++
						}
					}
					if len(m.input) < len(word) {
						for i := len(m.input); i < len(word); i++ {
							m.incorrectLetters[m.currentIdx][i] = true
							m.mistakes++
						}
					}
				} else {
					delete(m.incorrectLetters, m.currentIdx)
				}
				m.currentIdx++
				m.input = m.prevInputs[m.currentIdx]
				m.lastMistakePos = -1
				return m, nil
			}
			if len(msg.String()) == 1 && m.currentIdx < len(m.words) {
				if len(m.input) < len(m.words[m.currentIdx]) {
					m.input += msg.String()
					// Mistake tracking (real-time)
					pos := len(m.input) - 1
					if m.input[pos] != m.words[m.currentIdx][pos] {
						if m.incorrectLetters[m.currentIdx] == nil {
							m.incorrectLetters[m.currentIdx] = make(map[int]bool)
						}
						if !m.incorrectLetters[m.currentIdx][pos] {
							m.incorrectLetters[m.currentIdx][pos] = true
							m.mistakes++
						}
					} else if m.incorrectLetters[m.currentIdx] != nil && m.incorrectLetters[m.currentIdx][pos] {
						delete(m.incorrectLetters[m.currentIdx], pos)
					}
				}
			}
		case StateGameOver:
			if msg.String() == "enter" {
				m.state = StateWelcome
				m.currentIdx = 0
				m.input = ""
				m.finished = false
				m.mistakes = 0
				m.lastMistakePos = -1
				m.incorrectLetters = make(map[int]map[int]bool)
				m.finalElapsed = 0
				m.prevInputs = make(map[int]string)
				m.animFrame = 0
				m.barLetters = nil
				return m, welcomeAnimTick()
			}
		}
	case tickMsg:
		if m.state == StateGame && !m.finished {
			m.timeLeft -= time.Second
			if m.timeLeft <= 0 {
				m.finished = true
				if m.finalElapsed == 0 {
					m.finalElapsed = time.Since(m.startTime).Minutes()
				}
				m.state = StateGameOver
				return m, nil
			}
			return m, tick()
		}
	}
	return m, nil
}

type tickMsg struct{}

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}

func (m Model) View() string {
	switch m.state {
	case StateWelcome:
		return welcomeViewWithAnim(m, m.animFrame)
	case StateSettings:
		return settingsView(m)
	case StateGame:
		return gameView(m)
	case StateGameOver:
		return gameOverView(m)
	}
	return ""
}

func welcomeViewWithAnim(m Model, frame int) string {
	ascii := `
 _________  ___  ___  ___  _________    ___    ___ ________  _______      
|\___   ___\\  \|\  \|\  \|\___   ___\ |\  \  /  /|\   __  \|\  ___ \     
\|___ \  \_\ \  \\\  \ \  \|___ \  \_| \ \  \/  / | \  \|\  \ \   __/|    
     \ \  \ \ \  \\\  \ \  \   \ \  \   \ \    / / \ \   ____\ \  \_|/__  
      \ \  \ \ \  \\\  \ \  \   \ \  \   \/  /  /   \ \  \___|\ \  \_|\ \ 
       \ \__\ \ \_______\ \__\   \ \__\__/  / /      \ \__\    \ \_______\
        \|__|  \|_______|\|__|    \|__|\___/ /        \|__|     \|_______|
                                      \|___|/                             `
	padTop := "\n\n\n"
	padBottom := "\n\n\n"

	// Bouncing animation with random letters
	barWidth := 66
	blockSpacing := 6
	headPos := frame % ((barWidth - 2) * 2)
	if headPos >= barWidth-2 {
		headPos = (barWidth-2)*2 - headPos
	}
	bar := ""
	lettersStyle := lipgloss.NewStyle().Foreground(Catppuccin.Mauve)
	for i := 0; i < barWidth; i++ {
		if i == headPos {
			bar += lipgloss.NewStyle().Foreground(Catppuccin.Yellow).Render("█")
		} else if i%blockSpacing == 0 && i < headPos {
			bar += " "
		} else if i%blockSpacing == 0 && m.barLetters != nil && m.barLetters[i] != 0 {
			bar += lettersStyle.Render(string(m.barLetters[i]))
		} else {
			bar += " "
		}
	}
	bar = lipgloss.NewStyle().Width(barWidth).Align(lipgloss.Center).Render(bar)

	return WindowStyle.Width(100).Render(
		padTop +
			ascii +
			"\n" +
			bar +
			padBottom +
			lipgloss.NewStyle().Foreground(Catppuccin.Subtext0).Render("Press [Enter] or [Space] to start, [S] for settings") +
			"\n" +
			ButtonStyle.Render(" Start "),
	)
}

func welcomeView() string {
	return welcomeViewWithAnim(Model{}, 0)
}

func gameView(m Model) string {
	words := make([]string, len(m.words))
	for i, w := range m.words {
		if i < m.currentIdx {
			var b strings.Builder
			for j, r := range []rune(w) {
				if m.incorrectLetters[i] != nil && m.incorrectLetters[i][j] {
					b.WriteString(IncorrectStyle.Render(string(r)))
				} else {
					b.WriteString(TypedStyle.Render(string(r)))
				}
			}
			words[i] = b.String()
		} else if i == m.currentIdx {
			var b strings.Builder
			inputRunes := []rune(m.input)
			wordRunes := []rune(w)
			for j, r := range wordRunes {
				if j < len(inputRunes) {
					if inputRunes[j] == r {
						b.WriteString(TypedStyle.Render(string(r)))
					} else {
						b.WriteString(IncorrectStyle.Render(string(r)))
					}
				} else if j == len(inputRunes) {
					b.WriteString("|")
					b.WriteString(WordStyle.Render(string(r)))
				} else {
					b.WriteString(WordStyle.Render(string(r)))
				}
			}
			if len(inputRunes) == len(wordRunes) {
				b.WriteString("|")
			}
			words[i] = b.String()
		} else {
			words[i] = lipgloss.NewStyle().Foreground(Catppuccin.Overlay1).Render(w)
		}
	}
	wordsLine := lipgloss.NewStyle().Align(lipgloss.Center).Width(72).Render(strings.Join(words, " "))
	timer := lipgloss.NewStyle().Foreground(Catppuccin.Yellow).Render(fmt.Sprintf("%ds", int(m.timeLeft.Seconds())))
	padTop := "\n\n\n"
	padBottom := "\n\n\n"
	return WindowStyle.Width(100).Render(
		padTop +
			timer + "\n\n" + wordsLine + "\n\n" +
			lipgloss.NewStyle().Foreground(Catppuccin.Subtext0).Render("Type the words above. [Esc] to quit. [←][→] to move.") +
			padBottom,
	)
}

func gameOverView(m Model) string {
	minElapsed := m.finalElapsed
	if minElapsed == 0 {
		minElapsed = 0.5
		if !m.startTime.IsZero() {
			minElapsed = time.Since(m.startTime).Minutes()
			if minElapsed < 0.01 {
				minElapsed = 0.01
			}
		}
	}
	chars := 0
	for i := 0; i < m.currentIdx; i++ {
		w := m.words[i]
		for j := 0; j < len(w); j++ {
			if m.incorrectLetters[i] == nil || !m.incorrectLetters[i][j] {
				chars++
			}
		}
	}
	wpm := float64(chars) / 5.0 / minElapsed
	padTop := "\n\n\n"
	padBottom := "\n\n\n"
	return WindowStyle.Width(100).Render(
		padTop +
			TitleStyle.Render("Game Over!") +
			"\n" +
			lipgloss.NewStyle().Foreground(Catppuccin.Subtext0).Render(fmt.Sprintf("Words completed: %d", m.currentIdx)) +
			"\n" +
			lipgloss.NewStyle().Foreground(Catppuccin.Blue).Bold(true).Render(fmt.Sprintf("WPM: %.1f", wpm)) +
			"\n" +
			lipgloss.NewStyle().Foreground(Catppuccin.Red).Render(fmt.Sprintf("Mistakes: %d", m.mistakes)) +
			"\n\n" +
			ButtonStyle.Render(" Play Again (Enter) ") +
			padBottom,
	)
}

func settingsView(m Model) string {
	padTop := "\n\n\n"
	padBottom := "\n\n\n"
	return WindowStyle.Width(100).Render(
		padTop +
			TitleStyle.Render("Settings") +
			"\n\n" +
			lipgloss.NewStyle().Foreground(Catppuccin.Subtext0).Render("Time Limit: "+fmt.Sprintf("%d seconds", int(m.timeLimit.Seconds()))) +
			"\n" +
			lipgloss.NewStyle().Foreground(Catppuccin.Subtext0).Render("Use [↑][↓] to adjust, [Esc] to return") +
			padBottom,
	)
}
