<script>
  export let requestData;
  export let onRequestDataChange;

  function addHeader() {
    const newHeaders = [...requestData.headers, { key: '', value: '' }];
    onRequestDataChange({ ...requestData, headers: newHeaders });
  }

  function removeHeader(index) {
    const newHeaders = requestData.headers.filter((_, i) => i !== index);
    onRequestDataChange({ ...requestData, headers: newHeaders });
  }

  function updateHeader(index, field, value) {
    const newHeaders = [...requestData.headers];
    newHeaders[index] = { ...newHeaders[index], [field]: value };
    onRequestDataChange({ ...requestData, headers: newHeaders });
  }
</script>

<div class="flex-1 bg-white flex flex-col h-full">
  <div class="p-4 bg-white">
    <h2 class="text-lg font-semibold text-gray-900">Request</h2>
    <p class="text-sm text-gray-500 mt-1">Configure your HTTP request</p>
  </div>
  
  <div class="flex-1 px-4 pb-4 space-y-6 overflow-y-auto">
    <!-- Method and URL -->
    <div class="space-y-2">
      <label class="text-sm font-medium text-gray-700">Request URL</label>
      <div class="flex space-x-3">
        <select 
          bind:value={requestData.method}
          class="px-4 py-3 border border-gray-300 rounded-lg bg-white text-sm font-medium focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
        >
          <option value="GET">GET</option>
          <option value="POST">POST</option>
          <option value="PUT">PUT</option>
          <option value="DELETE">DELETE</option>
          <option value="PATCH">PATCH</option>
        </select>
        <input 
          type="text" 
          bind:value={requestData.url}
          placeholder="https://api.example.com/endpoint"
          class="flex-1 px-4 py-3 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
        />
      </div>
    </div>

    <!-- Headers -->
    <div class="space-y-3">
      <div class="flex items-center justify-between">
        <label class="text-sm font-medium text-gray-700">Headers</label>
        <button 
          on:click={addHeader}
          class="text-sm text-blue-600 hover:text-blue-700 font-medium flex items-center space-x-1"
        >
          <span>+</span>
          <span>Add Header</span>
        </button>
      </div>
      
      {#if requestData.headers.length === 0}
        <div class="text-center py-8 text-gray-500 border-2 border-dashed border-gray-200 rounded-lg">
          <div class="text-2xl mb-2">🔧</div>
          <p class="text-sm">No headers added</p>
          <p class="text-xs mt-1">Click "Add Header" to get started</p>
        </div>
      {:else}
        <div class="space-y-3">
          {#each requestData.headers as header, index}
            <div class="flex space-x-3 items-center">
              <input 
                type="text" 
                value={header.key}
                on:input={(e) => updateHeader(index, 'key', e.target.value)}
                placeholder="Content-Type"
                class="flex-1 px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
              <input 
                type="text" 
                value={header.value}
                on:input={(e) => updateHeader(index, 'value', e.target.value)}
                placeholder="application/json"
                class="flex-1 px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
              <button 
                on:click={() => removeHeader(index)}
                class="p-2 text-red-500 hover:text-red-700 hover:bg-red-50 rounded-md transition-colors"
                title="Remove header"
              >
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
                </svg>
              </button>
            </div>
          {/each}
        </div>
      {/if}
    </div>

    <!-- Body -->
    {#if requestData.method !== 'GET'}
      <div class="space-y-2">
        <label class="text-sm font-medium text-gray-700">Request Body</label>
        <textarea 
          bind:value={requestData.body}
          placeholder="Enter JSON, form data, or raw text"
          rows="8"
          class="w-full px-4 py-3 border border-gray-300 rounded-lg text-sm font-mono focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-none"
        ></textarea>
        <p class="text-xs text-gray-500">Enter JSON, form data, or raw text</p>
      </div>
    {/if}
  </div>
</div> 