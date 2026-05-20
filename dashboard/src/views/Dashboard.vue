<template>
  <div class="min-h-full bg-slate-50 flex flex-col font-sans">
    <!-- Navigation -->
    <nav class="bg-white border-b border-slate-200 sticky top-0 z-10">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between h-16">
          <div class="flex items-center">
            <div class="flex-shrink-0 flex items-center">
              <div class="h-10 w-10 bg-blue-600 rounded-xl flex items-center justify-center shadow-lg shadow-blue-100">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 7v10c0 2.21 3.582 4 8 4s8-1.79 8-4V7M4 7c0 2.21 3.582 4 8 4s8-1.79 8-4M4 7c0-2.21 3.582-4 8-4s8 1.79 8 4m0 5c0 2.21-3.582 4-8 4s-8-1.79-8-4" />
                </svg>
              </div>
              <span class="ml-3 text-xl font-bold text-slate-900 tracking-tight">GoStorage</span>
            </div>
          </div>
          <div class="flex items-center">
            <button 
              @click="handleLogout" 
              class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-bold rounded-xl text-white bg-red-600 hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500 shadow-lg shadow-red-100 transition-all duration-200"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
              </svg>
              Logout
            </button>
          </div>
        </div>
      </div>
    </nav>

    <!-- Content Area -->
    <main class="flex-grow flex items-center justify-center p-6">
      <div class="max-w-4xl w-full bg-white rounded-3xl shadow-2xl shadow-slate-200 border border-slate-100 p-8 sm:p-16 text-center">
        <div class="inline-flex items-center justify-center h-24 w-24 rounded-3xl bg-blue-50 text-blue-600 mb-8 border border-blue-100">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        </div>
        <h1 class="text-4xl font-extrabold text-slate-900 mb-6 tracking-tight">System Operational</h1>
        <p class="text-xl text-slate-600 max-w-2xl mx-auto leading-relaxed mb-12 font-medium">
          Welcome back! Your Object Storage node is healthy. We are currently finalizing the bucket exploration module for this dashboard.
        </p>
        
        <div class="grid grid-cols-1 sm:grid-cols-3 gap-6 max-w-3xl mx-auto">
          <div class="p-6 bg-slate-50 rounded-2xl border border-slate-200 text-left hover:border-blue-300 transition-colors duration-200">
            <span class="block text-xs font-bold text-slate-400 uppercase tracking-widest mb-2">API Status</span>
            <div class="flex items-center">
              <span class="h-3 w-3 bg-green-500 rounded-full animate-pulse mr-3"></span>
              <span class="text-lg font-bold text-slate-900">Connected</span>
            </div>
          </div>
          <div class="p-6 bg-slate-50 rounded-2xl border border-slate-200 text-left hover:border-blue-300 transition-colors duration-200">
            <span class="block text-xs font-bold text-slate-400 uppercase tracking-widest mb-2">Storage Path</span>
            <span class="text-lg font-bold text-slate-900">Local FS</span>
          </div>
          <div class="p-6 bg-slate-50 rounded-2xl border border-slate-200 text-left hover:border-blue-300 transition-colors duration-200">
            <span class="block text-xs font-bold text-slate-400 uppercase tracking-widest mb-2">Build</span>
            <span class="text-lg font-bold text-slate-900">1.0.0-PRO</span>
          </div>
        </div>
      </div>
    </main>

    <!-- Footer -->
    <footer class="py-6 text-center text-slate-400 text-sm font-medium">
      &copy; 2026 GoStorage - Secure & Fast Object Storage
    </footer>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
import axios from 'axios'

const router = useRouter()

const handleLogout = async () => {
  const token = localStorage.getItem('token')
  try {
    await axios.delete('/v1/auth', {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })
  } catch (err) {
    console.error('Logout failed on server, clearing local storage anyway', err)
  } finally {
    localStorage.removeItem('token')
    router.push('/login')
  }
}
</script>
