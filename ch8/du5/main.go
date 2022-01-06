package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var vFlag = flag.Bool("v", false, "show verbose progress messages")

type directoryInfo struct {
	Path string
	TotalSize int64
}

//!+
func main() {
	flag.Parse()
	roots := flag.Args()
	if len(roots) == 0 {
		roots = []string{"."}
	}
	fileSizes := make(chan int64)
	dirSizes := make(chan directoryInfo)
	subDirSizeChannel := make(chan int64)
	var n sync.WaitGroup
	for _, root := range roots {
		n.Add(1)
		go walkDir(root, &n, fileSizes, dirSizes, subDirSizeChannel)
	}
	go func() {
		<-subDirSizeChannel
		n.Wait()
		close(fileSizes)
		close(dirSizes)
	}()
	var tick <-chan time.Time
	if *vFlag {
		tick = time.Tick(500 * time.Millisecond)
	}
	var nfiles, nbytes int64
loop:
	for {
		select {
		case size, ok := <-fileSizes:
			if !ok {
				break loop
			}
			nfiles++
			nbytes += size
		case <-tick:
			printDiskUsage(nfiles, nbytes)
		case directorySize := <-dirSizes:
			printDirectorySize(directorySize)
		}
	}

	printDiskUsage(nfiles, nbytes)
}

func printDiskUsage(nfiles, nbytes int64) {
	fmt.Printf("%d files  %.1f GB\n", nfiles, float64(nbytes)/1e9)
}

func printDirectorySize(dirInfo directoryInfo) {
	fmt.Printf("%s: %.1f MB\n", dirInfo.Path, float64(dirInfo.TotalSize)/1e6)
}

// walkDir recursively walks the file tree rooted at dir
// and sends the size of each found file on fileSizes.
//!+walkDir
func walkDir(
	dir string, 
	n *sync.WaitGroup, 
	fileSizes chan<- int64, 
	dirSizes chan<- directoryInfo,
	dirSizeReport chan<- int64,
) {
	defer n.Done()
	defer close(dirSizeReport)
	var totalSize int64
	var subDirSizeChannels []chan int64
	dirEntities := dirents(dir)
	if dirEntities == nil {
		dirSizeReport <- 0
		return
	}
	for _, entry := range dirEntities {
		if entry.IsDir() {
			n.Add(1)
			subDirSizeChannel := make(chan int64)
			subDirSizeChannels = append(subDirSizeChannels, subDirSizeChannel)
			subdir := filepath.Join(dir, entry.Name())
			go walkDir(subdir, n, fileSizes, dirSizes, subDirSizeChannel)
		} else {
			fileSize := entry.Size()
			fileSizes <- fileSize
			totalSize += fileSize
		}
	}
	for _, channel := range subDirSizeChannels {
		subDirSize := <- channel
		totalSize += subDirSize
	}
	dirSizeReport <- totalSize
	dirSizes <- directoryInfo{dir, totalSize}
}

var sema = make(chan struct{}, 20)

func dirents(dir string) []os.FileInfo {
	sema <- struct{}{}
	defer func() { <-sema }()
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "du: %v\n", err)
		return nil
	}
	return entries
}
