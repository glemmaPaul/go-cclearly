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

<header class="bg-gray-900 border-b border-gray-800 px-6 py-4">
  <div class="flex items-center justify-between">
    <div class="flex items-center space-x-3">
      <h1 class="text-xl font-bold text-white font-mono">cclearly</h1>
      <span class="text-sm text-gray-300 bg-gray-800 px-3 py-1 rounded-full">cURL Inspector</span>
    </div>
    <div class="flex items-center space-x-3">
    </div>
  </div>
</header> 