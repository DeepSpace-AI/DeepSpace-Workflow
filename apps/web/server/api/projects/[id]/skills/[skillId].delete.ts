export default defineEventHandler(async (event) => {
  const { aiGateway } = useRuntimeConfig();
  if (!aiGateway?.url) {
    throw createError({ statusCode: 500, statusMessage: "Missing AI Gateway config" });
  }

  const id = getRouterParam(event, "id");
  const skillId = getRouterParam(event, "skillId");
  if (!id || !skillId) {
    throw createError({ statusCode: 400, statusMessage: "Missing project/skill id" });
  }

  const base = aiGateway.url.endsWith("/") ? aiGateway.url.slice(0, -1) : aiGateway.url;
  return await $fetch(`${base}/api/projects/${id}/skills/${skillId}`, {
    method: "DELETE",
    headers: {
      cookie: event.node.req.headers.cookie || "",
    },
  });
});
