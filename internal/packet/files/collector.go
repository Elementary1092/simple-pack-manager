package files

import (
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/Elementary1092/pm/internal/packet/parser"
)

const (
	maxWorkers = 5
)

var (
	ErrInternalError = errors.New("syscal failed")
	ErrInvalidPath   = errors.New("invalid path")
	ErrNoTargets     = errors.New("no input targets")
)

type finderInput struct {
	dir     string
	include string
	exclude string
}

// CollectLocalFileNames ignores symbolic and hard links
func CollectLocalFileNames(targets []*parser.Targets) (string, []string, error) {
	if len(targets) == 0 {
		return "", nil, ErrNoTargets
	}
	res := make([]string, 0)
	commonRoot := ""

	paths := make([]string, len(targets))
	for i := range paths {
		paths[i] = targets[i].Path
	}

	err := convertRelativePathsToAbsolutePaths(paths)
	if err != nil {
		return "", nil, err
	}

	// Assuming that to match all files in a directory pattern ends with '/*'
	commonRoot = filepath.Dir(findCommonPrefix(paths))

	var resCh = make(chan string)
	var inputCh = make(chan finderInput)
	var wg sync.WaitGroup
    var workers = len(targets)
    if maxWorkers < workers {
        workers = maxWorkers
    }
	// creating worker goroutines
	for i := 0; i < workers; i++ {
		go findAllFiles(inputCh, resCh, &wg)
		wg.Add(1)
	}

	go func() {
		for i, path := range paths {
			inputCh <- finderInput{
				dir:     filepath.Dir(path),
				include: filepath.Base(path),
				exclude: targets[i].Exclude,
			}
		}
		close(inputCh)
	}()

	go func() {
		wg.Wait()
		close(resCh)
	}()

	for path := range resCh {
		res = append(res, path)
	}

	return commonRoot, res, nil
}

// To implement files collection from subdirectory, I think, that a map should be created
// which maps directory to exclude and include patterns. Then, on encountering
// subdirectory in directory, check whether such directory is mapped to the patterns.
// If yes, it means that this directory is already in a queue.
// If no, add include and exclude (if it has wildcard matching) paths of the current directory to the queue.
func findAllFiles(input <-chan finderInput, res chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	for in := range input {
		dir, include, exclude := in.dir, in.include, in.exclude
		entries, err := os.ReadDir(dir)
		if err != nil {
			return
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			filePath := filepath.Join(dir, entry.Name())
			info, err := os.Lstat(filePath)
			if err != nil {
			    continue	
			}

			if !info.Mode().IsRegular() {
				continue
			}

			included, err := filepath.Match(include, info.Name())
			if err != nil || !included {
				continue
			}

			excluded, err := filepath.Match(exclude, info.Name())
			if err != nil || excluded {
				continue
			}

			res <- filePath
		}
	}
}

func convertRelativePathsToAbsolutePaths(paths []string) error {
	currPath, err := os.Getwd()
	if err != nil {
		return ErrInternalError
	}
	// could have implemented more efficient algorithm, but it would have been time consuming
	for i, path := range paths {
		if !filepath.IsAbs(path) {
			paths[i], err = filepath.Abs(filepath.Join(currPath, path))
			if err != nil {
				return ErrInvalidPath
			}
		}
	}

	return nil
}

func findCommonPrefix(strs []string) string {
	if len(strs) == 0 {
		return ""
	}

	// Creating new array inorder to not change original strs array
	var cp = make([]string, len(strs))
	copy(cp, strs)
	sort.Slice(cp, func(i, j int) bool {
		return len(cp[i]) < len(cp[j])
	})

	var commonStr strings.Builder
	for i, ch := range []byte(cp[0]) {
		matched := true
		for j := 1; j < len(cp); j++ {
			if ch != cp[j][i] {
				matched = false
				break
			}
		}

		if !matched {
			break
		}

		commonStr.WriteByte(ch)
	}

	return commonStr.String()
}
