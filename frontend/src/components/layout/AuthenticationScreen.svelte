<script>
  import { createEventDispatcher } from 'svelte';
  import Login from '../../Login.svelte';
  import Signup from '../../Signup.svelte';
  
  export let authMode = 'login';
  
  const dispatch = createEventDispatcher();
  
  function handleLoginSuccess(event) {
    dispatch('login-success', event.detail);
  }
  
  function handleSignupSuccess(event) {
    dispatch('signup-success', event.detail);
  }
  
  function switchToSignup() {
    dispatch('switch-to-signup');
  }
  
  function switchToLogin() {
    dispatch('switch-to-login');
  }
</script>

{#if authMode === 'login'}
  <Login 
    on:login-success={handleLoginSuccess} 
    on:switch-to-signup={switchToSignup}
  />
{:else if authMode === 'signup'}
  <Signup 
    on:signup-success={handleSignupSuccess}
    on:switch-to-login={switchToLogin}
  />
{/if}