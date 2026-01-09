<script lang="ts">
	import { Globe } from '@lucide/svelte';

	interface Props {
		config: Record<string, unknown>;
		onChange: (key: string, value: unknown) => void;
		disabled?: boolean;
	}

	let { config, onChange, disabled = false }: Props = $props();
</script>

<div class="flex flex-col gap-3">
	<div class="form-control">
		<label class="label">
			<span class="label-text flex items-center gap-2">
				<Globe class="size-4" />
				URL
			</span>
		</label>
		<input
			type="text"
			class="input-bordered input input-sm"
			placeholder="https://api.example.com/endpoint"
			value={(config.url as string) ?? ''}
			onchange={(e) => onChange('url', e.currentTarget.value)}
			{disabled}
		/>
	</div>

	<div class="form-control">
		<label class="label" for="http-method">
			<span class="label-text">Method</span>
		</label>
		<select
			id="http-method"
			class="select-bordered select select-sm"
			value={(config.method as string) ?? 'GET'}
			onchange={(e) => onChange('method', e.currentTarget.value)}
			{disabled}
		>
			<option value="GET">GET</option>
			<option value="POST">POST</option>
			<option value="PUT">PUT</option>
			<option value="DELETE">DELETE</option>
			<option value="PATCH">PATCH</option>
		</select>
	</div>

	<div class="form-control">
		<label class="label" for="http-headers">
			<span class="label-text">Headers (JSON)</span>
		</label>
		<textarea
			id="http-headers"
			class="textarea-bordered textarea h-20 font-mono text-sm"
			placeholder={'{"Authorization": "Bearer token"}'}
			value={(config.headers as string) ?? ''}
			onchange={(e) => onChange('headers', e.currentTarget.value)}
			{disabled}
		></textarea>
	</div>

	<div class="form-control">
		<label class="label" for="http-body">
			<span class="label-text">Request Body (JSON)</span>
		</label>
		<textarea
			id="http-body"
			class="textarea-bordered textarea h-24 font-mono text-sm"
			placeholder={'{"key": "value"}'}
			value={(config.body as string) ?? ''}
			onchange={(e) => onChange('body', e.currentTarget.value)}
			{disabled}
		></textarea>
	</div>

	<div class="form-control">
		<label class="label" for="http-timeout">
			<span class="label-text">Timeout (ms)</span>
		</label>
		<input
			id="http-timeout"
			type="number"
			class="input-bordered input input-sm"
			placeholder="30000"
			min="1000"
			max="120000"
			value={(config.timeoutMs as number) ?? 30000}
			onchange={(e) => onChange('timeoutMs', parseInt(e.currentTarget.value) || 30000)}
			{disabled}
		/>
	</div>
</div>
