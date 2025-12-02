# CLAUDE.md - AI Assistant Guide for codelab-tools

This document provides comprehensive guidance for AI assistants working with the codelab-tools repository. It covers codebase structure, development workflows, build systems, and key conventions.

## Table of Contents

- [Repository Overview](#repository-overview)
- [Directory Structure](#directory-structure)
- [Technologies & Dependencies](#technologies--dependencies)
- [Development Setup](#development-setup)
- [Build Systems](#build-systems)
- [Testing](#testing)
- [Common Development Tasks](#common-development-tasks)
- [Key Conventions](#key-conventions)
- [Architecture Patterns](#architecture-patterns)
- [Important Files](#important-files)
- [Git Workflow](#git-workflow)
- [Troubleshooting](#troubleshooting)

## Repository Overview

**Purpose**: Tools for authoring and serving interactive codelabs (tutorials)

**Main Components**:
1. **claat** - CLI tool (Go) for converting Google Docs/Markdown to codelabs
2. **codelab-elements** - Web components (JavaScript) for rendering codelabs
3. **site** - Landing page generator (Gulp/Polymer) for hosting codelab collections

**Module Path**: `github.com/chenglu/codelab-tools/claat`

**Current Version**: 3.0.7 (see claat/VERSION)

## Directory Structure

```
codelab-tools/
├── .github/workflows/     # CI/CD automation (Go tests, releases)
├── claat/                 # Go CLI tool (main component)
│   ├── cmd/              # Command implementations (export, serve, update)
│   ├── parser/           # Input parsers (Google Docs, Markdown)
│   ├── nodes/            # AST node types
│   ├── render/           # Output renderers (HTML, Markdown, Offline)
│   ├── types/            # Core data structures
│   ├── fetch/            # Resource fetching
│   ├── util/             # Utilities
│   ├── Makefile          # Build configuration
│   ├── go.mod            # Go dependencies
│   └── VERSION           # Current version
├── codelab-elements/      # Web components for rendering
│   ├── google-codelab/   # Main codelab component
│   ├── google-codelab-step/
│   ├── google-codelab-index/
│   └── */BUILD.bazel     # Component build files
├── site/                  # Landing page generator
│   ├── app/              # Frontend application
│   │   ├── views/        # Event/category-specific views
│   │   └── bower_components/  # Polymer dependencies
│   ├── gulpfile.js       # Build pipeline
│   └── package.json      # NPM dependencies
├── sample/                # Example codelab in Markdown
├── scripts/               # Build/release automation
├── third_party/           # External dependencies (Bazel)
├── BUILD.bazel            # Root build file
├── WORKSPACE              # Bazel workspace config
├── package.json           # NPM packaging (Bazel wrapper)
├── README.md              # User documentation
├── FORMAT-GUIDE.md        # Codelab authoring guide
└── CONTRIBUTING.md        # Contribution guidelines
```

## Technologies & Dependencies

### Programming Languages
- **Go 1.23.0+** (claat CLI, primary language)
- **JavaScript (ES6)** (web components, site)
- **SCSS/Sass** (styling)
- **Soy Templates** (Google Closure templating)

### Go Dependencies (claat/go.mod)
```go
require (
    github.com/google/go-cmp v0.6.0          // Testing utilities
    github.com/stoewer/go-strcase v1.2.0     // String case conversion
    github.com/x1ddos/csslex v0.0.0-...      // CSS lexing
    github.com/yuin/goldmark v1.4.13         // Markdown parsing
    golang.org/x/net v0.38.0                 // HTTP utilities
    golang.org/x/oauth2 v0.1.0               // Google Drive auth
)
```

### JavaScript Libraries
- Google Closure Library & Compiler (web components)
- Polymer (legacy, used in site)
- Gulp 4.0 (site build system)
- Babel (ES6 transpilation)
- Prettify.js (syntax highlighting)

### Build Tools
- **Bazel 0.18.1** (web components) - See `.bazelversion`
- **Make** (claat CLI)
- **Go Modules** (dependency management)
- **Gulp** (site generator)

## Development Setup

### Prerequisites

1. **For claat (Go CLI)**:
   ```bash
   # Install Go 1.23.0 or later
   # https://golang.org/dl/

   # Verify installation
   go version  # Should be 1.23.0+
   ```

2. **For codelab-elements (Web Components)**:
   ```bash
   # Install Bazel 0.18.1
   # https://docs.bazel.build/versions/0.18.0/install.html

   # Or use npm (package.json includes @bazel/bazel)
   npm install
   ```

3. **For site (Landing Page)**:
   ```bash
   cd site
   npm install
   ```

### Initial Setup

```bash
# Clone repository
git clone https://github.com/chenglu/codelab-tools.git
cd codelab-tools

# Build claat CLI
cd claat
make
# Binary created at: claat/bin/claat

# Test installation
./bin/claat version
```

## Build Systems

### claat (Make + Go Modules)

**Location**: `/home/user/codelab-tools/claat/Makefile`

**Key Targets**:
```bash
make              # Build local binary (claat/bin/claat)
make test         # Run Go tests
make lint         # Run golint
make serve        # Build and serve codelabs locally
make release      # Cross-compile for linux/darwin/windows (amd64, 386)
make clean        # Remove build artifacts
```

**Build Output**:
- Development: `claat/bin/claat`
- Release: `claat/bin/claat-{os}-{arch}`

**Version Injection**:
```makefile
VERSION := $(shell cat VERSION)
LDFLAGS := -ldflags "-X main.version=$(VERSION)"
```

**Manual Build**:
```bash
cd claat
go build -o bin/claat ./
go build -ldflags "-X main.version=3.0.7" -o bin/claat ./
```

### codelab-elements (Bazel)

**Location**: Root directory + `codelab-elements/` subdirectories

**Key Commands**:
```bash
# Build web components bundle
bazel build :bundle
# Output: bazel-bin/bundle.tar

# Build NPM distribution
bazel build npm_dist
# Output: bazel-genfiles/npm_dist.zip

# Run specific component tests
bazel test --test_output=all //codelab-elements/demo:hello_test
bazel test //codelab-elements/google-codelab-step:google_codelab_step_test

# Run dev server
bazel run //codelab-elements/tools:server
# Serves on http://localhost:8080

# Clean build artifacts
bazel clean
```

**Via NPM Scripts** (uses @bazel/bazel):
```bash
npm run build     # bazel build npm_dist
npm run test      # bazel test
npm run clean     # bazel clean
```

**Component Build Pattern**:
Each component has a `BUILD.bazel` with:
- `closure_js_library` - Source code
- `closure_js_binary` - Compiled JS
- `sass_binary` - Compiled CSS
- `closure_js_test` - Tests

### site (Gulp)

**Location**: `/home/user/codelab-tools/site/`

**Key Commands**:
```bash
cd site

# Development server with live reload
gulp serve

# Test production build locally
gulp serve:dist

# Build for production
gulp dist
# Output: site/dist/

# Deploy to staging
gulp publish:staging:views --staging-bucket=gs://your-bucket

# Deploy to production
gulp publish:prod:codelabs --prod-bucket=gs://your-bucket
```

**Filters** (for selective builds):
```bash
# Build only specific views
gulp serve --views-filter='^event-'

# Build only specific codelabs
gulp dist --codelabs-filter='android|ios'
```

## Testing

### Go Tests (claat)

**Convention**: `*_test.go` files alongside source

**Run Tests**:
```bash
cd claat

# All tests
go test ./...

# Verbose output
go test -v ./...

# Specific package
go test ./parser/md/...

# With coverage
go test -cover ./...

# Via Makefile
make test
```

**Test Structure**:
```
claat/
├── cmd/export_test.go
├── nodes/
│   ├── code_test.go
│   ├── image_test.go
│   └── ...
├── parser/
│   ├── gdoc/*_test.go
│   └── md/*_test.go
└── render/*_test.go
```

### JavaScript Tests (codelab-elements)

**Framework**: Web Component Tester via Bazel

**Run Tests**:
```bash
# Specific test
bazel test --test_output=all //codelab-elements/demo:hello_test

# All tests in package
bazel test --test_output=errors //codelab-elements/...

# Via NPM
npm run test
```

**Writing Tests**:
Tests use Chromium via rules_webtesting. Example:
```javascript
// google_codelab_step_test.js
suite('google-codelab-step', function() {
  test('instantiating the element works', function() {
    const element = document.createElement('google-codelab-step');
    assert.equal(element.is, 'google-codelab-step');
  });
});
```

## Common Development Tasks

### Working with claat CLI

**1. Export Codelab from Google Doc**:
```bash
cd claat
./bin/claat export <google-doc-id>
# Output: ./<codelab-id>/index.html
```

**2. Export from Markdown**:
```bash
./bin/claat export path/to/document.md
```

**3. Export to Different Format**:
```bash
./bin/claat -f markdown export <source>
./bin/claat -f offline export <source>
```

**4. Serve Codelabs Locally**:
```bash
./bin/claat serve
# Opens http://localhost:9090

# Custom port
./bin/claat serve -addr :8080
```

**5. Update Existing Codelab**:
```bash
./bin/claat update <codelab-directory>
```

### Working with Web Components

**1. Develop Component**:
```bash
# Make changes to component source
# e.g., codelab-elements/google-codelab/google_codelab.js

# Build component
bazel build //codelab-elements/google-codelab:google_codelab_bin

# Test component
bazel test //codelab-elements/google-codelab:google_codelab_test

# Run dev server to preview
bazel run //codelab-elements/tools:server
```

**2. Build NPM Package**:
```bash
npm run build
# Creates: bazel-genfiles/npm_dist.zip

# Publish (requires permissions)
npm run pub
```

**3. Add New Component**:
```bash
# 1. Create directory
mkdir codelab-elements/google-codelab-newfeature

# 2. Create files:
#    - google_codelab_newfeature.js
#    - google_codelab_newfeature_def.js
#    - google_codelab_newfeature.scss
#    - google_codelab_newfeature_test.js
#    - BUILD.bazel

# 3. Add to parent BUILD.bazel filegroup
```

### Working with Site Generator

**1. Create New View**:
```bash
cd site

# 1. Create directory
mkdir -p app/views/my-event

# 2. Create view.json
cat > app/views/my-event/view.json << EOF
{
  "title": "My Event Codelabs",
  "description": "Codelabs for My Event 2025",
  "tags": ["event", "my-event"],
  "categories": ["Web", "Mobile", "Cloud"],
  "logoUrl": "/images/logo.png",
  "url": "my-event"
}
EOF

# 3. (Optional) Create custom styles
touch app/views/my-event/style.css

# 4. Preview
gulp serve
```

**2. Deploy Codelabs**:
```bash
# Build
gulp dist

# Deploy to GCS (example)
gulp publish:prod:codelabs \
  --staging-bucket=gs://staging-bucket \
  --prod-bucket=gs://production-bucket
```

## Key Conventions

### Go Code Style (claat)

**Package Organization**:
```
claat/
├── main.go              # Entry point, CLI flags
├── cmd/                 # Commands (CmdExport, CmdServe, CmdUpdate)
├── parser/              # Input parsers
│   ├── common.go       # Parser registry
│   ├── gdoc/           # Google Docs implementation
│   └── md/             # Markdown implementation
├── nodes/               # AST node types
├── render/              # Output renderers
├── types/               # Data structures (Codelab, Step, Meta)
└── fetch/               # Resource fetching
```

**Naming Conventions**:
- **Files**: `lowercase_with_underscores.go`
- **Tests**: `*_test.go`
- **Packages**: Short, lowercase (cmd, parser, nodes)
- **Types**: PascalCase (Codelab, Step, NodeText)
- **Node Types**: Prefix with "Node" (NodeCode, NodeImage)
- **Interfaces**: Simple names (Parser, Node, Renderer)
- **Exported**: Capital first letter
- **Private**: lowercase first letter

**Architecture Pattern**:
- **Plugin-based parsers**: Register via `parser.Register("name", parser)`
- **AST**: Content as tree of `nodes.Node` interfaces
- **Template rendering**: Go templates in `render/` directory

**Example**:
```go
// parser/md/parser.go
func init() {
    parser.Register("md", NewParser())
}

type Parser struct {
    // ...
}

func (p *Parser) Parse(r io.Reader, opts *parser.Options) (*types.Codelab, error) {
    // Returns AST
}
```

### Web Component Structure

**Component Files**:
```
google-codelab-component/
├── BUILD.bazel                         # Build rules
├── google_codelab_component.js         # Main logic
├── google_codelab_component_def.js     # Element registration
├── google_codelab_component.scss       # Styles
├── google_codelab_component.soy        # Templates (optional)
└── google_codelab_component_test.js    # Tests
```

**Naming**:
- **HTML Elements**: kebab-case (`<google-codelab-step>`)
- **Files**: snake_case (`google_codelab_step.js`)
- **Classes**: PascalCase (GoogleCodelabStep)
- **Bazel targets**: snake_case with `_bin` for binaries

**BUILD.bazel Pattern**:
```python
closure_js_library(
    name = "google_codelab_component",
    srcs = ["google_codelab_component.js"],
)

closure_js_binary(
    name = "google_codelab_component_bin",
    deps = [":google_codelab_component"],
)

sass_binary(
    name = "google_codelab_component_scss_bin",
    src = "google_codelab_component.scss",
)

closure_js_test(
    name = "google_codelab_component_test",
    srcs = ["google_codelab_component_test.js"],
    deps = [":google_codelab_component"],
)
```

### Codelab Metadata

**Markdown Format** (at top of `.md` file):
```markdown
summary: Brief description of the codelab
id: unique-codelab-id
categories: Web, Mobile
environments: Web, Kiosk
status: Published
feedback link: https://github.com/user/repo/issues
analytics account: UA-XXXXX-Y
authors: Author Name
tags: tag1, tag2

# Codelab Title

## Step 1: Introduction
Duration: 5:00

Step content here...
```

**Google Docs Format**:
Two-column table before first Heading 1:
```
| Summary          | Brief description       |
| URL              | unique-codelab-id       |
| Category         | Web                     |
| Environment      | Web, Kiosk              |
| Status           | Published               |
| Feedback Link    | https://...             |
| Analytics Account| UA-XXXXX-Y              |
```

**Required Metadata**:
- `summary` - Short description
- `id` (Markdown) or `URL` (Google Docs) - Unique identifier
- `categories` or `Category` - Platform/topic grouping

**Optional Metadata**:
- `environments` - Target environments (default: "Web, Kiosk")
- `status` - Draft, Published, Deprecated, Hidden
- `feedback link` - Bug report URL
- `analytics account` - Custom Google Analytics ID
- `authors` - Author names
- `tags` - Search/filtering tags

### File Path Conventions

**Absolute Paths**: Always use absolute paths in code
```go
// Good
filepath := "/home/user/codelab-tools/claat/..."

// Bad
filepath := "../claat/..."
```

**Output Structure**:
```
output-directory/
└── codelab-id/
    ├── index.html          # Main HTML
    ├── codelab.json        # Metadata
    └── img/                # Images
```

**Generated Files**:
- HTML: `<codelab-id>/index.html`
- Markdown: `<codelab-id>/codelab.md`
- JSON: `<codelab-id>/codelab.json`

## Architecture Patterns

### claat CLI Architecture

**Command Pattern**:
```go
// main.go - Entry point
func main() {
    // Parse flags
    flag.Parse()

    // Route to command
    switch flag.Arg(0) {
    case "export":
        cmd.CmdExport(...)
    case "serve":
        cmd.CmdServe(...)
    case "update":
        cmd.CmdUpdate(...)
    }
}
```

**Parser Plugin System**:
```go
// parser/common.go
var parsers = make(map[string]Parser)

func Register(name string, p Parser) {
    parsers[name] = p
}

// parser/md/parser.go
func init() {
    parser.Register("md", NewParser())
}
```

**AST Node System**:
```go
// nodes/nodes.go
type Node interface {
    Type() NodeType
    Empty() bool
    // ...
}

// Implementations:
type NodeText struct { Value string }
type NodeCode struct { Code string; Lang string }
type NodeImage struct { Src string; Alt string }
// etc.
```

**Template-Based Rendering**:
```go
// render/html.go
func Execute(c *types.Codelab, w io.Writer) error {
    tmpl := template.Must(template.ParseFiles("template.html"))
    return tmpl.Execute(w, c)
}
```

### Web Component Architecture

**Custom Elements v1**:
```javascript
// google_codelab_component.js
class GoogleCodelabComponent extends HTMLElement {
  constructor() {
    super();
  }

  connectedCallback() {
    // Initialize
  }
}

// google_codelab_component_def.js
customElements.define('google-codelab-component', GoogleCodelabComponent);
```

**Event-Driven Communication**:
```javascript
// Dispatch custom event
this.dispatchEvent(new CustomEvent('codelab-action', {
  detail: { action: 'next' },
  bubbles: true
}));

// Listen for event
element.addEventListener('codelab-action', (e) => {
  // Handle event
});
```

### Site Generator Architecture

**View-Based Organization**:
```
site/app/views/
├── default/                 # Default view
│   └── view.json
├── event-name/              # Event-specific view
│   ├── view.json
│   └── style.css
└── category-name/           # Category view
    └── view.json
```

**view.json Schema**:
```json
{
  "title": "View Title",
  "description": "View description",
  "tags": ["tag1", "tag2"],
  "categories": ["Category1", "Category2"],
  "logoUrl": "/images/logo.png",
  "url": "view-url-suffix",
  "ga": "UA-XXXXX-Y"
}
```

**Gulp Task Pipeline**:
```javascript
// gulpfile.js
gulp.task('serve', gulp.series('build', 'watch', 'server'));
gulp.task('dist', gulp.series('clean', 'build', 'minify'));
gulp.task('publish', gulp.series('dist', 'upload'));
```

## Important Files

### Configuration Files

| File | Purpose |
|------|---------|
| `claat/go.mod` | Go module dependencies |
| `claat/go.sum` | Dependency checksums |
| `claat/VERSION` | Current version (3.0.7) |
| `claat/Makefile` | Build configuration |
| `.bazelversion` | Required Bazel version (0.18.1) |
| `WORKSPACE` | Bazel workspace with external deps |
| `BUILD.bazel` | Root Bazel build targets |
| `package.json` | NPM packaging, Bazel wrapper |
| `site/package.json` | Site build dependencies |
| `site/gulpfile.js` | Complete build pipeline |

### CI/CD Files

| File | Purpose |
|------|---------|
| `.github/workflows/release.yml` | Automated releases on main branch |
| `.github/workflows/go.yml` | Go CI (vet, build, test) |
| `scripts/bump_version.sh` | Version bumping script |

### Documentation Files

| File | Purpose |
|------|---------|
| `README.md` | User documentation |
| `FORMAT-GUIDE.md` | Codelab authoring guide |
| `CONTRIBUTING.md` | Contribution guidelines |
| `claat/README.md` | claat CLI documentation |
| `site/README.md` | Site generator documentation |
| `CLAUDE.md` | This file - AI assistant guide |

### Key Source Files

**claat**:
- `claat/main.go` - Entry point
- `claat/cmd/export.go` - Export command
- `claat/cmd/serve.go` - Serve command
- `claat/parser/md/parser.go` - Markdown parser
- `claat/parser/gdoc/parse.go` - Google Docs parser
- `claat/render/html.go` - HTML renderer
- `claat/types/codelab.go` - Core data structures

**codelab-elements**:
- `codelab-elements/google-codelab/google_codelab.js` - Main component
- `codelab-elements/google-codelab-step/google_codelab_step.js` - Step component
- `codelab-elements/google-codelab-index/google_codelab_index.js` - Index component

**site**:
- `site/gulpfile.js` - Build pipeline
- `site/app/index.html` - Main HTML
- `site/app/scripts/app.js` - Application logic

## Git Workflow

### Branch Naming

**Current Branch**: `claude/claude-md-mip3zk6leudad6qp-01KoK8e8eJyodkZvVMC4b1z3`

**Pattern**: `claude/claude-md-<session-id>-<unique-id>`

**Important**: When pushing, branch must:
- Start with `claude/`
- End with matching session ID
- Otherwise, push fails with 403 HTTP code

### Commit Messages

**Pattern**:
```
type: brief description

[optional body]

[optional footer]
```

**Types**: `chore`, `feat`, `fix`, `docs`, `test`, `refactor`

**Special Markers**:
- `[skip release]` - Prevents automatic version bump on main branch

**Examples**:
```
chore: bump claat to v3.0.7 [skip release]
feat: add support for custom templates
fix: resolve parsing error with nested lists
docs: update FORMAT-GUIDE with new metadata fields
```

### Release Process

**Automatic** (on main branch):
1. Push to main branch (without `[skip release]`)
2. GitHub Actions runs `.github/workflows/release.yml`:
   - Runs `scripts/bump_version.sh`
   - Updates `claat/VERSION`
   - Builds cross-platform binaries (linux, darwin, windows)
   - Commits version bump
   - Creates git tag (e.g., `v3.0.8`)
   - Publishes GitHub release with binaries

**Manual** (on feature branch):
```bash
cd claat

# 1. Update VERSION file
echo "3.0.8" > VERSION

# 2. Build release binaries
make release
# Creates: bin/claat-{linux,darwin,windows}-{amd64,386}

# 3. Test binaries
./bin/claat-linux-amd64 version

# 4. Commit and push
git add VERSION
git commit -m "chore: bump claat to v3.0.8"
git push origin <branch-name>
```

### Git Push/Pull Best Practices

**Push with Retry**:
```bash
# Use -u flag for first push
git push -u origin claude/claude-md-...

# If network error, retry with exponential backoff:
# Wait 2s, try again
# Wait 4s, try again
# Wait 8s, try again
# Wait 16s, try again (max 4 retries)
```

**Fetch/Pull**:
```bash
# Fetch specific branch
git fetch origin <branch-name>

# Pull specific branch
git pull origin <branch-name>

# Apply same retry logic for network errors
```

**Never**:
- Push to main/master without permission
- Use `git push --force` (especially to main)
- Skip hooks with `--no-verify` unless explicitly requested
- Amend other developers' commits

## Troubleshooting

### claat Build Issues

**Issue**: `go: module not found`
```bash
# Solution: Download dependencies
cd claat
go mod download
go mod tidy
```

**Issue**: `make: command not found`
```bash
# Solution: Build manually with go
cd claat
go build -o bin/claat ./
```

**Issue**: Version not showing correctly
```bash
# Solution: Build with version flag
cd claat
VERSION=$(cat VERSION)
go build -ldflags "-X main.version=$VERSION" -o bin/claat ./
```

### Bazel Build Issues

**Issue**: `bazel: command not found`
```bash
# Solution: Use NPM wrapper
npm install
npm run build
```

**Issue**: Wrong Bazel version
```bash
# Check required version
cat .bazelversion  # Shows 0.18.1

# Install specific version
# https://github.com/bazelbuild/bazel/releases/tag/0.18.1
```

**Issue**: Build fails with "missing dependency"
```bash
# Solution: Clean and rebuild
bazel clean --expunge
bazel build :bundle
```

### Testing Issues

**Issue**: Go tests fail with import errors
```bash
# Solution: Ensure you're in correct directory
cd claat
go test ./...

# Not from root:
# go test ./claat/...  # Wrong
```

**Issue**: JavaScript tests won't run
```bash
# Solution: Use full Bazel path
bazel test --test_output=all //codelab-elements/demo:hello_test

# Not just:
# bazel test hello_test  # Wrong
```

### Runtime Issues

**Issue**: `claat serve` port already in use
```bash
# Solution: Use different port
./bin/claat serve -addr :8080
```

**Issue**: Cannot access Google Doc (authentication error)
```bash
# Solution: claat will prompt for OAuth authentication
# Follow the URL and authorize access
# Credentials stored in ~/.config/claat/
```

**Issue**: Exported codelab missing images
```bash
# Solution: Check image permissions in Google Doc
# Images must be publicly accessible or shared
```

### Site Build Issues

**Issue**: `gulp: command not found`
```bash
# Solution: Install dependencies
cd site
npm install

# Run with npx
npx gulp serve
```

**Issue**: View not appearing
```bash
# Solution: Verify view.json format
cat app/views/my-view/view.json

# Must be valid JSON with required fields
```

**Issue**: Deploy fails with GCS error
```bash
# Solution: Verify bucket permissions
gsutil ls gs://your-bucket

# Check gulp command syntax
gulp publish:prod:codelabs --prod-bucket=gs://your-bucket
```

### Common Error Messages

| Error | Cause | Solution |
|-------|-------|----------|
| `405 Not Allowed` | Wrong HTTP method | Check API endpoint |
| `403 Forbidden` | Push to wrong branch | Verify branch name starts with `claude/` |
| `module not found` | Missing dependencies | Run `go mod download` |
| `bazel: unknown target` | Wrong target name | Use `//package:target` format |
| `port already in use` | Server already running | Change port with `-addr` flag |
| `template not found` | Missing template file | Check render/ directory |

## Additional Resources

### Official Documentation
- **Google Codelabs**: https://g.co/codelabs
- **Codelab Authors Group**: https://groups.google.com/forum/#!forum/codelab-authors
- **GitHub Issues**: https://github.com/chenglu/codelab-tools/issues
- **Releases**: https://github.com/chenglu/codelab-tools/releases

### Tutorials
- **Creating Codelabs**: https://medium.com/@zarinlo/publish-technical-tutorials-in-google-codelab-format-b07ef76972cd
- **FORMAT-GUIDE.md**: Complete formatting reference for authoring codelabs
- **claat/README.md**: CLI tool documentation

### Internal Documentation
- **README.md**: User-facing overview
- **CONTRIBUTING.md**: Contribution guidelines
- **claat/parser/md/README.md**: Markdown parser specifics

## Quick Reference

### Common Commands Cheat Sheet

```bash
# claat CLI
cd claat
make                                    # Build
make test                               # Test
./bin/claat export <doc-id>            # Export from Google Doc
./bin/claat export document.md          # Export from Markdown
./bin/claat serve                       # Serve locally

# codelab-elements
bazel build :bundle                     # Build bundle
bazel test //codelab-elements/...       # Test all
bazel run //codelab-elements/tools:server  # Dev server
npm run build                           # Build via NPM

# site
cd site
gulp serve                              # Dev server
gulp dist                               # Production build
gulp publish:prod:codelabs --prod-bucket=gs://bucket  # Deploy

# Git
git push -u origin <branch>            # Push with upstream
git fetch origin <branch>               # Fetch branch
git log --oneline -10                   # Recent commits
```

### File Locations Cheat Sheet

```
claat binary:      claat/bin/claat
claat version:     claat/VERSION
go dependencies:   claat/go.mod
build config:      claat/Makefile
web components:    codelab-elements/
bazel config:      WORKSPACE, BUILD.bazel, .bazelversion
site generator:    site/
site config:       site/gulpfile.js, site/package.json
views:             site/app/views/
CI/CD:             .github/workflows/
docs:              README.md, FORMAT-GUIDE.md, CONTRIBUTING.md
```

---

**Document Version**: 1.0
**Last Updated**: 2025-12-02
**Repository Version**: claat v3.0.7

This guide is maintained for AI assistants working with the codelab-tools codebase. For user-facing documentation, see README.md and FORMAT-GUIDE.md.
