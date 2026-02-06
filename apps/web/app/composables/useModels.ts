type ModelItem = {
  id: string;
  object?: string;
  owned_by?: string;
};

export function useModels() {
  const requestHeaders = useRequestHeaders(["cookie"]);
  const { data, pending, error, refresh } = useAsyncData("models", () =>
    $fetch<{ data: ModelItem[] }>("/api/models", {
      headers: requestHeaders,
    }).catch(() => ({ data: [] })),
  );

  const items = computed(() => {
    const result = data.value?.data ?? [];
    return result;
  });

  const menuItems = items.value.map((model) => model.id);
  const imageMenuItems = ref<string[]>([]);

  const defaultModel = computed(() => menuItems[0] || "deepseek-chat");
  const defaultImageModel = computed(() => imageMenuItems.value[0] || "");

  return {
    items,
    menuItems,
    imageMenuItems,
    defaultModel,
    defaultImageModel,
    pending,
    error,
    refresh,
  };
}
