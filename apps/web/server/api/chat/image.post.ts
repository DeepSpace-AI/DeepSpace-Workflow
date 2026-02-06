import { getGatewayBase } from "../../utils/gateway";

type ImageSize = "1024x1024" | "1024x1792" | "1792x1024";
type ImageQuality = "standard" | "hd";
type ImageStyle = "vivid" | "natural";

const DEFAULT_MODEL = "dall-e-3";
const DEFAULT_SIZE: ImageSize = "1024x1024";
const DEFAULT_QUALITY: ImageQuality = "standard";

function normalizeSize(value: unknown): ImageSize {
  return value === "1024x1792" || value === "1792x1024" || value === "1024x1024"
    ? value
    : DEFAULT_SIZE;
}

function normalizeQuality(value: unknown): ImageQuality {
  return value === "hd" || value === "standard" ? value : DEFAULT_QUALITY;
}

function normalizeStyle(value: unknown): ImageStyle | undefined {
  return value === "vivid" || value === "natural" ? value : undefined;
}

export default defineEventHandler(async (event) => {
  const { aiGateway } = useRuntimeConfig();
  if (!aiGateway?.url) {
    throw createError({ statusCode: 500, statusMessage: "Missing AI Gateway config" });
  }

  let body: any;
  try {
    body = await readBody(event);
  } catch {
    throw createError({ statusCode: 400, statusMessage: "Invalid body" });
  }

  const prompt = typeof body?.prompt === "string" ? body.prompt.trim() : "";
  if (!prompt) {
    throw createError({ statusCode: 400, statusMessage: "Prompt required" });
  }

  const model = typeof body?.model === "string" && body.model.trim() ? body.model.trim() : DEFAULT_MODEL;
  const size = normalizeSize(body?.size);
  const quality = normalizeQuality(body?.quality);
  const style = normalizeStyle(body?.style);

  const base = getGatewayBase(aiGateway.url);
  const cookie = event.node.req.headers.cookie || "";

  try {
    const walletResp = await $fetch<{ wallet?: any }>(`${base}/api/billing/wallet`, {
      headers: { cookie },
    });
    const wallet = walletResp?.wallet ?? {};
    const balance =
      typeof wallet?.balance === "number"
        ? wallet.balance
        : typeof wallet?.Balance === "number"
          ? wallet.Balance
          : 0;
    if (balance <= 0) {
      throw createError({ statusCode: 402, statusMessage: "Payment Required" });
    }
  } catch (err: any) {
    if (err?.statusCode === 402 || err?.status === 402) {
      throw createError({ statusCode: 402, statusMessage: "Payment Required" });
    }
    // Wallet check failures are ignored; gateway still enforces balance.
  }

  try {
    const response = await $fetch<{
      data?: Array<{ url?: string; revised_prompt?: string }>;
    }>(`${base}/v1/images/generations`, {
      method: "POST",
      headers: { cookie },
      body: {
        model,
        prompt,
        size,
        quality,
        style,
        n: 1,
        response_format: "url",
      },
    });

    const item = response?.data?.[0];
    const imageUrl = typeof item?.url === "string" ? item.url : "";
    if (!imageUrl) {
      throw createError({
        statusCode: 502,
        statusMessage: "Image generation returned empty URL",
      });
    }

    return {
      imageUrl,
      revisedPrompt:
        typeof item?.revised_prompt === "string" && item.revised_prompt.trim()
          ? item.revised_prompt.trim()
          : undefined,
      model,
    };
  } catch (error: any) {
    const status =
      error?.statusCode || error?.status || error?.response?.status || 500;
    const statusMessage =
      error?.statusMessage ||
      error?.data?.error?.message ||
      error?.message ||
      "Image generation failed";
    throw createError({ statusCode: status, statusMessage });
  }
});
