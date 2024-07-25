document.getElementById('convert-form').addEventListener('submit', async function(event) {
    event.preventDefault();
    const docUrl = document.getElementById('doc-url').value;
    const response = await fetch('/convert', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: `docURL=${encodeURIComponent(docUrl)}`,
    });

    if (response.ok) {
        const blob = await response.blob();
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.style.display = 'none';
        a.href = url;
        a.download = 'converted.md';
        document.body.appendChild(a);
        a.click();
        window.URL.revokeObjectURL(url);
    } else {
        alert('Failed to convert Google Docs to markdown.');
    }
});
