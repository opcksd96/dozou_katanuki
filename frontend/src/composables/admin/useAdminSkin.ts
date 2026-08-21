// frontend/src/composables/admin/useAdminSkin.ts (100行以下)
import { ref } from 'vue';

const getApp = () => (window as any)?.go?.main?.App;

export function useAdminSkin() {
  const skinCSS = ref('');
  const isSkinLoading = ref(false);
  const isSkinSaved = ref(false);

  const fetchSkinCSS = async (platform = 'twitter') => {
    isSkinLoading.value = true;
    const app = getApp();
    if (app?.GetSkinCSS) {
      skinCSS.value = await app.GetSkinCSS(platform);
    }
    isSkinLoading.value = false;
  };

  const saveSkinCSS = async (platform = 'twitter', css: string) => {
    isSkinLoading.value = true;
    const app = getApp();
    if (app?.SaveSkinCSS) {
      await app.SaveSkinCSS(platform, css);
      isSkinSaved.value = true;
      setTimeout(() => { isSkinSaved.value = false; }, 3000);
    }
    isSkinLoading.value = false;
  };

  return { skinCSS, isSkinLoading, isSkinSaved, fetchSkinCSS, saveSkinCSS };
}
