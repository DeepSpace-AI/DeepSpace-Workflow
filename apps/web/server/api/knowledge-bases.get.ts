export default defineEventHandler(async (event) => {
  const { aiGateway } = useRuntimeConfig();
  if (!aiGateway?.url || !aiGateway?.apiKey) {
    throw createError({ statusCode: 500, statusMessage: "Missing AI Gateway config" });
  }

  const base = aiGateway.url.endsWith("/") ? aiGateway.url.slice(0, -1) : aiGateway.url;
  const query = getQuery(event);
  return await $fetch(`${base}/api/knowledge-bases`, {
    headers: {
      Authorization: `Bearer ${aiGateway.apiKey}`,
    },
    query,
  });
});
