package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/dustin/go-humanize"
)

// WriteCounter counts the number of bytes written to it. It implements to the io.Writer interface
// and we can pass this into io.TeeReader() which will report progress on each write cycle.
type WriteCounter struct {
	Total uint64
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.PrintProgress()
	return n, nil
}

func (wc WriteCounter) PrintProgress() {
	// Clear the line by using a character return to go back to the start and remove
	// the remaining characters by filling it with spaces
	fmt.Printf("\r%s", strings.Repeat(" ", 35))

	// Return again and print current status of download
	// We use the humanize package to print the bytes in a meaningful way (e.g. 10 MB)
	fmt.Printf("\rDownloading... %s complete", humanize.Bytes(wc.Total))
}

func getRevision(lastChange string) string {
	revision := ""
	resp, err := http.Get(lastChange)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		revision = string(bodyBytes)
	}
	return revision
}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory. We pass an io.TeeReader
// into Copy() to report progress on the download.
func DownloadFile(filepath string, url string) error {

	// Create the file, but give it a tmp file extension, this means we won't overwrite a
	// file until it's downloaded, but we'll remove the tmp extension once downloaded.
	out, err := os.Create(filepath + ".tmp")
	if err != nil {
		return err
	}

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		out.Close()
		return err
	}
	defer resp.Body.Close()

	// Create our progress reporter and pass it to be used alongside our writer
	counter := &WriteCounter{}
	if _, err = io.Copy(out, io.TeeReader(resp.Body, counter)); err != nil {
		out.Close()
		return err
	}

	// The progress use the same line so print a new line once it's finished downloading
	fmt.Print("\n")

	// Close the file without defer so it can happen before Rename()
	out.Close()

	if err = os.Rename(filepath+".tmp", filepath); err != nil {
		return err
	}
	return nil
}

func main() {

	fileURL := ""
	platform := runtime.GOOS
	//Get latest revision snapshot for chromium.
	lastChange := "https://www.googleapis.com/download/storage/v1/b/chromium-browser-snapshots/o/Linux_x64%2FLAST_CHANGE?alt=media"
	revisionNo := getRevision(lastChange)
	fmt.Print("Go runs on ")
	switch pl := platform; pl {
	case "darwin":
		fmt.Println("Platform Mac OS.")
		fileURL = "https://storage.googleapis.com/chromium-browser-snapshots/Mac/" + revisionNo + "/chrome-mac.zip"
		//filepath = ""
	case "linux":
		fmt.Println("Platform Linux.")
		fileURL = "https://storage.googleapis.com/chromium-browser-snapshots/Linux_x64/" + revisionNo + "/chrome-linux.zip"
	case "windows":
		fmt.Println("Platform Windows.")
		fileURL = "https://storage.googleapis.com/chromium-browser-snapshots/Win_x64/" + revisionNo + "/chrome-win.zip"
	default:
		// freebsd, openbsd
		fmt.Printf("%s.\n", pl)
		fmt.Println("Unsupported OS platform.")
		os.Exit(3)
	}

	fmt.Println("Download Started")
	if err := DownloadFile("chromium-"+platform+"-x64_latest_"+revisionNo+".zip", fileURL); err != nil {
		panic(err)
	}
	fmt.Println("Download Finished")
}
