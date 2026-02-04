import { getGatewayBase } from "../../../utils/gateway";

export default defineEventHandler(async (event) => {
  const { aiGateway } = useRuntimeConfig();
  if (!aiGateway?.url) {
    throw createError({ statusCode: 500, statusMessage: "Missing AI Gateway config" });
  }

  const id = getRouterParam(event, "id");
  if (!id) {
    throw createError({ statusCode: 400, statusMessage: "Missing project id" });
  }

  const body = await readBody(event);
  const base = getGatewayBase(aiGateway.url);
  return await $fetch(`${base}/api/projects/${id}/conversations`, {
    method: "POST",
    headers: {
      cookie: event.node.req.headers.cookie || "",
    },
    body,
  });
});
