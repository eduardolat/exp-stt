<script lang="ts">
	import { MessageSquare } from '@lucide/svelte';

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
				<MessageSquare class="size-4" />
				Title (optional)
			</span>
		</label>
		<input
			type="text"
			class="input-bordered input input-sm"
			placeholder="Notification Title"
			value={(config.title as string) ?? ''}
			onchange={(e) => onChange('title', e.currentTarget.value)}
			{disabled}
		/>
	</div>

	<div class="form-control">
		<label class="label" for="notify-message">
			<span class="label-text">Message</span>
		</label>
		<textarea
			id="notify-message"
			class="textarea-bordered textarea h-20"
			placeholder="Enter notification message"
			value={(config.message as string) ?? ''}
			onchange={(e) => onChange('message', e.currentTarget.value)}
			{disabled}
		></textarea>
		<label class="label">
			<span class="label-text-alt opacity-70">Tip: Use variables like {'{{nodes.id.field}}'}</span>
		</label>
	</div>
</div>
