import axios from "axios"

export const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || "http://localhost:8080/api",
  withCredentials: true,
})

function getCookie(name: string): string | undefined {
  const value = `; ${document.cookie}`
  const parts = value.split(`; ${name}=`)
  if (parts.length === 2) {
    return parts.pop()?.split(";").shift()
  }
  return undefined
}

api.interceptors.request.use((config) => {
  const csrfToken = getCookie("csrf_token")
  if (csrfToken) {
    config.headers["X-CSRF-Token"] = csrfToken
  }
  return config
})
