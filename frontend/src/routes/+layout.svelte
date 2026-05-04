<script lang="ts">
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import { auth } from '$lib/auth';
	import { onMount } from 'svelte';
	import { showPinList } from '$lib/stores';
	import '../app.css';

	let { children } = $props();

	let showSearch = $state(false);
	let showAdd = $state(false);

	onMount(() => {
		if (!$auth && page.url.pathname !== '/login' && page.url.pathname !== '/change-password') {
			goto('/login');
		}
	});
</script>

<div class="flex h-screen w-screen flex-col overflow-hidden bg-gray-950">
	<button
		type="button"
		onclick={() => {
			showPinList.set(true);
			setTimeout(() => document.getElementById('pin-list')?.focus(), 50);
		}}
		class="skip-link"
	>
		Pomiń mapę i przejdź do listy punktów
	</button>

	<!-- Navbar -->
	<nav
		class="z-50 flex shrink-0 items-center justify-between border-b border-gray-800 bg-gray-950 px-4 py-2"
	>
		<div class="flex items-center gap-2">
			<a href="/" class="text-sm font-bold tracking-widest text-white">PUNKT</a>

			<button
				type="button"
				onclick={() => (showSearch = true)}
				class="flex items-center gap-1.5 rounded-lg bg-gray-800 px-3 py-1.5 text-sm text-gray-300 transition-colors hover:bg-gray-700"
			>
				<svg class="h-3.5 w-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
					></path>
				</svg>
				Szukaj
			</button>

			<button
				type="button"
				onclick={() => (showAdd = true)}
				class="flex items-center gap-1.5 rounded-lg bg-gray-800 px-3 py-1.5 text-sm text-gray-300 transition-colors hover:bg-gray-700"
			>
				<svg class="h-3.5 w-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"
					></path>
				</svg>
				Dodaj
			</button>
		</div>

		<div class="flex items-center gap-2">
			<a
				href="/stats"
				class="rounded-lg px-4 py-1.5 text-sm font-medium transition-colors {page.url.pathname ===
				'/stats'
					? 'bg-blue-600 text-white'
					: 'bg-gray-800 text-gray-300 hover:bg-gray-700'}"
			>
				Statystyki
			</a>

			<a
				href="/street"
				class="rounded-lg px-4 py-1.5 text-sm font-medium transition-colors {page.url.pathname ===
				'/street'
					? 'bg-green-600 text-white'
					: 'bg-gray-800 text-gray-300 hover:bg-gray-700'}"
			>
				Streetwork
			</a>
			<span class="text-xs text-gray-600">© OSM</span>
		</div>
	</nav>

	<!-- Page content -->
	<div class="relative flex-1 overflow-hidden">
		{@render children()}
	</div>

	<!-- Search modal -->
	{#if showSearch}
		<dialog
			open
			class="fixed inset-0 z-50 m-0 flex h-full w-full max-w-none items-start justify-center border-none bg-black/60 px-4 pt-20"
			onclose={() => (showSearch = false)}
		>
			<div class="w-full max-w-md rounded-2xl border border-gray-700 bg-gray-900 shadow-2xl">
				<div class="flex gap-2 border-b border-gray-800 p-4">
					<input
						type="text"
						placeholder="Szukaj pinu..."
						class="w-full rounded-lg bg-gray-800 px-4 py-2.5 text-sm text-white placeholder-gray-500 outline-none focus:ring-2 focus:ring-blue-500"
					/>
					<button
						type="button"
						onclick={() => (showSearch = false)}
						class="px-2 text-gray-500 hover:text-white">✕</button
					>
				</div>
				<div class="max-h-80 overflow-y-auto p-2">
					<p class="py-8 text-center text-sm text-gray-500">Zacznij pisać aby szukać...</p>
				</div>
			</div>
		</dialog>
	{/if}

	<!-- Add pin modal -->
	{#if showAdd}
		<dialog
			open
			class="fixed inset-0 z-50 m-0 flex h-full w-full max-w-none items-end justify-center border-none bg-black/60"
			onclose={() => (showAdd = false)}
		>
			<div
				class="w-full max-w-lg rounded-t-2xl border-t border-gray-700 bg-gray-900 p-6 shadow-2xl"
			>
				<h2 class="mb-4 font-semibold text-white">Dodaj pin</h2>
				<p class="text-sm text-gray-400">Tapnij na mapę żeby wybrać lokalizację...</p>
				<button
					type="button"
					onclick={() => (showAdd = false)}
					class="mt-4 w-full rounded-xl bg-gray-800 py-2.5 text-sm text-gray-300"
				>
					Anuluj
				</button>
			</div>
		</dialog>
	{/if}
</div>
