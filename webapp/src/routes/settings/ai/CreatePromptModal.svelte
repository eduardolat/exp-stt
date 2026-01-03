<script lang="ts">
	import { Modal, Card } from '$lib/components';
	import { Check, ArrowLeft, Search } from '@lucide/svelte';
	import { promptPresets, type PromptPreset } from '$lib/presets/prompts';
	import type { Prompt } from '$lib/client.gen';

	interface Props {
		onSave: (prompt: Prompt) => void;
	}

	let { onSave }: Props = $props();

	let modal: Modal | undefined = $state();
	let currentStep: 1 | 2 = $state(1);
	let promptName = $state('');
	let promptBody = $state('');
	let searchQuery = $state('');

	let filteredPresets = $derived.by(() => {
		if (!searchQuery.trim()) return promptPresets;
		const terms = searchQuery.toLowerCase().split(/\s+/).filter(Boolean);
		return promptPresets.filter((preset) => {
			const searchText = `${preset.name} ${preset.description}`.toLowerCase();
			return terms.every((term) => searchText.includes(term));
		});
	});

	export function open() {
		currentStep = 1;
		promptName = '';
		promptBody = '';
		searchQuery = '';
		modal?.open();
	}

	function selectPreset(preset: PromptPreset) {
		promptName = preset.name === 'Start from Scratch' ? '' : preset.name;
		promptBody = preset.body;
		currentStep = 2;
	}

	function goBack() {
		currentStep = 1;
	}

	function handleSave() {
		if (!promptName.trim() || !promptBody.trim()) return;

		const newPrompt: Prompt = {
			id: crypto.randomUUID(),
			name: promptName.trim(),
			body: promptBody.trim()
		};

		onSave(newPrompt);
		modal?.close();
	}

	function handleCancel() {
		modal?.close();
	}

	let modalTitle = $derived(
		currentStep === 1 ? 'Choose a Starting Point' : 'Customize Your Prompt'
	);
</script>

<Modal bind:this={modal} title={modalTitle} size="lg" disableBackdropClose>
	{#snippet children()}
		<div class="flex h-125 flex-col">
			{#if currentStep === 1}
				<div class="flex flex-1 flex-col gap-4 overflow-hidden">
					<label class="input w-full">
						<Search class="size-4" />
						<input
							type="text"
							class="grow"
							placeholder="Search presets..."
							bind:value={searchQuery}
						/>
					</label>

					<div class="flex-1 space-y-2 overflow-y-auto">
						{#each filteredPresets as preset (preset.id)}
							{@const PresetIcon = preset.icon}
							<Card interactive darker onclick={() => selectPreset(preset)}>
								<div class="flex items-center gap-3 p-3">
									<div
										class="flex size-10 shrink-0 items-center justify-center rounded-lg bg-base-300"
									>
										<PresetIcon class="size-5" />
									</div>
									<div class="flex-1">
										<p class="font-medium">{preset.name}</p>
										<p class="text-sm opacity-70">{preset.description}</p>
									</div>
								</div>
							</Card>
						{/each}

						{#if filteredPresets.length === 0}
							<p class="py-8 text-center text-sm opacity-70">No presets match your search</p>
						{/if}
					</div>
				</div>
			{/if}

			{#if currentStep === 2}
				<div class="flex flex-1 flex-col space-y-4 overflow-y-auto">
					<fieldset class="fieldset">
						<label class="label" for="newPromptName">
							<span class="label-text">Name</span>
						</label>
						<input
							type="text"
							id="newPromptName"
							class="input w-full"
							placeholder="My Custom Prompt"
							bind:value={promptName}
						/>
					</fieldset>
					<fieldset class="fieldset flex flex-1 flex-col">
						<label class="label" for="newPromptBody">
							<span class="label-text">
								Prompt Template
								<span class="opacity-70">
									(use {'${output}'} to insert the unprocessed transcription)
								</span>
							</span>
						</label>
						<textarea
							id="newPromptBody"
							class="textarea w-full flex-1"
							placeholder="Enter your prompt template..."
							bind:value={promptBody}
						></textarea>
					</fieldset>
				</div>
			{/if}
		</div>
	{/snippet}
	{#snippet actions()}
		{#if currentStep === 1}
			<button class="btn" onclick={handleCancel}>Cancel</button>
		{/if}
		{#if currentStep === 2}
			<button class="btn" onclick={goBack}>
				<ArrowLeft class="size-4" />
				Back
			</button>
			<button class="btn btn-primary" onclick={handleSave}>
				<Check class="size-4" />
				Save Prompt
			</button>
		{/if}
	{/snippet}
</Modal>
