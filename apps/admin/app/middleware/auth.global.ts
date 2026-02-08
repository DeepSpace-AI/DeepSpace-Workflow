export default defineNuxtRouteMiddleware(async (to) => {
  const publicRoutes = new Set(['/sign-in'])
  if (publicRoutes.has(to.path)) return

  const sessionCookie = useCookie('dsp_session')
  if (sessionCookie.value) return

  try {
    const me = await $fetch('/api/auth/me')
    if (!me) return navigateTo('/sign-in')
  } catch {
    return navigateTo('/sign-in')
  }
})
