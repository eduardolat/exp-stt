<script lang="ts">
	import { Modal } from '$lib/components';
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
	let editingPrompt: Prompt | null = $state(null);

	function startEdit() {
		editingPrompt = { ...prompt };
		modal?.open();
	}

	function cancelEdit() {
		editingPrompt = null;
		modal?.close();
	}

	function saveEdit() {
		if (!editingPrompt || !editingPrompt.name.trim() || !editingPrompt.body.trim()) return;
		onUpdate(editingPrompt);
		cancelEdit();
	}

	function handleDelete() {
		if (!confirm('Delete this prompt?')) return;
		onDelete(prompt.id);
	}
</script>

<div
	class="flex items-center rounded-lg border transition-colors {isSelected
		? 'border-primary bg-primary/10'
		: 'border-base-300'}"
>
	<button
		class="flex flex-1 cursor-pointer items-center gap-3 py-3 pr-2 pl-3 text-left"
		onclick={() => onSelect(prompt.id)}
	>
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
	</button>
	<div class="flex shrink-0 gap-1 py-3 pr-3">
		<button class="btn btn-square btn-ghost btn-sm" onclick={startEdit} title="Edit">
			<Pencil class="size-3.5" />
		</button>
		<button
			class="btn btn-square text-error btn-ghost btn-sm"
			onclick={handleDelete}
			title="Delete"
		>
			<Trash2 class="size-3.5" />
		</button>
	</div>
</div>

<Modal bind:this={modal} title="Edit Prompt: {prompt.name}" size="lg">
	{#snippet children()}
		{#if editingPrompt}
			<fieldset class="fieldset">
				<label class="label" for="editPromptName-{prompt.id}">
					<span class="label-text">Name</span>
				</label>
				<input
					type="text"
					id="editPromptName-{prompt.id}"
					class="input w-full"
					bind:value={editingPrompt.name}
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
					bind:value={editingPrompt.body}
				></textarea>
			</fieldset>
		{/if}
	{/snippet}
	{#snippet actions()}
		<button class="btn btn-primary" onclick={saveEdit}>
			<Check class="size-4" />
			Save
		</button>
		<button class="btn" onclick={cancelEdit}>Cancel</button>
	{/snippet}
</Modal>
