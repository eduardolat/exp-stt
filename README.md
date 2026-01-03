# Tribar Voice

Tribar Voice is a desktop application designed to make speech-to-text simple, private, and powerful. It listens to your voice and converts it into text in real-time, running entirely on your computer without sending audio to the cloud.

Whether you are drafting emails, writing code, or just prefer speaking over typing, Tribar Voice integrates seamlessly into your workflow with smart clipboard features and optional AI tools to polish your text.

## Features

- **100% Offline & Private**: Powered by the Parakeet model running locally via ONNX Runtime. Your voice data stays on your machine.
- **AI Enhancement**: Connect to OpenAI-compatible APIs to automatically fix grammar, punctuation, and style. This is optional and fully configurable.
- **Smart Clipboard Modes**:
  - **Copy Only**: Just copies the transcription to your clipboard.
  - **Copy & Paste**: Copies and immediately pastes the text into your active application.
  - **Ghost Paste**: Pastes the text without overwriting what you currently have in your clipboard.
- **Modern Web Dashboard**: Configure everything from a clean, user-friendly interface. Manage your history, tweak AI settings, and customize notifications.
- **History & Playback**: Every transcription is saved with its audio. You can search your history and replay recordings anytime.
- **System Tray Integration**: Quietly runs in the background. Access controls or open the dashboard via the system tray icon.

## Supported Operating Systems

Tribar Voice is built on cross-platform technologies (Go and Web) and is designed to work on:

- **Linux** (First-class support, includes a `flake.nix` for NixOS users)
- **Windows**
- **macOS**

## Installation

Tribar Voice is currently distributed as source code. You'll need a few tools to build it.

### Prerequisites

- **Go** (version 1.25.4 or later)
- **PNPM** (for building the web interface)
- **Task** (a task runner/make alternative)
- **UFO RPC** (`urpc` CLI tool, required for code generation)
  - You can find the binary in the [uforpc repository](https://github.com/uforg/uforpc) releases.

### Building from Source

1.  **Clone the repository:**

    ```bash
    git clone https://github.com/varavelio/tribar.git
    cd tribar
    ```

2.  **Install dependencies:**

    ```bash
    task deps
    ```

3.  **Build the application:**

    ```bash
    task build
    ```

4.  **Run it:**
    ```bash
    ./dist/tribar
    ```

**Note for NixOS users:** Just run `nix develop` to get a shell with all dependencies pre-configured, then run `task run`.

## Usage

Once running, Tribar Voice sits in your system tray. You can:

- **Open the Dashboard:** Click the tray icon and select "Open Web UI" (or go to `http://localhost:13000` by default).
- **Configure Hotkeys:** Set up a global hotkey to start/stop recording in the settings.
- **Review History:** Browse past transcriptions in the dashboard.

## Feedback

We want Tribar Voice to be the best tool for your workflow. If you have ideas, find bugs, or just want to share how you use it, feel free to open an issue or contribute code.
