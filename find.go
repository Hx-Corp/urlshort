package main

import (
	"bufio"
	"fmt"
	"os"
	//"path/filepath"
	"strings"
)

// processFind searches URLs for *any* of the given keywords.
// It returns a map where keys are matching URLs and values are true.
func processFind(urls []string, keywords string) map[string]bool {
	if keywords == "" {
		return make(map[string]bool) // Return empty map if no keywords
	}
	keywordList := parseKeywords(keywords)
	matchingUrls := make(map[string]bool)

	for _, url := range urls {
		for _, keyword := range keywordList {
			if strings.Contains(url, keyword) {
				matchingUrls[url] = true
				break // Move to the next URL once a keyword matches
			}
		}
	}
	return matchingUrls
}

// processFindX searches URLs where *all* of the given keywords are present.
// It returns a map where keys are matching URLs and values are true.
func processFindX(urls []string, keywords string) map[string]bool {
	if keywords == "" {
		return make(map[string]bool) // Return empty map if no keywords
	}
	keywordList := parseKeywords(keywords)
	matchingUrls := make(map[string]bool)

	for _, url := range urls {
		allMatch := true
		for _, keyword := range keywordList {
			if !strings.Contains(url, keyword) {
				allMatch = false
				break // Stop checking keywords for this URL if one doesn't match
			}
		}
		if allMatch {
			matchingUrls[url] = true
		}
	}
	return matchingUrls
}

// parseKeywords splits the keyword string by commas and trims spaces.
func parseKeywords(keywords string) []string {
	rawKeywords := strings.Split(keywords, ",")
	var cleanedKeywords []string
	for _, k := range rawKeywords {
		trimmed := strings.TrimSpace(k)
		if trimmed != "" {
			cleanedKeywords = append(cleanedKeywords, trimmed)
		}
	}
	return cleanedKeywords
}

// generateOutputFileName creates a filename based on the mode ("Find" or "FindX") and keywords.
func generateOutputFileName(mode string, keywords string) string {
	keywordList := parseKeywords(keywords)
	safeKeywords := strings.Join(keywordList, "-")
	// Basic sanitization (replace common problematic characters)
	safeKeywords = strings.ReplaceAll(safeKeywords, "/", "_")
	safeKeywords = strings.ReplaceAll(safeKeywords, "\\", "_")
	safeKeywords = strings.ReplaceAll(safeKeywords, ":", "_")
	// Add more replacements if needed

	return fmt.Sprintf("%s-%s.txt", mode, safeKeywords)
}

// saveUrlsToFile writes the provided URLs (from the map keys) to the specified file.
// It returns an error if writing fails.
func saveUrlsToFile(filePath string, urlsToSave map[string]bool, quietMode bool) error {
	if len(urlsToSave) == 0 {
		if !quietMode {
			fmt.Printf("%s[*] No URLs matched the criteria for file '%s'. File not created.%s\n", colorYellow, filePath, colorReset)
		}
		return nil // Not an error, just nothing to write
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("creating file '%s': %w", filePath, err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	count := 0
	for url := range urlsToSave {
		_, err := writer.WriteString(url + "\n")
		if err != nil {
			// Try to flush before returning error to write partial data if possible
			_ = writer.Flush()
			return fmt.Errorf("writing to file '%s': %w", filePath, err)
		}
		count++
	}
	err = writer.Flush()
	if err != nil {
		return fmt.Errorf("flushing file '%s': %w", filePath, err)
	}

	// Success message even in quiet mode for file operations
	fmt.Printf("%s[+] Successfully wrote %d URLs to %s%s\n", colorGreen+bold, count, filePath, colorReset)
	return nil
}
