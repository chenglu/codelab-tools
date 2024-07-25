import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/codelabs-cn/codelab-tools/claat/cmd"

	// allow parsers to register themselves
	_ "github.com/codelabs-cn/codelab-tools/claat/parser/gdoc"
	_ "github.com/codelabs-cn/codelab-tools/claat/parser/md"
	_ "github.com/codelabs-cn/codelab-tools/claat/web"
)

var (
	version string // set by linker -X

	// Flags.
	addr         = flag.String("addr", "localhost:9090", "hostname and port to bind web server to")
	authToken    = flag.String("auth", "", "OAuth2 Bearer token; alternative credentials override.")
	expenv       = flag.String("e", "web", "codelab environment")
	extra        = flag.String("extra", "", "Additional arguments to pass to format templates. JSON object of string,string key values.")
	globalGA     = flag.String("ga", "UA-49880327-14", "global Google Analytics account")
	output       = flag.String("o", ".", "output directory or '-' for stdout")
	passMetadata = flag.String("pass_metadata", "", "Metadata fields to pass through to the output. Comma-delimited list of field names.")
	prefix       = flag.String("prefix", "https://static.codelabs.cn", "URL prefix for html format")
	tmplout      = flag.String("f", "html", "output format")
)

func main() {
	log.SetFlags(0)
	rand.Seed(time.Now().UnixNano())
	if len(os.Args) == 1 {
		log.Fatalf("Need subcommand. Try '-h' for options.")
	}
	if os.Args[1] == "-h" || os.Args[1] == "--help" {
		usage()
		return
	}

	flag.Usage = usage
	flag.CommandLine.Parse(os.Args[2:])

	extraVars, err := ParseExtraVars(*extra)
	if err != nil {
		os.Exit(1)
	}

	pm := parsePassMetadata(*passMetadata)

	exitCode := 0
	switch os.Args[1] {
	case "export":
		exitCode = cmd.CmdExport(cmd.CmdExportOptions{
			AuthToken:    *authToken,
			Expenv:       *expenv,
			ExtraVars:    extraVars,
			GlobalGA:     *globalGA,
			Output:       *output,
			PassMetadata: pm,
			Prefix:       *prefix,
			Srcs:         flag.Args(),
			Tmplout:      *tmplout,
		})
	case "serve":
		exitCode = cmd.CmdServe(*addr)
	case "update":
		exitCode = cmd.CmdUpdate(cmd.CmdUpdateOptions{
			AuthToken:    *authToken,
			ExtraVars:    extraVars,
			GlobalGA:     *globalGA,
			PassMetadata: pm,
			Prefix:       *prefix,
		})
	case "web":
		cmd.ServeWebInterface()
	case "help":
		usage()
	case "version":
		fmt.Println(version)
	default:
		log.Fatalf("Unknown subcommand. Try '-h' for options.")
	}

	os.Exit(exitCode)
}

// parsePassMetadata parses metadata fields to parse that are not explicitly handled elsewhere.
// It expects the fields to be passed in as a comma separated list (extraneous spaces are autoremoved), and returns a set of strings.
func parsePassMetadata(passMeta string) map[string]bool {
	fields := map[string]bool{}
	for _, v := range strings.Split(passMeta, ",") {
		fields[strings.ToLower(strings.TrimSpace(v))] = true
	}
	return fields
}

// ParseExtraVars parses extra template variables from command line.
// extra is any additional arguments to pass to format templates. Should be formatted as JSON objects of string:string KV pairs.
func ParseExtraVars(extra string) (map[string]string, error) {
	vars := map[string]string{}
	if extra == "" {
		return vars, nil
	}
	b := []byte(extra)
	err := json.Unmarshal(b, &vars)
	if err != nil {
		log.Printf("Error parsing additional template data: %v", err)
		return nil, err
	}
	return vars, nil
}

// usage prints usageText and program arguments to stderr.
func usage() {
	fmt.Fprint(os.Stderr, usageText)
	flag.PrintDefaults()
}

const usageText = `Usage: claat <cmd> [options] src [src ...]
Available commands are: export, serve, update, web, version.
## Export command
Export takes one or more 'src' documents and converts them
to the format specified with -f option.
The following formats are built-in:
- html (Polymer-based app)
- md (Markdown)
- offline (plain HTML markup for offline consumption)
Note that the built-in templates of the formats are not guaranteed to be stable.
They can be found in https://github.com/codelabs-cn/codelab-tools/tree/master/claat/render.
Please avoid using default templates in production. Use your own copies.
To use a custom format, specify a local file path to a Go template file.
More info on Go templates: https://golang.org/pkg/text/template/.
Each 'src' can be either a remote HTTP resource or a local file.
Source formats currently supported are:
- Google Doc (Codelab Format, go/codelab-guide)
- Markdown
When 'src' is a Google Doc, it must be specified as a doc ID,
omitting https://docs.google.com/... part.
Instead of writing to an output directory, use "-o -" to specify
stdout. In this case images and metadata are not exported.
When writing to a directory, existing files will be overwritten.
The program exits with non-zero code if at least one src could not be exported.
## Serve command
Serve provides a simple web server for viewing exported codelabs.
It takes no arguments and presents the current directory contents.
Clicking on a directory representing an exported codelab will load
all the required dependencies and render the generated codelab as
it would appear in production.
The serve command takes a -addr host:port option, to specify the
desired hostname or IP address and port number to bind to.
## Update command
Update scans one or more 'src' local directories for codelab.json metadata
files, recursively. A directory containing the metadata file is expected
to be a codelab previously created with the export command.
Current directory is assumed if no 'src' argument is given.
Each found codelab is then re-exported using parameters from the metadata file.
Unused codelab assets will be deleted, as well as the entire codelab directory,
if codelab ID has changed since last update or export.
In the latter case, where codelab ID has changed, the new directory
will be placed alongside the old one. In other words, it will have the same ancestor
as the old one.
While -prefix and -ga can override existing codelab metadata, the other
arguments have no effect during update.
The program does not follow symbolic links and exits with non-zero code
if no metadata found or at least one src could not be updated.
## Web command
Web provides a standalone web interface for converting Google Docs to markdown format.
It takes no arguments and starts a web server on the specified port.
The web command takes a -addr host:port option, to specify the
desired hostname or IP address and port number to bind to.
## Flags
`
