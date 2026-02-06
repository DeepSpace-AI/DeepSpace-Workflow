type WorkflowItem = {
  id: number;
  name: string;
  description?: string | null;
  steps?: unknown;
  created_at?: string;
  updated_at?: string;
};

export function useProjectWorkflows(projectId: string) {
  const requestHeaders = useRequestHeaders(["cookie"]);
  const { data, refresh } = useAsyncData(
    () => `project-workflows-${projectId}`,
    () =>
      $fetch<{ items: WorkflowItem[] }>(`/api/projects/${projectId}/workflows`, {
        headers: requestHeaders,
      }).catch(() => ({ items: [] }))
  );

  const items = computed(() => data.value?.items ?? []);
  const options = computed(() =>
    items.value.map((item) => ({ label: item.name, value: String(item.id) }))
  );
  const itemMap = computed(() => {
    const map = new Map<string, WorkflowItem>();
    items.value.forEach((item) => map.set(String(item.id), item));
    return map;
  });

  return {
    items,
    options,
    itemMap,
    refresh,
  };
}
