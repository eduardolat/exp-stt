<script lang="ts">
	import { Sparkles } from '@lucide/svelte';

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
				<Sparkles class="size-4" />
				Custom Prompt (optional)
			</span>
		</label>
		<textarea
			class="textarea-bordered textarea h-32"
			placeholder="Leave empty to use the prompt from Settings. Or enter a custom prompt here."
			value={(config.prompt as string) ?? ''}
			onchange={(e) => onChange('prompt', e.currentTarget.value)}
			{disabled}
		></textarea>
		<label class="label">
			<span class="label-text-alt opacity-70">Uses AI settings from Settings → AI</span>
		</label>
	</div>

	<div class="rounded-lg bg-info/10 p-3 text-xs text-info">
		💡 The input text will be sent to your configured LLM with this prompt.
	</div>
</div>
