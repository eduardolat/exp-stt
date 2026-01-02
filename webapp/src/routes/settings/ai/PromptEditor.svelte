<script lang="ts">
	import { Card, Modal } from '$lib/components';
	import { Pencil, Trash2, Check } from '@lucide/svelte';
	import type { Prompt } from '$lib/client.gen';

	interface Props {
		prompt: Prompt;
		isSelected: boolean;
		onSelect: (id: string) => void;
		onUpdate: (prompt: Prompt) => void;
		onDelete: (id: string) => void;
	}

	let { prompt, isSelected, onSelect, onUpdate, onDelete }: Props = $props();

	let modal: Modal | undefined = $state();
	let editName = $state('');
	let editBody = $state('');

	function startEdit() {
		editName = prompt.name;
		editBody = prompt.body;
		modal?.open();
	}

	function cancelEdit() {
		modal?.close();
	}

	function saveEdit() {
		if (!editName.trim() || !editBody.trim()) return;

		onUpdate({
			id: prompt.id,
			name: editName.trim(),
			body: editBody.trim()
		});

		modal?.close();
	}

	function handleDelete() {
		if (!confirm('Delete this prompt?')) return;
		onDelete(prompt.id);
	}
</script>

<Card interactive darker active={isSelected} onclick={() => onSelect(prompt.id)}>
	<div class="flex items-center">
		<div class="flex flex-1 items-center gap-3 py-3 pr-2 pl-3">
			<input
				type="radio"
				name="selectedPrompt"
				class="pointer-events-none radio radio-sm radio-primary"
				checked={isSelected}
				tabindex="-1"
			/>
			<div class="flex-1">
				<p class="font-medium">{prompt.name}</p>
				<p class="line-clamp-1 text-xs opacity-70">{prompt.body}</p>
			</div>
		</div>
		<div class="flex shrink-0 gap-1 py-3 pr-3">
			<button
				class="btn btn-square btn-ghost btn-sm"
				onclick={(e) => {
					e.stopPropagation();
					startEdit();
				}}
				title="Edit"
			>
				<Pencil class="size-3.5" />
			</button>
			<button
				class="btn btn-square text-error btn-ghost btn-sm"
				onclick={(e) => {
					e.stopPropagation();
					handleDelete();
				}}
				title="Delete"
			>
				<Trash2 class="size-3.5" />
			</button>
		</div>
	</div>
</Card>

<Modal bind:this={modal} title="Edit Prompt: {prompt.name}" size="lg">
	{#snippet children()}
		<fieldset class="fieldset">
			<label class="label" for="editPromptName-{prompt.id}">
				<span class="label-text">Name</span>
			</label>
			<input
				type="text"
				id="editPromptName-{prompt.id}"
				class="input w-full"
				bind:value={editName}
			/>
		</fieldset>
		<fieldset class="fieldset">
			<label class="label" for="editPromptBody-{prompt.id}">
				<span class="label-text">Prompt Template</span>
			</label>
			<textarea
				id="editPromptBody-{prompt.id}"
				rows={12}
				class="textarea w-full"
				bind:value={editBody}
			></textarea>
		</fieldset>
	{/snippet}
	{#snippet actions()}
		<button class="btn btn-primary" onclick={saveEdit}>
			<Check class="size-4" />
			Save
		</button>
		<button class="btn" onclick={cancelEdit}>Cancel</button>
	{/snippet}
</Modal>
