import { ref } from 'vue';

export function useAdminRelationExplorer() {
  const targetTweetId = ref('');
  const targetTweetText = ref('');
  const detectedHandles = ref<string[]>([]);
  const selectedHandles = ref<string[]>([]);
  const selectedEngine = ref('sotwe');
  const enableWayback = ref(true);
  const isExploring = ref(false);
  const logLines = ref<string[]>(['[READY] Relation Explorer standby.']);

  const setTargetTweet = (id: string, text: string, handles: string[]) => {
    targetTweetId.value = id;
    targetTweetText.value = text;
    detectedHandles.value = handles;
    selectedHandles.value = [...handles];
  };

  const toggleHandle = (handle: string) => {
    const idx = selectedHandles.value.indexOf(handle);
    if (idx >= 0) selectedHandles.value.splice(idx, 1);
    else selectedHandles.value.push(handle);
  };

  const startExploration = async () => {
    if (!selectedHandles.value.length) return;
    isExploring.value = true;
    logLines.value.push(`[DISPATCH] Starting ${selectedEngine.value} exploration for @${selectedHandles.value.join(', @')}`);
    // ここで Go / Python のジョブキックAPIを呼出
  };

  return {
    targetTweetId,
    targetTweetText,
    detectedHandles,
    selectedHandles,
    selectedEngine,
    enableWayback,
    isExploring,
    logLines,
    setTargetTweet,
    toggleHandle,
    startExploration
  };
}
