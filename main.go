package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

// ANSI Color/Style constants
const (
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
	colorReset  = "\033[0m"
	bold        = "\033[1m"
	italic      = "\033[3m"
)

func main() {
	// Define command-line flags
	inputFile := flag.String("f", "", "Input file containing URLs (required)")
	outputFile := flag.String("o", "", "Output file to write shortened URLs")
	delimiters := flag.String("x", "=", "Delimiters to use for shortening (comma separated)")
	noDuplicates := flag.Bool("D", false, "Remove duplicate URLs")
	quietMode := flag.Bool("Q", false, "Quiet mode (suppress banner and URL output to console, only show errors and final success message)")
	splitPath := flag.Bool("p", false, "Split URLs at path segments (/)")
	appendString := flag.String("a", "", "String to append to each generated variation")
	appendFile := flag.String("F", "", "File containing strings to append (one per line, overrides -a)")
	help := flag.Bool("h", false, "Show help message")

	// Set custom usage message
	flag.Usage = showHelp

	// Parse the flags
	flag.Parse()

    // Show help and exit if -h is provided *after* parsing
	if *help {
        showBanner() // Show banner even when showing help
		showHelp()
		os.Exit(0)
	}

	// Show the banner only if not in quiet mode
	if !*quietMode {
		showBanner()
	}

	// Validate input: Input file is required
	if *inputFile == "" {
		// Print error to stderr
		fmt.Fprintf(os.Stderr, "\n%sError: Input file (-f) is required.%s\n", colorRed+bold, colorReset)
		showHelp() // Show help message on error
		os.Exit(1)
	}

	// Read input file
	urls, err := readLines(*inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%sError reading input file '%s': %v%s\n", colorRed+bold, *inputFile, err, colorReset)
		os.Exit(1)
	}
	if !*quietMode && len(urls) > 0 {
		fmt.Printf("%s[*] Read %d URLs from %s%s\n", colorGreen, len(urls), *inputFile, colorReset)
	} else if len(urls) == 0 {
		fmt.Printf("%s[*] Input file '%s' is empty or contains no valid lines.%s\n", colorYellow, *inputFile, colorReset)
		os.Exit(0) // Exit gracefully if input is empty
	}

	// Read append file if specified
	var appendStrings []string
	if *appendFile != "" {
		appendStrings, err = readLines(*appendFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%sError reading append file '%s': %v%s\n", colorRed+bold, *appendFile, err, colorReset)
			os.Exit(1)
		}
		if !*quietMode && len(appendStrings) > 0 {
			fmt.Printf("%s[*] Read %d strings to append from %s%s\n", colorGreen, len(appendStrings), *appendFile, colorReset)
		}
	}

	// Process URLs
	if !*quietMode {
		fmt.Printf("%s[*] Processing URLs...%s\n", colorCyan, colorReset)
	}
	shortenedURLs := processURLs(urls, *delimiters, *noDuplicates, *splitPath, *appendString, appendStrings)

	// Output to console if not in quiet mode
	if !*quietMode {
		fmt.Printf("%s[*] Generated %d variations:%s\n", colorCyan, len(shortenedURLs), colorReset)
		for _, url := range shortenedURLs {
			fmt.Println(url)
		}
	}

	// Write to output file if specified
	if *outputFile != "" {
		err = writeLines(*outputFile, shortenedURLs)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%sError writing to output file '%s': %v%s\n", colorRed+bold, *outputFile, err, colorReset)
			os.Exit(1)
		}
		// Success message is shown even in quiet mode if output file is used
		fmt.Printf("%s[+] Successfully wrote %d URLs to %s%s\n", colorGreen+bold, len(shortenedURLs), *outputFile, colorReset)
	} else if *quietMode {
		// If quiet mode and no output file, mention completion.
		fmt.Printf("%s[+] Processing complete. %d variations generated (output suppressed).%s\n", colorGreen+bold, len(shortenedURLs), colorReset)
	} else if len(shortenedURLs) > 0 && !*quietMode{
        // If not quiet and no output file, just indicate console output is done.
        fmt.Printf("%s[+] Output displayed above.%s\n", colorGreen+bold, colorReset)
    }
}

// Displays the application banner
func showBanner() {
    // Top border
    fmt.Printf("%s╔════════════════════════════════════════════════════════════════════════════════════════════════════╗%s\n", colorWhite, colorReset)
    
    // Logo lines (6 lines)
    fmt.Printf("%s║%s  ██╗  ██╗██████╗ ██╗      ███████╗██╗  ██╗ ██████╗ ██████╗ ████████╗                               %s║%s\n", colorWhite, bold+colorRed, colorReset+colorWhite, colorReset)
    fmt.Printf("%s║%s  ██║  ██║██╔══██╗██║      ██╔════╝██║  ██║██╔═══██╗██╔══██╗╚══██╔══╝                               %s║%s\n", colorWhite, bold+colorRed, colorReset+colorWhite, colorReset)
    fmt.Printf("%s║%s  ██║  ██║██████╔╝██║█████╗███████╗███████║██║   ██║██████╔╝   ██║                                  %s║%s\n", colorWhite, bold+colorRed, colorReset+colorWhite, colorReset)
    fmt.Printf("%s║%s  ██║  ██║██╔══██╗██║╚════╝╚════██║██╔══██║██║   ██║██╔══██╗   ██║                                  %s║%s\n", colorWhite, bold+colorRed, colorReset+colorWhite, colorReset)
    fmt.Printf("%s║%s  ╚██████╔╝██║  ██║███████╗███████║██║  ██║╚██████╔╝██║  ██║   ██║                                  %s║%s\n", colorWhite, bold+colorRed, colorReset+colorWhite, colorReset)
    fmt.Printf("%s║%s   ╚═════╝ ╚═╝  ╚═╝╚══════╝╚══════╝╚═╝  ╚═╝ ╚═════╝ ╚═╝  ╚═╝   ╚═╝                                  %s║%s\n", colorWhite, bold+colorRed, colorReset+colorWhite, colorReset)
    
    // Separator
    fmt.Printf("%s╠════════════════════════════════════════════════════════════════════════════════════════════════════╣%s\n", colorWhite, colorReset)
    
    // Tool Type line
    fmt.Printf("%s║ %s💡 Tool Type:%s %sAdvanced URL Shortener & Parameter Generator%s                                         %s║%s\n", 
        colorWhite, colorCyan+bold, colorReset, colorCyan+italic, colorReset, colorWhite, colorReset)
    
    // Use Case line
    fmt.Printf("%s║ %s💡 Use Case:%s  %sSecurity Testing • Web Dev Utility • Payload Injector%s                                %s║%s\n", 
        colorWhite, colorYellow+bold, colorReset, colorYellow+italic, colorReset, colorWhite, colorReset)
    
    // Separator
    fmt.Printf("%s╠════════════════════════════════════════════════════════════════════════════════════════════════════╣%s\n", colorWhite, colorReset)
    
    // Developer/Version/License line
    fmt.Printf("%s║ %s👾 Developed by:%s %sTeam HyperGod-X%s   %s📦 Version:%s %s1.0.0%s   %s📝 License:%s %sMIT%s                             %s║%s\n", 
        colorWhite, colorPurple+bold, colorReset, colorPurple+italic, colorReset, 
        colorBlue+bold, colorReset, colorBlue+bold, colorReset, 
        colorGreen+bold, colorReset, colorGreen+bold, colorReset, 
        colorWhite, colorReset)
    
    // Bottom border
    fmt.Printf("%s╚════════════════════════════════════════════════════════════════════════════════════════════════════╝%s\n", colorWhite, colorReset)
}


// Displays the help message for command-line arguments
func showHelp() {
    // Ensure help format aligns with flags
	fmt.Printf("\n%sUsage:%s\n", bold, colorReset)
	fmt.Println("  urlshort -f <input-file> [options]")
	fmt.Println("\n%sOptions:%s", bold, colorReset) // Use Printf for colors here too
	fmt.Println("  -f string     Input file containing URLs (required)")
	fmt.Println("  -o string     Output file to write shortened URLs")
	fmt.Println("  -x string     Delimiters to use for splitting parameters (comma separated, e.g., \"&,=\") (default \"=\")")
	fmt.Println("  -p            Split URLs at path segments (/) as well")
	fmt.Println("  -a string     String to append to each generated variation")
	fmt.Println("  -F string     File containing strings to append (one per line, overrides -a)")
	fmt.Println("  -D            Remove duplicate generated URLs")
	fmt.Println("  -Q            Quiet mode (suppress banner and URL output to console, only show errors and final success message)")
	fmt.Println("  -h            Show this help message")
	fmt.Println("\n%sExample:%s", bold, colorReset) // Use Printf for colors
	fmt.Println("  urlshort -f urls.txt -o shortened.txt -x \"&,=\" -p -F payloads.txt -D")
}


// Reads all non-empty lines from a file into a slice of strings.
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text()) // Trim whitespace
		if line != "" {                           // Skip empty lines
			lines = append(lines, line)
		}
	}
	return lines, scanner.Err()
}

// Writes a slice of strings to a file, each on a new line.
func writeLines(path string, lines []string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return err // Return error immediately if write fails
		}
	}
	// Flush ensures all buffered data is written to the file
	return writer.Flush()
}

// Processes URLs based on specified options: delimiters, duplicates, path splitting, appends.
func processURLs(urls []string, delimiters string, noDuplicates bool, splitPath bool, appendString string, appendStrings []string) []string {
	// Prepare delimiters
	rawDelimList := strings.Split(delimiters, ",")
    // Filter out empty strings that might result from trailing commas, etc.
    delimList := []string{}
    for _, d := range rawDelimList {
        if d != "" {
            delimList = append(delimList, d)
        }
    }

	// Add path delimiter if requested, ensuring it's not duplicated if already present
	if splitPath {
		found := false
		for _, d := range delimList {
			if d == "/" {
				found = true
				break
			}
		}
		if !found {
			delimList = append(delimList, "/")
		}
	}
     if len(delimList) == 0 {
        // Should not happen with default "=", but handle edge case
        fmt.Fprintf(os.Stderr, "%sWarning: No valid delimiters specified. Only applying appends.%s\n", colorYellow, colorReset)
     }


	// Use map for efficient duplicate checking if needed
	seen := make(map[string]bool)
	var result []string

	for _, url := range urls {
		// Generate base variations based on delimiters
		baseVariations := generateVariations(url, delimList)

		for _, variation := range baseVariations {
			// Apply append operations to each base variation
			finalURLs := applyAppends(variation, appendString, appendStrings)

			// Add final URLs to the result list, handling duplicates if requested
			for _, finalURL := range finalURLs {
				if noDuplicates {
					if !seen[finalURL] {
						seen[finalURL] = true
						result = append(result, finalURL)
					}
				} else {
					result = append(result, finalURL)
				}
			}
		}
	}

	return result
}

// Generates variations of a URL by splitting it at given delimiters
// and taking prefixes ending at each delimiter instance. Includes the original URL.
func generateVariations(url string, delimiters []string) []string {
	// Use map to automatically handle duplicates during generation
	variations := make(map[string]bool)
	variations[url] = true // Always include the original URL

    if len(delimiters) == 0 {
         // If no delimiters, just return the original URL
         return []string{url}
    }

	// Use a queue for breadth-first processing of variations and delimiters
	queue := []string{url}
	processed := make(map[string]bool) // Keep track of processed URLs to avoid loops with recursive delimiters
    processed[url] = true

	head := 0
	for head < len(queue) {
		currentURL := queue[head]
		head++

		for _, delim := range delimiters {
            // Ensure delimiter is not empty
            if delim == "" {
                continue
            }
			parts := strings.Split(currentURL, delim)
			if len(parts) <= 1 { // No delimiter found or only one part
				continue
			}

			// Generate prefixes ending with the delimiter
			for i := 1; i < len(parts); i++ {
				prefix := strings.Join(parts[:i], delim)
                // Only add the delimiter if the original string had content after it at this point
                // Or, more simply, always add it if splitting occurred.
                prefix += delim

				if !variations[prefix] {
					variations[prefix] = true
                    if !processed[prefix]{
                        queue = append(queue, prefix) // Add to queue for further splitting
                        processed[prefix] = true
                    }
				}
			}
            // Add the full split parts as well if they weren't the original
            fullSplitRejoined := strings.Join(parts, delim) // This should be same as currentURL
             if !variations[fullSplitRejoined] {
                 variations[fullSplitRejoined] = true
                  if !processed[fullSplitRejoined]{
                        queue = append(queue, fullSplitRejoined)
                        processed[fullSplitRejoined] = true
                    }
             }
		}
	}

	// Convert map keys to slice for the final result
	result := make([]string, 0, len(variations))
	for v := range variations {
		result = append(result, v)
	}

	return result
}


// Appends strings to a base URL. Prioritizes strings from a file over a single string.
func applyAppends(baseURL string, appendString string, appendStrings []string) []string {
	// If no append operation is needed, return the base URL directly in a slice
	if appendString == "" && len(appendStrings) == 0 {
		return []string{baseURL}
	}

	var result []string

	if len(appendStrings) > 0 {
		// If a list of strings to append is provided (from -F file), use them.
		result = make([]string, len(appendStrings))
		for i, s := range appendStrings {
			result[i] = baseURL + s
		}
	} else {
        // Check appendString again, as len(appendStrings) could be 0
        if appendString != "" {
		    // Otherwise, if a single append string is provided (from -a), use it.
		    result = []string{baseURL + appendString}
        } else {
            // This case means appendString was "" and appendStrings was empty, handled above.
            // But for robustness, return the original if somehow reached.
             result = []string{baseURL}
        }
	}

	return result
}
