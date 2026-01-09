<script lang="ts">
	import { Clock } from '@lucide/svelte';

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
			<Clock class="size-4" />
			Delay (milliseconds)
		</span>
	</label>
	<input
		type="number"
		class="input-bordered input input-sm"
		placeholder="1000"
		min="0"
		max="300000"
		value={(config.delayMs as number) ?? 1000}
		onchange={(e) => onChange('delayMs', parseInt(e.currentTarget.value) || 0)}
		{disabled}
	/>
	<label class="label">
		<span class="label-text-alt opacity-70">Max: 5 minutes (300000ms)</span>
	</label>
</div>
