package cmd

import (
	"fmt"
	"os"

	"github.com/eiannone/keyboard"
	"github.com/enescakir/emoji"
)

func (ui *UI) CheckDataset(dataset string) (bool, error) {
	fmt.Printf("|---> inspecting dataset: %s (c or q to exit)\n", dataset)

	fmt.Printf("|---> open '%s' and check if content is correct...\n", dataset)

	err := keyboard.Open()
	if err != nil {
		return false, err
	}
	defer keyboard.Close()

	for {
		fmt.Println("|---> output correct? (y/n):")
		ui.l.Trace("reading keyboard input")
		char, _, err := keyboard.GetKey()
		if err != nil {
			return false, err
		}

		switch char {
		case 'y', 'Y':
			fmt.Printf(
				"|---> %v dataset '%s' maked OK\n",
				emoji.CheckMarkButton,
				dataset)
			ui.l.Trace("input was 'y', returning true, nil")
			return true, nil
		case 'n', 'N':
			fmt.Printf(
				"|---> %v dataset '%s' maked ERR\n",
				emoji.CrossMarkButton,
				dataset)
			ui.l.Trace("input was 'n', returning false, nil")
			return false, nil
		case 'c', 'C', 'q', 'Q':
			os.Exit(0)
		default:
			fmt.Printf("|---> invalid input; options are: y, Y, n, N; c, q to quit")
		}
	}
}
