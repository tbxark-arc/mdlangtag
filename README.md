# mdlangtag

A powerful CLI tool that automatically detects programming languages in fenced code blocks within Markdown files and fills in the language info strings. Perfect for maintaining consistent and properly-tagged code blocks across your documentation.

## Features

- 🔍 **Automatic Language Detection**: Uses [chroma](https://github.com/alecthomas/chroma) to intelligently detect code block languages
- 📝 **Batch Processing**: Process single files or entire directories recursively
- ⚡ **Concurrent Processing**: Handle multiple files in parallel for improved performance
- 🎯 **Smart Updates**: Optionally force-update existing language tags or preserve them
- 🛡️ **Safe Operations**: Preview changes via stdout before writing to files
- 🔧 **Flexible Configuration**: Support for fallback languages and minimum line thresholds

## Install

```bash
go install github.com/TBXark/mdlangtag@latest
```

## Usage

### Basic Usage

```bash
# Process a single file and preview changes
mdlangtag README.md

# Process all Markdown files in a directory
mdlangtag docs/

# Write changes directly to files
mdlangtag -w docs/
```

### Options

```bash
-w, --write              Write result back to files (default: false)
    --stdout             Print output to stdout (default: false)
    --force              Overwrite existing language info (default: false)
    --default string     Fallback language when detection fails
    --min-lines int      Skip blocks with fewer than this many lines (default: 0)
-v, --verbose            Enable verbose logging (default: false)
-j, --concurrency int    Number of files to process concurrently (default: 1)
```

### Examples

```bash
# Preview changes for a directory without modifying files
mdlangtag docs/

# Write changes to all markdown files with verbose output
mdlangtag -w -v docs/

# Force update all language tags and skip small code blocks
mdlangtag -w --force --min-lines 3 docs/

# Set a fallback language for undetectable blocks
mdlangtag -w --default "text" README.md

# Process with multiple concurrent workers
mdlangtag -w -j 4 docs/
```

## How It Works

1. **Parsing**: Scans Markdown files for fenced code blocks (triple backticks)
2. **Detection**: Analyzes code content to identify the programming language
3. **Updating**: Fills in language info strings while preserving formatting
4. **Output**: Either displays changes or writes them back to files

## Requirements

- Go 1.25 or later

## License

See [LICENSE](LICENSE) file for details.
