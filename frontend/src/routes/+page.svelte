<script lang="ts">
	import { onMount } from 'svelte';
	import Map from '$lib/Map.svelte';
	import { getPins } from '$lib/api';
	import type { Pin } from '$lib/types';
	let pins: Pin[] = $state([]);
	let showList = $state(false);

	let pinListEl: HTMLDivElement;

	onMount(async () => {
		try {
			pins = await getPins();
		} catch (e) {
			console.error(e);
		}
	});
	function skipToList() {
		showList = true;
		setTimeout(() => pinListEl?.focus(), 50);
	}
</script>

<button
	type="button"
	onclick={skipToList}
	class="sr-only text-sm font-medium focus:not-sr-only focus:absolute focus:top-2 focus:left-2 focus:z-50 focus:rounded-lg focus:bg-blue-600 focus:px-4 focus:py-2 focus:text-white"
>
	Pomiń mapę i przejdź do listy punktów
</button>
<div
	class="absolute inset-0"
	role="application"
	aria-label="Mapa punktów streetworkowych"
	aria-describedby="map-desc"
>
	<p id="map-desc" class="sr-only">
		Interaktywna mapa z punktami streetworkowymi. Użyj listy poniżej mapy, aby przeglądać punkty za
		pomocą klawiatury.
	</p>
	<Map />
</div>
<div id="pin-list" tabindex="-1" class="absolute right-4 bottom-4 z-50" bind:this={pinListEl}>
	<button
		type="button"
		onclick={() => (showList = !showList)}
		class="rounded-xl border border-gray-700 bg-gray-900 px-4 py-2 text-sm text-white shadow-lg outline-none focus:ring-2 focus:ring-blue-500"
		aria-expanded={showList}
		aria-controls="pin-list-items"
	>
		{showList ? 'Ukryj listę punktów' : 'Pokaż listę punktów'}
	</button>
	{#if showList}
		<div
			id="pin-list-items"
			class="absolute right-0 bottom-12 max-h-96 w-72 overflow-y-auto rounded-2xl border border-gray-700 bg-gray-900 shadow-2xl"
			role="region"
			aria-label="Lista wszystkich punktów"
		>
			{#if pins.length === 0}
				<p class="p-4 text-sm text-gray-500">Ładowanie punktów...</p>
			{:else}
				<ul role="list" class="flex flex-col gap-1 p-2">
					{#each pins as pin}
						<li role="listitem">
							<a
								href="/pin/{pin.id}"
								class="flex flex-col rounded-xl px-3 py-2.5 transition-colors hover:bg-gray-800 focus:bg-gray-800 focus:ring-2 focus:ring-blue-500 focus:outline-none"
							>
								<span class="text-sm font-medium text-white">{pin.title}</span>
								<span class="mt-0.5 text-xs text-gray-400">
									<span class="sr-only">Kategoria:</span>{pin.category}{pin.description
										? ` — ${pin.description}`
										: ''}
								</span>
							</a>
						</li>
					{/each}
				</ul>
			{/if}
		</div>
	{/if}
</div>
