<script lang="ts">
	import { Pause, Play } from '@lucide/svelte';

	interface Props {
		label: string;
		soundType: 'record' | 'success' | 'error';
		options: string[];
		value: string;
		onchange: (value: string) => void;
	}

	let { label, soundType, options, value, onchange }: Props = $props();

	let isPlaying = $state(false);
	let currentAudio: HTMLAudioElement | null = null;

	function getSoundUrl(id: string): string {
		const apiBaseUrl = localStorage.getItem('tribar_api_base_url') || '';
		const baseUrl = apiBaseUrl ? apiBaseUrl.replace('/urpc', '') : '/api/v1';
		return `${baseUrl}/sound/${soundType}/${id}`;
	}

	function playSound() {
		if (currentAudio) {
			currentAudio.pause();
			currentAudio = null;
		}

		const audio = new Audio(getSoundUrl(value));
		currentAudio = audio;
		isPlaying = true;

		audio.onended = () => {
			isPlaying = false;
			currentAudio = null;
		};

		audio.onerror = () => {
			isPlaying = false;
			currentAudio = null;
		};

		audio.play();
	}

	function stopSound() {
		if (currentAudio) {
			currentAudio.pause();
			currentAudio.currentTime = 0;
			currentAudio = null;
		}
		isPlaying = false;
	}
</script>

<fieldset class="fieldset">
	<label class="label" for={`sound-${soundType}`}>
		<span class="label-text">{label}</span>
	</label>
	<div class="flex gap-2">
		<select
			id={`sound-${soundType}`}
			class="select flex-1 select-sm"
			{value}
			onchange={(e) => onchange(e.currentTarget.value)}
		>
			{#each options as option (option)}
				<option value={option}>Sound {option}</option>
			{/each}
		</select>
		<button
			type="button"
			class="btn btn-square btn-sm"
			onclick={() => (isPlaying ? stopSound() : playSound())}
			title={isPlaying ? 'Stop' : 'Preview'}
			aria-label={isPlaying ? 'Stop sound preview' : 'Preview sound'}
		>
			{#if isPlaying}
				<Pause class="size-4" />
			{:else}
				<Play class="size-4" />
			{/if}
		</button>
	</div>
</fieldset>
