export default defineEventHandler(async (event) => {
  const { aiGateway } = useRuntimeConfig();
  if (!aiGateway?.url || !aiGateway?.apiKey) {
    throw createError({ statusCode: 500, statusMessage: "Missing AI Gateway config" });
  }

  const id = getRouterParam(event, "id");
  const docId = getRouterParam(event, "docId");
  if (!id || !docId) {
    throw createError({ statusCode: 400, statusMessage: "Missing knowledge base id or document id" });
  }

  const base = aiGateway.url.endsWith("/") ? aiGateway.url.slice(0, -1) : aiGateway.url;
  return await $fetch(`${base}/api/knowledge-bases/${id}/documents/${docId}`, {
    method: "DELETE",
    headers: {
      Authorization: `Bearer ${aiGateway.apiKey}`,
    },
  });
});
