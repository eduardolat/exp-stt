<!--
	Animated logo component that mirrors the systray animation behavior.
	Animation positions cycle in a ping-pong pattern: middle → right → middle → left → middle
	For static states (unloaded, loaded), the position is always 'middle'.
-->

<script lang="ts" module>
	export type AppStatus =
		| 'unknown'
		| 'unloaded'
		| 'downloading'
		| 'loading'
		| 'loaded'
		| 'listening'
		| 'transcribing'
		| 'post_processing';
</script>

<script lang="ts">
	import { untrack } from 'svelte';

	type LogoPosition = 'left' | 'middle' | 'right';
	type LogoColor = 'gray' | 'amber' | 'white' | 'pink' | 'blue' | 'green';

	type Props = {
		status: AppStatus;
		size?: number;
		class?: string;
	};

	const FRAME_DURATION_MS = 200;

	const STATUS_COLOR: Record<AppStatus, LogoColor> = {
		unknown: 'gray',
		unloaded: 'gray',
		downloading: 'amber',
		loading: 'amber',
		loaded: 'white',
		listening: 'pink',
		transcribing: 'blue',
		post_processing: 'green'
	};

	const STATIC_STATUSES: AppStatus[] = ['unknown', 'unloaded', 'loaded'];

	let { status, size = 24, class: className = '' }: Props = $props();

	let position: LogoPosition = $state('middle');
	let backward = false;
	let timer: ReturnType<typeof setTimeout> | null = null;

	const logoPath = $derived(`/logo/svg/black-${STATUS_COLOR[status] ?? 'gray'}-${position}.svg`);

	function stopAnimation() {
		if (timer !== null) {
			clearTimeout(timer);
			timer = null;
		}
	}

	function getNextPosition(current: LogoPosition): [LogoPosition, boolean] {
		switch (current) {
			case 'middle':
				return backward ? ['left', backward] : ['right', backward];
			case 'right':
				return ['middle', true];
			case 'left':
				return ['middle', false];
			default:
				return ['middle', false];
		}
	}

	function tick() {
		const currentStatus = untrack(() => status);

		if (STATIC_STATUSES.includes(currentStatus)) {
			position = 'middle';
			backward = false;
			return;
		}

		const [nextPos, nextBackward] = getNextPosition(untrack(() => position));
		position = nextPos;
		backward = nextBackward;
		timer = setTimeout(tick, FRAME_DURATION_MS);
	}

	$effect(() => {
		status; // Track status changes
		untrack(() => {
			stopAnimation();
			tick();
		});
		return () => stopAnimation();
	});
</script>

<img src={logoPath} alt="Tribar Logo" width={size} height={size} class={className} />
