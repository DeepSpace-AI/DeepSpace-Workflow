import { getGatewayBase } from "../../utils/gateway";

export default defineEventHandler(async (event) => {
  const { aiGateway } = useRuntimeConfig();
  if (!aiGateway?.url) {
    throw createError({ statusCode: 500, statusMessage: "Missing AI Gateway config" });
  }

  const base = getGatewayBase(aiGateway.url);
  const query = getQuery(event);
  return await $fetch(`${base}/api/billing/usage`, {
    headers: {
      cookie: event.node.req.headers.cookie || "",
    },
    query,
  });
});
