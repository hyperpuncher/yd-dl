package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

const apiBase = "https://cloud-api.yandex.net/v1/disk/public/resources"

type apiResponse struct {
	Type     string `json:"type"`
	Name     string `json:"name"`
	Path     string `json:"path"`
	Size     int64  `json:"size"`
	File     string `json:"file"` // direct download URL for files
	Embedded *struct {
		Limit  int           `json:"limit"`
		Offset int           `json:"offset"`
		Total  int           `json:"total"`
		Items  []apiResponse `json:"items"`
	} `json:"_embedded"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: yd-dl <yandex-disk-public-link>")
		os.Exit(1)
	}
	link := os.Args[1]
	if err := downloadAll(link, "."); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func downloadAll(link, dest string) error {
	root, err := fetchResource(link, "/")
	if err != nil {
		return err
	}
	return walkResources(link, root, dest)
}

// recursive walk + flat single-goroutine download, add concurrency if many small files become slow
func walkResources(link string, res apiResponse, dest string) error {
	if res.Embedded == nil || len(res.Embedded.Items) == 0 {
		return nil
	}
	for _, item := range res.Embedded.Items {
		if item.Type == "dir" {
			sub, err := fetchResource(link, item.Path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "skip dir %s: %v\n", item.Path, err)
				continue
			}
			dirDest := filepath.Join(dest, safeName(item.Name))
			if err := os.MkdirAll(dirDest, 0755); err != nil {
				return err
			}
			if err := walkResources(link, sub, dirDest); err != nil {
				return err
			}
		} else if item.Type == "file" {
			fmt.Printf("%s  ", item.Name)
			if err := downloadFile(item, dest); err != nil {
				fmt.Fprintf(os.Stderr, "skip %s: %v\n", item.Name, err)
				continue
			}
			fmt.Println("ok")
		}
	}
	return nil
}

func fetchResource(link string, resourcePath string) (apiResponse, error) {
	u, err := url.Parse(apiBase)
	if err != nil {
		return apiResponse{}, err
	}
	q := u.Query()
	q.Set("public_key", link)
	if resourcePath != "/" {
		q.Set("path", resourcePath)
	}
	// no pagination — Yandex Disk defaults to limit=20, but for public shares
	// with many files we'd need offset loops. YAGNI for now.
	q.Set("limit", "10000")
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return apiResponse{}, fmt.Errorf("fetch %s: %w", resourcePath, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return apiResponse{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return apiResponse{}, fmt.Errorf("API %d: %s", resp.StatusCode, string(body))
	}
	var r apiResponse
	if err := json.Unmarshal(body, &r); err != nil {
		return apiResponse{}, err
	}
	return r, nil
}

func downloadFile(item apiResponse, dest string) error {
	// file field on the item itself, skip download API round-trip
	if item.File == "" {
		return fmt.Errorf("no download URL for %s", item.Name)
	}
	resp, err := http.Get(item.File)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download %s: HTTP %d", item.Name, resp.StatusCode)
	}
	outPath := filepath.Join(dest, safeName(item.Name))
	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, resp.Body)
	return err
}

// strip chars forbidden on any OS
func safeName(name string) string {
	for _, c := range []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|", "\x00"} {
		name = strings.ReplaceAll(name, c, "_")
	}
	return name
}
