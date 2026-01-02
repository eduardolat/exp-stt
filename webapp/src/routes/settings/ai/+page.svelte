<script lang="ts">
	import { store } from '$lib/store.svelte';
	import { Card, Modal } from '$lib/components';
	import PromptEditor from './PromptEditor.svelte';
	import { Sparkles, Plus, Check, Bot, Key, Link, FileText } from '@lucide/svelte';
	import type { Prompt } from '$lib/client.gen';

	let newPromptName = $state('');
	let newPromptBody = $state('');

	let createModal: Modal | undefined = $state();

	function startCreate() {
		newPromptName = '';
		newPromptBody =
			'You are a helpful assistant. Clean up and improve the following transcription:\n\n${output}\n\nProvide only the improved text without any additional commentary.';
		createModal?.open();
	}

	function cancelCreate() {
		newPromptName = '';
		newPromptBody = '';
		createModal?.close();
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

	function updatePrompt(updatedPrompt: Prompt) {
		const updatedPrompts = store.prompts.map((p) =>
			p.id === updatedPrompt.id ? updatedPrompt : p
		);
		store.updateSettings({ prompts: updatedPrompts });
	}

	function deletePrompt(id: string) {
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

<div class="space-y-4">
	<Card class="card-body">
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
	</Card>

	{#if store.settings.postProcessEnabled}
		<Card class="card-body">
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
		</Card>

		<Card class="card-body">
			<div class="flex items-center justify-between">
				<h3 class="card-title text-base">
					<FileText class="size-4" />
					Prompts
				</h3>
				<button class="btn btn-outline btn-sm" onclick={startCreate}>
					<Plus class="size-3.5" />
					New Prompt
				</button>
			</div>

			{#if store.prompts.length === 0}
				<p class="py-4 text-center text-sm opacity-70">
					No prompts yet. Create one to get started.
				</p>
			{/if}

			{#if store.prompts.length > 0}
				<div class="mt-4 space-y-2">
					{#each store.prompts as prompt (prompt.id)}
						<PromptEditor
							{prompt}
							isSelected={store.settings.postProcessPromptId === prompt.id}
							onSelect={selectPrompt}
							onUpdate={updatePrompt}
							onDelete={deletePrompt}
						/>
					{/each}
				</div>
			{/if}
		</Card>
	{/if}
</div>

<Modal bind:this={createModal} title="Create New Prompt" size="lg">
	{#snippet children()}
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
				rows={12}
				class="textarea w-full"
				placeholder="Enter your prompt template..."
				bind:value={newPromptBody}
			></textarea>
		</fieldset>
	{/snippet}
	{#snippet actions()}
		<button class="btn btn-primary" onclick={saveNewPrompt}>
			<Check class="size-4" />
			Save
		</button>
		<button class="btn" onclick={cancelCreate}>Cancel</button>
	{/snippet}
</Modal>
