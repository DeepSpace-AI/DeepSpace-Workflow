export default defineNuxtRouteMiddleware(async (to) => {
  const publicRoutes = new Set([
    "/",
    "/sign-in",
    "/sign-up",
    "/sign-out",
    "/privacy-policy",
    "/terms-of-service",
  ]);
  if (publicRoutes.has(to.path)) return;

  try {
    const auth = await $fetch<{ user_id?: number | null } | null>("/api/auth/me");
    if (!auth?.user_id) return navigateTo("/sign-in");
  } catch {
    return navigateTo("/sign-in");
  }

  if (to.path === "/profile") return;

  try {
    const me = await $fetch<{
      profile?: {
        display_name?: string | null;
        full_name?: string | null;
        DisplayName?: string | null;
        FullName?: string | null;
      };
    } | null>("/api/users/me");
    const profile = me?.profile || {};
    const displayName = String(profile.display_name || profile.DisplayName || "").trim();
    const fullName = String(profile.full_name || profile.FullName || "").trim();
    if (!displayName || !fullName) {
      return navigateTo({
        path: "/profile",
        query: {
          force_init: "1",
          redirect: to.fullPath,
        },
      });
    }
  } catch {
    // 资料接口异常时不应误判为登出，放行当前页面
  }
});
