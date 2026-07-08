<h1 align="center">
  <img src="static/httpx-logo.png" alt="httpx-cr" width="200px">
  <br>
</h1>

<p align="center">
<a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/license-MIT-_red.svg"></a>
<a href="https://github.com/dace-cyber/httpx-cr/releases"><img src="https://img.shields.io/github/release/projectdiscovery/httpx.svg"></a>
</p>

<p align="center">
  <a href="#copyright-extraction">Copyright Flag</a> •
  <a href="#installation">Installation</a> •
  <a href="#usage">Usage</a> •
  <a href="#credits">Credits</a>
</p>

**httpx-cr** is a fork of [projectdiscovery/httpx](https://github.com/projectdiscovery/httpx) with an added **`-copyright`** flag that extracts copyright notices from web pages.

`httpx` is a fast and multi-purpose HTTP toolkit that allows running multiple probes using the [retryablehttp](https://github.com/projectdiscovery/retryablehttp-go) library.

# Copyright Extraction

The `-copyright` flag automatically detects and extracts copyright information from responses. It searches three sources in order:

| Priority | Source | Example |
|----------|--------|---------|
| 1 | `<meta name="copyright">` / `<meta name="rights">` | `<meta name="copyright" content="© 2026 Company">` |
| 2 | `<footer>` / `<small>` elements containing `©` or "copyright" | `<footer>© 2026 GitHub, Inc.</footer>` |
| 3 | Any `©` or "copyright" text in the page body (fallback) | `© 2026 Example Corp` |

# Installation

### From source (Go required)

```bash
go install github.com/dace-cyber/httpx-cr/cmd/httpx@latest
```

### Build from source

```bash
git clone https://github.com/dace-cyber/httpx-cr.git
cd httpx-cr
go build -o httpx-cr ./cmd/httpx/
```

### Download binary

Pre-built binaries are available on the [releases page](https://github.com/dace-cyber/httpx-cr/releases).

# Usage

```bash
# Single URL
echo "https://example.com" | httpx-cr -copyright

# Multiple URLs from file
httpx-cr -l urls.txt -copyright

# With JSON output
echo "https://github.com" | httpx-cr -copyright -j

# Combined with other probes
echo "https://github.com" | httpx-cr -copyright -title -sc -td
```

### Example output

#### CLI
```
https://github.com [© 2026 GitHub, Inc.]
```

#### JSON
```json
{"url":"https://github.com","copyright":"© 2026 GitHub, Inc.","title":"GitHub ..."}
```

#### CSV
```csv
url,copyright,title
https://github.com,© 2026 GitHub, Inc.,GitHub ...
```

### All original httpx probes still available

| Probe | Flag | Probe | Flag |
|-------|------|-------|------|
| URL | (default) | IP | `-ip` |
| Title | `-title` | CNAME | `-cname` |
| Status Code | `-sc` | CDN | `-cdn` |
| Content Length | `-cl` | Response Time | `-rt` |
| Copyright | `-copyright` | Technology | `-td` |

For the full list of original httpx flags, see the [official documentation](https://docs.projectdiscovery.io/tools/httpx/).

# Credits

- **Original httpx**: [projectdiscovery/httpx](https://github.com/projectdiscovery/httpx)
- **Copyright flag**: Added by [dace-cyber](https://github.com/dace-cyber)
