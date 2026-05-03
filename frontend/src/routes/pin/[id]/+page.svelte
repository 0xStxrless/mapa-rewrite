<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import {
		getPin,
		getVisitsByPin,
		getWorkers,
		getCategories,
		createVisit,
		updateVisit,
		deleteVisit,
		deletePin,
		updatePin
	} from '$lib/api';
	import type { Pin, Visit, Category } from '$lib/types';

	const id = Number(page.params.id);

	let pin: Pin | null = $state(null);
	let visits: Visit[] = $state([]);
	let workers: string[] = $state([]);
	let categories: Category[] = $state([]);
	let loading = $state(true);
	let categoryColor = $state('#888');

	let workersInput = $state('');

	// add visit
	let selectedWorkers: string[] = $state([]);
	let note = $state('');
	let backdate = $state('');
	let showBackdate = $state(false);
	let submitting = $state(false);

	// edit visit
	let editingVisit: Visit | null = $state(null);
	let editNote = $state('');
	let editName = $state('');

	// edit pin
	let showEditPin = $state(false);
	let editTitle = $state('');
	let editDescription = $state('');
	let editCategory = $state('');
	let editingPin = $state(false);

	// delete
	let confirmDelete = $state(false);

	const quickNotes = ['Nikogo', 'Pusto', 'Brak aktywności', 'Ślady, nikt'];

	onMount(async () => {
		try {
			const [p, v, w, cats] = await Promise.all([
				getPin(id),
				getVisitsByPin(id),
				getWorkers(),
				getCategories()
			]);
			pin = p;
			visits = v ?? [];
			workers = w ?? [];
			categories = cats ?? [];
			const cat = cats.find((c) => c.name === p.category);
			categoryColor = cat?.color ?? '#888';

			editTitle = p.title;
			editDescription = p.description ?? '';
			editCategory = p.category;
		} catch (e) {
			console.error(e);
		} finally {
			loading = false;
		}
	});

	function toggleWorker(w: string) {
		const current = workersInput
			.split(',')
			.map((s) => s.trim())
			.filter(Boolean);
		if (current.includes(w)) {
			workersInput = current.filter((x) => x !== w).join(', ');
		} else {
			workersInput = [...current, w].join(', ');
		}
	}

	function appendQuickNote(n: string) {
		note = note ? `${note}. ${n}` : n;
	}

	async function submitVisit() {
		if (!workersInput.trim() || !pin) return;
		submitting = true;
		try {
			const newVisit = await createVisit({
				pin_id: pin.id,
				name: workersInput.trim(),
				note: note || undefined,
				...(backdate ? { visited_at: new Date(backdate).toISOString() } : {})
			});
			visits = [newVisit, ...visits];
			workersInput = '';
			note = '';
			backdate = '';
			showBackdate = false;
		} finally {
			submitting = false;
		}
	}

	function isWorkerSelected(w: string) {
		return workersInput
			.split(',')
			.map((s) => s.trim())
			.includes(w);
	}

	async function handleDeletePin() {
		if (!pin) return;
		await deletePin(pin.id);
		goto('/');
	}

	function startEditVisit(v: Visit) {
		editingVisit = v;
		editNote = v.note ?? '';
		editName = v.name;
	}

	async function submitEditVisit() {
		if (!editingVisit) return;
		const updated = await updateVisit({ ...editingVisit, note: editNote, name: editName });
		visits = visits.map((v) => (v.id === updated.id ? updated : v));
		editingVisit = null;
	}

	async function handleDeleteVisit(visitId: number) {
		await deleteVisit(visitId);
		visits = visits.filter((v) => v.id !== visitId);
	}

	async function submitEditPin() {
		if (!pin) return;
		editingPin = true;
		try {
			const updated = await updatePin(pin.id, {
				title: editTitle,
				description: editDescription,
				category: editCategory,
				lat: pin.lat,
				lng: pin.lng,
				image_url: pin.image_url ?? undefined
			});
			pin = updated;
			const cat = categories.find((c) => c.name === updated.category);
			categoryColor = cat?.color ?? '#888';
			showEditPin = false;
		} finally {
			editingPin = false;
		}
	}

	function formatDate(d: string) {
		return new Date(d).toLocaleString('pl-PL', {
			day: '2-digit',
			month: '2-digit',
			year: '2-digit',
			hour: '2-digit',
			minute: '2-digit'
		});
	}
</script>

<div class="min-h-screen bg-gray-950 text-white">
	<!-- Topbar -->
	<div
		class="sticky top-0 z-10 flex items-center gap-3 border-b border-gray-800 bg-gray-950 px-4 py-3"
	>
		<button
			type="button"
			onclick={() => goto('/')}
			class="text-gray-400 transition-colors hover:text-white"
		>
			<svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"
				></path>
			</svg>
		</button>
		<h1 class="flex-1 truncate text-base font-semibold text-white">{pin?.title ?? '...'}</h1>
		<span class="text-sm text-gray-500">{visits.length} wizyt</span>
	</div>

	{#if loading}
		<div class="flex h-40 items-center justify-center">
			<p class="text-sm text-gray-500">Ładowanie...</p>
		</div>
	{:else if pin}
		<!-- Image -->
		{#if pin.image_url}
			<img src={pin.image_url} alt={pin.title} class="h-48 w-full object-cover" />
		{/if}

		<!-- Pin info -->
		<div class="border-b border-gray-800 px-4 py-4">
			<div class="mb-1 flex items-center gap-2">
				<span class="h-2.5 w-2.5 shrink-0 rounded-full" style="background: {categoryColor}"></span>
				<span class="text-sm text-gray-400">{pin.category}</span>
			</div>
			{#if pin.description}
				<p class="text-sm text-gray-300">{pin.description}</p>
			{/if}

			<!-- Actions -->
			<div class="mt-4 flex gap-2">
				{#if confirmDelete}
					<button
						type="button"
						onclick={handleDeletePin}
						class="flex-1 rounded-xl bg-red-600 py-2.5 text-sm font-medium text-white"
						>Na pewno usuń</button
					>
					<button
						type="button"
						onclick={() => (confirmDelete = false)}
						class="flex-1 rounded-xl bg-gray-800 py-2.5 text-sm text-gray-300">Anuluj</button
					>
				{:else}
					<button
						type="button"
						onclick={() => (confirmDelete = true)}
						class="flex-1 rounded-xl bg-red-600/80 py-2.5 text-sm font-medium text-white"
						>Usuń</button
					>
					<button
						type="button"
						onclick={() => (showEditPin = !showEditPin)}
						class="flex-1 rounded-xl bg-gray-700 py-2.5 text-sm font-medium text-white"
						>Edytuj</button
					>
				{/if}
			</div>
			<button
				type="button"
				class="mt-2 w-full rounded-xl bg-blue-600 py-2.5 text-sm font-medium text-white transition-colors hover:bg-blue-500"
				>+ Dodaj na patrol</button
			>
		</div>

		<!-- Edit pin form -->
		{#if showEditPin}
			<div class="border-b border-gray-800 bg-gray-900 px-4 py-4">
				<p class="mb-3 text-xs font-semibold tracking-wider text-gray-500 uppercase">Edytuj pin</p>
				<div class="flex flex-col gap-3">
					<input
						type="text"
						bind:value={editTitle}
						placeholder="Nazwa"
						class="rounded-xl bg-gray-800 px-4 py-2.5 text-sm text-white outline-none focus:ring-2 focus:ring-blue-500"
					/>
					<textarea
						bind:value={editDescription}
						placeholder="Opis"
						rows="2"
						class="resize-none rounded-xl bg-gray-800 px-4 py-2.5 text-sm text-white outline-none focus:ring-2 focus:ring-blue-500"
					></textarea>
					<select
						bind:value={editCategory}
						class="rounded-xl bg-gray-800 px-4 py-2.5 text-sm text-white outline-none focus:ring-2 focus:ring-blue-500"
					>
						{#each categories as cat}
							<option value={cat.name}>{cat.name}</option>
						{/each}
					</select>
					<div class="flex gap-2">
						<button
							type="button"
							onclick={submitEditPin}
							disabled={editingPin}
							class="flex-1 rounded-xl bg-blue-600 py-2.5 text-sm font-medium text-white hover:bg-blue-500 disabled:opacity-50"
						>
							{editingPin ? 'Zapisywanie...' : 'Zapisz'}
						</button>
						<button
							type="button"
							onclick={() => (showEditPin = false)}
							class="flex-1 rounded-xl bg-gray-800 py-2.5 text-sm text-gray-300">Anuluj</button
						>
					</div>
				</div>
			</div>
		{/if}

		<!-- Visit history -->
		<div class="border-b border-gray-800 px-4 py-4">
			<p class="mb-3 text-xs font-semibold tracking-wider text-gray-500 uppercase">
				Historia odwiedzin
			</p>
			{#if visits.length === 0}
				<p class="text-sm text-gray-600">Brak odwiedzin</p>
			{:else}
				<div class="flex flex-col gap-2">
					{#each visits as visit}
						{#if editingVisit?.id === visit.id}
							<div class="flex flex-col gap-2 rounded-xl bg-gray-800 p-3">
								<input
									type="text"
									bind:value={editName}
									class="rounded-lg bg-gray-700 px-3 py-1.5 text-sm text-white outline-none"
									placeholder="Pracownicy"
								/>
								<textarea
									bind:value={editNote}
									class="resize-none rounded-xl bg-gray-700 px-3 py-1.5 text-sm text-white outline-none"
									rows="2"
									placeholder="Notatka"
								></textarea>
								<div class="flex gap-2">
									<button
										type="button"
										onclick={submitEditVisit}
										class="flex-1 rounded-lg bg-blue-600 py-1.5 text-sm text-white">Zapisz</button
									>
									<button
										type="button"
										onclick={() => (editingVisit = null)}
										class="flex-1 rounded-lg bg-gray-700 py-1.5 text-sm text-gray-300"
										>Anuluj</button
									>
								</div>
							</div>
						{:else}
							<div class="rounded-xl bg-gray-800 p-3">
								<div class="flex items-start justify-between gap-2">
									<div class="flex-1">
										<p class="text-sm text-white">
											<span class="font-medium">{visit.name}</span>{#if visit.note}
												– {visit.note}{/if}
										</p>
										<p class="mt-0.5 text-xs text-gray-500">{formatDate(visit.visited_at)}</p>
									</div>
									<div class="flex shrink-0 gap-1">
										<button
											type="button"
											onclick={() => startEditVisit(visit)}
											class="p-1 text-gray-500 transition-colors hover:text-white">✎</button
										>
										<button
											type="button"
											onclick={() => handleDeleteVisit(visit.id)}
											class="p-1 text-gray-500 transition-colors hover:text-red-400">✕</button
										>
									</div>
								</div>
							</div>
						{/if}
					{/each}
				</div>
			{/if}
		</div>

		<!-- Add visit -->
		<div class="px-4 py-4">
			<p class="mb-3 text-xs font-semibold tracking-wider text-gray-500 uppercase">
				Dodaj aktualizację
			</p>

			<div class="mb-3 flex flex-wrap gap-2"></div>
			<div class="mb-3 flex flex-wrap gap-2">
				{#each workers as w}
					<button
						type="button"
						onclick={() => toggleWorker(w)}
						class="rounded-xl px-3 py-1.5 text-sm font-medium transition-colors {isWorkerSelected(w)
							? 'bg-blue-600 text-white'
							: 'bg-gray-800 text-gray-300 hover:bg-gray-700'}">{w}</button
					>
				{/each}
			</div>

			<textarea
				bind:value={workersInput}
				class="mb-2 w-full resize-none overflow-hidden rounded-xl bg-gray-800 px-4 py-2.5 text-sm text-white placeholder-gray-500 outline-none focus:ring-2 focus:ring-blue-500"
				rows="1"
				oninput={(e) => {
					const t = e.currentTarget;
					t.style.height = 'auto';
					t.style.height = t.scrollHeight + 'px';
				}}
				placeholder="Pracownicy (opcjonalnie)"
			></textarea>

			<textarea
				bind:value={note}
				class="mb-2 w-full resize-none rounded-xl bg-gray-800 px-4 py-2.5 text-sm text-white placeholder-gray-500 outline-none focus:ring-2 focus:ring-blue-500"
				rows="2"
				placeholder="Notatka (opcjonalnie)"
			></textarea>
			<div class="mb-3 flex flex-wrap gap-2">
				{#each quickNotes as qn}
					<button
						type="button"
						onclick={() => appendQuickNote(qn)}
						class="rounded-xl bg-gray-800 px-3 py-1.5 text-sm text-gray-300 transition-colors hover:bg-gray-700"
						>{qn}</button
					>
				{/each}
			</div>

			<button
				type="button"
				onclick={() => (showBackdate = !showBackdate)}
				class="mb-2 w-full rounded-xl bg-gray-800 py-2 text-sm text-gray-300 transition-colors hover:bg-gray-700"
				>Dodaj wizytę wstecz</button
			>

			{#if showBackdate}
				<input
					type="datetime-local"
					bind:value={backdate}
					class="mb-2 w-full rounded-xl bg-gray-800 px-4 py-2.5 text-sm text-white outline-none focus:ring-2 focus:ring-blue-500"
				/>
			{/if}

			<button
				type="button"
				class="mb-3 w-full rounded-xl bg-gray-800 py-2 text-sm text-gray-300 transition-colors hover:bg-gray-700"
				>Dodaj zdjęcie</button
			>

			<button
				type="button"
				onclick={submitVisit}
				disabled={submitting || workersInput.trim() === ''}
				class="w-full rounded-xl bg-blue-600 py-3 text-sm font-semibold text-white transition-colors hover:bg-blue-500 disabled:opacity-40"
				>{submitting ? 'Zapisywanie...' : 'Dodaj aktualizację'}</button
			>
		</div>
	{/if}
</div>
