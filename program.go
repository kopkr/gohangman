// ----- Word list provided by -----
// http://www.desiquintans.com/nounlist

// ----- Get single unbuffered key input and hide typed keys -----
// https://stackoverflow.com/questions/14094190/function-similar-to-getchar#answer-17278730

// ----- On Ctrl+C interrupt (to return key show on terminal afterwards) -----
// https://stackoverflow.com/questions/11268943/is-it-possible-to-capture-a-ctrlc-signal-and-run-a-cleanup-function-in-a-defe
// https://gobyexample.com/signals
// https://gobyexample.com/exit

// ----- Clear screen -----
// https://rosettacode.org/wiki/Terminal_control/Clear_the_screen#Go
// https://stackoverflow.com/questions/22891644/how-can-i-clear-the-terminal-screen-in-go
// https://golang.org/pkg/runtime/

// ---- Misc -----
// https://stackoverflow.com/questions/10485743/contains-method-for-a-slice

package main

import (
	"fmt"
	"os"
	"os/exec"
	"log"
	"strings"
	"os/signal"
	"syscall"
	"runtime"
)

// Check if slice (word and guessed lists) has a specific character.
func hasKey(slice []byte, key byte) bool {
   for _, a := range slice {
      if strings.ToUpper(string(a)) == strings.ToUpper(string(key)) {
         return true
      }
   }
   return false
}

// Main gameplay loop.
func playGame(diff int) {
	// TODO: words from wordfile
	wordlist, err := os.Open("words.txt")
	if err != nil {log.Fatal(err)}
	defer wordlist.Close();

	// --------- Prototype ---------- Does not use wordlist
	word := strings.ToUpper("Dictionary")
	letters := []byte(word)
	var guessed []byte;
	var message string = "How's it hanging?";
	var status int = 0;
	//var won bool = false;
	for {
		clear()
		banner(false)
		hangman(message,status)

		if status >= 8 {
			message="RIP, dude."
			restart(message, status, diff)
		}

		//fmt.Println("Status:",status,"Difficulty:",diff)
		fmt.Print("                 ")
		var hidden []byte;
		for i:=0;i<len(letters);i++{
			if hasKey(guessed, letters[i]) {
				fmt.Print(strings.ToUpper(string(letters[i])))
				fmt.Print(" ")
			} else if letters[i] == '-' {
				guessed = append(guessed, '-')
				fmt.Print("- ")
			} else {
				fmt.Print("_ ")
				hidden = append(hidden, '_')
			}
		}

		if !hasKey(hidden, '_') {
			message="You survived. Play again?"
			restart(message,status,diff)
		}
		fmt.Println()
		fmt.Println()
		fmt.Println("                       ( Press key )")
		fmt.Println()
		fmt.Println()

		var keys []byte;
		isLetter := false;
		isNumber := false;
		keys = make([]byte, 1)
		os.Stdin.Read(keys)
		key := keys[0]

		switch {
			case key >= 'a' && key <= 'z' || key >= 'A' && key <='Z':
				isLetter = true
			case key >= '0' && key <= '9':
				isNumber = true
			case key == 27:
				main()
			default:
		}

		guess:=strings.ToUpper(string(key))

		switch {
			case hasKey(letters, key) && hasKey(guessed, key) && isLetter:
				message="You already got \""+guess+"\", bro."
				status++;
			case hasKey(letters, key) && !hasKey(guessed, key) && isLetter:
				message="Yep! It has \""+guess+"\"."
				guessed = append(guessed, key)
			case !isLetter && !isNumber:
				message="Not a letter!"
			case isNumber:
				message="That's... a number..."
			default:
				message="Nope! No \""+guess+"\"."
				status++;
		}
	}
}

// Restart or return to main menu once the game after game is finished.
func restart(message string, status, diff int) {
	clear()
	banner(false)
	hangman(message,status)

	fmt.Println("\n")
	fmt.Print("                         ( Press R to restart ) \n")
	fmt.Print("                         ( Press Esc or Q to go to menu ) ")

	for {
		var keys []byte;
		keys = make([]byte, 1)
		os.Stdin.Read(keys)
		key := keys[0]
		switch {
			case key == 'R' || key == 'r':
				status=0
				message=""
				playGame(diff)
			case key == 27 || key == 'Q' || key == 'q':
				main()
			default:
		}
	}
}

// Main menu selection screen for difficulty, help and exiting.
func chooseDifficulty() int{
	var sel byte
	var diff int
	help := false
	diffSel:for {
		clear()
		banner(true)
		fmt.Println()
		fmt.Println("  Choose difficulty:")
		fmt.Println("  ----------------")
		fmt.Println("  1 - Normal")
		fmt.Println("  ----------------")
		fmt.Println("  H - Help")
		fmt.Println("  Esc - Exit")
		if help {
			fmt.Println()
			fmt.Println("  Help:")
			fmt.Println("  Press Esc anytime to end program.")
		}
		var key []byte = make([]byte, 1)
		os.Stdin.Read(key)
		sel  = key[0]
		switch {
			case sel == '1':
				diff=1
				break diffSel
			case sel == 'h'|| sel == 'H':
				help = true;
			case sel == 27:
				escape()
			default:
		}
	}
	return diff
}

// Prints "Hangman" banner.
func banner(menu bool) {
	fmt.Println(" _________________________________________________________________")
	fmt.Println("    _    _   ")
	fmt.Println("   | |  | |")
	fmt.Println("   | |__| |   __ _   _ __     __ _   _ __ ___     __ _   _ __")
	fmt.Println("   |  __  |  / _` | | '_ \\   / _` | | '_ ` _ \\   / _` | | '_ \\")
	fmt.Println("   | |  | | | (_| | | | | | | (_| | | | | | | | | (_| | | | | |")
	fmt.Println("   |_|  |_|  \\__,_| |_| |_|  \\__, | |_| |_| |_|  \\__,_| |_| |_|")
	fmt.Println("                              __/ |")
	if menu {
		fmt.Println("                             |___/")
	} else {
		fmt.Println("                             |___/           Press Esc to quit ")
	}
	fmt.Println(" _________________________________________________________________")
}

// Prints hangman art depending on game status.
func hangman(message string, status int) {
	fmt.Println()
	fmt.Println()
	switch {
		case status == 0:
			fmt.Println()
			spaces(message)
			fmt.Println(message)
			fmt.Println("                          \\")
			fmt.Println()
			fmt.Println("                       ¯\\_(ツ)_/¯")
			fmt.Println("                            |")
			fmt.Println("                            |")
			fmt.Println("                           /\\")
			fmt.Println("                          |  \\")
			fmt.Println("                         ☐☐☐☐☐☐☐☐☐┉┉┉┉┉┉┉☐☐☐☐☐☐☐☐☐")
			fmt.Println()
		case status == 1:
			fmt.Println("                                    _____________")
			spaces(message)
			fmt.Println(message,"       |_____________|")
			fmt.Println("                          \\                \\\\ ||")
			fmt.Println("                                            \\\\||")
			fmt.Println("                       ¯\\_(ツ)_/¯            \\||")
			fmt.Println("                            |                 ||")
			fmt.Println("                            |                 ||")
			fmt.Println("                           /\\                 ||")
			fmt.Println("                          |  \\                ||")
			fmt.Println("                         ☐☐☐☐☐☐☐☐☐┉┉┉┉┉┉┉☐☐☐☐☐☐☐☐☐")
			fmt.Println()
		case status == 2:
			fmt.Println("                                    _____________")
			spaces(message)
			fmt.Println(message,"       |_____________|")
			fmt.Println("                         \\           _┃    \\\\ ||")
			fmt.Println("                                    ('')    \\\\||")
			fmt.Println("                       ¯\\_(ツ)_/¯            \\||")
			fmt.Println("                            |                 ||")
			fmt.Println("                            |                 ||")
			fmt.Println("                           /\\                 ||")
			fmt.Println("                          |  \\                ||")
			fmt.Println("                         ☐☐☐☐☐☐☐☐☐┉┉┉┉┉┉┉☐☐☐☐☐☐☐☐☐")
			fmt.Println()
		case status == 3:
			fmt.Println("                                    _____________")
			spaces(message)
			fmt.Println(message,"       |_____________|")
			fmt.Println("                         \\           _┃    \\\\ ||")
			fmt.Println("                                    ('')    \\\\||")
			fmt.Println("                       ¯\\_(ツ)_/¯    |       \\||")
			fmt.Println("                            |        |        ||")
			fmt.Println("                            |                 ||")
			fmt.Println("                           /\\                 ||")
			fmt.Println("                          |  \\                ||")
			fmt.Println("                         ☐☐☐☐☐☐☐☐☐┉┉┉┉┉┉┉☐☐☐☐☐☐☐☐☐")
			fmt.Println()
		case status == 4:
			fmt.Println("                                    _____________")
			spaces(message)
			fmt.Println(message,"       |_____________|")
			fmt.Println("                         \\           _┃    \\\\ ||")
			fmt.Println("                                    ('')    \\\\||")
			fmt.Println("                       ¯\\_(ツ)_/¯    |\\      \\||")
			fmt.Println("                            |        | \\      ||")
			fmt.Println("                            |                 ||")
			fmt.Println("                           /\\                 ||")
			fmt.Println("                          |  \\                ||")
			fmt.Println("                         ☐☐☐☐☐☐☐☐☐┉┉┉┉┉┉┉☐☐☐☐☐☐☐☐☐")
			fmt.Println()
		case status == 5:
			fmt.Println("                                    _____________")
			spaces(message)
			fmt.Println(message,"       |_____________|")
			fmt.Println("                         \\           _┃    \\\\ ||")
			fmt.Println("                                    ('')    \\\\||")
			fmt.Println("                       ¯\\_(ツ)_/¯   /|\\      \\||")
			fmt.Println("                            |      / | \\      ||")
			fmt.Println("                            |                 ||")
			fmt.Println("                           /\\                 ||")
			fmt.Println("                          |  \\                ||")
			fmt.Println("                         ☐☐☐☐☐☐☐☐☐┉┉┉┉┉┉┉☐☐☐☐☐☐☐☐☐")
			fmt.Println()
		case status == 6:
			fmt.Println("                                    _____________")
			spaces(message)
			fmt.Println(message,"       |_____________|")
			fmt.Println("                         \\           _┃    \\\\ ||")
			fmt.Println("                                    ('')    \\\\||")
			fmt.Println("                       ¯\\_(ツ)_/¯   /|\\      \\||")
			fmt.Println("                            |      / | \\      ||")
			fmt.Println("                            |         \\       ||")
			fmt.Println("                           /\\          \\      ||")
			fmt.Println("                          |  \\                ||")
			fmt.Println("                         ☐☐☐☐☐☐☐☐☐┉┉┉┉┉┉┉☐☐☐☐☐☐☐☐☐")
			fmt.Println()
		case status == 7:
			fmt.Println("                                    _____________")
			spaces(message)
			fmt.Println(message,"       |_____________|")
			fmt.Println("                         \\           _┃    \\\\ ||")
			fmt.Println("                                    ('')    \\\\||")
			fmt.Println("                       ¯\\_(ツ)_/¯   /|\\      \\||")
			fmt.Println("                            |      / | \\      ||")
			fmt.Println("                            |       / \\       ||")
			fmt.Println("                           /\\      /   \\      ||")
			fmt.Println("                          |  \\                ||")
			fmt.Println("                         ☐☐☐☐☐☐☐☐☐┉┉┉┉┉┉┉☐☐☐☐☐☐☐☐☐")
			fmt.Println()
		case status == 8:
			fmt.Println("                                    _____________")
			spaces(message)
			fmt.Println(message,"       |_____________|")
			fmt.Println("                         \\           _┃    \\\\ ||")
			fmt.Println("                                    ('')    \\\\||")
			fmt.Println("                       ¯\\_(ツ)_/¯   /|\\      \\||")
			fmt.Println("                            |      / | \\      ||")
			fmt.Println("                            |       / \\       ||")
			fmt.Println("                           /\\      /   \\      ||")
			fmt.Println("                          |  \\                ||")
			fmt.Println("                         ☐☐☐☐☐☐☐☐☐|     |☐☐☐☐☐☐☐☐☐")
			fmt.Println("                                  |     |")
	}
	fmt.Println()

}





// Counts how many spaces need to be inserted before message for Unicode art formatting to stay aesthetic
func spaces(msg string) {
	space:=27-len(msg)
	for i:=0;i < space;i++ {
		fmt.Print(" ")
	}
}

// Main function
func main() {
	diff := chooseDifficulty()
	playGame(diff)
}

// Init function, runs first
// Establishes signals (so you can safely interrupt program with ctrl+C)
// Establishes terminal related stuff (get unbuffered key press. No need to press Enter after every line).
func init() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go onCtrlC(sigs,done)
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
}

// Handles interrupting with Ctrl+C
func onCtrlC(sigs chan os.Signal, done chan bool) {
	sig := <-sigs
	fmt.Println()
	fmt.Println(sig)
	done <- true
	cleanRun()
	clear()
	os.Exit(1)
}

// Shut it down.
func escape() {
	fmt.Println("Exiting...")
	cleanRun()
	clear()
	os.Exit(0)
}

// Re-enables command prompt to show key presses again, so you don't have to reopen terminal.
func cleanRun() {
	exec.Command("stty", "-F", "/dev/tty", "echo").Run()
}

// Clears screen for tidyness. OS specific clear commands (Windows, Linux, Mac)
func clear() {
	goos := runtime.GOOS
	switch {
	case goos == "linux" || goos == "darwin":
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	case goos == "windows":
		cmd := exec.Command("cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	default:
		fmt.Println("Screen clear not supported on your OS.")
		fmt.Println("Please contact author.")
		fmt.Println()
	}
}
