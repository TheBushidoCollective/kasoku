package analytics

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Analytics struct {
	path string
	mu   sync.Mutex
	data *AnalyticsData
}

type AnalyticsData struct {
	TotalCacheHits   int64                    `json:"total_cache_hits"`
	TotalCacheMisses int64                    `json:"total_cache_misses"`
	TotalTimeSaved   time.Duration            `json:"total_time_saved"`
	TotalCommands    int64                    `json:"total_commands"`
	CommandStats     map[string]*CommandStats `json:"command_stats"`
	DailyStats       map[string]*DailyStats   `json:"daily_stats"`
}

type CommandStats struct {
	Command       string        `json:"command"`
	CacheHits     int64         `json:"cache_hits"`
	CacheMisses   int64         `json:"cache_misses"`
	TimeSaved     time.Duration `json:"time_saved"`
	TotalDuration time.Duration `json:"total_duration"`
	LastRun       time.Time     `json:"last_run"`
}

type DailyStats struct {
	Date        string        `json:"date"`
	CacheHits   int64         `json:"cache_hits"`
	CacheMisses int64         `json:"cache_misses"`
	TimeSaved   time.Duration `json:"time_saved"`
	CommandsRun int64         `json:"commands_run"`
}

type ExecutionEvent struct {
	Command   string
	CacheHit  bool
	Duration  time.Duration
	TimeSaved time.Duration
	Timestamp time.Time
}

func New(basePath string) (*Analytics, error) {
	if basePath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		basePath = filepath.Join(homeDir, ".brisk", "analytics")
	}

	basePath = os.ExpandEnv(basePath)

	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create analytics directory: %w", err)
	}

	analyticsPath := filepath.Join(basePath, "analytics.json")

	a := &Analytics{
		path: analyticsPath,
		data: &AnalyticsData{
			CommandStats: make(map[string]*CommandStats),
			DailyStats:   make(map[string]*DailyStats),
		},
	}

	if err := a.load(); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to load analytics: %w", err)
	}

	return a, nil
}

func (a *Analytics) RecordExecution(event ExecutionEvent) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.data.TotalCommands++

	if event.CacheHit {
		a.data.TotalCacheHits++
		a.data.TotalTimeSaved += event.TimeSaved
	} else {
		a.data.TotalCacheMisses++
	}

	cmdStats, ok := a.data.CommandStats[event.Command]
	if !ok {
		cmdStats = &CommandStats{
			Command: event.Command,
		}
		a.data.CommandStats[event.Command] = cmdStats
	}

	if event.CacheHit {
		cmdStats.CacheHits++
		cmdStats.TimeSaved += event.TimeSaved
	} else {
		cmdStats.CacheMisses++
	}
	cmdStats.TotalDuration += event.Duration
	cmdStats.LastRun = event.Timestamp

	dateStr := event.Timestamp.Format("2006-01-02")
	dailyStats, ok := a.data.DailyStats[dateStr]
	if !ok {
		dailyStats = &DailyStats{
			Date: dateStr,
		}
		a.data.DailyStats[dateStr] = dailyStats
	}

	if event.CacheHit {
		dailyStats.CacheHits++
		dailyStats.TimeSaved += event.TimeSaved
	} else {
		dailyStats.CacheMisses++
	}
	dailyStats.CommandsRun++

	return a.save()
}

func (a *Analytics) GetStats() *AnalyticsData {
	a.mu.Lock()
	defer a.mu.Unlock()

	dataCopy := *a.data
	dataCopy.CommandStats = make(map[string]*CommandStats, len(a.data.CommandStats))
	for k, v := range a.data.CommandStats {
		statsCopy := *v
		dataCopy.CommandStats[k] = &statsCopy
	}
	dataCopy.DailyStats = make(map[string]*DailyStats, len(a.data.DailyStats))
	for k, v := range a.data.DailyStats {
		statsCopy := *v
		dataCopy.DailyStats[k] = &statsCopy
	}

	return &dataCopy
}

func (a *Analytics) GetCommandStats() []*CommandStats {
	a.mu.Lock()
	defer a.mu.Unlock()

	stats := make([]*CommandStats, 0, len(a.data.CommandStats))
	for _, cmdStats := range a.data.CommandStats {
		statsCopy := *cmdStats
		stats = append(stats, &statsCopy)
	}

	return stats
}

func (a *Analytics) GetDailyStats(days int) []*DailyStats {
	a.mu.Lock()
	defer a.mu.Unlock()

	stats := make([]*DailyStats, 0)
	now := time.Now()

	for i := 0; i < days; i++ {
		date := now.AddDate(0, 0, -i)
		dateStr := date.Format("2006-01-02")

		if dailyStats, ok := a.data.DailyStats[dateStr]; ok {
			statsCopy := *dailyStats
			stats = append(stats, &statsCopy)
		} else {
			stats = append(stats, &DailyStats{
				Date: dateStr,
			})
		}
	}

	return stats
}

func (a *Analytics) load() error {
	data, err := os.ReadFile(a.path)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, a.data)
}

func (a *Analytics) save() error {
	data, err := json.MarshalIndent(a.data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal analytics: %w", err)
	}

	tmpPath := a.path + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write analytics: %w", err)
	}

	if err := os.Rename(tmpPath, a.path); err != nil {
		return fmt.Errorf("failed to rename analytics file: %w", err)
	}

	return nil
}

func (a *Analytics) PrintSummary() {
	stats := a.GetStats()

	fmt.Println("\n📊 Brisk Analytics Summary")
	fmt.Println("═══════════════════════════")
	fmt.Printf("Total commands run:  %d\n", stats.TotalCommands)
	fmt.Printf("Cache hits:          %d (%.1f%%)\n",
		stats.TotalCacheHits,
		float64(stats.TotalCacheHits)/float64(stats.TotalCommands)*100)
	fmt.Printf("Cache misses:        %d (%.1f%%)\n",
		stats.TotalCacheMisses,
		float64(stats.TotalCacheMisses)/float64(stats.TotalCommands)*100)
	fmt.Printf("Total time saved:    %v\n", stats.TotalTimeSaved)

	if len(stats.CommandStats) > 0 {
		fmt.Println("\n🏆 Top Commands by Time Saved:")
		for _, cmdStats := range stats.CommandStats {
			if cmdStats.TimeSaved > 0 {
				avgDuration := cmdStats.TotalDuration / time.Duration(cmdStats.CacheHits+cmdStats.CacheMisses)
				fmt.Printf("  %s: %v saved (avg: %v, hits: %d, misses: %d)\n",
					cmdStats.Command,
					cmdStats.TimeSaved,
					avgDuration,
					cmdStats.CacheHits,
					cmdStats.CacheMisses)
			}
		}
	}
}
