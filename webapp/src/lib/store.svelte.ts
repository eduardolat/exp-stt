import {
	NewClient,
	type Client,
	type State,
	type Settings,
	type HistoryEntry,
	type AudioDevice,
	type DownloadProgress,
	type SystemInfo,
	type Prompt,
	type Shortcut
} from './client.gen';

const API_BASE_URL = '/api/v1/urpc';

function createDefaultState(): State {
	return {
		status: 'unknown',
		downloadProgress: { fileName: '', downloaded: 0, total: 0, percent: 0 },
		devices: { inputDevices: [], outputDevices: [] },
		history: [],
		systemInfo: { os: 'unknown', arch: 'unknown', displayServer: 'unknown' }
	};
}

function createDefaultSettings(): Settings {
	return {
		schemaVersion: 1,
		notifyOnError: true,
		notifyOnStart: false,
		notifyOnFinish: false,
		soundFeedbackEnable: true,
		soundFeedbackRecordId: 'default',
		soundFeedbackSuccessId: 'default',
		soundFeedbackErrorId: 'default',
		soundFeedbackVolume: 80,
		inputDevice: 'default',
		outputMode: 'copy_only',
		outputTrailingSpace: true,
		postProcessEnabled: false,
		postProcessBaseUrl: '',
		postProcessApiKey: '',
		postProcessModel: '',
		postProcessPromptId: '',
		prompts: [],
		historyLimit: 100,
		modelUnloadEnable: false,
		modelUnloadSeconds: 300,
		shortcutToggle: { modifiers: [], key: '' }
	};
}

class Store {
	private client: Client;
	private eventStreamCancel: (() => void) | null = null;

	state: State = $state(createDefaultState());
	settings: Settings = $state(createDefaultSettings());
	isConnected = $state(false);
	isLoading = $state(true);
	error: string | null = $state(null);

	get status() {
		return this.state.status;
	}

	get downloadProgress(): DownloadProgress {
		return this.state.downloadProgress;
	}

	get inputDevices(): AudioDevice[] {
		return this.state.devices.inputDevices;
	}

	get outputDevices(): AudioDevice[] {
		return this.state.devices.outputDevices;
	}

	get history(): HistoryEntry[] {
		return this.state.history;
	}

	get systemInfo(): SystemInfo {
		return this.state.systemInfo;
	}

	get prompts(): Prompt[] {
		return this.settings.prompts;
	}

	get shortcut(): Shortcut {
		return this.settings.shortcutToggle;
	}

	get isRecording(): boolean {
		return this.status === 'listening';
	}

	get isProcessing(): boolean {
		return (
			this.status === 'transcribing' ||
			this.status === 'post_processing' ||
			this.status === 'downloading' ||
			this.status === 'loading'
		);
	}

	get isReady(): boolean {
		return this.status === 'loaded';
	}

	get statusLabel(): string {
		const labels: Record<string, string> = {
			unknown: 'Unknown',
			unloaded: 'Model Unloaded',
			downloading: 'Downloading Model',
			loading: 'Loading Model',
			loaded: 'Ready',
			listening: 'Recording',
			transcribing: 'Transcribing',
			post_processing: 'Post-Processing'
		};
		return labels[this.status] ?? 'Unknown';
	}

	constructor() {
		this.client = NewClient(API_BASE_URL).build();
	}

	async initialize(): Promise<void> {
		this.isLoading = true;
		this.error = null;

		try {
			await Promise.all([this.fetchState(), this.fetchSettings()]);
			this.startEventStream();
			this.isConnected = true;
		} catch (err) {
			this.error = err instanceof Error ? err.message : 'Failed to initialize';
			console.error('Store initialization failed:', err);
		} finally {
			this.isLoading = false;
		}
	}

	private async fetchState(): Promise<void> {
		const { state } = await this.client.procs.stateGet().execute({});
		this.state = state;
	}

	private async fetchSettings(): Promise<void> {
		const { settings } = await this.client.procs.settingsGet().execute({});
		this.settings = settings;
	}

	private startEventStream(): void {
		const { stream, cancel } = this.client.streams
			.listenForEvents()
			.withReconnect({ maxAttempts: 10, initialDelayMs: 1000, maxDelayMs: 10000 })
			.execute({});

		this.eventStreamCancel = cancel;

		(async () => {
			for await (const event of stream) {
				if (!event.ok) {
					console.error('Event stream error:', event.error);
					this.isConnected = false;
					continue;
				}

				this.isConnected = true;

				switch (event.output.eventType) {
					case 'stateUpdated':
						if (event.output.stateUpdated) {
							this.state = event.output.stateUpdated;
						}
						break;
					case 'settingsUpdated':
						if (event.output.settingsUpdated) {
							this.settings = event.output.settingsUpdated;
						}
						break;
					case 'ping':
						break;
				}
			}
		})();
	}

	async toggleRecording(): Promise<void> {
		try {
			await this.client.procs.recordingToggle().execute({});
		} catch (err) {
			console.error('Failed to toggle recording:', err);
			throw err;
		}
	}

	async updateSettings(newSettings: Partial<Settings>): Promise<void> {
		const merged: Settings = { ...this.settings, ...newSettings };
		try {
			await this.client.procs.settingsUpdate().execute({ settings: merged });
		} catch (err) {
			console.error('Failed to update settings:', err);
		}
	}

	async updateShortcut(shortcut: Shortcut): Promise<void> {
		try {
			await this.client.procs.shortcutToggleUpdate().execute({ shortcut });
		} catch (err) {
			console.error('Failed to update shortcut:', err);
			throw err;
		}
	}

	async deleteHistoryEntry(id: string): Promise<void> {
		try {
			await this.client.procs.historyDeleteEntry().execute({ id });
		} catch (err) {
			console.error('Failed to delete history entry:', err);
			throw err;
		}
	}

	async clearHistory(): Promise<void> {
		try {
			await this.client.procs.historyClear().execute({});
		} catch (err) {
			console.error('Failed to clear history:', err);
			throw err;
		}
	}

	destroy(): void {
		if (this.eventStreamCancel) {
			this.eventStreamCancel();
			this.eventStreamCancel = null;
		}
	}
}

export const store = new Store();
