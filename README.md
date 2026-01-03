<div align="center">
  <img src="assets/logo/png/black-white-logo-128.png" alt="Tribar Voice Logo" width="128" />
  <h1>Tribar Voice</h1>
</div>

Tribar Voice is a desktop application designed to make speech-to-text simple, private, and powerful. It listens to your voice and converts it into text in real-time, running entirely on your computer without sending audio to the cloud.

Whether you are drafting emails, writing code, or just prefer speaking over typing, Tribar Voice integrates seamlessly into your workflow with smart clipboard features and optional AI tools to polish your text.

## Features

- **100% Offline & Private**: Powered by the Parakeet model running locally via ONNX Runtime. Your voice data stays on your machine.
- **AI Enhancement**: Connect to OpenAI-compatible APIs to automatically fix grammar, punctuation, and style. You can connect it to local LLMs for a 100% private and offline experience, or disable it completely.
- **Smart Clipboard Modes**:
  - **Copy Only**: Just copies the transcription to your clipboard.
  - **Copy & Paste**: Copies and immediately pastes the text into your active application.
  - **Ghost Paste**: Pastes the text without overwriting what you currently have in your clipboard.
- **Modern Web Dashboard**: Configure everything from a clean, user-friendly interface. Manage your history, tweak AI settings, and customize notifications.
- **History & Playback**: Every transcription is saved with its audio. You can search your history and replay recordings anytime.
- **System Tray Integration**: Quietly runs in the background. Access controls or open the dashboard via the system tray icon.

## Supported Operating Systems

Tribar Voice is built on cross-platform technologies (Go and Web) and is designed to work on:

- **Linux**
- **Windows**
- **macOS**

## Installation

You can download the latest pre-compiled binary for your operating system from the [Releases page](https://github.com/varavelio/tribar/releases).

If you prefer to build from source, follow the instructions below.

### Building from Source

#### Prerequisites

This project includes a `flake.nix` file that automatically installs all necessary dependencies. If you are not using Nix, you will need to install the required tools manually.

#### Steps

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

**Note for Nix users:** Just run `nix develop` to get a shell with all dependencies pre-configured, then run `task run`.

## Usage

Once running, Tribar Voice sits in your system tray. You can:

- **Open the Dashboard:** Click the tray icon and select "Open Web UI".
- **Configure Hotkeys:** Set up a global hotkey to start/stop recording in the settings.
- **Review History:** Browse past transcriptions in the dashboard.

## Feedback

We want Tribar Voice to be the best tool for your workflow. If you have ideas, find bugs, or just want to share how you use it, feel free to open an issue or contribute code.

## License

This project is 100% free and open source software, licensed under the [MIT License](LICENSE).
