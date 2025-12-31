# LLM Agent instructions

## Summary

This repository contains a Speech To Text application named Tribar Voice. It works by executing a machine learning model called Parakeet using the Onnx Runtime port with Golang. It also have a webapp built with SvelteKit to control the application.

## General instructions

- Always use Context7 MCP (context7) when you need setup or configuration steps, or
  library/API documentation. This means you should automatically use the Context7 MCP
  tools to resolve library id and get library docs without me having to explicitly ask. (With the exception of Svelte and Svelte Kit, see below)

- Whenever you work with Svelte, use Svelte’s MCP (svelte) to validate that what you’re writing is correct.

- Before starting any task, you must run the tree command to view the file structure we are working with in the most up-to-date way possible. Make sure to use that tree command, but don’t read things like node_modules or things like that so you don’t read unnecessary content.

- When you write code, make sure it is production-quality, readable, and maintainable. And it is well documented. Everything in an idiomatic way.

- Regarding comments in the code, I want you not to overuse them, not to add comments that add no value to the project. They should be solely for documentation purposes and not redundant, since the code must be good enough to be self-explanatory.

- Whenever you install any dependency using pnpm, make sure to install it at its exact version (`pnpm add <package> --save-exact`). That allows us to have reproducibility in the future.

- Before finishing any task, double-check the code you created. That way you can notice any errors you might have missed during its construction. Also check the editor's errors to see Type Safety issues and other problems.

- When you write conditionals, avoid using `else` block, this ensures the code is readable and prevents nesting-hell, this applies for all programming languages, even Svelte templates. Instead of `else` use inverted conditions in `if` blocks.

- Use `@lucide/svelte` icons whenever you need icons and you are working on the webapp.

- Before you finish any task, go to the root of the project and run the `task ci` command. That way all the code you wrote will be type-checked, tested and linted to catch errors as soon as possible. Fix any type errors you find before finishing the task. It includes checks for both the Go application and the Svelte webapp.

## Shared architecture

The UFO RPC schema and config file is located here:

- schema/schema.urpc
- schema/uforpc.toml

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

#### History

Source: `internal/history`

The `history` package implements a minimalist flat-file database system for transcriptions. It uses UUIDv7 for K-sortable unique identifiers and stores each entry as a pair of files (`.json` for metadata and `.wav` for audio). It handles atomic writes, automatic pruning based on user settings, and provides a thread-safe in-memory cache that is populated asynchronously at startup.

#### App State

Source: `internal/state`

Manages the global application state (status) in a thread-safe way, providing access to other packages. It also acts as a bridge to the `history` package, allowing other components to query and manage past transcriptions.

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

#### Server

Source: `internal/server`

Provides an HTTP API and a SvelteKit web UI to control and monitor the application. The Web UI includes a configuration manager to update user preferences, which are persisted via the `config` package.

It receives the state to react to changes (read-only) and the Engine to perform actions, as all interactions must be handled by the orchestrator (engine).

The server is built with the `github.com/labstack/echo/v4` framework.

The HTTP API is built with UFO RPC (a framework to build HTTP based RPC APIs), if you need documentation of UFO RPC you can find it here:

https://raw.githubusercontent.com/uforg/uforpc/refs/heads/main/docs/src/content/docs/reference/about.md
https://raw.githubusercontent.com/uforg/uforpc/refs/heads/main/docs/src/content/docs/reference/request-lifecycle.md
https://raw.githubusercontent.com/uforg/uforpc/refs/heads/main/docs/src/content/docs/reference/urpc-spec.md

Please when you create a handler for a procedure or stream, create a new file inside `internal/server/api` exclusively for that handler, use the `proc_{snake_case_name}.go` and `stream_{snake_case_name}.go` file names, read some of the existing files to understand the structure. After creating the handler please register it in the `internal/server/api/router.go` file, read this router and two or three other handlers to understand the structure.

#### Systray

Source: `internal/systray`

A system tray interface that displays app status and provides quick controls.

It receives the state to react to changes (read-only) and the Engine to perform actions, as all interactions must be handled by the orchestrator (engine).

It also has a button to open the Web UI in the user's default browser.

## Web UI

The Web UI source code is inside the `webapp/` directory, it is built with SvelteKit. It uses PNPM as the package manager.

### Svelte MCP

You are able to use the Svelte MCP server, where you have access to comprehensive Svelte 5 and SvelteKit documentation. Here's how to use the available tools effectively:

#### Available MCP Tools:

##### 1. list-sections

Use this FIRST to discover all available documentation sections. Returns a structured list with titles, use_cases, and paths.
When asked about Svelte or SvelteKit topics, ALWAYS use this tool at the start of the chat to find relevant sections.

##### 2. get-documentation

Retrieves full documentation content for specific sections. Accepts single or multiple sections.
After calling the list-sections tool, you MUST analyze the returned documentation sections (especially the use_cases field) and then use the get-documentation tool to fetch ALL documentation sections that are relevant for the user's task.

##### 3. svelte-autofixer

Analyzes Svelte code and returns issues and suggestions.
You MUST use this tool whenever writing Svelte code before sending it to the user. Keep calling it until no issues or suggestions are returned.

##### 4. playground-link

Generates a Svelte Playground link with the provided code.
After completing the code, ask the user if they want a playground link. Only call this tool after user confirmation and NEVER if code was written to files in their project.

## Style

This project uses Tailwind CSS v4 for styling. But it uses Basecoat CSS (`basecoat-css`) which is a framework on top of Tailwind CSS, it provides a set of pre-styled class based components similar to Shad CN but without the need to use React and its more convenient and framework agnostic, please use it instead of Tailwind CSS classes when possible. Only use Tailwind CSS classes when you really need to, otherwise use the Basecoat CSS classes. You can find more documentation about Basecoat CSS here:

- https://basecoatui.com/
- https://github.com/hunvreus/basecoat

Search for the previous documentation in Context7 MCP to find documentation about of the available components and how to use them.

Regarding the style, it should be desktop only (no responsive design) because the software is intended to be used only on desktops so we can simplify the design and make it desktop only. It can have a maximum width of 1024px which is configured in the layout and it should be visually comfortable to use and horizontally centered.

Make the style minimalistic and simple, without any unnecessary elements or distractions but without losing the functionality or user experience which is the absolute priority.

## Icons

Use `@lucide/svelte` icons whenever you need icons and you are working on the webapp. For example if you need `Phone` icon from lucide you can use this:

```svelte
<script>
	import { Phone } from '@lucide/svelte';
</script>

<!-- Simple -->
<Phone />
<!-- Or with custom classes -->
<Phone class="size-10 text-red-500" />
```
