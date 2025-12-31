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

<svelte:head>
	<title>AI Settings - Tribar Voice</title>
</svelte:head>

<div class="space-y-6">
	<div class="card bg-base-200">
		<div class="card-body">
			<div class="flex items-center justify-between">
				<div class="flex items-center gap-3">
					<div class="flex size-10 items-center justify-center rounded-lg bg-secondary">
						<Sparkles class="size-5 text-secondary-content" />
					</div>
					<div>
						<h3 class="font-medium">AI Post-Processing</h3>
						<p class="text-xs opacity-70">Enhance transcriptions with LLM</p>
					</div>
				</div>
				<input
					type="checkbox"
					class="toggle"
					checked={store.settings.postProcessEnabled}
					onchange={(e) => store.updateSettings({ postProcessEnabled: e.currentTarget.checked })}
				/>
			</div>
		</div>
	</div>

	{#if store.settings.postProcessEnabled}
		<div class="card bg-base-200">
			<div class="card-body">
				<h3 class="card-title text-base">
					<Bot class="size-4" />
					API Configuration
				</h3>
				<div class="space-y-4">
					<fieldset class="fieldset">
						<label class="label" for="baseUrl">
							<span class="label-text flex items-center gap-1.5">
								<Link class="size-3.5" />
								Base URL
							</span>
						</label>
						<input
							type="url"
							id="baseUrl"
							class="input w-full"
							placeholder="https://api.openai.com/v1"
							value={store.settings.postProcessBaseUrl}
							onblur={(e) => store.updateSettings({ postProcessBaseUrl: e.currentTarget.value })}
						/>
						<p class="label">
							<span class="label-text-alt">OpenAI-compatible API endpoint</span>
						</p>
					</fieldset>

					<fieldset class="fieldset">
						<label class="label" for="apiKey">
							<span class="label-text flex items-center gap-1.5">
								<Key class="size-3.5" />
								API Key
							</span>
						</label>
						<input
							type="password"
							id="apiKey"
							class="input w-full"
							placeholder="sk-..."
							value={store.settings.postProcessApiKey}
							onblur={(e) => store.updateSettings({ postProcessApiKey: e.currentTarget.value })}
						/>
					</fieldset>

					<fieldset class="fieldset">
						<label class="label" for="model">
							<span class="label-text">Model</span>
						</label>
						<input
							type="text"
							id="model"
							class="input w-full"
							placeholder="gpt-4o-mini"
							value={store.settings.postProcessModel}
							onblur={(e) => store.updateSettings({ postProcessModel: e.currentTarget.value })}
						/>
					</fieldset>
				</div>
			</div>
		</div>

		<div class="card bg-base-200">
			<div class="card-body">
				<div class="flex items-center justify-between">
					<h3 class="card-title text-base">
						<FileText class="size-4" />
						Prompts
					</h3>
					{#if !isCreating}
						<button class="btn btn-outline btn-sm" onclick={startCreate}>
							<Plus class="size-3.5" />
							New Prompt
						</button>
					{/if}
				</div>

				{#if isCreating}
					<div class="mt-4 rounded-lg border border-base-300 p-4">
						<div class="space-y-3">
							<fieldset class="fieldset">
								<label class="label" for="newPromptName">
									<span class="label-text">Name</span>
								</label>
								<input
									type="text"
									id="newPromptName"
									class="input w-full"
									placeholder="My Custom Prompt"
									bind:value={newPromptName}
								/>
							</fieldset>
							<fieldset class="fieldset">
								<label class="label" for="newPromptBody">
									<span class="label-text">
										Prompt Template
										<span class="opacity-70">(use {'${output}'} for transcription)</span>
									</span>
								</label>
								<textarea
									id="newPromptBody"
									rows={5}
									class="textarea w-full"
									placeholder="Enter your prompt template..."
									bind:value={newPromptBody}
								></textarea>
							</fieldset>
						</div>
						<div class="mt-3 flex gap-2">
							<button class="btn btn-sm btn-primary" onclick={saveNewPrompt}>
								<Check class="size-3.5" />
								Save
							</button>
							<button class="btn btn-ghost btn-sm" onclick={cancelCreate}>
								<X class="size-3.5" />
								Cancel
							</button>
						</div>
					</div>
				{/if}

				{#if store.prompts.length === 0}
					<p class="py-4 text-center text-sm opacity-70">
						No prompts yet. Create one to get started.
					</p>
				{/if}

				{#if store.prompts.length > 0}
					<div class="mt-4 space-y-2">
						{#each store.prompts as prompt (prompt.id)}
							{#if editingPrompt?.id === prompt.id}
								<div class="rounded-lg border border-base-300 p-4">
									<div class="space-y-3">
										<fieldset class="fieldset">
											<label class="label" for="editPromptName">
												<span class="label-text">Name</span>
											</label>
											<input
												type="text"
												id="editPromptName"
												class="input w-full"
												bind:value={editingPrompt.name}
											/>
										</fieldset>
										<fieldset class="fieldset">
											<label class="label" for="editPromptBody">
												<span class="label-text">Prompt Template</span>
											</label>
											<textarea
												id="editPromptBody"
												rows={5}
												class="textarea w-full"
												bind:value={editingPrompt.body}
											></textarea>
										</fieldset>
									</div>
									<div class="mt-3 flex gap-2">
										<button class="btn btn-sm btn-primary" onclick={saveEdit}>
											<Check class="size-3.5" />
											Save
										</button>
										<button class="btn btn-ghost btn-sm" onclick={cancelEdit}>
											<X class="size-3.5" />
											Cancel
										</button>
									</div>
								</div>
							{:else}
								{@const isSelected = store.settings.postProcessPromptId === prompt.id}
								<div
									class="flex items-center justify-between rounded-lg border p-3 transition-colors {isSelected
										? 'border-primary bg-primary/10'
										: 'border-base-300'}"
								>
									<button
										class="flex flex-1 items-center gap-3 text-left"
										onclick={() => selectPrompt(prompt.id)}
									>
										<input
											type="radio"
											name="selectedPrompt"
											class="radio radio-sm radio-primary"
											checked={isSelected}
											onchange={() => selectPrompt(prompt.id)}
										/>
										<div>
											<p class="font-medium">{prompt.name}</p>
											<p class="line-clamp-1 text-xs opacity-70">{prompt.body}</p>
										</div>
									</button>
									<div class="flex gap-1">
										<button
											class="btn btn-ghost btn-sm"
											onclick={() => startEdit(prompt)}
											title="Edit"
										>
											<Pencil class="size-3.5" />
										</button>
										<button
											class="btn text-error btn-ghost btn-sm"
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
			</div>
		</div>
	{/if}
</div>
