<script lang="ts">
	import { Code } from '@lucide/svelte';

	interface Props {
		config: Record<string, unknown>;
		onChange: (key: string, value: unknown) => void;
		disabled?: boolean;
	}

	let { config, onChange, disabled = false }: Props = $props();
</script>

<div class="form-control">
	<label class="label">
		<span class="label-text flex items-center gap-2">
			<Code class="size-4" />
			JavaScript Code
		</span>
	</label>
	<textarea
		class="textarea-bordered textarea h-48 font-mono text-sm"
		placeholder="// Your JavaScript code here&#10;// Use 'input' for incoming text&#10;// Return a value&#10;&#10;return input.toUpperCase();"
		value={(config.code as string) ?? ''}
		onchange={(e) => onChange('code', e.currentTarget.value)}
		{disabled}
	></textarea>
	<label class="label">
		<span class="label-text-alt opacity-70">Use 'input' for incoming text. Return a value.</span>
	</label>
</div>
