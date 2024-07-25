package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/codelabs-cn/codelab-tools/claat/cmd"
)

func convertToMarkdown(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	docURL := r.FormValue("docURL")
	if docURL == "" {
		http.Error(w, "Missing Google Docs URL", http.StatusBadRequest)
		return
	}

	opts := cmd.CmdExportOptions{
		Expenv:  "web",
		Output:  "-",
		Tmplout: "md",
	}

	meta, err := cmd.ExportCodelabToMarkdown(docURL, nil, opts)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to convert Google Docs to markdown: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.md", meta.ID))
	w.Header().Set("Content-Type", "text/markdown")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(meta.Content))
}

func serveWebInterface() {
	http.HandleFunc("/convert", convertToMarkdown)
	http.Handle("/", http.FileServer(http.Dir("./static")))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Starting server on port %s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}
