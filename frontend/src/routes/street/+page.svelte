<script lang="ts">
	import { onMount } from 'svelte';
	import { getPatrolPlans } from '$lib/api';
	import type { PatrolPlan } from '$lib/types';

	let plans: PatrolPlan[] = $state([]);
	let loading = $state(true);

	onMount(async () => {
		plans = await getPatrolPlans();
		loading = false;
	});
</script>

<div class="h-full overflow-y-auto bg-gray-950 p-4 text-white">
	<h1 class="mb-4 text-lg font-bold">Streetwork</h1>

	{#if loading}
		<p class="text-sm text-gray-500">Ładowanie...</p>
	{:else}
		<div class="grid grid-cols-1 gap-3">
			{#each plans as plan}
				<div class="rounded-xl border border-gray-800 bg-gray-900 p-4">
					<p class="font-medium">{plan.name}</p>
					<p class="text-sm text-gray-500">{plan.date}</p>
				</div>
			{/each}
		</div>
	{/if}
</div>
