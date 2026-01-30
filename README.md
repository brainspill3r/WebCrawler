# 🕷️ Go Web Crawler

A high-performance, concurrent web crawler built in Go that extracts structured data from websites and exports it to CSV format for easy analysis.

## ✨ Features

- 🚀 **Concurrent Crawling** - Configurable concurrency levels for optimal performance
- 🎯 **Domain-Focused** - Stays within the target domain boundaries
- 📊 **Rich Data Extraction** - Captures page titles, content, links, and images
- 📈 **CSV Export** - Exports structured data for analysis in spreadsheets
- 🛑 **Smart Limits** - Configurable maximum page limits to prevent runaway crawls
- 🔒 **Thread-Safe** - Built with Go's concurrency best practices
- ⚡ **Fast & Efficient** - Leverages goroutines for parallel processing

## 🛠️ Installation

### Prerequisites
- Go 1.19 or higher
- Internet connection for crawling

### Setup
1. Clone the repository:
```bash
git clone <your-repo-url>
cd WebCrawler
```

2. Install dependencies:
```bash
go mod init webcrawler
go get github.com/PuerkitoBio/goquery
```

3. Build the crawler:
```bash
go build -o crawler
```

## 🚀 Usage

### Basic Command
```bash
./crawler <URL> <maxConcurrency> <maxPages>
```

### Parameters
- **URL** - The website to crawl (must include `http://` or `https://`)
- **maxConcurrency** - Number of concurrent requests (1-10 recommended)
- **maxPages** - Maximum number of pages to crawl (prevents runaway crawls)

### Examples

#### 📝 Small Website Crawl
```bash
./crawler "https://example.com" 2 10
```

#### 🏢 Medium Blog Crawl
```bash
./crawler "https://blog.boot.dev/" 3 25
```

#### 🌐 Large Site Crawl
```bash
./crawler "https://wagslane.dev" 5 50
```

### Using `go run`
You can also run directly without building:
```bash
go run . "https://example.com" 3 10
```

## 📋 Output

### Console Output
```
starting crawl of: https://blog.boot.dev/
maxConcurrency: 3
maxPages: 25
crawling: https://blog.boot.dev/
crawling: https://blog.boot.dev/golang/
crawling: https://blog.boot.dev/python/
...

Generating CSV report: report.csv
Crawl completed: 25 pages found (max: 25)
Report saved to report.csv
```

### CSV Export
The crawler generates a `report.csv` file with the following columns:

| Column | Description | Example |
|--------|-------------|---------|
| `page_url` | Full URL of the crawled page | `https://blog.boot.dev/golang/` |
| `h1` | Main heading of the page | `"Learn Go Programming"` |
| `first_paragraph` | First paragraph of content | `"Go is a powerful language..."` |
| `outgoing_link_urls` | All links found on the page | `https://go.dev;https://golang.org` |
| `image_urls` | All images found on the page | `logo.png;banner.jpg` |

**Note:** Multiple links and images are separated by semicolons (`;`)

## �� Project Structure

```
WebCrawler/
├── main.go                     # 🏠 Main application entry point
├── concurrent_crawler.go       # 🕸️ Core crawling logic and concurrency
├── csv_report.go              # 📊 CSV export functionality  
├── extract_page_data.go       # 🔍 Page data extraction
├── get_urls.go                # 🔗 URL and image extraction
├── get_html.go                # 🌐 HTTP client for fetching pages
├── normalize_url.go           # 🧹 URL normalization utilities
├── *_test.go                  # 🧪 Test files
├── go.mod                     # 📦 Go module definition
└── README.md                  # 📖 This file
```

## 🧠 How It Works

1. **🎯 Target Selection** - Starts with the provided URL and parses the domain
2. **�� Concurrent Crawling** - Spawns goroutines up to the concurrency limit  
3. **🕷️ Page Processing** - For each page:
   - Fetches HTML content
   - Extracts H1, first paragraph, links, and images
   - Finds new URLs to crawl
4. **🛑 Smart Limiting** - Stops when max pages reached or no more pages found
5. **📊 Data Export** - Saves all structured data to CSV format

## ⚙️ Technical Details

### Concurrency Model
- Uses **buffered channels** to limit concurrent requests
- **Mutex-protected** shared data structures
- **WaitGroups** ensure all goroutines complete before exit

### Data Extraction
- **HTML parsing** with goquery (jQuery-like selectors)
- **URL normalization** for deduplication
- **Relative to absolute** URL conversion
- **Domain boundary** enforcement

### Performance
- **Non-blocking I/O** with goroutines
- **Memory efficient** page data storage
- **Configurable concurrency** for different network conditions

## 🧪 Testing

Run the test suite:
```bash
go test -v
```

Test specific functionality:
```bash
go test -run TestAddPageVisit
go test -run TestWriteCSVReport
```

## 📊 Example Results

After crawling `https://blog.boot.dev/` with 25 pages:

```csv
page_url,h1,first_paragraph,outgoing_link_urls,image_urls
https://blog.boot.dev/,Boot.dev Blog,Learn backend development with our tutorials,https://boot.dev;https://boot.dev/courses,https://blog.boot.dev/logo.png
https://blog.boot.dev/golang/,Go Tutorials,Master Go programming language,https://go.dev;https://boot.dev/courses/go,https://blog.boot.dev/go-logo.png;https://blog.boot.dev/gopher.jpg
```

## 🚀 Performance Tips

### Concurrency Settings
- **Start with 1-2** for initial testing
- **Use 3-5** for most websites  
- **Up to 10** for high-performance scenarios
- **Be respectful** - don't overwhelm target servers

### Page Limits
- **10-25 pages** for small sites
- **50-100 pages** for medium sites
- **500+ pages** for large sites (be careful!)


