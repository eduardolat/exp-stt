<script lang="ts">
	import { Terminal } from '@lucide/svelte';

	interface Props {
		config: Record<string, unknown>;
		onChange: (key: string, value: unknown) => void;
		disabled?: boolean;
	}

	let { config, onChange, disabled = false }: Props = $props();
</script>

<div class="flex flex-col gap-3">
	<div class="form-control">
		<label class="label">
			<span class="label-text flex items-center gap-2">
				<Terminal class="size-4" />
				Command
			</span>
		</label>
		<input
			type="text"
			class="input-bordered input input-sm font-mono"
			placeholder="echo 'Hello World'"
			value={(config.command as string) ?? ''}
			onchange={(e) => onChange('command', e.currentTarget.value)}
			{disabled}
		/>
	</div>

	<div class="form-control">
		<label class="label" for="terminal-args">
			<span class="label-text">Arguments (optional)</span>
		</label>
		<input
			id="terminal-args"
			type="text"
			class="input-bordered input input-sm font-mono"
			placeholder="--flag value"
			value={(config.args as string) ?? ''}
			onchange={(e) => onChange('args', e.currentTarget.value)}
			{disabled}
		/>
	</div>

	<div class="rounded-lg border border-error/50 bg-error/10 p-3 text-xs">
		<p class="font-semibold text-error">⚠️ Security Warning</p>
		<p class="mt-1 text-error/80">
			Commands are executed directly on your system with your user permissions. Only use this node
			if you understand the risks. If you imported this workflow, verify the command carefully
			before running.
		</p>
	</div>
</div>
