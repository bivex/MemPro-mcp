package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

var (
	defaultJSONPath = `C:\Program Files\PureDevSoftware\MemPro\MemProReader\test_memory_analysis.json`
)

func main() {
	// Create MCP server
	s := server.NewMCPServer(
		"MemPro Memory Analyzer",
		"1.0.0",
		server.WithResourceCapabilities(true, false),
	)

	// Add tools for memory analysis
	setupTools(s)

	// Add resources for quick data access
	setupResources(s)

	// Start server using stdio transport
	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func setupTools(s *server.MCPServer) {
	// Tool 1: Analyze Memory Leaks
	analyzeLeaksTool := mcp.NewTool("analyze_leaks",
		mcp.WithDescription("Analyzes memory leaks from MemPro JSON data and returns prioritized list of issues"),
		mcp.WithString("json_path",
			mcp.Description("Path to MemPro JSON analysis file"),
		),
	)

	s.AddTool(analyzeLeaksTool, handleAnalyzeLeaks)

	// Tool 2: Get Memory Summary
	summarizeTool := mcp.NewTool("get_summary",
		mcp.WithDescription("Provides overall memory usage summary including leak percentage and fragmentation"),
		mcp.WithString("json_path",
			mcp.Description("Path to MemPro JSON analysis file"),
		),
	)

	s.AddTool(summarizeTool, handleGetSummary)

	// Tool 3: Get Top Leakers
	topLeakersTool := mcp.NewTool("get_top_leakers",
		mcp.WithDescription("Returns the top N functions causing the most memory leaks"),
		mcp.WithString("json_path",
			mcp.Description("Path to MemPro JSON analysis file"),
		),
		mcp.WithNumber("count",
			mcp.Description("Number of top leakers to return (default: 10)"),
		),
	)

	s.AddTool(topLeakersTool, handleGetTopLeakers)

	// Tool 4: Analyze Fragmentation
	fragmentationTool := mcp.NewTool("analyze_fragmentation",
		mcp.WithDescription("Analyzes memory fragmentation and provides recommendations"),
		mcp.WithString("json_path",
			mcp.Description("Path to MemPro JSON analysis file"),
		),
	)

	s.AddTool(fragmentationTool, handleAnalyzeFragmentation)

	// Tool 5: Find Large Allocations
	largeAllocsTool := mcp.NewTool("find_large_allocations",
		mcp.WithDescription("Identifies unusually large memory allocations that may need optimization"),
		mcp.WithString("json_path",
			mcp.Description("Path to MemPro JSON analysis file"),
		),
	)

	s.AddTool(largeAllocsTool, handleFindLargeAllocations)

	// Tool 6: Get All Issues
	allIssues := mcp.NewTool("get_all_issues",
		mcp.WithDescription("Returns comprehensive analysis of all memory issues including leaks, fragmentation, and large allocations"),
		mcp.WithString("json_path",
			mcp.Description("Path to MemPro JSON analysis file"),
		),
	)

	s.AddTool(allIssues, handleGetAllIssues)
}

func setupResources(s *server.MCPServer) {
	// Resource: Quick stats
	statsResource := mcp.NewResource(
		"mempro://stats",
		"Memory Statistics",
		mcp.WithResourceDescription("Quick memory statistics from the most recent analysis"),
		mcp.WithMIMEType("application/json"),
	)

	s.AddResource(statsResource, func(request mcp.ReadResourceRequest) ([]interface{}, error) {
		analyzer, err := NewMemoryAnalyzer(defaultJSONPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load analyzer: %w", err)
		}

		leakPercentage := 0.0
		if analyzer.data.TotalSize > 0 {
			leakPercentage = float64(analyzer.data.LeakSize) / float64(analyzer.data.TotalSize) * 100
		}

		stats := map[string]interface{}{
			"session":            analyzer.data.SessionName,
			"total_allocations":  analyzer.data.TotalAllocations,
			"total_size":         analyzer.data.TotalSize,
			"leak_count":         analyzer.data.LeakCount,
			"leak_size":          analyzer.data.LeakSize,
			"fragmentation":      analyzer.data.MemoryFragmentation,
			"leak_percentage":    leakPercentage,
		}

		jsonData, err := json.MarshalIndent(stats, "", "  ")
		if err != nil {
			return nil, err
		}

		textContent := mcp.TextResourceContents{
			ResourceContents: mcp.ResourceContents{
				URI:      "mempro://stats",
				MIMEType: "application/json",
			},
			Text: string(jsonData),
		}

		return []interface{}{textContent}, nil
	})
}

// Tool handlers

func handleAnalyzeLeaks(args map[string]interface{}) (*mcp.CallToolResult, error) {
	jsonPath := getJSONPath(args)

	analyzer, err := NewMemoryAnalyzer(jsonPath)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to analyze: %v", err)), nil
	}

	issues := analyzer.AnalyzeLeaks()
	result, err := json.MarshalIndent(issues, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to format results: %v", err)), nil
	}

	return mcp.NewToolResultText(string(result)), nil
}

func handleGetSummary(args map[string]interface{}) (*mcp.CallToolResult, error) {
	jsonPath := getJSONPath(args)

	analyzer, err := NewMemoryAnalyzer(jsonPath)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to analyze: %v", err)), nil
	}

	summary := analyzer.GetSummary()
	return mcp.NewToolResultText(summary), nil
}

func handleGetTopLeakers(args map[string]interface{}) (*mcp.CallToolResult, error) {
	jsonPath := getJSONPath(args)

	count := 10
	if countArg, ok := args["count"].(float64); ok {
		count = int(countArg)
	}

	analyzer, err := NewMemoryAnalyzer(jsonPath)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to analyze: %v", err)), nil
	}

	topLeakers := analyzer.GetTopLeakers(count)
	return mcp.NewToolResultText(topLeakers), nil
}

func handleAnalyzeFragmentation(args map[string]interface{}) (*mcp.CallToolResult, error) {
	jsonPath := getJSONPath(args)

	analyzer, err := NewMemoryAnalyzer(jsonPath)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to analyze: %v", err)), nil
	}

	issues := analyzer.AnalyzeFragmentation()
	result, err := json.MarshalIndent(issues, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to format results: %v", err)), nil
	}

	return mcp.NewToolResultText(string(result)), nil
}

func handleFindLargeAllocations(args map[string]interface{}) (*mcp.CallToolResult, error) {
	jsonPath := getJSONPath(args)

	analyzer, err := NewMemoryAnalyzer(jsonPath)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to analyze: %v", err)), nil
	}

	issues := analyzer.AnalyzeLargeAllocations()
	result, err := json.MarshalIndent(issues, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to format results: %v", err)), nil
	}

	return mcp.NewToolResultText(string(result)), nil
}

func handleGetAllIssues(args map[string]interface{}) (*mcp.CallToolResult, error) {
	jsonPath := getJSONPath(args)

	analyzer, err := NewMemoryAnalyzer(jsonPath)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to analyze: %v", err)), nil
	}

	allIssues := struct {
		Summary       string          `json:"summary"`
		Leaks         []MemoryIssue   `json:"leaks"`
		Fragmentation []MemoryIssue   `json:"fragmentation"`
		LargeAllocs   []MemoryIssue   `json:"large_allocations"`
	}{
		Summary:       analyzer.GetSummary(),
		Leaks:         analyzer.AnalyzeLeaks(),
		Fragmentation: analyzer.AnalyzeFragmentation(),
		LargeAllocs:   analyzer.AnalyzeLargeAllocations(),
	}

	result, err := json.MarshalIndent(allIssues, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to format results: %v", err)), nil
	}

	return mcp.NewToolResultText(string(result)), nil
}

// Helper function to get JSON path from arguments or use default
func getJSONPath(args map[string]interface{}) string {
	if path, ok := args["json_path"].(string); ok && path != "" {
		return path
	}

	// Check if environment variable is set
	if envPath := os.Getenv("MEMPRO_JSON_PATH"); envPath != "" {
		return envPath
	}

	return defaultJSONPath
}
