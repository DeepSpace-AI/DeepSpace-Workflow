type ModelItem = {
  id: string;
  name?: string;
  capabilities?: string[];
  object?: string;
  owned_by?: string;
};

export function useModels() {
  const requestHeaders = useRequestHeaders(["cookie"]);
  const { data, pending, error, refresh } = useAsyncData("models", () =>
    $fetch<{ items?: ModelItem[]; data?: ModelItem[] }>("/api/models", {
      headers: requestHeaders,
    }).catch(() => ({ items: [], data: [] })),
  );

  const getItems = () => {
    const result = data.value?.items ?? data.value?.data ?? [];
    return result;
  };

  const getMenuItems = () =>
    Array.from(
      new Set(
        getItems()
          .map((model) => model.name || model.id)
          .filter((model): model is string => Boolean(model && model.trim())),
      ),
    );

  const getImageMenuItems = () =>
    Array.from(
      new Set(
        getItems()
          .filter((model) => Array.isArray(model.capabilities) && model.capabilities.includes("image"))
          .map((model) => model.name || model.id)
          .filter((model): model is string => Boolean(model && model.trim())),
      ),
    );

  const defaultModel = {
    get value() {
      return getMenuItems()[0] || "deepseek-chat";
    },
  };

  const defaultImageModel = {
    get value() {
      return getImageMenuItems()[0] || "";
    },
  };

  return {
    get items() {
      return getItems();
    },
    get menuItems() {
      return getMenuItems();
    },
    get imageMenuItems() {
      return getImageMenuItems();
    },
    defaultModel,
    defaultImageModel,
    pending,
    error,
    refresh,
  };
}
