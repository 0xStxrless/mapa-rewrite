<script lang="ts">
  import { goto } from '$app/navigation'
  import { login } from '$lib/auth'

  let email = $state('')
  let password = $state('')
  let error = $state('')
  let loading = $state(false)

  async function handleLogin() {
    error = ''
    loading = true
    try {
      const data = await login(email, password)
      if (data.must_change_password) {
        goto('/change-password')
      } else {
        goto('/')
      }
    } catch (e) {
      error = 'Nieprawidłowy email lub hasło'
    } finally {
      loading = false
    }
  }
</script>

<div class="min-h-screen bg-gray-950 flex items-center justify-center px-4">
  <div class="w-full max-w-sm">
    <h1 class="text-white text-2xl font-bold mb-8 text-center tracking-widest">PUNKT</h1>

    <div class="bg-gray-900 rounded-2xl p-6 border border-gray-800 flex flex-col gap-4">
      <div class="flex flex-col gap-1.5">
        <label class="text-gray-400 text-xs" for="email">Email</label>
        <input
          id="email"
          type="email"
          bind:value={email}
          class="bg-gray-800 text-white rounded-lg px-4 py-2.5 text-sm outline-none focus:ring-2 focus:ring-blue-500"
          placeholder="worker@example.com"
        />
      </div>

      <div class="flex flex-col gap-1.5">
        <label class="text-gray-400 text-xs" for="password">Hasło</label>
        <input
          id="password"
          type="password"
          bind:value={password}
          onkeydown={(e) => e.key === 'Enter' && handleLogin()}
          class="bg-gray-800 text-white rounded-lg px-4 py-2.5 text-sm outline-none focus:ring-2 focus:ring-blue-500"
          placeholder="••••••••"
        />
      </div>

      {#if error}
        <p class="text-red-400 text-xs">{error}</p>
      {/if}

      <button
        type="button"
        onclick={handleLogin}
        disabled={loading}
        class="w-full py-2.5 rounded-xl bg-blue-600 hover:bg-blue-500 disabled:opacity-50 text-white text-sm font-medium transition-colors"
      >
        {loading ? 'Logowanie...' : 'Zaloguj się'}
      </button>
    </div>
  </div>
</div>