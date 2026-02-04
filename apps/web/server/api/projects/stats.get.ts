import { getGatewayBase } from "#server/utils/gateway";

export default defineEventHandler(async (event) => {
  const { aiGateway } = useRuntimeConfig();

  const base = getGatewayBase(aiGateway.url);

  const res = await fetch(`${base}/api/projects/stats`, {
    headers: {
      cookie: event.node.req.headers.cookie || "",
    },
  });

  const data = await res.json();
  if (!res.ok) {
    const msg =
      typeof data?.error === "string"
        ? data.error
        : data?.error?.message || "Failed to fetch project stats";
    throw createError({ statusCode: res.status, statusMessage: msg });
  }

  return data;
});
