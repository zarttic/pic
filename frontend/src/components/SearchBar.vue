<template>
  <div class="search-bar">
    <div class="search-input-wrapper">
      <input
        v-model="searchQuery"
        type="text"
        :placeholder="placeholder"
        @keyup.enter="handleSearch"
        @keyup.esc="clearSearch"
        class="search-input"
      />
      <button v-if="searchQuery" @click="clearSearch" class="clear-button" title="清除搜索">
        ✕
      </button>
      <button @click="handleSearch" class="search-button" title="搜索">
        🔍
      </button>
    </div>

    <div v-if="showFilters" class="filters">
      <select v-model="filters.year" @change="handleFilterChange" class="filter-select">
        <option value="">所有年份</option>
        <option v-for="year in availableYears" :key="year" :value="year">
          {{ year }}
        </option>
      </select>

      <select v-model="filters.camera" @change="handleFilterChange" class="filter-select">
        <option value="">所有相机</option>
        <option v-for="camera in availableCameras" :key="camera" :value="camera">
          {{ camera }}
        </option>
      </select>

      <label class="featured-filter">
        <input
          v-model="filters.featured"
          type="checkbox"
          @change="handleFilterChange"
        />
        仅显示精选
      </label>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'

const props = defineProps({
  placeholder: {
    type: String,
    default: '搜索照片...'
  },
  showFilters: {
    type: Boolean,
    default: false
  },
  availableYears: {
    type: Array,
    default: () => []
  },
  availableCameras: {
    type: Array,
    default: () => []
  },
  modelValue: {
    type: Object,
    default: () => ({
      search: '',
      year: '',
      camera: '',
      featured: false
    })
  }
})

const emit = defineEmits(['update:modelValue', 'search'])

const searchQuery = ref(props.modelValue.search || '')
const filters = ref({
  year: props.modelValue.year || '',
  camera: props.modelValue.camera || '',
  featured: props.modelValue.featured || false
})

const handleSearch = () => {
  emitSearch()
}

const clearSearch = () => {
  searchQuery.value = ''
  emitSearch()
}

const handleFilterChange = () => {
  emitSearch()
}

const emitSearch = () => {
  const searchParams = {
    search: searchQuery.value,
    ...filters.value
  }
  emit('update:modelValue', searchParams)
  emit('search', searchParams)
}

watch(
  () => props.modelValue,
  (newValue) => {
    searchQuery.value = newValue.search || ''
    filters.value = {
      year: newValue.year || '',
      camera: newValue.camera || '',
      featured: newValue.featured || false
    }
  },
  { deep: true }
)
</script>

<style scoped>
.search-bar {
  width: 100%;
  max-width: 800px;
  margin: 0 auto;
}

.search-input-wrapper {
  display: flex;
  gap: 8px;
  margin-bottom: 16px;
}

.search-input {
  flex: 1;
  padding: 12px 16px;
  border: 2px solid #e0e0e0;
  border-radius: 8px;
  font-size: 16px;
  transition: border-color 0.3s;
}

.search-input:focus {
  outline: none;
  border-color: #667eea;
}

.clear-button,
.search-button {
  padding: 12px 16px;
  border: none;
  border-radius: 8px;
  font-size: 16px;
  cursor: pointer;
  transition: background-color 0.3s, transform 0.2s;
}

.clear-button {
  background-color: #f5f5f5;
}

.clear-button:hover {
  background-color: #e0e0e0;
}

.search-button {
  background-color: #667eea;
  color: white;
}

.search-button:hover {
  background-color: #5568d3;
  transform: translateY(-2px);
}

.filters {
  display: flex;
  gap: 16px;
  flex-wrap: wrap;
  align-items: center;
}

.filter-select {
  padding: 8px 12px;
  border: 1px solid #e0e0e0;
  border-radius: 6px;
  font-size: 14px;
  cursor: pointer;
  transition: border-color 0.3s;
}

.filter-select:focus {
  outline: none;
  border-color: #667eea;
}

.featured-filter {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  cursor: pointer;
}

.featured-filter input[type='checkbox'] {
  width: 18px;
  height: 18px;
  cursor: pointer;
}

@media (max-width: 768px) {
  .filters {
    flex-direction: column;
    align-items: stretch;
  }

  .filter-select {
    width: 100%;
  }
}
</style>
