package engine

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type loadFn func() error

func SyncRules(loadRules loadFn) error {
	err := loadRules()
	if err != nil {
		log.Fatalf("Failed to load rules, error %v", err)
		return fmt.Errorf("error loading rules: %v", err)
	}
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		for range ticker.C {
			err := loadRules()
			if err != nil {
				log.Fatalf("Failed to load rules, error: %v", err)
			}
		}
	}()
	return nil
}

type FileVersion struct {
	FileName string
	Version  int
	Hash     string
	Change   bool
}

var (
	fileVersions = make(map[string]*FileVersion)
	mu           sync.RWMutex
	once sync.Once
)

func loadRuleFiles(dir string, extension string) ([]FileVersion, error) {
	var result []FileVersion
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), extension) {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			hash := sha256.Sum256(content)
			hashStr := hex.EncodeToString(hash[:])

			mu.Lock()
			if fv, ok := fileVersions[info.Name()]; ok {
				if fv.Hash != hashStr {
					fv.Version++
					fv.Hash = hashStr
					fv.Change = true
				} else {
					fv.Change = false
				}
			} else {
				fv = &FileVersion{
					FileName: info.Name(),
					Version:  0,
					Hash:     hashStr,
					Change:   true,
				}
				fileVersions[info.Name()] = fv
			}
			mu.Unlock()
			result = append(result, *fileVersions[info.Name()])
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}
