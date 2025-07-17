<script>
  export let onSendRequest;
  
  function handleSendClick() {
    console.log('Header: Send Request button clicked');
    if (onSendRequest) {
      onSendRequest();
    } else {
      console.error('onSendRequest function not provided');
    }
  }
</script>

<header class="bg-white border-b border-gray-200 px-6 py-4 shadow-sm">
  <div class="flex items-center justify-between">
    <div class="flex items-center space-x-3">
      <h1 class="text-xl font-bold text-gray-900">cclearly</h1>
      <span class="text-sm text-gray-500 bg-gray-100 px-2 py-1 rounded-full">cURL Inspector</span>
    </div>
    <div class="flex items-center space-x-3">
      <button 
        on:click={handleSendClick}
        class="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors duration-200 font-medium shadow-sm hover:shadow-md"
      >
        Send Request
      </button>
      <button 
        on:click={() => console.log('Test button clicked')}
        class="px-4 py-2 bg-gray-500 text-white rounded-lg hover:bg-gray-600 transition-colors duration-200 text-sm"
      >
        Test
      </button>
    </div>
  </div>
</header> 