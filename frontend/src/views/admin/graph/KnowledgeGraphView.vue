<template>
  <div class="knowledge-graph-container">
    <div class="graph-header">
      <h2>Accounts Knowledge Graph</h2>
    </div>
    <div class="graph-layout">
      <!-- グラフ描画エリア (Cytoscape.js 等を想定) -->
      <div id="cy" class="cy-container" ref="cyRef"></div>
      
      <!-- サイドバー -->
      <GraphInspectorSidebar 
        v-if="selectedEdge" 
        :edge="selectedEdge"
        @close="selectedEdge = null"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import GraphInspectorSidebar from './GraphInspectorSidebar.vue'

const cyRef = ref<HTMLElement | null>(null)
const selectedEdge = ref<any>(null)

onMounted(() => {
  // TODO: Initialize Cytoscape and fetch data via window.go.wails.RelationApp.GetKnowledgeGraph
  console.log('Knowledge Graph mounted')
  // Mock data for initial render
  if (cyRef.value) {
    cyRef.value.innerHTML = '<div style="padding: 20px;">Graph rendering area (Cytoscape.js)</div>'
  }
})
</script>

<style scoped>
.knowledge-graph-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
}
.graph-header {
  padding: 1rem;
  background: #2c3e50;
  color: white;
}
.graph-layout {
  display: flex;
  flex: 1;
  overflow: hidden;
}
.cy-container {
  flex: 1;
  background: #ecf0f1;
  position: relative;
}
</style>
