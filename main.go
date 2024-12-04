package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func readURLsFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var urls []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		urls = append(urls, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return urls, nil
}

func downloadImage(url string, dest string) error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request for URL %s: %v", url, err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "image/webp,image/apng,image/*,*/*;q=0.8")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to fetch image from %s: %v", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to download image from %s: Status Code %d", url, resp.StatusCode)
	}

	file, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %v", dest, err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save image to %s: %v", dest, err)
	}

	fmt.Printf("Downloaded image from %s to %s\n", url, dest)
	return nil
}

func ensureDirectoryExists(directory string) error {
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		err := os.MkdirAll(directory, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create directory %s: %v", directory, err)
		}
	}
	return nil
}

func main() {
	urlFile := "urls_drawings.txt"
	directory := "downloaded_images"

	err := ensureDirectoryExists(directory)
	if err != nil {
		fmt.Println("Error ensuring directory:", err)
		return
	}

	urls, err := readURLsFromFile(urlFile)
	if err != nil {
		fmt.Println("Error reading URLs from file:", err)
		return
	}

	for i, url := range urls {
		ext := filepath.Ext(url)
		if ext == "" {
			ext = ".jpg"
		}
		fileName := fmt.Sprintf("%s/image_%d%s", directory, i+1, ext)

		err := downloadImage(url, fileName)
		if err != nil {
			fmt.Println("Error downloading image:", err)
		}
	}
}
