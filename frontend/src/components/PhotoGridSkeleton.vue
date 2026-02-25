<template>
  <div class="photo-grid-skeleton">
    <LoadingSkeleton
      v-for="i in count"
      :key="i"
      type="photo-card"
    />
  </div>
</template>

<script setup>
import LoadingSkeleton from './LoadingSkeleton.vue'

defineProps({
  count: {
    type: Number,
    default: 12
  }
})
</script>

<style scoped>
.photo-grid-skeleton {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 24px;
  padding: 20px 0;
}

@media (max-width: 768px) {
  .photo-grid-skeleton {
    grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
    gap: 16px;
  }
}

@media (max-width: 480px) {
  .photo-grid-skeleton {
    grid-template-columns: 1fr;
  }
}
</style>
