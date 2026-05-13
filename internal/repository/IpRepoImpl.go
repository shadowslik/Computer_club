package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	ipfilter "github.com/shadowslik/go-ipfilter"
	"go.uber.org/zap"
)

type ipFile struct {
	IPs []string `json:"ip"`
}

type IpRepoImpl struct {
	mu       sync.RWMutex
	rawList  []string
	matcher  *ipfilter.Matcher
	fileName string
	watcher  *fsnotify.Watcher
	log      *zap.Logger
}

func NewIpRepoImpl(file string, cacheSize int, cacheTTL time.Duration, log *zap.Logger) (*IpRepoImpl, error) {
	repo := &IpRepoImpl{
		fileName: file,
		matcher:  ipfilter.New(cacheSize, cacheTTL),
		log:      log,
	}
	if err := repo.loadFromFile(); err != nil {
		return nil, err
	}
	if err := repo.startWatcher(); err != nil {
		log.Warn("cannot start file watcher", zap.String("file", file), zap.Error(err))
	}
	return repo, nil
}

func (repo *IpRepoImpl) loadFromFile() error {
	entries, err := readIPFile(repo.fileName)
	if err != nil {
		return err
	}
	if err := repo.matcher.Reset(entries); err != nil {
		return fmt.Errorf("ipfilter reset: %w", err)
	}
	repo.mu.Lock()
	repo.rawList = entries
	repo.mu.Unlock()
	return nil
}

func (repo *IpRepoImpl) startWatcher() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	repo.watcher = watcher
	if err := watcher.Add(repo.fileName); err != nil {
		return err
	}
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&(fsnotify.Write|fsnotify.Create) != 0 {
					repo.log.Info("IP list changed, reloading", zap.String("file", repo.fileName))
					if err := repo.loadFromFile(); err != nil {
						repo.log.Error("failed to reload IP list", zap.String("file", repo.fileName), zap.Error(err))
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				repo.log.Error("file watcher error", zap.String("file", repo.fileName), zap.Error(err))
			}
		}
	}()
	return nil
}

func (repo *IpRepoImpl) GetAll() (map[string]bool, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()
	result := make(map[string]bool, len(repo.rawList))
	for _, s := range repo.rawList {
		result[s] = true
	}
	return result, nil
}

func (repo *IpRepoImpl) IsAllowed(_ context.Context, ipStr string) (bool, error) {
	return repo.matcher.Match(ipStr)
}

func (repo *IpRepoImpl) Add(_ context.Context, ip string) error {
	if err := repo.matcher.Add(ip); err != nil {
		return err
	}
	repo.mu.Lock()
	repo.rawList = append(repo.rawList, ip)
	err := repo.saveToFile()
	repo.mu.Unlock()
	return err
}

func (repo *IpRepoImpl) Remove(_ context.Context, ip string) error {
	repo.matcher.Remove(ip)

	repo.mu.Lock()
	filtered := repo.rawList[:0]
	for _, s := range repo.rawList {
		if s != ip {
			filtered = append(filtered, s)
		}
	}
	repo.rawList = filtered
	err := repo.saveToFile()
	repo.mu.Unlock()
	return err
}

func (repo *IpRepoImpl) saveToFile() error {
	out := ipFile{IPs: repo.rawList}
	data, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return err
	}
	tmp := repo.fileName + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return err
	}
	return os.Rename(tmp, repo.fileName)
}

func readIPFile(file string) ([]string, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	if len(data) == 0 {
		return nil, nil
	}
	var parsed ipFile
	if err := json.Unmarshal(data, &parsed); err != nil {
		return nil, fmt.Errorf("invalid JSON in %s: %w", file, err)
	}
	return parsed.IPs, nil
}
