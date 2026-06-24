package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

// main is the entry point of the program.
// It orchestrates the password generation process.
func main() {
	// Seed the random number generator with the current time
	// to ensure different results on each run.
	rand.Seed(time.Now().UnixNano())

	// Print a welcome message to the user.
	fmt.Println("Welcome to the Go Password Generator!")
	fmt.Println("------------------------------------")

	// Prompt the user to enter the desired password length.
	fmt.Print("Enter the length of the password: ")

	// Read the user input from standard input.
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	// Trim whitespace and convert the input to an integer.
	input = input[:len(input)-1] // Remove the trailing newline
	length, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("Invalid input. Please enter a number.")
		return
	}

	// Check if the length is valid (greater than 0).
	if length <= 0 {
		fmt.Println("Password length must be greater than 0.")
		return
	}

	// Generate the password using the specified length.
	password := generatePassword(length)

	// Print the generated password.
	fmt.Println("Generated Password:", password)
}

// generatePassword generates a random password of the specified length.
// It uses a combination of uppercase letters, lowercase letters, digits, and special characters.
func generatePassword(length int) string {
	// Define the character sets for the password.
	// These sets include uppercase letters, lowercase letters, digits, and special characters.
	uppercase := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lowercase := "abcdefghijklmnopqrstuvwxyz"
	digits := "0123456789"
	specialChars := "!@#$%^&*()_+-=[]{}|;:,.<>?"

	// Combine all character sets into a single string for random selection.
	allChars := uppercase + lowercase + digits + specialChars

	// Initialize a byte slice to hold the password characters.
	password := make([]byte, length)

	// Generate random characters for the password.
	for i := 0; i < length; i++ {
		// Select a random index from the combined character set.
		randomIndex := rand.Intn(len(allChars))
		// Append the character at the random index to the password.
		password[i] = allChars[randomIndex]
	}

	// Convert the byte slice to a string and return it.
	return string(password)
}