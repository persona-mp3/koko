# Koko — Simple HTTP request helper

Purpose
- `koko` is a small command-line tool for making quick HTTP GET and POST requests against local development servers. It reads a `spec.json` file for POST request bodies and uses local pagers (`jq`, `bat`, or `less`) to render server responses.

How it works
- The program accepts flags to choose GET vs POST, the endpoint, content type, and whether to follow redirects. For POSTs it reads the request body from `spec.json` in the current directory.

Requirements
- Go (project uses `go.mod` with `go 1.25.6`)
- Optional tools (recommended for best output experience): `jq`, `bat`, and `less` (homebrew paths are used in the code: `/opt/homebrew/bin/jq`, `/opt/homebrew/bin/bat`, `/usr/bin/less`).

Files
- `main.go`: CLI flags, request orchestration, and output handling.
- `post.go`: POST request helpers and `spec.json` reader.
- `pager.go`: Chooses and runs a pager based on response content-type.
- `spec.json`: JSON file used as POST request body (example provided in repository).

Usage
- Run the program directly with `go run .` or use the Makefile target:

```
make run
```

Flags
- `-post` : use POST instead of GET.
- `-ep`   : endpoint to call (default `http://localhost:3000`).
- `-ct`   : content type for POST body (default `application/json`).
- `-redirect` : follow redirects (by default the tool does not follow redirects when using the NoRedirect mode).

Examples
- GET (default):

```
go run .
```

- GET against a custom endpoint:

```
go run . -ep http://localhost:8080/health
```

- POST using `spec.json` as the request body (default content type `application/json`):

```
go run . -post
```

- POST and follow redirects:

```
go run . -post -redirect
```

Using `spec.json`
- When making POST requests, the program reads `spec.json` from the current working directory and sends its contents as the request body. If `spec.json` is empty or minimal JSON (example in repo), that content will be sent as-is.

Output and paging
- Responses with a JSON content-type are piped to `jq` (configured path in code). HTML and other text responses default to `bat`. If those tools are not found, responses will fall back to `less`.

Notes and limitations
- `readSpecFile` currently returns an empty body instead of an error when `spec.json` cannot be read — be sure a `spec.json` file exists in the working directory when making POST requests.
- The code logs and exits if the remote server refuses the connection.
- Content-type detection is basic (string containment check).

Next steps
- (Optional) Adjust pager paths for your platform or install `jq`/`bat` via Homebrew for best results.

License
- No license specified.

