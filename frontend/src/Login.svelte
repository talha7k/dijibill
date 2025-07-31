<script>
  import { createEventDispatcher } from 'svelte';
  import { Login } from '../wailsjs/go/main/App';

  const dispatch = createEventDispatcher();

  let username = '';
  let password = '';
  let loading = false;
  let error = '';

  async function handleLogin() {
    if (!username || !password) {
      error = 'Please enter both username and password';
      return;
    }

    loading = true;
    error = '';

    try {
      const authContext = await Login(username, password);
      dispatch('login-success', authContext);
    } catch (err) {
      error = err.message || 'Login failed';
    } finally {
      loading = false;
    }
  }

  function handleKeyPress(event) {
    if (event.key === 'Enter') {
      handleLogin();
    }
  }
</script>

<div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 to-indigo-100">
  <div class="max-w-md w-full space-y-8 p-8 bg-white rounded-xl shadow-lg">
    <div class="text-center">
      <h2 class="text-3xl font-bold text-gray-900 mb-2">DijiBill</h2>
      <p class="text-gray-600">Sign in to your account</p>
    </div>

    <form on:submit|preventDefault={handleLogin} class="space-y-6">
      {#if error}
        <div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg">
          {error}
        </div>
      {/if}

      <div>
        <label for="username" class="block text-sm font-medium text-gray-700 mb-2">
          Username
        </label>
        <input
          id="username"
          type="text"
          bind:value={username}
          on:keypress={handleKeyPress}
          disabled={loading}
          class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent disabled:bg-gray-100"
          placeholder="Enter your username"
          required
        />
      </div>

      <div>
        <label for="password" class="block text-sm font-medium text-gray-700 mb-2">
          Password
        </label>
        <input
          id="password"
          type="password"
          bind:value={password}
          on:keypress={handleKeyPress}
          disabled={loading}
          class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent disabled:bg-gray-100"
          placeholder="Enter your password"
          required
        />
      </div>

      <button
        type="submit"
        disabled={loading}
        class="w-full flex justify-center py-2 px-4 border border-transparent rounded-lg shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
      >
        {#if loading}
          <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          Signing in...
        {:else}
          Sign in
        {/if}
      </button>
    </form>

    <div class="text-center text-sm text-gray-600">
      <p>Default credentials:</p>
      <p class="font-mono text-xs bg-gray-100 p-2 rounded mt-1">
        Username: admin<br>
        Password: admin123
      </p>
    </div>
  </div>
</div>