package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
)

type ipFile struct {
	IPs []string `json:"ip"`
}

type IpRepoImpl struct {
	mu       sync.Mutex
	store    map[string]bool
	fileName string
}

func NewIpRepoImpl(file string) (*IpRepoImpl, error) {
	store, err := jsonFileToMap(file)
	if err != nil {
		return nil, err
	}
	return &IpRepoImpl{
		store:    store,
		fileName: file,
	}, nil
}

func jsonFileToMap(file string) (map[string]bool, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]bool), nil
		}
		return nil, err
	}

	var parsed ipFile
	if err := json.Unmarshal(data, &parsed); err != nil {
		return nil, err
	}

	result := make(map[string]bool)
	for _, ip := range parsed.IPs {
		result[ip] = true
	}
	return result, nil
}

func (repo *IpRepoImpl) GetAll() (map[string]bool, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	result := make(map[string]bool)
	for ip, v := range repo.store {
		result[ip] = v
	}
	return result, nil

}

func (repo *IpRepoImpl) IsAllowed(ctx context.Context, ipStr string) (bool, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false, fmt.Errorf("invalid IP: %s", ipStr)
	}

	for entry := range repo.store {
		if ipMatchesEntry(ip, entry) {
			return true, nil
		}
	}
	return false, nil
}

func ipMatchesEntry(ip net.IP, entry string) bool {
	entry = strings.TrimSpace(entry)
	if strings.Contains(entry, "/") {
		_, network, err := net.ParseCIDR(entry)
		return err == nil && network.Contains(ip)
	}
	if strings.Contains(entry, "-") {
		parts := strings.SplitN(entry, "-", 2)
		if len(parts) != 2 {
			return false
		}
		start := net.ParseIP(strings.TrimSpace(parts[0]))
		end := net.ParseIP(strings.TrimSpace(parts[1]))
		if start == nil || end == nil {
			return false
		}
		return bytes.Compare(ip, start) >= 0 && bytes.Compare(ip, end) <= 0
	}

	return entry == ip.String()
}

func (repo *IpRepoImpl) Add(ctx context.Context, ip string) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	repo.store[ip] = true
	return repo.saveToFile()
}

func (repo *IpRepoImpl) Remove(ctx context.Context, ip string) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	delete(repo.store, ip)
	return repo.saveToFile()
}

func (repo *IpRepoImpl) saveToFile() error {
	ips := make([]string, 0, len(repo.store))
	for ip := range repo.store {
		ips = append(ips, ip)
	}

	out := ipFile{IPs: ips}
	jsonData, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return err
	}

	tmpFile := repo.fileName + ".tmp"
	if err := os.WriteFile(tmpFile, jsonData, 0644); err != nil {
		log.Printf("Ошибка записи временного файла %s: %v", tmpFile, err)
		return err
	}

	if err := os.Rename(tmpFile, repo.fileName); err != nil {
		log.Printf("Ошибка переименования %s -> %s: %v", tmpFile, repo.fileName, err)
		return err
	}

	log.Printf("Успешно сохранено %d IP в файл %s", len(ips), repo.fileName)
	return nil
}
