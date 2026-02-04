export default defineEventHandler(async (event) => {
  const { aiGateway } = useRuntimeConfig();
  if (!aiGateway?.url || !aiGateway?.apiKey) {
    throw createError({ statusCode: 500, statusMessage: "Missing AI Gateway config" });
  }

  const id = getRouterParam(event, "id");
  if (!id) {
    throw createError({ statusCode: 400, statusMessage: "Missing knowledge base id" });
  }

  const parts = await readMultipartFormData(event);
  if (!parts || parts.length === 0) {
    throw createError({ statusCode: 400, statusMessage: "Missing file" });
  }

  const form = new FormData();
  for (const part of parts) {
    if (part.filename) {
      const blob = new Blob([part.data], { type: part.type || "application/octet-stream" });
      form.append(part.name, blob, part.filename);
      continue;
    }
    if (part.data) {
      form.append(part.name, part.data.toString());
    }
  }

  const base = aiGateway.url.endsWith("/") ? aiGateway.url.slice(0, -1) : aiGateway.url;
  return await $fetch(`${base}/api/knowledge-bases/${id}/documents`, {
    method: "POST",
    headers: {
      Authorization: `Bearer ${aiGateway.apiKey}`,
    },
    body: form,
  });
});
