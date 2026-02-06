export default defineEventHandler(async (event) => {
  const { aiGateway } = useRuntimeConfig();
  if (!aiGateway?.url) {
    throw createError({ statusCode: 500, statusMessage: "Missing AI Gateway config" });
  }

  const id = getRouterParam(event, "id");
  const workflowId = getRouterParam(event, "workflowId");
  if (!id || !workflowId) {
    throw createError({ statusCode: 400, statusMessage: "Missing project/workflow id" });
  }

  const base = aiGateway.url.endsWith("/") ? aiGateway.url.slice(0, -1) : aiGateway.url;
  return await $fetch(`${base}/api/projects/${id}/workflows/${workflowId}`, {
    method: "DELETE",
    headers: {
      cookie: event.node.req.headers.cookie || "",
    },
  });
});
