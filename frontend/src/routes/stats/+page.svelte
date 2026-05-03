<script lang="ts">
	import { onMount } from 'svelte';
	import { getAllStats } from '$lib/api';
	import type { StreetworkStat } from '$lib/types';

	let stats: StreetworkStat[] = $state([]);
	let loading = $state(true);

	onMount(async () => {
		stats = await getAllStats();
		loading = false;
	});
</script>

<div class="h-full overflow-y-auto bg-gray-950 p-4 text-white">
	<h1 class="mb-4 text-lg font-bold">Statystyki</h1>

	{#if loading}
		<p class="text-sm text-gray-500">Ładowanie...</p>
	{:else}
		<div class="grid grid-cols-1 gap-3">
			{#each stats as stat}
				<div class="rounded-xl border border-gray-800 bg-gray-900 p-4">
					<div class="mb-3 flex items-center gap-3">
						{#if stat.avatar}
							<img src={stat.avatar} alt={stat.worker_name} class="h-8 w-8 rounded-full" />
						{:else}
							<div
								class="flex h-8 w-8 items-center justify-center rounded-full text-xs font-bold"
								style="background: {stat.bg_color ?? '#374151'}"
							>
								{stat.worker_name[0]}
							</div>
						{/if}
						<div>
							<p class="text-sm font-medium">{stat.worker_name}</p>
							<p class="text-xs text-gray-500">{stat.month}</p>
						</div>
					</div>
					<div class="grid grid-cols-3 gap-2 text-center">
						<div class="rounded-lg bg-gray-800 p-2">
							<p class="text-lg font-bold">{stat.interactions}</p>
							<p class="text-xs text-gray-400">Kontakty</p>
						</div>
						<div class="rounded-lg bg-gray-800 p-2">
							<p class="text-lg font-bold">{stat.new_contacts}</p>
							<p class="text-xs text-gray-400">Nowe</p>
						</div>
						<div class="rounded-lg bg-gray-800 p-2">
							<p class="text-lg font-bold">{stat.interventions}</p>
							<p class="text-xs text-gray-400">Interwencje</p>
						</div>
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>
