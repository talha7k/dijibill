<script>
  import FormField from './FormField.svelte'
  
  export let customerForm = {
    name: '',
    name_arabic: '',
    vat_number: '',
    email: '',
    phone: '',
    address: '',
    address_arabic: '',
    city: '',
    city_arabic: '',
    country: '',
    country_arabic: ''
  }
  
  export let errors = {}

  // Country options
  const countryOptions = [
    { value: 'SA', label: 'Saudi Arabia' },
    { value: 'AE', label: 'United Arab Emirates' },
    { value: 'KW', label: 'Kuwait' },
    { value: 'QA', label: 'Qatar' },
    { value: 'BH', label: 'Bahrain' },
    { value: 'OM', label: 'Oman' },
    { value: 'JO', label: 'Jordan' },
    { value: 'LB', label: 'Lebanon' },
    { value: 'EG', label: 'Egypt' },
    { value: 'US', label: 'United States' },
    { value: 'GB', label: 'United Kingdom' },
    { value: 'CA', label: 'Canada' },
    { value: 'AU', label: 'Australia' }
  ]

  // Set default country to Saudi Arabia
  if (!customerForm.country) {
    customerForm.country = 'SA'
    customerForm.country_arabic = 'المملكة العربية السعودية'
  }

  // Update Arabic country name when country changes
  function updateCountryArabic(country) {
    const arabicNames = {
      'SA': 'المملكة العربية السعودية',
      'AE': 'الإمارات العربية المتحدة',
      'KW': 'الكويت',
      'QA': 'قطر',
      'BH': 'البحرين',
      'OM': 'عمان',
      'JO': 'الأردن',
      'LB': 'لبنان',
      'EG': 'مصر',
      'US': 'الولايات المتحدة',
      'GB': 'المملكة المتحدة',
      'CA': 'كندا',
      'AU': 'أستراليا'
    }
    customerForm.country_arabic = arabicNames[country] || ''
  }

  $: if (customerForm.country) {
    updateCountryArabic(customerForm.country)
  }
</script>

<div class="space-y-6">
  <!-- Basic Information -->
  <div class="glass-card">
    <div class="p-6">
      <h3 class="heading-3 mb-6">Basic Information</h3>
      
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <FormField
          label="Customer Name"
          labelArabic="اسم العميل"
          type="text"
          bind:value={customerForm.name}
          placeholder="Enter customer name"
          required={true}
          error={errors.name}
          id="customer-name"
          name="name"
        />

        <FormField
          label="Customer Name (Arabic)"
          labelArabic="اسم العميل بالعربية"
          type="text"
          bind:value={customerForm.name_arabic}
          placeholder="أدخل اسم العميل"
          dir="rtl"
          error={errors.name_arabic}
          id="customer-name-arabic"
          name="name_arabic"
        />

        <FormField
          label="VAT Number"
          labelArabic="الرقم الضريبي"
          type="text"
          bind:value={customerForm.vat_number}
          placeholder="Enter VAT number"
          error={errors.vat_number}
          id="customer-vat"
          name="vat_number"
        />

        <FormField
          label="Email Address"
          labelArabic="البريد الإلكتروني"
          type="email"
          bind:value={customerForm.email}
          placeholder="customer@example.com"
          error={errors.email}
          id="customer-email"
          name="email"
        />

        <FormField
          label="Phone Number"
          labelArabic="رقم الهاتف"
          type="tel"
          bind:value={customerForm.phone}
          placeholder="+966 50 123 4567"
          error={errors.phone}
          id="customer-phone"
          name="phone"
        />
      </div>
    </div>
  </div>

  <!-- Address Information -->
  <div class="glass-card">
    <div class="p-6">
      <h3 class="heading-3 mb-6">Address Information</h3>
      
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <FormField
          label="Address"
          labelArabic="العنوان"
          type="textarea"
          bind:value={customerForm.address}
          placeholder="Enter full address"
          rows={3}
          error={errors.address}
          id="customer-address"
          name="address"
        />

        <FormField
          label="Address (Arabic)"
          labelArabic="العنوان بالعربية"
          type="textarea"
          bind:value={customerForm.address_arabic}
          placeholder="أدخل العنوان الكامل"
          dir="rtl"
          rows={3}
          error={errors.address_arabic}
          id="customer-address-arabic"
          name="address_arabic"
        />

        <FormField
          label="City"
          labelArabic="المدينة"
          type="text"
          bind:value={customerForm.city}
          placeholder="Enter city name"
          error={errors.city}
          id="customer-city"
          name="city"
        />

        <FormField
          label="City (Arabic)"
          labelArabic="المدينة بالعربية"
          type="text"
          bind:value={customerForm.city_arabic}
          placeholder="أدخل اسم المدينة"
          dir="rtl"
          error={errors.city_arabic}
          id="customer-city-arabic"
          name="city_arabic"
        />

        <FormField
          label="Country"
          labelArabic="البلد"
          type="select"
          bind:value={customerForm.country}
          placeholder="Select country"
          options={countryOptions}
          error={errors.country}
          id="customer-country"
          name="country"
        />

        <FormField
          label="Country (Arabic)"
          labelArabic="البلد بالعربية"
          type="text"
          bind:value={customerForm.country_arabic}
          placeholder="أدخل اسم البلد"
          dir="rtl"
          error={errors.country_arabic}
          id="customer-country-arabic"
          name="country_arabic"
        />
      </div>
    </div>
  </div>
</div>