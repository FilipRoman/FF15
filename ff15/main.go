package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// findFiles will search the given directory for a file with the given name
// The searchName will be trimmed of it's extension and made lowercase for comparison
func findFiles(dir string, search string) ([]string, error) {
	var files []string
	searchName := strings.ToLower(filepath.Base(search))
	searchName = strings.TrimSuffix(searchName, filepath.Ext(searchName))

	// Create a slice of strings for the loading message
	// This is used to create a "Loading" animation
	loadingMessages := []string{"Searching |", "Searching /", "Searching -", "Searching \\"}

	// This index is used to cycle through the loading messages
	var loadingIndex int

	// Create a done channel so that we can stop the goroutine
	done := make(chan struct{})

	// Start a goroutine for the loading message
	go func() {
		// Run this goroutine until the done channel is closed
		for {
			select {
			case <-done:
				return
			default:
				// Get the current loading message
				loadingMessage := loadingMessages[loadingIndex]

				// Print the loading message and overwrite the previous line
				fmt.Printf("\r%s", loadingMessage)

				// Increment the loading index
				loadingIndex = (loadingIndex + 1) % len(loadingMessages)

				// Sleep for 500 milliseconds
				time.Sleep(500 * time.Millisecond)
			}
		}
	}()

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		// If there is an error, return it
		if err != nil {
			// If the error is a permission error, return nil
			// This is so that the program doesn't crash if it can't access a directory
			if os.IsPermission(err) {
				return nil
			}

			return err
		}

		// If the info is a directory, return nil
		if info.IsDir() {
			return nil
		}

		// Get the name of the file in lowercase
		fileName := strings.ToLower(info.Name())

		// If the file name contains the search name, add it to the files slice
		if strings.Contains(fileName, searchName) {
			files = append(files, path)
		}

		return nil
	})

	// Close the done channel to stop the goroutine
	close(done)

	// Print a newline character to overwrite the loading message
	fmt.Printf("\r%s", strings.Repeat(" ", len(loadingMessages[0])))

	return files, err
}

func main() {
	var dir string
	ff := `
   ___  ___  _  ____  
  / __\/ __\/ || ___| 
 / _\ / _\  | ||___ \ 
/ /  / /    | | ___) |
\/   \/     |_||____/ 
                      
`
	// Print the logo
	fmt.Println(ff)

	// Get the directory path and check if it exists
	for {
		fmt.Print("\033[32mEnter the directory path: (ex. D:\\Games\\) \033[0m")
		fmt.Scan(&dir)
		if _, err := os.Stat(dir); err != nil {
			if os.IsNotExist(err) {
				// If the directory does not exist, print a message and loop again
				fmt.Print("\033[H\033[2J")
				fmt.Println("\033[31mDirectory does not exist\033[0m")

			} else {
				// If there is an error accessing the directory, print a message and loop again
				fmt.Print("\033[H\033[2J")
				fmt.Println("\033[31mError accessing directory\033[0m")

			}
		} else {
			// If the directory exists, break out of the loop
			break
		}
	}
	var wait string
	var search string
	var oporex int
	// Loop until the user chooses to exit
	for {

		fmt.Print("\033[32mEnter the file name: (ex. virmox.txt)  \033[0m")
		fmt.Scan(&search)

		// Find all files in the directory with the given name
		files, err := findFiles(dir, search)
		if err != nil {
			// If there is an error, print a message and loop again
			fmt.Print("\033[H\033[2J")
			fmt.Println("Could not find files ")
			fmt.Scan(&wait)
		}

		// If no files were found, print a message and loop again
		if len(files) == 0 {
			fmt.Print("\033[H\033[2J")
			fmt.Println("Exactly named file not found ")
			fmt.Scan(&wait)
		} else {
			// If files were found, print a message and the files
			fmt.Print("\033[H\033[2J")
			fmt.Println("Files found: ")
			for _, file := range files {
				fmt.Println(file)
			}
		}

		fmt.Println(files)

		// Ask the user if they want to open the file or exit
		fmt.Print("\033[32m[1] to open the file \n[2] to exit\n\033[0m")
		fmt.Scan(&oporex)
		switch oporex {
		case 1:
			// Loop through the files and open the one with the exact name
			for _, file := range files {
				fileName := filepath.Base(file)
				if strings.ToLower(fileName) == strings.ToLower(search) {
					filePath := file
					dirPath := filepath.Dir(filePath)
					// Open the folder with the file selected
					cmd := exec.Command("explorer", "/select,", filePath)
					err := cmd.Start()
					if err != nil {
						// If there is an error, print a message and loop again
						fmt.Print("\033[H\033[2J")
						fmt.Println("Error opening folder: ", err, "(", dirPath, ")")
						fmt.Scan(&wait)
					} else {
						// If the folder is opened successfully, print a message and loop again
						fmt.Print("\033[H\033[2J")
						fmt.Println("Opening folder...")
						fmt.Scan(&wait)
					}
				}
			}
		case 2:
			// If the user chooses to exit, print a message, wait a second and exit
			fmt.Print("\033[H\033[2J")
			fmt.Println("Goodbye")
			time.Sleep(time.Second)
			os.Exit(0)
		}
	}
}
