<script lang="ts">
	import { store } from '$lib/store.svelte';
	import { Card, Modal, PageHeader } from '$lib/components';
	import PromptEditor from './PromptEditor.svelte';
	import ModelPicker from './ModelPicker.svelte';
	import BaseUrlPicker from './BaseUrlPicker.svelte';
	import {
		Sparkles,
		Plus,
		Check,
		Bot,
		Key,
		Link,
		FileText,
		Zap,
		ZapOff,
		Search,
		FileDigit
	} from '@lucide/svelte';
	import type { Prompt } from '$lib/client.gen';

	let newPromptName = $state('');
	let newPromptBody = $state('');
	let modelPicker: { open: () => void } | undefined = $state();
	let baseUrlPicker: { open: () => void } | undefined = $state();

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

<PageHeader
	icon={Sparkles}
	title="AI Enhancement"
	description="Configure AI post-processing for transcriptions"
>
	{#snippet actions()}
		<label class="label cursor-pointer gap-2">
			<span class="text-sm">Enable AI Enhancement</span>
			<input
				type="checkbox"
				class="toggle"
				checked={store.settings.postProcessEnabled}
				onchange={(e) => store.updateSettings({ postProcessEnabled: e.currentTarget.checked })}
			/>
		</label>
	{/snippet}
</PageHeader>

{#if store.settings.postProcessEnabled}
	<div class="space-y-4">
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
					<div class="join w-full">
						<input
							type="url"
							id="baseUrl"
							class="input join-item w-full"
							placeholder="https://api.openai.com/v1"
							value={store.settings.postProcessBaseUrl}
							onblur={(e) => store.updateSettings({ postProcessBaseUrl: e.currentTarget.value })}
						/>
						<button
							class="btn join-item border-base-content/20 btn-outline"
							onclick={() => baseUrlPicker?.open()}
							title="Browse providers"
						>
							<Search class="size-4" />
						</button>
					</div>
					<p class="label">
						<span class="label-text-alt">Enter manually or browse providers</span>
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
						autocomplete="off"
						class="input w-full"
						placeholder="sk-..."
						value={store.settings.postProcessApiKey}
						onblur={(e) => store.updateSettings({ postProcessApiKey: e.currentTarget.value })}
					/>
					<p class="label">
						<span class="label-text-alt">Get your API key from your provider's dashboard</span>
					</p>
				</fieldset>

				<fieldset class="fieldset">
					<label class="label" for="model">
						<span class="label-text flex items-center gap-1.5">
							<FileDigit class="size-3.5" />
							Model
						</span>
					</label>
					<div class="join w-full">
						<input
							type="text"
							id="model"
							class="input join-item w-full"
							placeholder="gpt-4o-mini"
							value={store.settings.postProcessModel}
							onblur={(e) => store.updateSettings({ postProcessModel: e.currentTarget.value })}
						/>
						<button
							class="btn join-item border-base-content/20 btn-outline"
							onclick={() => modelPicker?.open()}
							title="Browse available models"
						>
							<Search class="size-4" />
						</button>
					</div>
					<p class="label">
						<span class="label-text-alt">Enter manually or browse available models</span>
					</p>
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
	</div>
{:else}
	<Card class="card-body">
		<div class="flex flex-col items-center justify-center py-8 text-center">
			<ZapOff class="mb-4 size-10 opacity-50" />

			<h3 class="mb-2 text-lg font-medium">AI Enhancement Disabled</h3>
			<p class="mb-6 max-w-md text-sm opacity-70">
				AI Enhancement uses an OpenAI compatible API (either cloud based or running locally) to
				automatically improve, clean up, and format your transcriptions after the transcription is
				completed.
			</p>
			<button
				class="btn btn-primary"
				onclick={() => store.updateSettings({ postProcessEnabled: true })}
			>
				<Zap class="size-4" />
				Enable AI Enhancement
			</button>
		</div>
	</Card>
{/if}

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

<ModelPicker
	bind:this={modelPicker}
	baseUrl={store.settings.postProcessBaseUrl}
	apiKey={store.settings.postProcessApiKey}
	onSelect={(id) => store.updateSettings({ postProcessModel: id })}
/>

<BaseUrlPicker
	bind:this={baseUrlPicker}
	onSelect={(url) => store.updateSettings({ postProcessBaseUrl: url })}
/>
