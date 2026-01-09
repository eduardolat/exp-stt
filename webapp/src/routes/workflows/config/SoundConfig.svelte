<script lang="ts">
	import { Volume2 } from '@lucide/svelte';

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
			<Volume2 class="size-4" />
			Sound Type
		</span>
	</label>
	<select
		class="select-bordered select select-sm"
		value={(config.soundType as string) ?? 'success'}
		onchange={(e) => onChange('soundType', e.currentTarget.value)}
		{disabled}
	>
		<option value="success">Success</option>
		<option value="error">Error</option>
		<option value="record">Record Start</option>
	</select>
	<label class="label">
		<span class="label-text-alt opacity-70">Uses sounds configured in Settings</span>
	</label>
</div>
