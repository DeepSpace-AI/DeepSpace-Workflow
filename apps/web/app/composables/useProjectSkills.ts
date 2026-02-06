type SkillItem = {
  id: number;
  name: string;
  description?: string | null;
  prompt?: string | null;
  created_at?: string;
  updated_at?: string;
};

export function useProjectSkills(projectId: string) {
  const requestHeaders = useRequestHeaders(["cookie"]);
  const { data, refresh } = useAsyncData(
    () => `project-skills-${projectId}`,
    () =>
      $fetch<{ items: SkillItem[] }>(`/api/projects/${projectId}/skills`, {
        headers: requestHeaders,
      }).catch(() => ({ items: [] }))
  );

  const items = computed(() => data.value?.items ?? []);
  const options = computed(() =>
    items.value.map((item) => ({ label: item.name, value: String(item.id) }))
  );
  const itemMap = computed(() => {
    const map = new Map<string, SkillItem>();
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
