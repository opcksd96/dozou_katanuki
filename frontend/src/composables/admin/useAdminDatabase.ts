// frontend/src/composables/admin/useAdminDatabase.ts (100行以下)
import { ref } from 'vue';

const getApp = () => (window as any)?.go?.main?.App;

export function useAdminDatabase() {
  const searchResults = ref<any[]>([]);
  const totalCount = ref(0);
  const isSearchLoading = ref(false);
  const searchQuery = ref('');
  const searchAccount = ref('all');
  const searchFilter = ref('all');
  const page = ref(1);
  const limit = 20;

  const searchArticles = async () => {
    isSearchLoading.value = true;
    const app = getApp();
    if (app?.SearchArticles) {
      const offset = (page.value - 1) * limit;
      const res = await app.SearchArticles(searchQuery.value, searchAccount.value, searchFilter.value, limit, offset);
      searchResults.value = res?.Items || [];
      totalCount.value = res?.Total || 0;
    }
    isSearchLoading.value = false;
  };

  const saveTranslation = async (articleId: string, ja: string, en: string, zh: string) => {
    const app = getApp();
    if (app?.UpdateArticleTranslations) {
      await app.UpdateArticleTranslations(articleId, ja, en, zh);
      await searchArticles();
    }
  };

  const retryMedia = async (mediaId: string) => {
    const app = getApp();
    if (app?.RetryMediaDownload) {
      await app.RetryMediaDownload(mediaId);
    }
  };

  return {
    searchResults, totalCount, isSearchLoading, searchQuery, searchAccount, searchFilter, page, limit,
    searchArticles, saveTranslation, retryMedia,
  };
}
