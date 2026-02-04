import { proxyRequest } from "h3";

export default defineEventHandler(async (event) => {
  const { aiGateway } = useRuntimeConfig();
  if (!aiGateway?.url) {
    throw createError({ statusCode: 500, statusMessage: "Missing AI Gateway config" });
  }

  const id = getRouterParam(event, "id");
  const docId = getRouterParam(event, "docId");
  if (!id || !docId) {
    throw createError({ statusCode: 400, statusMessage: "Missing knowledge base id or document id" });
  }

  const base = aiGateway.url.endsWith("/") ? aiGateway.url.slice(0, -1) : aiGateway.url;
  const query = getQuery(event);
  const url = new URL(`${base}/api/knowledge-bases/${id}/documents/${docId}/download`);
  for (const [key, value] of Object.entries(query)) {
    if (typeof value === "string" && value) {
      url.searchParams.set(key, value);
    }
  }

  return proxyRequest(event, url.toString(), {
    headers: {
      cookie: event.node.req.headers.cookie || "",
    },
  });
});
