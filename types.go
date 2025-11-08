package main

// MemProData represents the complete MemPro analysis JSON structure
type MemProData struct {
	SessionName          string        `json:"SessionName"`
	TotalSnapshots       int           `json:"TotalSnapshots"`
	TotalAllocations     int           `json:"TotalAllocations"`
	TotalSize            int64         `json:"TotalSize"`
	LeakCount            int           `json:"LeakCount"`
	LeakSize             int64         `json:"LeakSize"`
	MemoryFragmentation  float64       `json:"MemoryFragmentation"`
	CallTrees            []CallTree    `json:"CallTrees"`
	Functions            []Function    `json:"Functions"`
	Leaks                []Leak        `json:"Leaks"`
	PageViews            []PageView    `json:"PageViews"`
	Types                []AllocType   `json:"Types"`
}

// CallTree represents a call tree entry with allocation information
type CallTree struct {
	FunctionName    string     `json:"FunctionName"`
	FileName        string     `json:"FileName"`
	LineNumber      int        `json:"LineNumber"`
	AllocationCount int        `json:"AllocationCount"`
	TotalSize       int64      `json:"TotalSize"`
	SelfSize        int64      `json:"SelfSize"`
	InclusiveSize   int64      `json:"InclusiveSize"`
	Children        []CallTree `json:"Children"`
}

// Function represents function-level allocation statistics
type Function struct {
	FunctionName    string  `json:"FunctionName"`
	FileName        string  `json:"FileName"`
	LineNumber      int     `json:"LineNumber"`
	AllocationCount int     `json:"AllocationCount"`
	TotalSize       int64   `json:"TotalSize"`
	AverageSize     float64 `json:"AverageSize"`
	MinSize         int64   `json:"MinSize"`
	MaxSize         int64   `json:"MaxSize"`
	Percentage      float64 `json:"Percentage"`
}

// Leak represents a memory leak with suspect information
type Leak struct {
	FunctionName string  `json:"FunctionName"`
	FileName     string  `json:"FileName"`
	LineNumber   int     `json:"LineNumber"`
	LeakSize     int64   `json:"LeakSize"`
	LeakCount    int     `json:"LeakCount"`
	LeakScore    float64 `json:"LeakScore"`
	CallStack    string  `json:"CallStack"`
	IsSuspect    bool    `json:"IsSuspect"`
}

// PageView represents memory page usage information
type PageView struct {
	Address         int64  `json:"Address"`
	State           string `json:"State"`
	Type            string `json:"Type"`
	Protection      int    `json:"Protection"`
	StackId         int    `json:"StackId"`
	Usage           int    `json:"Usage"`
	AllocationCount int    `json:"AllocationCount"`
	TotalSize       int64  `json:"TotalSize"`
	FunctionName    string `json:"FunctionName"`
	CallStack       string `json:"CallStack"`
}

// AllocType represents allocation type statistics
type AllocType struct {
	TypeName           string  `json:"TypeName"`
	AllocationCount    int     `json:"AllocationCount"`
	TotalSize          int64   `json:"TotalSize"`
	AverageSize        float64 `json:"AverageSize"`
	MinSize            int64   `json:"MinSize"`
	MaxSize            int64   `json:"MaxSize"`
	Percentage         float64 `json:"Percentage"`
	MostCommonFunction string  `json:"MostCommonFunction"`
	MostCommonFile     string  `json:"MostCommonFile"`
	MostCommonLine     int     `json:"MostCommonLine"`
}

// MemoryIssue represents a detected memory issue for AI analysis
type MemoryIssue struct {
	Severity     string  `json:"severity"`     // Critical, High, Medium, Low
	Type         string  `json:"type"`         // Leak, Fragmentation, LargeAllocation
	Description  string  `json:"description"`
	FunctionName string  `json:"functionName"`
	FileName     string  `json:"fileName"`
	LineNumber   int     `json:"lineNumber"`
	Size         int64   `json:"size"`
	Count        int     `json:"count"`
	Score        float64 `json:"score"`
	Suggestion   string  `json:"suggestion"`
}
