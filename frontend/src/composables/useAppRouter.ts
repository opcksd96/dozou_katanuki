// frontend/src/composables/useAppRouter.ts (100行以下 - SPEC-PRINCIPLE-001)
import { onMounted, watch, type Ref } from 'vue';
import { useRoute } from 'vue-router';

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
  const route = useRoute();

  const parseUrlParams = (): RouteState => {
    const params = new URLSearchParams(window.location.search);
    const hash = window.location.hash.replace(/^#/, '');
    const hashQuery = hash.includes('?') ? hash.split('?')[1] : '';
    const hashParams = new URLSearchParams(hashQuery);
    return {
      account: (route?.query?.account as string) || params.get('account') || hashParams.get('account') || undefined,
      post: (route?.query?.post as string) || params.get('post') || hashParams.get('post') || undefined,
      adminTab: (route?.query?.tab as string) || params.get('tab') || hashParams.get('tab') || undefined,
      showAdmin: !!route?.query?.admin || params.has('admin') || hashParams.has('admin'),
    };
  };

  const applyRoute = (r: RouteState) => {
    if (r.showAdmin || r.adminTab) {
      callbacks.onOpenAdmin(r.adminTab || 'accounts');
    } else if (r.post) {
      callbacks.onSelectPost(r.post);
    } else if (r.account) {
      callbacks.onSelectAccount(r.account);
    }
  };

  onMounted(() => {
    applyRoute(parseUrlParams());
    window.addEventListener('beforeunload', () => {
      sessionStorage.setItem('katana_scroll_y', window.scrollY.toString());
    });
  });

  watch(() => route?.query?.post, (newPost) => {
    if (newPost && typeof newPost === 'string' && newPost !== state.activeArticleId.value) {
      callbacks.onSelectPost(newPost);
    }
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
