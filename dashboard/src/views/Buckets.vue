<template>
  <MainLayout>
    <div class="space-y-6">
      <!-- Header & Search -->
      <div class="bg-white rounded-3xl shadow-xl shadow-slate-200 border border-slate-100 p-6">
        <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
          <div>
            <h1 class="text-2xl font-extrabold text-slate-900 tracking-tight">Buckets</h1>
            <p class="text-slate-500 font-medium">Manage your storage containers</p>
          </div>
          <div class="flex flex-col sm:flex-row gap-3">
            <!-- Search -->
            <div class="relative">
              <span class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none text-slate-400">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                </svg>
              </span>
              <input 
                v-model="search"
                type="text" 
                placeholder="Search buckets..."
                class="block w-full pl-10 pr-4 py-2.5 border border-slate-200 rounded-2xl text-sm font-medium focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all duration-200"
                @input="handleSearch"
              />
            </div>
            <!-- Limit -->
            <select 
              v-model="limit"
              class="block w-full sm:w-auto px-4 py-2.5 border border-slate-200 rounded-2xl text-sm font-bold text-slate-700 focus:outline-none focus:ring-2 focus:ring-blue-500 transition-all duration-200"
            >
              <option :value="10">10 per page</option>
              <option :value="20">20 per page</option>
              <option :value="50">50 per page</option>
              <option :value="100">100 per page</option>
            </select>
          </div>
        </div>

        <div class="mt-4 flex flex-wrap gap-3 items-center border-t border-slate-50 pt-4">
          <span class="text-xs font-bold text-slate-400 uppercase tracking-widest">Sort by:</span>
          <select 
            v-model="sortBy"
            class="px-3 py-1.5 border border-slate-200 rounded-xl text-xs font-bold text-slate-700 focus:outline-none focus:ring-2 focus:ring-blue-500 transition-all duration-200"
          >
            <option value="name">Name</option>
            <option value="created_at">Creation Date</option>
          </select>
          <select 
            v-model="sortDir"
            class="px-3 py-1.5 border border-slate-200 rounded-xl text-xs font-bold text-slate-700 focus:outline-none focus:ring-2 focus:ring-blue-500 transition-all duration-200"
          >
            <option value="asc">Ascending</option>
            <option value="desc">Descending</option>
          </select>
        </div>
      </div>

      <!-- Table Card -->
      <div class="bg-white rounded-3xl shadow-xl shadow-slate-200 border border-slate-100 overflow-hidden">
        <div class="overflow-x-auto">
          <table class="min-w-full divide-y divide-slate-100">
            <thead class="bg-slate-50">
              <tr>
                <th scope="col" class="px-6 py-4 text-left text-xs font-bold text-slate-400 uppercase tracking-widest">Name</th>
                <th scope="col" class="px-6 py-4 text-left text-xs font-bold text-slate-400 uppercase tracking-widest">Visibility</th>
                <th scope="col" class="px-6 py-4 text-left text-xs font-bold text-slate-400 uppercase tracking-widest">Created At</th>
                <th scope="col" class="px-6 py-4 text-right text-xs font-bold text-slate-400 uppercase tracking-widest">Actions</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100">
              <tr v-if="loading" v-for="i in 3" :key="'skeleton-'+i" class="animate-pulse">
                <td class="px-6 py-4 whitespace-nowrap"><div class="h-4 bg-slate-100 rounded w-3/4"></div></td>
                <td class="px-6 py-4 whitespace-nowrap"><div class="h-6 bg-slate-100 rounded-full w-16"></div></td>
                <td class="px-6 py-4 whitespace-nowrap"><div class="h-4 bg-slate-100 rounded w-1/2"></div></td>
                <td class="px-6 py-4 whitespace-nowrap text-right"><div class="h-4 bg-slate-100 rounded w-1/4 ml-auto"></div></td>
              </tr>
              <tr v-else-if="buckets.length === 0">
                <td colspan="4" class="px-6 py-12 text-center text-slate-500 font-medium">
                  No buckets found.
                </td>
              </tr>
              <tr v-else v-for="bucket in buckets" :key="bucket.id" class="hover:bg-slate-50 transition-colors duration-150">
                <td class="px-6 py-4 whitespace-nowrap">
                  <span class="text-sm font-bold text-slate-900">{{ bucket.name }}</span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <span 
                    v-if="bucket.allow_public" 
                    class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-bold bg-green-100 text-green-700"
                  >
                    public
                  </span>
                  <span 
                    v-else 
                    class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-bold bg-red-100 text-red-700"
                  >
                    private
                  </span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <span class="text-sm text-slate-600 font-medium">{{ formatDate(bucket.created_at) }}</span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                  <div class="flex justify-end gap-2">
                    <!-- Edit -->
                    <button class="p-2 text-slate-400 hover:text-blue-600 hover:bg-blue-50 rounded-xl transition-all duration-200" title="Edit">
                      <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                      </svg>
                    </button>
                    <!-- Privacy Toggle -->
                    <button class="p-2 text-slate-400 hover:text-amber-600 hover:bg-amber-50 rounded-xl transition-all duration-200" title="Toggle Privacy">
                      <svg v-if="bucket.allow_public" xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 11V7a4 4 0 118 0m-4 8v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2z" />
                      </svg>
                      <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
                      </svg>
                    </button>
                    <!-- Delete -->
                    <button class="p-2 text-slate-400 hover:text-red-600 hover:bg-red-50 rounded-xl transition-all duration-200" title="Delete">
                      <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                      </svg>
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Pagination -->
        <div v-if="totalPages > 1" class="bg-slate-50 px-6 py-4 flex items-center justify-between border-t border-slate-100">
          <div class="hidden sm:flex-1 sm:flex sm:items-center sm:justify-between">
            <div>
              <p class="text-sm text-slate-600 font-medium">
                Showing page <span class="font-bold">{{ page }}</span> of <span class="font-bold">{{ totalPages }}</span>
              </p>
            </div>
            <div>
              <nav class="relative z-0 inline-flex rounded-xl shadow-sm -space-x-px" aria-label="Pagination">
                <!-- Previous -->
                <button 
                  v-if="page > 1"
                  @click="page--"
                  class="relative inline-flex items-center px-3 py-2 rounded-l-xl border border-slate-200 bg-white text-sm font-bold text-slate-500 hover:bg-slate-50 transition-colors duration-200"
                >
                  <span class="sr-only">Previous</span>
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
                  </svg>
                  <span class="ml-1">voltar</span>
                </button>

                <!-- Page Numbers -->
                <button 
                  v-for="p in visiblePages" 
                  :key="p"
                  @click="page = p"
                  class="relative inline-flex items-center px-4 py-2 border border-slate-200 text-sm font-bold transition-all duration-200"
                  :class="[p === page ? 'z-10 bg-blue-600 border-blue-600 text-white' : 'bg-white text-slate-500 hover:bg-slate-50']"
                >
                  {{ p }}
                </button>

                <!-- Next -->
                <button 
                  v-if="page < totalPages"
                  @click="page++"
                  class="relative inline-flex items-center px-3 py-2 rounded-r-xl border border-slate-200 bg-white text-sm font-bold text-slate-500 hover:bg-slate-50 transition-colors duration-200"
                >
                  <span class="mr-1">avançar</span>
                  <span class="sr-only">Next</span>
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
                  </svg>
                </button>
              </nav>
            </div>
          </div>
        </div>
      </div>
    </div>
  </MainLayout>
</template>

<script setup>
import { ref, watch, computed, onMounted } from 'vue'
import axios from 'axios'
import MainLayout from '../components/MainLayout.vue'

const buckets = ref([])
const loading = ref(false)
const search = ref('')
const page = ref(1)
const limit = ref(10)
const sortBy = ref('name')
const sortDir = ref('asc')
const totalPages = ref(1)

let searchTimeout = null

const handleSearch = () => {
  clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    page.value = 1
    fetchBuckets()
  }, 500)
}

const fetchBuckets = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const response = await axios.get('/v1/bucket', {
      params: {
        page: page.value,
        limit: limit.value,
        sort_by: sortBy.value,
        sort_dir: sortDir.value,
        search: search.value
      },
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })

    if (response.data.s) {
      const data = response.data.d
      buckets.value = data.data || []
      totalPages.value = data.total_pages || 1
    }
  } catch (err) {
    console.error('Failed to fetch buckets', err)
  } finally {
    loading.value = false
  }
}

const formatDate = (dateString) => {
  if (!dateString) return '-'
  const date = new Date(dateString)
  return new Intl.DateTimeFormat('pt-BR', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  }).format(date)
}

const visiblePages = computed(() => {
  const pages = []
  const start = Math.max(1, page.value - 2)
  const end = Math.min(totalPages.value, page.value + 2)
  
  for (let i = start; i <= end; i++) {
    pages.push(i)
  }
  return pages
})

watch([limit, sortBy, sortDir, page], () => {
  fetchBuckets()
})

onMounted(() => {
  fetchBuckets()
})
</script>
