package main

import (
    "context"
    "fmt"
    _ "embed"
    "strings"

    "github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// HideWindow hides the application window
func (a *App) HideWindow() {
	runtime.WindowHide(a.ctx)
}

// ShowWindow shows the application window
func (a *App) ShowWindow() {
	runtime.WindowShow(a.ctx)
}

// BIP39Logic represents BIP39 logic
type BIP39Logic struct {
	Wordlist []string
}

var bip39Logic *BIP39Logic

//go:embed wordlists/english.txt
var embeddedEnglishWordlist string

// InitBIP39 initializes BIP39 logic
func (a *App) InitBIP39() error {
	if bip39Logic != nil {
		return nil
	}

    // Parse embedded wordlist
    raw := embeddedEnglishWordlist
    if raw == "" {
        return fmt.Errorf("embedded wordlist is empty")
    }

    items := strings.Split(raw, "\n")
    wordlist := make([]string, 0, len(items))
    for _, w := range items {
        w = strings.TrimSpace(strings.TrimSuffix(w, "\r"))
        if w == "" {
            continue
        }
        wordlist = append(wordlist, w)
    }

    if len(wordlist) != 2048 {
        return fmt.Errorf("invalid wordlist length: expected 2048, got %d", len(wordlist))
    }

    bip39Logic = &BIP39Logic{Wordlist: wordlist}

	return nil
}

// GetWord returns the word at the given index
func (a *App) GetWord(index int) (string, error) {
	if bip39Logic == nil {
		if err := a.InitBIP39(); err != nil {
			return "", err
		}
	}

	if index < 0 || index >= len(bip39Logic.Wordlist) {
		return "", fmt.Errorf("invalid index: %d", index)
	}

	return bip39Logic.Wordlist[index], nil
}

// IsValidNumber checks if a number is valid (power of 2 from 1 to 1024)
func (a *App) IsValidNumber(num int) bool {
	validNumbers := []int{1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024}
	for _, valid := range validNumbers {
		if num == valid {
			return true
		}
	}
	return false
}

// CalculateWordIndex calculates word index from sum (sum - 1)
func (a *App) CalculateWordIndex(sum int) int {
	return sum - 1
}

// ValidateWordIndex validates if the calculated index is valid
func (a *App) ValidateWordIndex(index int) bool {
	return index >= 0 && index < 2048
}
