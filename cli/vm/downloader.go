package vm

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
)

type WriteCounter struct {
	Downloaded int64
	TotalSize  int64
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Downloaded += int64(n)
	wc.PrintProgress()
	return n, nil
}

func (wc WriteCounter) PrintProgress() {
	percentage := 100 * wc.Downloaded / wc.TotalSize
	fmt.Printf("\rDownloading... %d%% completed", percentage)
}

func DownloadFile(filepath string, url string) error {

	// Create the file, but give it a tmp file extension, this means we won't overwrite a
	// file until it's downloaded, but we'll remove the tmp extension once downloaded.
	out, err := os.Create(filepath + ".tmp")
	if err != nil {
		return err
	}
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		fmt.Printf("\nInterrupt signal received while downloading image, removing temp files\n")
		os.Remove(filepath + ".tmp")
		os.Exit(1)
	}()
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Disk image error: 404 not found")
	}
	defer resp.Body.Close()

	fmt.Printf("File size: %d bytes\n", resp.ContentLength)
	// Create our progress reporter and pass it to be used alongside our writer
	counter := &WriteCounter{TotalSize: resp.ContentLength}
	_, err = io.Copy(out, io.TeeReader(resp.Body, counter))
	if err != nil {
		return err
	}
	fmt.Println()

	err = os.Rename(filepath+".tmp", filepath)
	if err != nil {
		return err
	}

	return nil
}
