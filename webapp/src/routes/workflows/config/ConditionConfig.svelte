<script lang="ts">
	import { GitBranch } from '@lucide/svelte';

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
				<GitBranch class="size-4" />
				Field to Check
			</span>
		</label>
		<input
			type="text"
			class="input-bordered input input-sm"
			placeholder="settings.postProcessEnabled"
			value={(config.field as string) ?? ''}
			onchange={(e) => onChange('field', e.currentTarget.value)}
			{disabled}
		/>
		<label class="label">
			<span class="label-text-alt opacity-70">Use "settings.fieldName" for settings fields</span>
		</label>
	</div>

	<div class="form-control">
		<label class="label" for="condition-operator">
			<span class="label-text">Operator</span>
		</label>
		<select
			id="condition-operator"
			class="select-bordered select select-sm"
			value={(config.operator as string) ?? 'equals'}
			onchange={(e) => onChange('operator', e.currentTarget.value)}
			{disabled}
		>
			<option value="equals">Equals (==)</option>
			<option value="notEquals">Not Equals (!=)</option>
			<option value="contains">Contains</option>
			<option value="startsWith">Starts With</option>
			<option value="endsWith">Ends With</option>
			<option value="isEmpty">Is Empty</option>
			<option value="isNotEmpty">Is Not Empty</option>
			<option value="greaterThan">Greater Than (&gt;)</option>
			<option value="lessThan">Less Than (&lt;)</option>
		</select>
	</div>

	<div class="form-control">
		<label class="label" for="condition-value">
			<span class="label-text">Expected Value</span>
		</label>
		<input
			id="condition-value"
			type="text"
			class="input-bordered input input-sm"
			placeholder="true"
			value={(config.value as string) ?? ''}
			onchange={(e) => onChange('value', e.currentTarget.value)}
			{disabled}
		/>
	</div>
</div>
