export const getGatewayBase = (raw: string) => {
  let base = raw.endsWith("/") ? raw.slice(0, -1) : raw;

  const v1Index = base.indexOf("/v1");
  if (v1Index !== -1) {
    base = base.slice(0, v1Index);
  }

  if (base.endsWith("/api")) {
    base = base.slice(0, -4);
  }

  return base;
};
