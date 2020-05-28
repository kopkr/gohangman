// http://www.desiquintans.com/nounlist
// https://stackoverflow.com/questions/14094190/function-similar-to-getchar#answer-17278730
// https://stackoverflow.com/questions/11268943/is-it-possible-to-capture-a-ctrlc-signal-and-run-a-cleanup-function-in-a-defe
// https://gobyexample.com/signals
// https://gobyexample.com/exit

package main

import (
	"fmt"
	"os"
	"os/exec"
	"log"
	"strings"
	"os/signal"
	"syscall"
)

func hasKey(slice []byte, key byte) bool {
   for _, a := range slice {
      if strings.ToUpper(string(a)) == strings.ToUpper(string(key)) {
         return true
      }
   }
   return false
}

func playGame(diff byte) {
	fmt.Println()
	// TODO: words from wordfile
	wordlist, err := os.Open("words.txt")
	if err != nil {log.Fatal(err)}
	defer wordlist.Close();

	// --------- Prototype ---------- Does not use wordlist
	word := strings.ToUpper("goPhEr")
	wordB := []byte(word)
	var guessed []byte;
	for {
		fmt.Print("Contains letters: ")
		for i:=0;i<len(wordB);i++{
			if hasKey(guessed, wordB[i]) {
				fmt.Print(string(wordB[i])," ")
				//fmt.Print(strings.ToUpper(string(guessed)),"\n")
			} else {
				fmt.Print("_ ")
			}
		}

		fmt.Print("Input a letter: ")

		var key []byte;
		keySel:for {
			key = make([]byte, 1)
			os.Stdin.Read(key)
			if key[0] >= 'a' && key[0] <= 'z' || key[0] >= 'A' && key[0] <='Z' {
				break keySel
			} else if key[0] == 27{
				escape()
			}
		}
		fmt.Println(strings.ToUpper(string(key[0])))

		if hasKey(wordB, key[0]) {
			if hasKey(guessed, key[0]) == false {
				guessed = append(guessed, key[0])
			}
			//fmt.Println("Word contains:",strings.ToUpper(string(key[0])))
		} else {
			//fmt.Println("Word does not contain:",strings.ToUpper(string(key[0])))
		}
	}
}

func chooseDifficulty() byte{
	fmt.Println("Choose difficulty:")
	fmt.Println("1 - Normal")
	fmt.Println("----------------")
	fmt.Println("H - Help")
	fmt.Println("Esc - Exit")

	var sel byte
	diffSel:for {
		var key []byte = make([]byte, 1)
		os.Stdin.Read(key)
		sel  = key[0]
		switch {
			case sel == '1':
				fmt.Println("Selection: "+string(sel))
				break diffSel
			case sel == 'h'|| sel == 'H':
				fmt.Println("Press Esc anytime to end program.")
			case sel == 27:
				escape()
			default:
				fmt.Println("Incorrect selection")
		}
	}
	return sel
}

func introLogo() {
	fmt.Println("  _    _   ")
	fmt.Println(" | |  | |")
	fmt.Println(" | |__| |   __ _   _ __     __ _   _ __ ___     __ _   _ __")
	fmt.Println(" |  __  |  / _` | | '_ \\   / _` | | '_ ` _ \\   / _` | | '_ \\")
	fmt.Println(" | |  | | | (_| | | | | | | (_| | | | | | | | | (_| | | | | |")
	fmt.Println(" |_|  |_|  \\__,_| |_| |_|  \\__, | |_| |_| |_|  \\__,_| |_| |_|")
	fmt.Println("                            __/ |")
	fmt.Println("                           |___/")
	fmt.Println("_________________________________________________________________")
}

/*
	fmt.Printf("   _____________\n")
	fmt.Printf("  |_____________|\n")
	fmt.Printf("     ┃     \\\\ ||\n")
	fmt.Printf("     😐     \\\\||\n")
	fmt.Printf("    /|\\      \\||\n")
	fmt.Printf("   / | \\      ||\n")
	fmt.Printf("    / \\       ||\n")
	fmt.Printf("   /   \\      ||\n")
	fmt.Printf("              ||\n")
	fmt.Printf("☐☐☐☐☐☐☐☐☐☐☐☐☐☐☐☐☐☐\n")
*/

func main() {
	introLogo()
	diff := chooseDifficulty()
	playGame(diff)
}

func init() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go onCtrlC(sigs,done)
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
}

func contains(s []int, e int) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}

func onCtrlC(sigs chan os.Signal, done chan bool) {
	sig := <-sigs
	fmt.Println()
	fmt.Println(sig)
	done <- true
	clean()
	os.Exit(1)
}

func escape() {
	fmt.Println("Exiting...")
	clean()
	os.Exit(0)
}

func clean() {
	exec.Command("stty", "-F", "/dev/tty", "echo").Run()
}
