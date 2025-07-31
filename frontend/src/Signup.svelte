<script>
  import { createEventDispatcher } from 'svelte';
  import { Signup } from '../wailsjs/go/main/App';

  const dispatch = createEventDispatcher();

  // User fields
  let username = '';
  let email = '';
  let password = '';
  let confirmPassword = '';
  let firstName = '';
  let lastName = '';

  // Company fields
  let companyName = '';
  let companyNameAr = '';
  let vatNumber = '';
  let crNumber = '';
  let companyEmail = '';
  let companyPhone = '';
  let companyAddress = '';
  let companyAddressAr = '';
  let city = '';
  let cityAr = '';
  let country = '';
  let countryAr = '';

  let loading = false;
  let error = '';
  let currentStep = 1;

  function nextStep() {
    if (currentStep === 1) {
      // Validate user fields
      if (!username || !email || !password || !confirmPassword || !firstName || !lastName) {
        error = 'Please fill in all user information fields';
        return;
      }
      if (password !== confirmPassword) {
        error = 'Passwords do not match';
        return;
      }
      if (password.length < 6) {
        error = 'Password must be at least 6 characters long';
        return;
      }
    }
    error = '';
    currentStep++;
  }

  function prevStep() {
    error = '';
    currentStep--;
  }

  async function handleSignup() {
    if (!companyName) {
      error = 'Company name is required';
      return;
    }

    loading = true;
    error = '';

    try {
      const signupRequest = {
        username,
        email,
        password,
        first_name: firstName,
        last_name: lastName,
        company_name: companyName,
        company_name_arabic: companyNameAr,
        vat_number: vatNumber,
        cr_number: crNumber,
        company_email: companyEmail,
        company_phone: companyPhone,
        company_address: companyAddress,
        company_address_arabic: companyAddressAr,
        city,
        city_arabic: cityAr,
        country,
        country_arabic: countryAr
      };

      const authContext = await Signup(signupRequest);
      dispatch('signup-success', authContext);
    } catch (err) {
      error = err.message || 'Signup failed';
    } finally {
      loading = false;
    }
  }

  function switchToLogin() {
    dispatch('switch-to-login');
  }
</script>

<div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 to-indigo-100">
  <div class="max-w-2xl w-full space-y-8 p-8 bg-white rounded-xl shadow-lg">
    <div class="text-center">
      <h2 class="text-3xl font-bold text-gray-900 mb-2">Create Your Account</h2>
      <p class="text-gray-600">Join DijiBill and start managing your invoices</p>
    </div>

    <!-- Progress indicator -->
    <div class="flex items-center justify-center space-x-4 mb-8">
      <div class="flex items-center">
        <div class="w-8 h-8 rounded-full {currentStep >= 1 ? 'bg-blue-600 text-white' : 'bg-gray-200 text-gray-600'} flex items-center justify-center text-sm font-medium">
          1
        </div>
        <span class="ml-2 text-sm {currentStep >= 1 ? 'text-blue-600' : 'text-gray-500'}">User Info</span>
      </div>
      <div class="w-16 h-0.5 {currentStep >= 2 ? 'bg-blue-600' : 'bg-gray-200'}"></div>
      <div class="flex items-center">
        <div class="w-8 h-8 rounded-full {currentStep >= 2 ? 'bg-blue-600 text-white' : 'bg-gray-200 text-gray-600'} flex items-center justify-center text-sm font-medium">
          2
        </div>
        <span class="ml-2 text-sm {currentStep >= 2 ? 'text-blue-600' : 'text-gray-500'}">Company Info</span>
      </div>
    </div>

    {#if error}
      <div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg">
        {error}
      </div>
    {/if}

    {#if currentStep === 1}
      <!-- Step 1: User Information -->
      <form on:submit|preventDefault={nextStep} class="space-y-6">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div>
            <label for="firstName" class="block text-sm font-medium text-gray-700 mb-2">
              First Name *
            </label>
            <input
              id="firstName"
              type="text"
              bind:value={firstName}
              disabled={loading}
              class="w-full px-3 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent disabled:bg-gray-100"
              placeholder="Enter your first name"
              required
            />
          </div>

          <div>
            <label for="lastName" class="block text-sm font-medium text-gray-700 mb-2">
              Last Name *
            </label>
            <input
              id="lastName"
              type="text"
              bind:value={lastName}
              disabled={loading}
              class="w-full px-3 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent disabled:bg-gray-100"
              placeholder="Enter your last name"
              required
            />
          </div>
        </div>

        <div>
          <label for="username" class="block text-sm font-medium text-gray-700 mb-2">
            Username *
          </label>
          <input
            id="username"
            type="text"
            bind:value={username}
            disabled={loading}
            class="w-full px-3 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent disabled:bg-gray-100"
            placeholder="Choose a username"
            required
          />
        </div>

        <div>
          <label for="email" class="block text-sm font-medium text-gray-700 mb-2">
            Email Address *
          </label>
          <input
            id="email"
            type="email"
            bind:value={email}
            disabled={loading}
            class="w-full px-3 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent disabled:bg-gray-100"
            placeholder="Enter your email address"
            required
          />
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div>
            <label for="password" class="block text-sm font-medium text-gray-700 mb-2">
              Password *
            </label>
            <input
              id="password"
              type="password"
              bind:value={password}
              disabled={loading}
              class="w-full px-3 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent disabled:bg-gray-100"
              placeholder="Create a password"
              required
            />
          </div>

          <div>
            <label for="confirmPassword" class="block text-sm font-medium text-gray-700 mb-2">
              Confirm Password *
            </label>
            <input
              id="confirmPassword"
              type="password"
              bind:value={confirmPassword}
              disabled={loading}
              class="w-full px-3 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent disabled:bg-gray-100"
              placeholder="Confirm your password"
              required
            />
          </div>
        </div>

        <div class="flex justify-between">
          <button
            type="button"
            on:click={switchToLogin}
            class="text-blue-600 hover:text-blue-700 text-sm font-medium"
          >
            Already have an account? Sign in
          </button>
          <button
            type="submit"
            disabled={loading}
            class="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50"
          >
            Next
          </button>
        </div>
      </form>
    {:else if currentStep === 2}
      <!-- Step 2: Company Information -->
      <form on:submit|preventDefault={handleSignup} class="space-y-6">
        <div>
          <label for="companyName" class="block text-sm font-medium text-gray-700 mb-2">
            Company Name *
          </label>
          <input
            id="companyName"
            type="text"
            bind:value={companyName}
            disabled={loading}
            class="w-full px-3 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent disabled:bg-gray-100"
            placeholder="Enter your company name"
            required
          />
        </div>

        <div>
          <label for="companyNameAr" class="block text-sm font-medium text-gray-700 mb-2">
            Company Name (Arabic)
          </label>
          <input
            id="companyNameAr"
            type="text"
            bind:value={companyNameAr}
            disabled={loading}
            class="w-full px-3 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent disabled:bg-gray-100"
            placeholder="اسم الشركة بالعربية"
            dir="rtl"
          />
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div>
            <label for="vatNumber" class="block text-sm font-medium text-gray-700 mb-2">
              VAT Number
            </label>
            <input
              id="vatNumber"
              type="text"
              bind:value={vatNumber}
              disabled={loading}
              class="w-full px-3 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent disabled:bg-gray-100"
              placeholder="Enter VAT number"
            />
          </div>

          <div>
            <label for="crNumber" class="block text-sm font-medium text-gray-700 mb-2">
              CR Number
            </label>
            <input
              id="crNumber"
              type="text"
              bind:value={crNumber}
              disabled={loading}
              class="w-full px-3 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent disabled:bg-gray-100"
              placeholder="Enter CR number"
            />
          </div>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div>
            <label for="companyEmail" class="block text-sm font-medium text-gray-700 mb-2">
              Company Email
            </label>
            <input
              id="companyEmail"
              type="email"
              bind:value={companyEmail}
              disabled={loading}
              class="w-full px-3 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent disabled:bg-gray-100"
              placeholder="company@example.com"
            />
          </div>

          <div>
            <label for="companyPhone" class="block text-sm font-medium text-gray-700 mb-2">
              Company Phone
            </label>
            <input
              id="companyPhone"
              type="tel"
              bind:value={companyPhone}
              disabled={loading}
              class="w-full px-3 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent disabled:bg-gray-100"
              placeholder="+966 XX XXX XXXX"
            />
          </div>
        </div>

        <div>
          <label for="companyAddress" class="block text-sm font-medium text-gray-700 mb-2">
            Company Address
          </label>
          <textarea
            id="companyAddress"
            bind:value={companyAddress}
            disabled={loading}
            rows="3"
            class="w-full px-3 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent disabled:bg-gray-100"
            placeholder="Enter company address"
          ></textarea>
        </div>

        <div>
          <label for="companyAddressAr" class="block text-sm font-medium text-gray-700 mb-2">
            Company Address (Arabic)
          </label>
          <textarea
            id="companyAddressAr"
            bind:value={companyAddressAr}
            disabled={loading}
            rows="3"
            class="w-full px-3 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent disabled:bg-gray-100"
            placeholder="عنوان الشركة بالعربية"
            dir="rtl"
          ></textarea>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div>
            <label for="city" class="block text-sm font-medium text-gray-700 mb-2">
              City
            </label>
            <input
              id="city"
              type="text"
              bind:value={city}
              disabled={loading}
              class="w-full px-3 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent disabled:bg-gray-100"
              placeholder="Enter city"
            />
          </div>

          <div>
            <label for="cityAr" class="block text-sm font-medium text-gray-700 mb-2">
              City (Arabic)
            </label>
            <input
              id="cityAr"
              type="text"
              bind:value={cityAr}
              disabled={loading}
              class="w-full px-3 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent disabled:bg-gray-100"
              placeholder="المدينة بالعربية"
              dir="rtl"
            />
          </div>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div>
            <label for="country" class="block text-sm font-medium text-gray-700 mb-2">
              Country
            </label>
            <input
              id="country"
              type="text"
              bind:value={country}
              disabled={loading}
              class="w-full px-3 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent disabled:bg-gray-100"
              placeholder="Enter country"
            />
          </div>

          <div>
            <label for="countryAr" class="block text-sm font-medium text-gray-700 mb-2">
              Country (Arabic)
            </label>
            <input
              id="countryAr"
              type="text"
              bind:value={countryAr}
              disabled={loading}
              class="w-full px-3 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent disabled:bg-gray-100"
              placeholder="البلد بالعربية"
              dir="rtl"
            />
          </div>
        </div>

        <div class="flex justify-between">
          <button
            type="button"
            on:click={prevStep}
            disabled={loading}
            class="px-6 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-gray-500 disabled:opacity-50"
          >
            Back
          </button>
          <button
            type="submit"
            disabled={loading}
            class="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50 flex items-center"
          >
            {#if loading}
              <i class="fas fa-spinner fa-spin -ml-1 mr-3 h-5 w-5 text-white"></i>
              Creating Account...
            {:else}
              Create Account
            {/if}
          </button>
        </div>
      </form>
    {/if}
  </div>
</div>