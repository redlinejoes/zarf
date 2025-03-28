<!-- 
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2021-Present The Zarf Authors
 -->
<script lang="ts">
	import YamlCode from './yaml-code.svelte';

	import type { ZarfComponent } from '$lib/api-types';
	import { pkgComponentDeployStore } from '$lib/store';
	import { Accordion, IconButton, Typography } from '@ui';
	import ZarfChip from './zarf-chip.svelte';

	export let idx: number;
	export let readOnly: boolean = true;
	export let component: ZarfComponent;

	const toggleComponentDeployment = (list: number[], idx: number) => {
		const enabled = list.includes(idx);
		if (enabled) {
			list = [...list].filter((n) => n !== idx);
		} else {
			list = [...list, idx];
		}
		list.sort();
		pkgComponentDeployStore.set(list);
	};
	$: requiredText = (component.required && 'required') || 'optional';
	$: componentTitle = `${component.name} ${requiredText}`;
</script>

<Accordion id={`component-accordion-${idx}`} class="package-component-accordion">
	<div slot="headerContent" class="component-accordion-header">
		<div style="flex: 1; gap: 8px;">
			<Typography variant="subtitle2" element="div" class="component-title" title={componentTitle}>
				{component.name}
			</Typography>
			<ZarfChip>{requiredText}</ZarfChip>
		</div>
		<div style="flex: 3">
			<Typography
				variant="body2"
				color="text-secondary-on-dark"
				class="component-description"
				title={component.description || ''}
				element="div"
			>
				{component.description || ' '}
			</Typography>
		</div>
		{#if !readOnly}
			<div style={`gap: 5px; visibility: ${component.required ? 'hidden' : 'initial'}`}>
				<IconButton
					toggleable
					class="deploy-component-toggle"
					iconColor="inherit"
					iconContent="toggle_off"
					toggledIconColor="primary"
					id={`deploy-component-${idx}`}
					toggledIconContent="toggle_on"
					iconClass="material-symbols-outlined"
					disabled={component.required}
					toggled={$pkgComponentDeployStore.includes(idx)}
					on:click={() => toggleComponentDeployment($pkgComponentDeployStore, idx)}
				/>
				<Typography
					variant="body1"
					element="label"
					for={`deploy-component-${idx}`}
					style={component.required
						? 'color: var(--mdc-theme-text-secondary-on-light);'
						: 'var(--mdc-theme-text-disabled-on-light, rgba(0, 0, 0, 0.38))'}
				>
					Deploy
				</Typography>
			</div>
		{/if}
	</div>
	<YamlCode slot="content" code={component} />
</Accordion>

<style>
	:global(.accordion-header-wrapper) {
		height: 51px;
	}
	:global(.accordion-header) {
		width: 100%;
		display: flex;
		overflow: hidden;
	}
	:global(.component-accordion-header) {
		gap: 30px;
		width: 100%;
		display: flex;
		overflow: hidden;
		justify-content: space-between;
	}
	:global(.component-title .optional) {
		color: var(--mdc-theme-primary-dark);
	}
	.component-accordion-header div {
		display: flex;
		align-items: center;
		overflow: hidden;
		flex-wrap: nowrap;
		white-space: nowrap;
		text-overflow: ellipsis;
	}
	.component-accordion-header div :global(div) {
		overflow: hidden;
		flex-wrap: nowrap;
		white-space: nowrap;
		text-overflow: ellipsis;
	}
</style>
