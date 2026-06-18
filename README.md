# Youshitsune
A minimalist, high-performance personal portfolio and blog for me.

**[👉 Visit the Live Site](https://youshitsune.top/)**

## ⚡ Quick Start

Since this is a static-content site served by a Go backend, you can get it running locally in seconds:

```bash
# Install dependencies and run
go mod tidy
go run main.go
```
Then open `http://localhost:8080` in your browser.

## Features

- **Markdown Blog**: Full support for Markdown posts with YAML frontmatter for dates, tags, and descriptions.
- **Syntax Highlighting**: Integrated Monokai-style code highlighting via `goldmark-highlighting`.
- **Dynamic Tagging**: Automated tag indexing and filtering for easy discovery of related content.
- **Interactive UI**: A clean, responsive design featuring a particle-system background and a dedicated project showcase.
- **Fast Rendering**: Built with the Echo framework for minimal overhead and rapid page loads.

## Local Development

### Prerequisites
- **Language**: Go 1.20+
- **Dependencies**: The project uses `github.com/labstack/echo` for routing and `github.com/yuin/goldmark` for Markdown processing.

### Setup
1. **Clone the repository**
   ```bash
   git clone https://github.com/youshitsune/youshitsune.github.io.git
   cd youshitsune.github.io
   ```

2. **Install Go modules**
   ```bash
   go mod tidy
   ```

3. **Run the server**
   ```bash
   go run main.go
   ```

## How it Works

The site is engineered as a hybrid between a static site generator and a live web server. Instead of a heavy build step, it uses a **custom-built Markdown pipeline** in Go:

- **Pre-parsing**: At startup, the server scans the `posts/` directory, parses YAML frontmatter for metadata, and converts Markdown to HTML using `goldmark`.
- **In-Memory Cache**: The processed posts are stored in memory (maps and slices), ensuring that page requests are served with virtually zero latency.
- **Templating**: It utilizes Go's `text/template` engine to inject dynamic content into HTML layouts, allowing for a consistent look and feel across the home page, blog posts, and tag views.

## Credits

- [Echo Framework](https://echo.labstack.com/) - High performance web framework for Go.
- [Goldmark](https://github.com/yuin/goldmark) - Markdown parser that follows the CommonMark spec.
- [Particles.js](https://vincentgarreau.com/particles.js/) - For the interactive background effects.
