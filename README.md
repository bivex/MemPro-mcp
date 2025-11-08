# MemPro MCP Server

An MCP (Model Context Protocol) server for analyzing MemPro memory profiling data and detecting memory errors using AI-powered analysis.

## Overview

This server provides intelligent analysis of MemPro JSON exports, enabling AI assistants to quickly identify memory leaks, fragmentation issues, and large allocations in C++ applications.

## Features

### MCP Tools

1. **analyze_leaks** - Analyzes memory leaks and returns prioritized issues
   - Input: `json_path` (optional, defaults to test_memory_analysis.json)
   - Output: JSON array of memory leak issues with severity, descriptions, and suggestions

2. **get_summary** - Provides comprehensive memory usage summary
   - Input: `json_path` (optional)
   - Output: Text summary with key metrics and critical findings

3. **get_top_leakers** - Returns top N functions causing memory leaks
   - Input: `json_path` (optional), `count` (default: 10)
   - Output: Formatted list of top leakers with details

4. **analyze_fragmentation** - Analyzes memory fragmentation
   - Input: `json_path` (optional)
   - Output: Fragmentation issues and recommendations

5. **find_large_allocations** - Identifies unusually large allocations
   - Input: `json_path` (optional)
   - Output: List of large allocation issues

6. **get_all_issues** - Comprehensive analysis of all issues
   - Input: `json_path` (optional)
   - Output: Complete analysis including summary, leaks, fragmentation, and large allocations

### MCP Resources

- **mempro://stats** - Quick access to memory statistics in JSON format

## Installation

1. Ensure Go 1.22+ is installed
2. Navigate to the MemProMCP directory:
   ```bash
   cd "C:\Program Files\PureDevSoftware\MemPro\MemProReader\MemProMCP"
   ```

3. Download dependencies:
   ```bash
   go mod tidy
   ```

4. Build the server:
   ```bash
   go build -o mempro-mcp.exe
   ```

## Usage

### Running the Server

The server uses stdio transport for communication with MCP clients:

```bash
./mempro-mcp.exe
```

### Integration with Claude Desktop

Add to your Claude Desktop configuration (`claude_desktop_config.json`):

```json
{
  "mcpServers": {
    "mempro": {
      "command": "C:\\Program Files\\PureDevSoftware\\MemPro\\MemProReader\\MemProMCP\\mempro-mcp.exe",
      "env": {
        "MEMPRO_JSON_PATH": "C:\\Program Files\\PureDevSoftware\\MemPro\\MemProReader\\test_memory_analysis.json"
      }
    }
  }
}
```

### Environment Variables

- `MEMPRO_JSON_PATH` - Default path to MemPro JSON file (optional)

## Analysis Capabilities

### Memory Leak Detection

The analyzer identifies memory leaks based on:
- Leak size and count
- Leak score from MemPro
- Suspect flags
- Function context and call stacks

Severity levels:
- **Critical**: Suspect leaks > 100KB
- **High**: Suspect leaks or leaks > 50KB
- **Medium**: Leaks > 10KB or > 100 allocations
- **Low**: Other leaks

### Intelligent Suggestions

The analyzer provides context-aware suggestions:
- **Unknown functions**: Enable debug symbols
- **STL containers**: Check destructors and circular references
- **Main function**: Review allocation ownership, use RAII
- **General**: Use smart pointers (std::unique_ptr, std::shared_ptr)

### Fragmentation Analysis

Detects memory fragmentation issues:
- **High** (>80%): Severe fragmentation requiring immediate attention
- **Medium** (>50%): Moderate fragmentation to monitor

Recommendations include object pooling, memory arenas, and allocation pattern optimization.

### Large Allocation Detection

Identifies allocations that may benefit from optimization:
- Average size > 10KB
- Maximum size > 50KB (Medium) or > 100KB (High)

Suggests chunking, streaming, or incremental allocation strategies.

## Example Queries for AI

When using this server with an AI assistant:

1. "What are the most critical memory leaks in the application?"
2. "Analyze memory fragmentation and provide recommendations"
3. "Show me the top 5 functions causing memory leaks"
4. "Give me a comprehensive analysis of all memory issues"
5. "What large allocations should I optimize?"

## Data Structure

The server parses MemPro JSON with the following key sections:
- **CallTrees**: Hierarchical allocation call trees
- **Functions**: Function-level allocation statistics
- **Leaks**: Detected memory leaks with suspect flags
- **PageViews**: Memory page usage information
- **Types**: Allocation type statistics

## Development

### Project Structure

```
MemProMCP/
├── main.go       # MCP server setup and tool handlers
├── analyzer.go   # Memory analysis logic
├── types.go      # Data structures for MemPro JSON
├── go.mod        # Go module definition
└── README.md     # This file
```

### Adding New Tools

1. Define tool in `setupTools()` in main.go
2. Create handler function following the pattern
3. Implement analysis logic in analyzer.go
4. Update README with tool documentation

## License

This tool is designed to work with PureDevSoftware's MemPro memory profiler.

## Support

For issues or questions about MemPro, visit: https://www.puredevsoftware.com/mempro/

For MCP protocol information: https://modelcontextprotocol.io/
