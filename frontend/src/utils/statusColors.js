// Centralized status color utilities
export function getStatusColor(status) {
  if (status >= 200 && status < 300) return 'bg-emerald-900 text-emerald-300';
  if (status >= 400 && status < 500) return 'bg-amber-900 text-amber-300';
  if (status >= 500) return 'bg-red-900 text-red-300';
  return 'bg-gray-600 text-gray-200';
}

export function getStatusText(status) {
  if (status >= 200 && status < 300) return 'Success';
  if (status >= 400 && status < 500) return 'Client Error';
  if (status >= 500) return 'Server Error';
  return 'Unknown';
}

export function getStatusTextColor(status) {
  if (status >= 200 && status < 300) return 'text-emerald-400';
  if (status >= 400 && status < 500) return 'text-amber-400';
  if (status >= 500) return 'text-red-400';
  return 'text-gray-400';
} 