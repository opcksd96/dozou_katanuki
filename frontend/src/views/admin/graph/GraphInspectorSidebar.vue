<template>
  <div class="inspector-sidebar">
    <div class="sidebar-header">
      <h3>Edge Inspector</h3>
      <button @click="$emit('close')">Close</button>
    </div>
    <div class="sidebar-content" v-if="edge">
      <h4>Relation Details</h4>
      <ul>
        <li><strong>Source:</strong> {{ edge.source }}</li>
        <li><strong>Target:</strong> {{ edge.target }}</li>
        <li><strong>Type:</strong> {{ edge.type }}</li>
        <li><strong>Weight:</strong> {{ edge.weight }}</li>
      </ul>
      
      <h4>Evidences</h4>
      <div class="evidence-list">
        <div class="evidence-item" v-for="ev in evidences" :key="ev.id">
          <p>{{ ev.context_snippet }}</p>
          <button v-if="!ev.is_salvaged" @click="salvage(ev)">Salvage Media</button>
          <span v-else>Salvaged ✓</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const props = defineProps<{
  edge: any
}>()

defineEmits(['close'])

const evidences = ref<any[]>([
  { id: '1', context_snippet: 'Mock evidence snippet 1', is_salvaged: false },
  { id: '2', context_snippet: 'Mock evidence snippet 2', is_salvaged: true }
])

const salvage = (ev: any) => {
  console.log('Salvaging media for evidence:', ev.id)
  ev.is_salvaged = true
}
</script>

<style scoped>
.inspector-sidebar {
  width: 350px;
  background: white;
  border-left: 1px solid #bdc3c7;
  display: flex;
  flex-direction: column;
}
.sidebar-header {
  display: flex;
  justify-content: space-between;
  padding: 1rem;
  background: #ecf0f1;
  border-bottom: 1px solid #bdc3c7;
}
.sidebar-content {
  padding: 1rem;
  overflow-y: auto;
}
.evidence-item {
  border: 1px solid #eee;
  padding: 0.5rem;
  margin-bottom: 0.5rem;
}
</style>
