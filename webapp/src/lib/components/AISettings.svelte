<script lang="ts">
	import { store } from '$lib/store.svelte';
	import {
		Sparkles,
		Plus,
		Pencil,
		Trash2,
		Check,
		X,
		Bot,
		Key,
		Link,
		FileText
	} from '@lucide/svelte';
	import type { Prompt } from '$lib/client.gen';

	let editingPrompt = $state<Prompt | null>(null);
	let isCreating = $state(false);

	let newPromptName = $state('');
	let newPromptBody = $state('');

	function startCreate() {
		isCreating = true;
		newPromptName = '';
		newPromptBody =
			'You are a helpful assistant. Clean up and improve the following transcription:\n\n${output}\n\nProvide only the improved text without any additional commentary.';
	}

	function cancelCreate() {
		isCreating = false;
		newPromptName = '';
		newPromptBody = '';
	}

	function saveNewPrompt() {
		if (!newPromptName.trim() || !newPromptBody.trim()) return;

		const newPrompt: Prompt = {
			id: crypto.randomUUID(),
			name: newPromptName.trim(),
			body: newPromptBody.trim()
		};

		store.updateSettings({
			prompts: [...store.prompts, newPrompt]
		});

		cancelCreate();
	}

	function startEdit(prompt: Prompt) {
		editingPrompt = { ...prompt };
	}

	function cancelEdit() {
		editingPrompt = null;
	}

	function saveEdit() {
		if (!editingPrompt || !editingPrompt.name.trim() || !editingPrompt.body.trim()) return;

		const updatedPrompts = store.prompts.map((p) =>
			p.id === editingPrompt!.id ? editingPrompt! : p
		);

		store.updateSettings({ prompts: updatedPrompts });
		editingPrompt = null;
	}

	function deletePrompt(id: string) {
		if (!confirm('Delete this prompt?')) return;

		const updatedPrompts = store.prompts.filter((p) => p.id !== id);
		const updates: Record<string, unknown> = { prompts: updatedPrompts };

		if (store.settings.postProcessPromptId === id) {
			updates.postProcessPromptId = updatedPrompts[0]?.id ?? '';
		}

		store.updateSettings(updates);
	}

	function selectPrompt(id: string) {
		store.updateSettings({ postProcessPromptId: id });
	}
</script>

<div class="space-y-6">
	<!-- AI Toggle -->
	<section class="card p-4">
		<div class="flex items-center justify-between">
			<div class="flex items-center gap-3">
				<div class="flex size-10 items-center justify-center rounded-lg bg-purple-500/10">
					<Sparkles class="size-5 text-purple-600" />
				</div>
				<div>
					<h3 class="font-medium">AI Post-Processing</h3>
					<p class="text-xs text-muted-foreground">Enhance transcriptions with LLM</p>
				</div>
			</div>
			<input
				type="checkbox"
				role="switch"
				checked={store.settings.postProcessEnabled}
				onchange={(e) => store.updateSettings({ postProcessEnabled: e.currentTarget.checked })}
			/>
		</div>
	</section>

	{#if store.settings.postProcessEnabled}
		<!-- API Configuration -->
		<section class="card p-4">
			<h3 class="mb-4 flex items-center gap-2 font-medium">
				<Bot class="size-4" />
				API Configuration
			</h3>
			<div class="space-y-4">
				<div class="field">
					<label for="baseUrl" class="flex items-center gap-1.5 text-sm">
						<Link class="size-3.5" />
						Base URL
					</label>
					<input
						type="url"
						id="baseUrl"
						placeholder="https://api.openai.com/v1"
						value={store.settings.postProcessBaseUrl}
						onchange={(e) => store.updateSettings({ postProcessBaseUrl: e.currentTarget.value })}
						class="input"
					/>
					<p class="text-xs text-muted-foreground">OpenAI-compatible API endpoint</p>
				</div>

				<div class="field">
					<label for="apiKey" class="flex items-center gap-1.5 text-sm">
						<Key class="size-3.5" />
						API Key
					</label>
					<input
						type="password"
						id="apiKey"
						placeholder="sk-..."
						value={store.settings.postProcessApiKey}
						onchange={(e) => store.updateSettings({ postProcessApiKey: e.currentTarget.value })}
						class="input"
					/>
				</div>

				<div class="field">
					<label for="model" class="text-sm">Model</label>
					<input
						type="text"
						id="model"
						placeholder="gpt-4o-mini"
						value={store.settings.postProcessModel}
						onchange={(e) => store.updateSettings({ postProcessModel: e.currentTarget.value })}
						class="input"
					/>
				</div>
			</div>
		</section>

		<!-- Prompts -->
		<section class="card p-4">
			<div class="mb-4 flex items-center justify-between">
				<h3 class="flex items-center gap-2 font-medium">
					<FileText class="size-4" />
					Prompts
				</h3>
				{#if !isCreating}
					<button class="btn btn-sm btn-outline" onclick={startCreate}>
						<Plus class="mr-1.5 size-3.5" />
						New Prompt
					</button>
				{/if}
			</div>

			{#if isCreating}
				<div class="mb-4 rounded-lg border p-4">
					<div class="mb-3 space-y-3">
						<div class="field">
							<label for="newPromptName" class="text-sm">Name</label>
							<input
								type="text"
								id="newPromptName"
								placeholder="My Custom Prompt"
								bind:value={newPromptName}
								class="input"
							/>
						</div>
						<div class="field">
							<label for="newPromptBody" class="text-sm">
								Prompt Template
								<span class="text-muted-foreground">(use {'${output}'} for transcription)</span>
							</label>
							<textarea
								id="newPromptBody"
								rows={5}
								bind:value={newPromptBody}
								class="textarea"
								placeholder="Enter your prompt template..."
							></textarea>
						</div>
					</div>
					<div class="flex gap-2">
						<button class="btn btn-sm" onclick={saveNewPrompt}>
							<Check class="mr-1.5 size-3.5" />
							Save
						</button>
						<button class="btn btn-sm btn-outline" onclick={cancelCreate}>
							<X class="mr-1.5 size-3.5" />
							Cancel
						</button>
					</div>
				</div>
			{/if}

			{#if store.prompts.length === 0}
				<p class="py-4 text-center text-sm text-muted-foreground">
					No prompts yet. Create one to get started.
				</p>
			{:else}
				<div class="space-y-2">
					{#each store.prompts as prompt (prompt.id)}
						{#if editingPrompt?.id === prompt.id}
							<div class="rounded-lg border p-4">
								<div class="mb-3 space-y-3">
									<div class="field">
										<label for="editPromptName" class="text-sm">Name</label>
										<input
											type="text"
											id="editPromptName"
											bind:value={editingPrompt.name}
											class="input"
										/>
									</div>
									<div class="field">
										<label for="editPromptBody" class="text-sm">Prompt Template</label>
										<textarea
											id="editPromptBody"
											rows={5}
											bind:value={editingPrompt.body}
											class="textarea"
										></textarea>
									</div>
								</div>
								<div class="flex gap-2">
									<button class="btn btn-sm" onclick={saveEdit}>
										<Check class="mr-1.5 size-3.5" />
										Save
									</button>
									<button class="btn btn-sm btn-outline" onclick={cancelEdit}>
										<X class="mr-1.5 size-3.5" />
										Cancel
									</button>
								</div>
							</div>
						{:else}
							{@const isSelected = store.settings.postProcessPromptId === prompt.id}
							<div
								class="flex items-center justify-between rounded-lg border p-3 transition-colors {isSelected
									? 'border-primary bg-primary/5'
									: ''}"
							>
								<button
									class="flex flex-1 items-center gap-3 text-left"
									onclick={() => selectPrompt(prompt.id)}
								>
									<div
										class="flex size-4 items-center justify-center rounded-full border {isSelected
											? 'border-primary bg-primary'
											: ''}"
									>
										{#if isSelected}
											<Check class="size-2.5 text-primary-foreground" />
										{/if}
									</div>
									<div>
										<p class="font-medium">{prompt.name}</p>
										<p class="line-clamp-1 text-xs text-muted-foreground">{prompt.body}</p>
									</div>
								</button>
								<div class="flex gap-1">
									<button
										class="btn btn-sm btn-ghost"
										onclick={() => startEdit(prompt)}
										title="Edit"
									>
										<Pencil class="size-3.5" />
									</button>
									<button
										class="btn btn-sm btn-ghost text-destructive"
										onclick={() => deletePrompt(prompt.id)}
										title="Delete"
									>
										<Trash2 class="size-3.5" />
									</button>
								</div>
							</div>
						{/if}
					{/each}
				</div>
			{/if}
		</section>
	{/if}
</div>
