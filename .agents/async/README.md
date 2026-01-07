# Asynchronous LLM agents

This directory contains the definitions and state for **autonomous AI agents** that run on a schedule (e.g., via tools like Jules, Codex, etc).

Unlike standard LLM chat context, these agents are designed to operate independently, in the background, with no human interaction, perform specific maintenance or logic tasks, and persist their memory directly in the repository.

## Directory structure

We use a **"Folder-per-Agent"** architecture to keep logic and state encapsulated.

```text
.agents/async/
└── {{agent_name}}/
    ├── prompt.md       # The Source Code (Instructions for the LLM)
    ├── journal.md      # The Memory (Agent state, mutable by the agent)
```

### File definitions

- **`prompt.md`**:
  Contains the **Identity, Goal, and Instructions** for the agent. This is the "brain" of the agent. It defines _what_ the agent should do and _how_ it should do it.
- **`journal.md`**:
  Contains the **Persistent Memory** of the agent. The agent reads this file to understand what it did in previous runs and appends new findings or logs at the end. This allows the agent to have continuity across different execution sessions. The instructions of how to use this file should be included in the `prompt.md` file.

## How to configure the agents

To run these agents, you don't need to copy the prompt into your agent runner tool. Instead, we treat the runner as a generic execution runtime.

**Configuration steps:**

1. Create a scheduled task in your agent runner tool (e.g., Jules, Codex, etc).
2. Set the target repository to this repo.
3. Use the following **Bootstrap Prompt** as the system instruction.

### Bootstrap Prompt Template

Copy the text below. You only need to change the `AGENT_NAME` variable in the configuration block.

```text
## RUNTIME CONFIGURATION

AGENT_NAME:        "{{AGENT_NAME}}"
IDENTITY_FILE:     ".agents/async/{{AGENT_NAME}}/prompt.md"
MEMORY_FILE:       ".agents/async/{{AGENT_NAME}}/journal.md"
GLOBAL_RULES_FILE: "AGENTS.md"

## SYSTEM INSTRUCTIONS

You are an autonomous agent bound to this repository. Your behavior is defined strictly by the files specified in the CONFIGURATION block above.

Strictly follow this initialization and execution sequence:

1. **LOAD IDENTITY:** Read the file specified in `IDENTITY_FILE`. This defines your core personality, specific tasks, and the MANDATORY protocols for managing your memory.
2. **LOAD GLOBAL CONTEXT:** Read the file specified in `GLOBAL_RULES_FILE`. This contains general project rules and operational guidelines you must adhere to.
3. **LOAD MEMORY:** Read the file specified in `MEMORY_FILE` (if it exists). This is your long-term memory from previous runs.
4. **EXECUTE:** Rigorously execute the instructions found in the `IDENTITY_FILE`. **STRICT ADHERENCE** to those instructions is required, especially regarding the logic for reading, writing, and updating your `MEMORY_FILE`. You must not deviate from the specified journal management protocol.

Start working now.
```
