# 📝 Godo

> ⚡ A fast, lightweight, file-based CLI todo manager built with Go. No sync, no cloud — just simple todos that live locally and are versionable with Git.

![Go Version](https://img.shields.io/badge/Go-1.20+-brightgreen)
![Platform](https://img.shields.io/badge/platform-Windows%20%7C%20Linux%20%7C%20macOS-blue)
![License](https://img.shields.io/badge/license-MIT-green)

---

## ✨ Features

- ⚡ **Fast and Lightweight**: Zero dependencies. Blazing fast in terminal.
- 💾 **File-based**: Your todos are stored as simple local files (`.godo`)
- 🛠️ **Git-friendly**: Use Git to track changes to your todos.
- 📂 **Multiple lists**: Manage multiple todo files (projects or categories).
- 📦 **Cross-platform**: Works on Windows, Linux, and macOS.

---

## 📦 Installation

### Using `go install`

```bash
go install github.com/yourusername/godo@latest
````

Make sure your `$GOPATH/bin` is added to your `$PATH`.

---

## 🚀 Getting Started

```bash
godo help         # Show usage
godo new work     # Create a new todo file named 'work'
godo add work     # Add a new todo to 'work'
godo list work    # List all todos in 'work'
godo done work    # Mark a todo as completed
godo edit work    # Edit a todo entry
godo del work     # Delete a todo
```

---

## 🛠️ Commands Manual

| Command            | Description                               |
| ------------------ | ----------------------------------------- |
| `godo help`        | Show help info                            |
| `godo new <file>`  | Create a new todo file                    |
| `godo add <file>`  | Add a new todo item                       |
| `godo list <file>` | List all todo items                       |
| `godo done <file>` | Mark a todo as completed                  |
| `godo del <file>`  | Delete a todo item                        |
| `godo edit <file>` | Edit a todo entry                         |
| `godo files`       | Show all todo files                       |
| `godo open <file>` | Open a todo file in editor (if supported) |

Each todo file is stored as a plain `.godo` file under the `~/.godo/` directory (or the working directory depending on config).

---

## 📂 Example

```bash
$ godo new personal
Created todo file: personal.godo

$ godo add personal
> What do you want to add?
Buy groceries

$ godo list personal
1. [ ] Buy groceries

$ godo done personal
> Which todo ID to mark done?
1

$ godo list personal
1. [x] Buy groceries
```

---

## 🧩 Directory Structure

By default, your todo files are saved in:

```bash
~/.godo/
```

Each file is stored as:

```text
<filename>.godo
```

---

## 🧠 Why Godo?

* No cloud lock-in
* No internet dependency
* Great for developers
* Can be version-controlled
* Minimalist and efficient

---

## 🤝 Contributing

1. Fork the repo
2. Create your feature branch: `git checkout -b feat/new-command`
3. Commit your changes: `git commit -m "Add new feature"`
4. Push to the branch: `git push origin feat/new-command`
5. Open a Pull Request

---

## 📄 License

MIT © \[Your Name]

---

## 🔗 Related Projects

* [taskwarrior](https://taskwarrior.org/)
* [todo.txt](https://github.com/todotxt/todo.txt-cli)
* [yazi](https://github.com/sxyazi/yazi) – for file managers (naming inspo)

---
