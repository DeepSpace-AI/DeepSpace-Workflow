import { getGatewayBase } from "#server/utils/gateway";

export default defineEventHandler(async (event) => {
  const { aiGateway } = useRuntimeConfig();
  if (!aiGateway?.url) {
    throw createError({ statusCode: 500, statusMessage: "Missing AI Gateway config" });
  }

  const base = getGatewayBase(aiGateway.url);

  const res = await fetch(`${base}/api/users/me`, {
    headers: {
      cookie: event.node.req.headers.cookie || "",
    },
  });

  if (res.status === 401) {
    return null;
  }

  const data = await res.json();
  if (!res.ok) {
    const msg =
      typeof data?.error === "string"
        ? data.error
        : data?.error?.message || "User fetch failed";
    throw createError({ statusCode: res.status, statusMessage: msg });
  }

  return data;
});
