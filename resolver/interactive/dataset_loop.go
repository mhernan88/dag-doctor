package interactive

import (
	"fmt"
	"log/slog"

	"github.com/eiannone/keyboard"
	"github.com/enescakir/emoji"
)

func CheckDataset(dataset string, l *slog.Logger) (string, error) {
	fmt.Printf("|---> inspecting dataset: %s (c or q to exit)\n", dataset)
	fmt.Printf("|---> open '%s' and check if content is correct...\n", dataset)

	err := keyboard.Open()
	if err != nil {
		return "err", err
	}
	defer keyboard.Close()

	for {
		fmt.Println("|---> output correct? (y/n):")
		char, _, err := keyboard.GetKey()
		if err != nil {
			return "err", err
		}

		switch char {
		case 'y', 'Y':
			fmt.Printf(
				"|---> %v dataset '%s' maked OK\n",
				emoji.CheckMarkButton,
				dataset)
			l.Info("input was 'y', returning true, nil")
			return "ok", nil
		case 'n', 'N':
			fmt.Printf(
				"|---> %v dataset '%s' maked ERR\n",
				emoji.CrossMarkButton,
				dataset)
			l.Info("input was 'n', returning false, nil")
			return "err", nil
		case 'c', 'C', 'q', 'Q':
			return "aborted", nil
		default:
			fmt.Printf("|---> invalid input; options are: y, Y, n, N; c, q to quit")
		}
	}
}
