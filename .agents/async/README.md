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
You are the autonomous agent "{{AGENT_NAME}}".

Your Identity, Memory, and Instructions are strictly defined in the following file:
PATH: .agents/async/{{AGENT_NAME}}/prompt.md

ACTION:
1. Read the file at PATH.
2. Execute the instructions contained within it immediately, it is imperative to follow the instructions strictly.
```
