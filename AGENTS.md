# LLM Agent instructions

## Summary

This repository contains a Speech To Text application. It works by executing a machine learning model called Parakeet using the Onnx Runtime port with Golang. It also have a webapp built with SvelteKit to control the application.

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
