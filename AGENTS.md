# LLM Agent instructions

## Summary

This repository contains a Speech To Text application named Tribar Voice. It works by executing a machine learning model called Parakeet using the Onnx Runtime port with Golang. It also have a webapp built with SvelteKit to control the application.

## General instructions

- Always use Context7 MCP (context7) when you need setup or configuration steps, or
  library/API documentation. This means you should automatically use the Context7 MCP
  tools to resolve library id and get library docs without me having to explicitly ask. (With the exception of Svelte and Svelte Kit, see below)

- Whenever you work with Svelte, use Svelte’s MCP (svelte-mcp) to validate that what you’re writing is correct.

- Before starting any task, you must run the tree command to view the file structure we are working with in the most up-to-date way possible. Make sure to use that tree command, but don’t read things like node_modules or things like that so you don’t read unnecessary content.

- When you write code, make sure it is production-quality, readable, and maintainable. And it is well documented. Everything in an idiomatic way.

- Regarding comments in the code, I want you not to overuse them, not to add comments that add no value to the project. They should be solely for documentation purposes and not redundant, since the code must be good enough to be self-explanatory.

- Whenever you install any dependency using pnpm, make sure to install it at its exact version (`pnpm add <package> --save-exact`). That allows us to have reproducibility in the future.

- Before finishing any task, double-check the code you created. That way you can notice any errors you might have missed during its construction. Also check the editor's errors to see Type Safety issues and other problems.

- When you write conditionals, avoid using `else` block, this ensures the code is readable and prevents nesting-hell, this applies for all programming languages, even Svelte templates. Instead of `else` use inverted conditions in `if` blocks.

- Use `@lucide/svelte` icons whenever you need icons and you are working on the webapp.

- Before you finish any task, go to the root of the project and run the `task ci` command. That way all the code you wrote will be type-checked, tested and linted to catch errors as soon as possible. Fix any type errors you find before finishing the task.

## Backend architecture

The backend is divided into several packages inside `internal` which are orchestrated in the main process (`cmd/tribar/main.go`).

### Backend layers

The dependencies between packages follow an order and some must be created before others since they are received as parameters using dependency injection to maintain the order and testability of the project.

Below I list the main packages; this order must be respected because the first ones are injected into the following ones:

#### Main package

Source: `cmd/tribar/main.go`

The `main` package is responsible for bootstrapping the entire project, creating and injecting dependencies, and starting all the necessary project services.

Here is the graceful shutdown handled via context cancelation and signal listeners.

#### Logger

Source: `internal/logger`

The `logger` package is a utility for printing important data to STDOUT in a structured way; it must be created right after starting the program to allow capturing logs of absolutely everything else in the program.

#### Config

Source: `internal/config`

The `config` package contains global and general program settings such as name, version, etc. It ensures the existence of all required directories and manages a JSON configuration file that persists user preferences (notifications, sounds, AI settings, history limits), which can be updated via the Web UI.

#### Onnx Runtime

Source: `internal/onnx`

The `onnx` package, like the `config` package, is vital to the program, and if it fails, the program cannot continue. The function of this package is to place the shared libraries of the onnx runtime within the program's directories so that subsequent packages can use the onnx runtime without problems. These shared libraries are embedded in the program using `go embed` and extracted into its directory using this package.

#### App State

Source: `internal/state`

Manages the global application state (status, settings) in a thread-safe way, providing access to other packages. It also handles a configurable history of transcriptions and their corresponding audio files.

#### Recorder

Source: `internal/record`

Handles audio recording from the system's input device and saves the output as WAV files in the designated directory for further processing.

#### Transcriber

Source: `internal/transcribe`

Converts audio files into text using the Parakeet model via ONNX Runtime, handling the inference process and returning the raw transcription.

#### Post-processor

Source: `internal/postprocess`

Refines and enhances transcriptions using LLM-based AI processing to improve grammar, punctuation, and overall readability. It is disabled by default and supports OpenAI-compatible APIs with a prompt manager for predefined or custom enhancements.

#### Notify

Source: `internal/notify`

Sends desktop notifications to inform the user about important application events. By default, it only alerts on errors, but users can enable notifications for transcription start and completion.

#### Clipboard

Source: `internal/clipboard`

Responsible for outputting the final transcription. Supports three modes: `copy_only` (copies text to clipboard), `copy_paste` (copies and triggers paste), and `ghost_paste` (pastes without modifying clipboard by temporarily storing existing content).

#### Sound

Source: `internal/sound`

Plays audio cues to provide acoustic feedback for application events, helping the user know the app's status without looking at the screen. Cues for starting and finishing transcriptions are enabled by default but can be disabled by the user.

#### Engine

Source: `internal/engine`

The central orchestrator that connects all components and manages the workflow using dependency injection.

It is the only package allowed to modify the application state and receives all other functional packages as dependencies (except visualization layers like Systray or Server).

#### Systray

Source: `internal/systray`

A system tray interface that displays app status and provides quick controls.

It receives the state to react to changes (read-only) and the Engine to perform actions, as all interactions must be handled by the orchestrator (engine).

#### Server

Source: `internal/server`

Provides an HTTP API and a SvelteKit web UI to control and monitor the application. The Web UI includes a configuration manager to update user preferences, which are persisted via the `config` package.

It receives the state to react to changes (read-only) and the Engine to perform actions, as all interactions must be handled by the orchestrator (engine).
