import { getGatewayBase } from "#server/utils/gateway";

export default defineEventHandler(async (event) => {
  const { aiGateway } = useRuntimeConfig();
  if (!aiGateway?.url) {
    throw createError({ statusCode: 500, statusMessage: "Missing AI Gateway config" });
  }

  const body = await readBody(event);
  const base = getGatewayBase(aiGateway.url);

  const res = await fetch(`${base}/api/auth/login`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });

  const setCookie = res.headers.get("set-cookie");
  if (setCookie) {
    setHeader(event, "set-cookie", setCookie);
  }

  const data = await res.json();
  if (!res.ok) {
    const msg =
      typeof data?.error === "string"
        ? data.error
        : data?.error?.message || "Login failed";
    throw createError({ statusCode: res.status, statusMessage: msg });
  }

  return data;
});
