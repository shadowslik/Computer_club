package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	lruexpirable "github.com/hashicorp/golang-lru/v2/expirable"
	"github.com/yl2chen/cidranger"
	"go.uber.org/zap"
)

type ipFile struct {
	IPs []string `json:"ip"`
}

type ipRange struct {
	start net.IP
	end   net.IP
	raw   string
}

type IpRepoImpl struct {
	mu       sync.RWMutex
	rawList  []string
	exactIPs map[string]struct{}
	ranger   cidranger.Ranger
	ranges   []ipRange

	cache    *lruexpirable.LRU[string, bool]
	fileName string
	watcher  *fsnotify.Watcher
	log      *zap.Logger
}

func NewIpRepoImpl(file string, cacheSize int, cacheTTL time.Duration, log *zap.Logger) (*IpRepoImpl, error) {
	repo := &IpRepoImpl{
		fileName: file,
		cache:    lruexpirable.NewLRU[string, bool](cacheSize, nil, cacheTTL),
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
	rawList, err := readIPFile(repo.fileName)
	if err != nil {
		return err
	}
	exactIPs := make(map[string]struct{})
	ranger := cidranger.NewPCTrieRanger()
	var ranges []ipRange

	for _, entry := range rawList {
		entry = strings.TrimSpace(entry)
		if strings.Contains(entry, "/") {
			_, network, err := net.ParseCIDR(entry)
			if err == nil {
				ranger.Insert(cidranger.NewBasicRangerEntry(*network))
			}
		} else if strings.Contains(entry, "-") {
			parts := strings.SplitN(entry, "-", 2)
			if len(parts) == 2 {
				start := net.ParseIP(strings.TrimSpace(parts[0]))
				end := net.ParseIP(strings.TrimSpace(parts[1]))
				if start != nil && end != nil {
					ranges = append(ranges, ipRange{start: start, end: end, raw: entry})
				}
			}
		} else {
			if net.ParseIP(entry) != nil {
				exactIPs[entry] = struct{}{}
			}
		}
	}

	repo.mu.Lock()
	defer repo.mu.Unlock()
	repo.rawList = rawList
	repo.exactIPs = exactIPs
	repo.ranger = ranger
	repo.ranges = ranges
	repo.cache.Purge()
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

func (repo *IpRepoImpl) IsAllowed(ctx context.Context, ipStr string) (bool, error) {
	if cached, ok := repo.cache.Get(ipStr); ok {
		return cached, nil
	}

	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false, fmt.Errorf("invalid IP: %s", ipStr)
	}

	repo.mu.RLock()
	defer repo.mu.RUnlock()

	// 1. Exact match — O(1)
	normalised := ip.String()
	if _, ok := repo.exactIPs[normalised]; ok {
		repo.cache.Add(ipStr, true)
		return true, nil
	}

	// 2. CIDR via radix trie — O(log n)
	entries, err := repo.ranger.ContainingNetworks(ip)
	if err == nil && len(entries) > 0 {
		repo.cache.Add(ipStr, true)
		return true, nil
	}

	// 3. IP ranges — O(n ranges), typically very few
	for _, r := range repo.ranges {
		if bytes.Compare(ip, r.start) >= 0 && bytes.Compare(ip, r.end) <= 0 {
			repo.cache.Add(ipStr, true)
			return true, nil
		}
	}

	repo.cache.Add(ipStr, false)
	return false, nil
}

func (repo *IpRepoImpl) Add(ctx context.Context, ip string) error {
	repo.mu.Lock()
	repo.rawList = append(repo.rawList, ip)
	// rebuild index for the new entry
	if strings.Contains(ip, "/") {
		_, network, err := net.ParseCIDR(ip)
		if err == nil {
			repo.ranger.Insert(cidranger.NewBasicRangerEntry(*network))
		}
	} else if strings.Contains(ip, "-") {
		parts := strings.SplitN(ip, "-", 2)
		if len(parts) == 2 {
			start := net.ParseIP(strings.TrimSpace(parts[0]))
			end := net.ParseIP(strings.TrimSpace(parts[1]))
			if start != nil && end != nil {
				repo.ranges = append(repo.ranges, ipRange{start: start, end: end, raw: ip})
			}
		}
	} else {
		if parsed := net.ParseIP(ip); parsed != nil {
			repo.exactIPs[parsed.String()] = struct{}{}
		}
	}
	repo.cache.Purge()
	err := repo.saveToFile()
	repo.mu.Unlock()
	return err
}

func (repo *IpRepoImpl) Remove(ctx context.Context, ip string) error {
	repo.mu.Lock()
	filtered := repo.rawList[:0]
	for _, s := range repo.rawList {
		if s != ip {
			filtered = append(filtered, s)
		}
	}
	repo.rawList = filtered
	repo.cache.Purge()
	err := repo.saveToFile()
	repo.mu.Unlock()
	if err != nil {
		return err
	}
	// full rebuild needed for ranger (no remove in cidranger v1)
	return repo.loadFromFile()
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
