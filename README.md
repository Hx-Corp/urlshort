# 🔗 URLShort — Advanced URL Shortener & Parameter Generator

![Go Version](https://img.shields.io/badge/Go-1.17+-00ADD8?style=flat&logo=go)
![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)
![Platform](https://img.shields.io/badge/Platform-linux%20%7C%20macOS%20%7C%20windows-blue)

**URLShort** is a powerful and intelligent tool written in Go to split, shorten, and fuzz URLs for automation, testing, and security purposes.  
Built with ❤️ by [Neeraj Sah (Team HyperGod-X)](https://github.com/nxneeraj) to help you automate recon, parameter manipulation, and payload injection — quickly and effectively.

---

## 🙋‍♂️ Who Is It For?

- 🧑‍💻 **Web Developers** — Analyze and test URL behavior.
- 🐞 **Bug Bounty Hunters** — Generate fuzzable URL vectors fast.
- 🛡 **Security Professionals** — Inject payloads into endpoints for XSS/IDOR testing.
- 🧪 **QA Engineers** — Automate URL manipulation for test coverage.
- 🧰 **Red Teamers & 🔷 Blue Teamers / Recon Experts** — Enhance endpoint discovery and testing.

---

## 🧠 Why Use It?

- ✅ Generate recursive URL variations
- ✅ Append payloads from CLI or file
- ✅ Target both query strings and path segments
- ✅ Deduplicate, output to file, and integrate in scripts
- ✅ Ideal for recon, fuzzing, XSS testing, path traversal, etc.

---

## 🌟 Features

| 🔹 Feature                          | ✅ Status  |
|-----------------------------------|-----------|
| Split URLs by `=`, `&`, `/`, etc. | ✔️ |
| Recursive URL prefix generation   | ✔️ |
| Append payloads (`-a` or `-F`)    | ✔️ |
| Remove duplicates with `-D`       | ✔️ |
| Input/output file support         | ✔️ |
| Quiet mode for automation         | ✔️ |
| Cross-platform support            | ✔️ |
| Beautiful banner & color output   | ✔️ |

---

## ⚡️ Installation

### ✅ Option 1: Global Install via Go (Recommended)

```bash
go install github.com/Hx-Corp/urlshort@latest
```

> ⚠️ Ensure your `$GOPATH/bin` is in your system `PATH`:

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

Then use it globally:

```bash
urlshort -f urls.txt [options]
```

---

### ⚙️ Option 2: Manual Global Setup (via `setup.go`)

```bash
go build -o urlshort
go run setup.go
```

| OS         | Installs To                          |
|------------|--------------------------------------|
| Linux      | `/usr/local/bin/hxscanner`           |
| macOS      | `/usr/local/bin/hxscanner`           |
| Termux     | `$HOME/.termux/bin/hxscanner`        |
| Windows    | Start Menu as `hxscanner.exe`        |

🛡️ Requires **root/sudo** privileges on Unix/macOS.

---

## 🚀 Usage

```bash
urlshort -f <input-file> [options]
```

### 🔧 Command-line Flags

| Flag | Description |
|------|-------------|
| `-f` | Input file with URLs (**required**) |
| `-o` | Output file to save results |
| `-x` | Delimiters for splitting (e.g., `"=&"`) (default: `=`) |
| `-p` | Enable splitting on path segments (`/`) |
| `-a` | Append a string to each URL variation |
| `-F` | File of strings to append (overrides `-a`) |
| `-D` | Remove duplicate URLs |
| `-Q` | Quiet mode (suppress output, show only final messages) |
| `-h` | Display help message |

---

## 🧪 Example Use Case

```bash
urlshort -f urls.txt -o out.txt -x "=&" -p -F payloads.txt -D
```

- Splits URLs by `=`, `&`, `/`
- Appends payloads from `payloads.txt`
- Deduplicates results
- Writes final output to `out.txt`

---

## 📁 Sample Files

### `urls.txt`
```
https://example.com/page?user=admin&pass=1234
https://site.com/login.php?auth=token
```

### `payloads.txt`
```
<script>alert(1)</script>
%27-- OR 1=1
../
```

---

## 📤 Sample Output

```
https://example.com/
https://example.com/page?
https://example.com/page?user=
https://example.com/page?user=admin&
https://example.com/page?user=admin&pass=
https://example.com/page?user=admin&pass=1234<script>alert(1)</script>
https://site.com/
https://site.com/login.php?
https://site.com/login.php?auth=
https://site.com/login.php?auth=token%27-- OR 1=1
...
```

---

## 🧠 How It Works (Internals)

1. **Input Reading**  
   - Loads all non-empty lines from the file specified by `-f`.

2. **Splitting Logic**  
   - Uses delimiters from `-x` to split URL query strings.
   - Adds `/` splitting if `-p` is enabled.

3. **Variation Generation**  
   - Recursively builds prefix variations of URLs ending at each delimiter.

4. **Appending**  
   - Combines each variation with payloads using `-a` or `-F`.

5. **Deduplication**  
   - Optional: `-D` removes repeated entries using a map.

6. **Output**  
   - Writes to file with `-o` or prints to console (unless `-Q` is used).

---

## 🛠 Development

### Clone & Build

```bash
git clone https://github.com/nxneeraj/urlshort
cd urlshort
go build -o urlshort
./urlshort -h
```

---

## 📂 Project Structure

```bash
urlshort/
├── main.go       # Main logic & CLI interface
├── setup.go          # Optional global installer script
├── README.md         # Full documentation
```

---

## 📜 License

```text
MIT License

Copyright (c) 2025 Neeraj

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction...
```

---

## 🤝 Contributing

Contributions are warmly welcome! Here’s how you can help:

- Report bugs or edge cases via [GitHub Issues](https://github.com/nxneeraj/urlshort/issues)
- Suggest new features or CLI improvements
- Submit PRs to enhance logic, UX, or compatibility
- Share use cases, templates, and examples

---

## 📬 Contact

- 👤 **Author**: [Neeraj Sah](https://github.com/nxneeraj)
- 📧 **Email**: neerajsahnx@gmail.com
- 🏴‍☠️ **Org**: [Team HyperGod-X](https://github.com/hypergodx)

---

## ⭐ Support This Project

If this tool helped you:

- ⭐ Star this repo
- 🚀 Share it with fellow hackers
- 🧠 Mention it in your toolkit, blog, or course
- 🔁 Fork and make it even better!

> Build faster. Test smarter. Hack ethically.  
> With 💥 from Team HyperGod-X 👾
<p align="center"><strong> Keep Moving Forward </strong></p>
