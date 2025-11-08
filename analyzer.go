package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
)

// MemoryAnalyzer analyzes MemPro data and detects memory issues
type MemoryAnalyzer struct {
	data *MemProData
}

// NewMemoryAnalyzer creates a new analyzer from a JSON file
func NewMemoryAnalyzer(jsonPath string) (*MemoryAnalyzer, error) {
	fileData, err := os.ReadFile(jsonPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read JSON file: %w", err)
	}

	var data MemProData
	if err := json.Unmarshal(fileData, &data); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return &MemoryAnalyzer{data: &data}, nil
}

// AnalyzeLeaks detects and prioritizes memory leaks
func (ma *MemoryAnalyzer) AnalyzeLeaks() []MemoryIssue {
	var issues []MemoryIssue

	if ma == nil || ma.data == nil {
		return issues
	}

	for _, leak := range ma.data.Leaks {
		if leak.LeakSize == 0 && leak.LeakCount == 0 {
			continue
		}

		severity := ma.calculateLeakSeverity(leak)
		suggestion := ma.generateLeakSuggestion(leak)

		issues = append(issues, MemoryIssue{
			Severity:     severity,
			Type:         "MemoryLeak",
			Description:  ma.formatLeakDescription(leak),
			FunctionName: leak.FunctionName,
			FileName:     leak.FileName,
			LineNumber:   leak.LineNumber,
			Size:         leak.LeakSize,
			Count:        leak.LeakCount,
			Score:        leak.LeakScore,
			Suggestion:   suggestion,
		})
	}

	// Sort by severity and size
	severityOrder := map[string]int{"Critical": 0, "High": 1, "Medium": 2, "Low": 3}
	sort.Slice(issues, func(i, j int) bool {
		orderI, okI := severityOrder[issues[i].Severity]
		orderJ, okJ := severityOrder[issues[j].Severity]

		// Handle unknown severity levels by treating them as lowest priority
		if !okI {
			orderI = 999
		}
		if !okJ {
			orderJ = 999
		}

		if orderI != orderJ {
			return orderI < orderJ
		}
		return issues[i].Size > issues[j].Size
	})

	return issues
}

// AnalyzeFragmentation detects memory fragmentation issues
func (ma *MemoryAnalyzer) AnalyzeFragmentation() []MemoryIssue {
	var issues []MemoryIssue

	if ma == nil || ma.data == nil {
		return issues
	}

	if ma.data.MemoryFragmentation > 80.0 {
		issues = append(issues, MemoryIssue{
			Severity:    "High",
			Type:        "MemoryFragmentation",
			Description: fmt.Sprintf("Memory fragmentation is at %.2f%%, which indicates severe fragmentation", ma.data.MemoryFragmentation),
			Size:        ma.data.TotalSize,
			Count:       ma.data.TotalAllocations,
			Score:       ma.data.MemoryFragmentation,
			Suggestion:  "Consider implementing object pooling or using memory arenas to reduce fragmentation. Review allocation patterns and consolidate small allocations where possible.",
		})
	} else if ma.data.MemoryFragmentation > 50.0 {
		issues = append(issues, MemoryIssue{
			Severity:    "Medium",
			Type:        "MemoryFragmentation",
			Description: fmt.Sprintf("Memory fragmentation is at %.2f%%, which may impact performance", ma.data.MemoryFragmentation),
			Size:        ma.data.TotalSize,
			Count:       ma.data.TotalAllocations,
			Score:       ma.data.MemoryFragmentation,
			Suggestion:  "Monitor fragmentation levels and consider optimizing allocation patterns if fragmentation increases.",
		})
	}

	return issues
}

// AnalyzeLargeAllocations finds unusually large allocations
func (ma *MemoryAnalyzer) AnalyzeLargeAllocations() []MemoryIssue {
	var issues []MemoryIssue

	if ma == nil || ma.data == nil {
		return issues
	}

	for _, fn := range ma.data.Functions {
		if fn.AverageSize > 10000 || fn.MaxSize > 50000 {
			severity := "Medium"
			if fn.MaxSize > 100000 {
				severity = "High"
			}

			issues = append(issues, MemoryIssue{
				Severity:     severity,
				Type:         "LargeAllocation",
				Description:  fmt.Sprintf("Large allocation detected: average %.0f bytes, max %d bytes across %d allocations", fn.AverageSize, fn.MaxSize, fn.AllocationCount),
				FunctionName: fn.FunctionName,
				FileName:     fn.FileName,
				LineNumber:   fn.LineNumber,
				Size:         fn.TotalSize,
				Count:        fn.AllocationCount,
				Score:        float64(fn.MaxSize),
				Suggestion:   "Review if large allocations can be split into smaller chunks or allocated incrementally. Consider using streaming or chunked processing for large data.",
			})
		}
	}

	return issues
}

// GetSummary provides an overall summary of memory usage
func (ma *MemoryAnalyzer) GetSummary() string {
	if ma == nil || ma.data == nil {
		return "Error: No data available for analysis"
	}

	leakPercentage := 0.0
	if ma.data.TotalSize > 0 {
		leakPercentage = float64(ma.data.LeakSize) / float64(ma.data.TotalSize) * 100
	}

	summary := fmt.Sprintf(`Memory Analysis Summary
======================
Session: %s
Total Allocations: %d
Total Size: %d bytes (%.2f MB)
Leak Count: %d
Leak Size: %d bytes (%.2f MB)
Leak Percentage: %.2f%%
Memory Fragmentation: %.2f%%

Critical Findings:
`, ma.data.SessionName, ma.data.TotalAllocations, ma.data.TotalSize,
		float64(ma.data.TotalSize)/1024/1024,
		ma.data.LeakCount, ma.data.LeakSize,
		float64(ma.data.LeakSize)/1024/1024,
		leakPercentage, ma.data.MemoryFragmentation)

	if leakPercentage > 50 {
		summary += "- CRITICAL: Over 50% of allocated memory is leaked!\n"
	}
	if ma.data.MemoryFragmentation > 80 {
		summary += "- HIGH: Severe memory fragmentation detected\n"
	}

	suspectLeaks := 0
	for _, leak := range ma.data.Leaks {
		if leak.IsSuspect {
			suspectLeaks++
		}
	}
	if suspectLeaks > 0 {
		summary += fmt.Sprintf("- %d suspect leak locations identified\n", suspectLeaks)
	}

	return summary
}

// GetTopLeakers returns the top N functions by leak size
func (ma *MemoryAnalyzer) GetTopLeakers(n int) string {
	if ma == nil || ma.data == nil {
		return "Error: No data available for analysis"
	}

	leaks := make([]Leak, len(ma.data.Leaks))
	copy(leaks, ma.data.Leaks)

	sort.Slice(leaks, func(i, j int) bool {
		return leaks[i].LeakSize > leaks[j].LeakSize
	})

	if n > len(leaks) {
		n = len(leaks)
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("Top %d Memory Leakers:\n", n))
	result.WriteString("====================\n\n")

	for i := 0; i < n; i++ {
		leak := leaks[i]
		result.WriteString(fmt.Sprintf("%d. %s\n", i+1, leak.FunctionName))
		result.WriteString(fmt.Sprintf("   Leak Size: %d bytes (%.2f KB)\n", leak.LeakSize, float64(leak.LeakSize)/1024))
		result.WriteString(fmt.Sprintf("   Leak Count: %d allocations\n", leak.LeakCount))
		result.WriteString(fmt.Sprintf("   Leak Score: %.2f\n", leak.LeakScore))
		result.WriteString(fmt.Sprintf("   Suspect: %v\n", leak.IsSuspect))
		if leak.FileName != "" {
			result.WriteString(fmt.Sprintf("   Location: %s:%d\n", leak.FileName, leak.LineNumber))
		}
		if leak.CallStack != "" {
			result.WriteString(fmt.Sprintf("   CallStack: %s\n", leak.CallStack))
		}
		result.WriteString("\n")
	}

	return result.String()
}

// Helper functions

func (ma *MemoryAnalyzer) calculateLeakSeverity(leak Leak) string {
	if leak.IsSuspect && leak.LeakSize > 100000 {
		return "Critical"
	}
	if leak.IsSuspect || leak.LeakSize > 50000 {
		return "High"
	}
	if leak.LeakSize > 10000 || leak.LeakCount > 100 {
		return "Medium"
	}
	return "Low"
}

func (ma *MemoryAnalyzer) formatLeakDescription(leak Leak) string {
	if strings.Contains(leak.FunctionName, "Unknown Function") {
		return fmt.Sprintf("Unknown function leaked %d bytes across %d allocations. This may indicate missing debug symbols or dynamically loaded code.",
			leak.LeakSize, leak.LeakCount)
	}
	return fmt.Sprintf("Function leaked %d bytes across %d allocations", leak.LeakSize, leak.LeakCount)
}

func (ma *MemoryAnalyzer) generateLeakSuggestion(leak Leak) string {
	if strings.Contains(leak.FunctionName, "Unknown Function") {
		return "Enable debug symbols and rebuild with full symbol information to identify the exact source of this leak. Check for third-party libraries or dynamically loaded modules."
	}

	if strings.Contains(leak.FunctionName, "std::_Allocate") || strings.Contains(leak.FunctionName, "std::vector") {
		return "STL container leak detected. Ensure proper cleanup in destructors, check for circular references, and verify that containers are properly cleared before going out of scope."
	}

	if strings.Contains(leak.FunctionName, "main") {
		return "Leak originated from main function. Review allocation ownership and ensure all allocated resources are freed before program exit. Consider using RAII or smart pointers."
	}

	return "Review allocation patterns in this function and ensure all allocated memory is properly deallocated. Consider using smart pointers (std::unique_ptr, std::shared_ptr) or RAII patterns."
}
