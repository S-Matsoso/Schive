# Schive

A cross-platform desktop application for university students to organise their academic life — modules, files, timetables, assignments, flashcards, and backups — all in one place.

Built with [Wails](https://wails.io/) (Go + Vue 3).

![License](https://img.shields.io/badge/license-MIT-green)
![Go](https://img.shields.io/badge/Go-1.26-00ADD8?logo=go&logoColor=white)
![Vue](https://img.shields.io/badge/Vue-3-4FC08D?logo=vuedotjs&logoColor=white)
![Platform](https://img.shields.io/badge/platform-Linux%20%7C%20Windows-informational)

---

## ✨ Features

### 📚 Module Management
- Create, edit, and delete academic modules (course code, name, colour)
- Module code and name are required (defensive validation)
- Each module gets its own isolated workspace for files, folders, assignments, and flashcards

### 📁 File & Folder Organisation
- Upload files (documents, PDFs, images, spreadsheets, etc.) into any module
- Create folders inside modules to organise content (e.g. Notes, Assignments, Past Papers)
- Drag-and-drop file uploads
- Move files between folders
- Rename files in-place
- In-app file preview for text-based files
- Grid and list view toggle
- Last-modified date shown on editable document types (`.doc`, `.docx`, `.odt`, etc.)

### 📅 Timetable
- Weekly class schedule with a visual calendar (powered by [vue-cal](https://github.com/antoniandre/vue-cal))
- Add, edit, and delete class slots (day, start hour, duration, room)
- Assignment deadlines are automatically overlaid on the calendar
- Export your timetable as a PDF via the browser print dialog

### ✅ Tasks & Assignments
- Create assignments with title, description, priority (low / medium / high), status, and due date
- Track assignments per module or view all from the dashboard
- Filter by status: All, Pending, In Progress, Completed
- Mark assignments as done directly from the dashboard
- Upcoming deadlines section on the home screen with priority colour coding

### 🃏 Flashcards
- Create question-and-answer flashcards per module
- Interactive study mode: flip cards, mark as mastered
- Track study progress

### 🏠 Dashboard
- At-a-glance stats: total modules, today's classes, pending tasks, total files
- Today's class schedule
- Upcoming deadlines with quick "Mark Done" action
- Module quick-access grid
- Recent files list with one-click open

### 🔍 Global Search
- Spotlight-style search (`Ctrl+K`) across modules, files, assignments, and timetable slots

### 💾 Backup
- One-click backup that creates a `.zip` archive of your database and uploaded files

---

## 🏗️ Tech Stack

| Layer       | Technology                                                   |
| ----------- | ------------------------------------------------------------ |
| Framework   | [Wails v2](https://wails.io/)                                |
| Backend     | Go 1.26, [GORM](https://gorm.io/), SQLite                   |
| Frontend    | Vue 3 (Composition API), [Vuetify 4](https://vuetifyjs.com/) |
| Build tool  | [Vite](https://vitejs.dev/)                                  |
| Calendar    | [vue-cal](https://github.com/antoniandre/vue-cal)            |
| Icons       | [Material Design Icons](https://materialdesignicons.com/)    |

---

## 📂 Project Structure

```
Schive/
├── main.go                  # Wails entry point
├── app.go                   # App struct, service/controller wiring
├── wails.json               # Wails project config
├── go.mod / go.sum          # Go dependencies
│
├── models/                  # GORM data models
│   ├── Module.go
│   ├── File.go
│   ├── Folder.go
│   ├── ClassSlot.go
│   ├── Assignment.go
│   └── Flashcard.go
│
├── store/
│   └── db.go                # Database initialisation & migrations
│
├── services/                # Business logic layer
│   ├── ModuleService.go
│   ├── FileService.go
│   ├── TimeTableService.go
│   ├── AssignmentService.go
│   ├── FlashcardService.go
│   ├── BackupService.go
│   └── *_test.go            # Unit tests
│
├── controllers/             # Wails-bound API controllers
│   ├── ModuleController.go
│   ├── FileController.go
│   ├── TimeTableController.go
│   ├── AssignmentController.go
│   ├── FlashcardController.go
│   └── BackupController.go
│
├── frontend/
│   ├── src/
│   │   ├── main.js          # Vue + Vuetify setup
│   │   ├── style.css        # Global styles & CSS variables (light/dark)
│   │   ├── views/
│   │   │   ├── Home.vue           # Dashboard
│   │   │   ├── ModuleDetails.vue  # Module workspace (files, assignments, flashcards)
│   │   │   └── Timetable.vue      # Weekly timetable & calendar
│   │   ├── components/
│   │   │   └── SideBar.vue        # Navigation, search, theme toggle, backup
│   │   └── router/
│   │       └── index.js           # Vue Router config
│   └── package.json
│
├── build/
│   └── appicon.png          # Application icon
│
└── .github/
    └── workflows/
        └── release.yml      # CI: build & publish GitHub releases (Linux + Windows)
```

---

## 🚀 Getting Started

### Prerequisites

- **Go** ≥ 1.26
- **Node.js** ≥ 18
- **Wails CLI** v2 — install with:
  ```bash
  go install github.com/wailsapp/wails/v2/cmd/wails@latest
  ```
- **Linux only:** WebKit2GTK development libraries
  ```bash
  # Ubuntu / Debian / Mint
  sudo apt install libgtk-3-dev libwebkit2gtk-4.1-dev

  # Fedora
  sudo dnf install gtk3-devel webkit2gtk4.1-devel
  ```

### Clone & Run

```bash
git clone https://github.com/<your-username>/schive.git
cd schive
```

#### Development mode (hot reload)

```bash
wails dev
```

This starts the Go backend and a Vite dev server with hot module replacement for the frontend.

#### Production build

```bash
# Linux
wails build -tags "webkit2_41"

# Windows
wails build
```

The compiled binary is output to `build/bin/`.

### Running tests

```bash
go test -v ./...
```

---

## 🔄 CI / CD

The project includes a GitHub Actions workflow (`.github/workflows/release.yml`) that:

1. Builds the application for **Linux (amd64)** and **Windows (amd64)** on every push to `main`
2. Packages the binaries into `.tar.gz` (Linux) and `.zip` (Windows)
3. Publishes a GitHub Release with the built artifacts and auto-generated release notes

---

## 🤝 Contributing

Contributions are welcome! Here's how to get started:

1. **Fork** the repository
2. **Create a feature branch** — `git checkout -b feat/my-feature`
3. **Commit your changes** — `git commit -m "feat: add my feature"`
4. **Push** — `git push origin feat/my-feature`
5. **Open a Pull Request**

Please make sure:
- All existing tests pass (`go test -v ./...`)
- The frontend builds without errors (`npm --prefix frontend run build`)
- New features include appropriate tests where possible

---

## 📄 License

This project is licensed under the **MIT License** — see the [LICENSE](LICENSE) file for details.

---

## 👤 Author

**Seboko Matsoso** — [smatsoso246@gmail.com](mailto:smatsoso246@gmail.com)
