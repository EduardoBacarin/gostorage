<template>
  <MainLayout>
    <div class="bg-white rounded-3xl shadow-2xl shadow-slate-200 border border-slate-100 p-8 sm:p-16 text-center">
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
  </MainLayout>
</template>

<script setup>
import MainLayout from '../components/MainLayout.vue'
</script>

