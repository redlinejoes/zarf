<script lang="ts">
	import { Typography, TextField, ListItem, ListItemAdornment } from '@ui';
	import { Packages } from '$lib/api';
	import ZarfDialog from '../zarf-dialog.svelte';
	import ButtonDense from '../button-dense.svelte';
	import Spinner from '../spinner.svelte';
	import type { DeployedPackage } from '$lib/api-types';

	export let pkg: DeployedPackage;

	export let toggleDialog: () => void;
	let inputValue: string;
	let errorMessage: string;
	let happyZarf: boolean;
	let titleText: string;
	let zarfAlt: string;

	let removing = false;

	function closeDialog() {
		inputValue = '';
		errorMessage = '';
		removing = false;
		toggleDialog();
	}

	async function removePkg(): Promise<void> {
		removing = true;
		Packages.remove(pkg.name)
			.then(closeDialog)
			.catch((e) => {
				errorMessage = e.message;
				removing = false;
			});
	}

	$: {
		if (!errorMessage) {
			happyZarf = true;
			titleText = 'Remove Package from Cluster';
			zarfAlt = 'Succeeded in removing a package from the cluster.';
		} else {
			happyZarf = false;
			titleText = `Failed to remove package ${pkg.name}`;
			zarfAlt = 'Failed to remove a package from the cluster.';
		}
	}
</script>

<ZarfDialog clickAway={!removing} bind:toggleDialog {happyZarf} {titleText} {zarfAlt}>
	{#if !errorMessage}
		{#if removing}
			<div style="display: flex; justify-content: center; align-items: center: width: 100%">
				<Spinner color="blue-200" diameter="50px" />
			</div>
		{:else}
			<Typography variant="body2" color="text-secondary-on-dark">
				Type the name of the package and click remove to delete the package and all of it’s
				resources from the cluster. This action step cannot be undone.
			</Typography>
			<Typography variant="subtitle1">{pkg.name}</Typography>
			<TextField
				variant="outlined"
				label="Package to Delete"
				color="primary"
				bind:value={inputValue}
				helperText="Type the name of the package."
			/>
		{/if}
	{:else}
		<Typography variant="body2" color="text-secondary-on-dark">
			{errorMessage}
		</Typography>
	{/if}
	<svelte:fragment slot="actions">
		{#if !errorMessage}
			<ButtonDense backgroundColor="white" variant="outlined" on:click={closeDialog}>
				Cancel
			</ButtonDense>
			<ButtonDense
				disabled={pkg.name !== inputValue}
				variant="flat"
				textColor="text-primary-on-light"
				backgroundColor="grey-300"
				on:click={removePkg}
			>
				Remove Package
			</ButtonDense>
		{:else}
			<ButtonDense backgroundColor="white" variant="outlined" on:click={closeDialog}>
				Close
			</ButtonDense>
		{/if}
	</svelte:fragment>
</ZarfDialog>
