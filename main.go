package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	//"path/filepath" // Added for potential future path handling, used in find.go
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

	// --- New Flags ---
	findKeywords := flag.String("find", "", "Keywords to find (comma separated). Matching URLs are highlighted green and saved to Find-<keywords>.txt")
	findXKeywords := flag.String("findX", "", "Keywords that *all* must exist in URL (comma separated). Matching URLs are highlighted green and saved to FindX-<keywords>.txt")
	// --- End New Flags ---

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

	// Process URLs (Original Shortening/Variation Logic)
	if !*quietMode {
		fmt.Printf("%s[*] Processing URLs...%s\n", colorCyan, colorReset)
	}
	shortenedURLs := processURLs(urls, *delimiters, *noDuplicates, *splitPath, *appendString, appendStrings)

	// --- Start New Find/FindX Processing ---
	var foundUrlsMap map[string]bool
	var foundXUrlsMap map[string]bool
	var findOutputFile string
	var findXOutputFile string
	var findErr error // Use a separate error variable for find operations

	// Process --find
	if *findKeywords != "" {
		if !*quietMode {
			fmt.Printf("%s[*] Finding URLs containing any of: [%s]%s\n", colorCyan, *findKeywords, colorReset)
		}
		// Pass shortenedURLs to find functions
		foundUrlsMap = processFind(shortenedURLs, *findKeywords)
		if len(foundUrlsMap) > 0 {
			findOutputFile = generateOutputFileName("Find", *findKeywords)
			findErr = saveUrlsToFile(findOutputFile, foundUrlsMap, *quietMode) // Call function from find.go
			if findErr != nil {
				fmt.Fprintf(os.Stderr, "%sError saving --find results: %v%s\n", colorRed+bold, findErr, colorReset)
				// Decide if you want to os.Exit(1) here or just continue
			}
		} else if !*quietMode {
			fmt.Printf("%s[*] No URLs matched --find criteria.%s\n", colorYellow, colorReset)
		}
	} else {
		foundUrlsMap = make(map[string]bool) // Initialize empty map if flag not used
	}

	// Process --findX
	if *findXKeywords != "" {
		if !*quietMode {
			fmt.Printf("%s[*] Finding URLs containing all of: [%s]%s\n", colorCyan, *findXKeywords, colorReset)
		}
		// Pass shortenedURLs to find functions
		foundXUrlsMap = processFindX(shortenedURLs, *findXKeywords)
		if len(foundXUrlsMap) > 0 {
			findXOutputFile = generateOutputFileName("FindX", *findXKeywords)
			findErr = saveUrlsToFile(findXOutputFile, foundXUrlsMap, *quietMode) // Call function from find.go
			if findErr != nil {
				fmt.Fprintf(os.Stderr, "%sError saving --findX results: %v%s\n", colorRed+bold, findErr, colorReset)
				// Decide if you want to os.Exit(1) here or just continue
			}
		} else if !*quietMode {
			fmt.Printf("%s[*] No URLs matched --findX criteria.%s\n", colorYellow, colorReset)
		}
	} else {
		foundXUrlsMap = make(map[string]bool) // Initialize empty map if flag not used
	}
	// --- End New Find/FindX Processing ---

	// Output to console if not in quiet mode (MODIFIED LOOP)
	if !*quietMode {
		fmt.Printf("%s[*] Generated %d variations:%s\n", colorCyan, len(shortenedURLs), colorReset)
		for _, url := range shortenedURLs {
			// Check if the URL matched either find condition
			isFindMatch := foundUrlsMap[url]  // Check if key exists (true if matched)
			isFindXMatch := foundXUrlsMap[url] // Check if key exists (true if matched)

			if isFindMatch || isFindXMatch {
				// Print in green if it matched either
				fmt.Printf("%s%s%s\n", colorGreen, url, colorReset)
			} else {
				// Print normally otherwise
				fmt.Println(url)
			}
		}
	}

	// Write to output file if specified (using the ORIGINAL shortenedURLs list)
	// Note: The --o flag saves *all* generated URLs, not just the found ones.
	if *outputFile != "" {
		err = writeLines(*outputFile, shortenedURLs) // Use original writeLines
		if err != nil {
			fmt.Fprintf(os.Stderr, "%sError writing to output file '%s': %v%s\n", colorRed+bold, *outputFile, err, colorReset)
			os.Exit(1) // Exit on primary output file error
		}
		// Success message shown even in quiet mode if output file is used
		fmt.Printf("%s[+] Successfully wrote %d URLs to %s%s\n", colorGreen+bold, len(shortenedURLs), *outputFile, colorReset)
	} else if *quietMode {
		// Adjusted quiet message to mention find results if any were saved
		findMsg := ""
		if findOutputFile != "" {
			findMsg += fmt.Sprintf(" Saved --find results to %s.", findOutputFile)
		}
		if findXOutputFile != "" {
			findMsg += fmt.Sprintf(" Saved --findX results to %s.", findXOutputFile)
		}
		if findMsg == "" && (*findKeywords != "" || *findXKeywords != "") {
			findMsg = " No matching URLs found for --find/--findX."
		}

		fmt.Printf("%s[+] Processing complete. %d variations generated (output suppressed).%s%s\n", colorGreen+bold, len(shortenedURLs), findMsg, colorReset) // adjusted message
	} else if len(shortenedURLs) > 0 && !*quietMode {
		// Indicate console output done
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
	fmt.Printf("%s║ %s👾 Developed by:%s %sTeam HyperGod-X%s   %s📦 Version:%s %s1.1.0%s   %s📝 License:%s %sMIT%s                             %s║%s\n", // Consider bumping version
		colorWhite, colorPurple+bold, colorReset, colorPurple+italic, colorReset,
		colorBlue+bold, colorReset, colorBlue+bold, colorReset, // Updated Version to 1.1.0
		colorGreen+bold, colorReset, colorGreen+bold, colorReset,
		colorWhite, colorReset)

	// Bottom border
	fmt.Printf("%s╚════════════════════════════════════════════════════════════════════════════════════════════════════╝%s\n", colorWhite, colorReset)
}

// Displays the help message for command-line arguments (UPDATED)
func showHelp() {
	// Ensure help format aligns with flags
	fmt.Printf("\n%sUsage:%s\n", bold, colorReset)
	fmt.Println("  urlshort -f <input-file> [options]")
	fmt.Printf("\n%sOptions:%s\n", bold, colorReset) // Use Printf for colors here too
	fmt.Println("  -f string     Input file containing URLs (required)")
	fmt.Println("  -o string     Output file to write shortened URLs")
	fmt.Println("  -x string     Delimiters to use for splitting parameters (comma separated, e.g., \"&,=\") (default \"=\")")
	fmt.Println("  -p            Split URLs at path segments (/) as well")
	fmt.Println("  -a string     String to append to each generated variation")
	fmt.Println("  -F string     File containing strings to append (one per line, overrides -a)")
	fmt.Println("  -D            Remove duplicate generated URLs")
	fmt.Println("  -Q            Quiet mode (suppress banner and URL output to console, only show errors and final success message)")
	// --- Additions for Find/FindX ---
	fmt.Println("  --find string Keywords to find (comma separated). Highlights matches and saves to Find-<keywords>.txt")
	fmt.Println("  --findX string Keywords where *all* must exist in URL (comma separated). Highlights matches and saves to FindX-<keywords>.txt")
	// --- End Additions ---
	fmt.Println("  -h            Show this help message")
	fmt.Printf("\n%sExamples:%s\n", bold, colorReset) // Use Printf for colors
	fmt.Println("  urlshort -f urls.txt -o shortened.txt -x \"&,=\" -p -F payloads.txt -D")
	fmt.Println("  urlshort -f urls.txt --find \"wp-json,api\"") // Added example for find
	fmt.Println("  urlshort -f urls.txt --findX \"user,token\" -o results.txt") // Added example for findX
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
		trimmed := strings.TrimSpace(d)
		if trimmed != "" {
			delimList = append(delimList, trimmed)
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
	processed := make(map[string]bool) // Keep track of processed URLs to avoid loops
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
			currentPrefix := ""
			for i := 0; i < len(parts)-1; i++ {
				currentPrefix += parts[i] + delim
				if !variations[currentPrefix] {
					variations[currentPrefix] = true
					if !processed[currentPrefix] {
						queue = append(queue, currentPrefix) // Add to queue for further splitting
						processed[currentPrefix] = true
					}
				}
			}
			// Check the full original string split/rejoined as well (though it should be currentURL)
			fullSplitRejoined := strings.Join(parts, delim)
			if !variations[fullSplitRejoined] {
				variations[fullSplitRejoined] = true
				if !processed[fullSplitRejoined] {
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

// NOTE: The functions processFind, processFindX, parseKeywords,
// generateOutputFileName, and saveUrlsToFile should be in the 'find.go' file
// alongside this 'main.go' file in the same package ('main').
