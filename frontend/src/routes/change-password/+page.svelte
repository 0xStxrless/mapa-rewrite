<script lang="ts">
  import { goto } from '$app/navigation'
  import { changePassword } from '$lib/api'

  let oldPassword = $state('')
  let newPassword = $state('')
  let error = $state('')
  let loading = $state(false)

  async function handleChange() {
    error = ''
    loading = true
    try {
      await changePassword({ old_password: oldPassword, new_password: newPassword })
      goto('/')
    } catch (e) {
      error = 'Nie udało się zmienić hasła'
    } finally {
      loading = false
    }
  }
</script>

<div class="min-h-screen bg-gray-950 flex items-center justify-center px-4">
  <div class="w-full max-w-sm">
    <h1 class="text-white text-2xl font-bold mb-2 text-center tracking-widest">PUNKT</h1>
    <p class="text-gray-500 text-sm text-center mb-8">Musisz zmienić hasło przed kontynuowaniem</p>

    <div class="bg-gray-900 rounded-2xl p-6 border border-gray-800 flex flex-col gap-4">
      <div class="flex flex-col gap-1.5">
        <label class="text-gray-400 text-xs" for="old">Obecne hasło</label>
        <input
          id="old"
          type="password"
          bind:value={oldPassword}
          class="bg-gray-800 text-white rounded-lg px-4 py-2.5 text-sm outline-none focus:ring-2 focus:ring-blue-500"
        />
      </div>

      <div class="flex flex-col gap-1.5">
        <label class="text-gray-400 text-xs" for="new">Nowe hasło</label>
        <input
          id="new"
          type="password"
          bind:value={newPassword}
          class="bg-gray-800 text-white rounded-lg px-4 py-2.5 text-sm outline-none focus:ring-2 focus:ring-blue-500"
        />
      </div>

      {#if error}
        <p class="text-red-400 text-xs">{error}</p>
      {/if}

      <button
        type="button"
        onclick={handleChange}
        disabled={loading}
        class="w-full py-2.5 rounded-xl bg-blue-600 hover:bg-blue-500 disabled:opacity-50 text-white text-sm font-medium transition-colors"
      >
        {loading ? 'Zapisywanie...' : 'Zmień hasło'}
      </button>
    </div>
  </div>
</div>