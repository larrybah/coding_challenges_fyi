#include <stdio.h>   // For input/output functions (e.g., printf, fgetc)
#include <ctype.h>   // For character classification (e.g., isspace)
#include <stddef.h>  // For NULL and size_t

// Function to count lines, words, and characters in a file or stdin
void word_count(FILE *file, long *lines, long *words, long *chars) {
    int c;          // Variable to store each character read
    int in_word = 0; // Flag to track if we are inside a word

    // Initialize counters
    *lines = 0;
    *words = 0;
    *chars = 0;

    // Read the file character by character
    while ((c = fgetc(file)) != EOF) {
        // Increment character count
        (*chars)++;

        // Check for newline to count lines
        if (c == '\n') {
            (*lines)++;
        }

        // Check if the current character is a whitespace
        if (isspace(c)) {
            in_word = 0; // Not inside a word
        }
        // If the current character is not a whitespace and we were not already in a word
        else if (!in_word) {
            in_word = 1; // We are now inside a word
            (*words)++;  // Increment word count
        }
    }

    // If the file does not end with a newline, increment line count
    if (c != '\n' && *chars != 0) {
        (*lines)++;
    }
}

int main(int argc, char *argv[]) {
    FILE *file;      // File pointer
    long lines = 0;  // Line counter
    long words = 0;  // Word counter
    long chars = 0;  // Character counter

    // Check if a filename is provided as an argument
    if (argc == 1) {
        // No filename provided, read from stdin
        word_count(stdin, &lines, &words, &chars);
        printf("%ld %ld %ld\n", lines, words, chars);
    }
    else {
        // Open the file provided as an argument
        file = fopen(argv[1], "r");
        if (file == NULL) {
            fprintf(stderr, "Error: Cannot open file '%s'\n", argv[1]);
            return 1;
        }

        // Count lines, words, and characters in the file
        word_count(file, &lines, &words, &chars);
        printf("%ld %ld %ld %s\n", lines, words, chars, argv[1]);

        // Close the file
        fclose(file);
    }

    return 0;
}