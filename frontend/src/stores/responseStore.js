import { writable } from 'svelte/store';

export const responseStore = writable({
  status: null,
  body: '',
  formattedBody: '',
  responseType: 'raw',
  headers: []
}); 