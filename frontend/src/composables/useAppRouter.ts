// frontend/src/composables/useAppRouter.ts (100行以下 - SPEC-PRINCIPLE-001)
import { onMounted, watch, type Ref } from 'vue';

export interface RouteState {
  account?: string;
  post?: string;
  adminTab?: string;
  showAdmin?: boolean;
}

export function useAppRouter(
  state: {
    activeArticleId: Ref<string | null>;
    selectedAccount: Ref<string>;
    isAdminOpen: Ref<boolean>;
    isTimelineLoading: Ref<boolean>;
    isDetailLoading: Ref<boolean>;
  },
  callbacks: {
    onSelectAccount: (username: string) => void;
    onSelectPost: (articleId: string) => void;
    onOpenAdmin: (tab?: string) => void;
  }
) {
  const parseUrlParams = (): RouteState => {
    const params = new URLSearchParams(window.location.search);
    const hash = window.location.hash.replace('#', '');
    const hashParams = new URLSearchParams(hash);
    return {
      account: params.get('account') || hashParams.get('account') || undefined,
      post: params.get('post') || hashParams.get('post') || undefined,
      adminTab: params.get('tab') || hashParams.get('tab') || undefined,
      showAdmin: params.has('admin') || hashParams.has('admin'),
    };
  };

  const updateUrl = () => {
    const url = new URL(window.location.href);
    if (state.activeArticleId.value) url.searchParams.set('post', state.activeArticleId.value);
    else url.searchParams.delete('post');

    if (state.selectedAccount.value !== 'all') url.searchParams.set('account', state.selectedAccount.value);
    else url.searchParams.delete('account');

    if (state.isAdminOpen.value) url.searchParams.set('admin', 'true');
    else url.searchParams.delete('admin');

    window.history.replaceState({}, '', url);
  };

  onMounted(() => {
    const route = parseUrlParams();
    if (route.showAdmin || route.adminTab) {
      callbacks.onOpenAdmin(route.adminTab || 'accounts');
    } else if (route.post) {
      callbacks.onSelectPost(route.post);
    } else if (route.account) {
      callbacks.onSelectAccount(route.account);
    }

    window.addEventListener('beforeunload', () => {
      sessionStorage.setItem('katana_scroll_y', window.scrollY.toString());
    });
  });

  watch([state.activeArticleId, state.selectedAccount, state.isAdminOpen], () => {
    updateUrl();
  });

  watch([state.isTimelineLoading, state.isDetailLoading], ([timelineLoad, detailLoad]) => {
    if (!timelineLoad && !detailLoad) {
      const sy = sessionStorage.getItem('katana_scroll_y');
      if (sy) {
        setTimeout(() => {
          window.scrollTo({ top: parseInt(sy, 10), behavior: 'instant' });
          sessionStorage.removeItem('katana_scroll_y');
        }, 100);
      }
    }
  });
}
