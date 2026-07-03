import { api } from "../services/api"
import type { AxiosRequestConfig } from "axios"

export const http = {
  async get<TResponse>(url: string, config?: AxiosRequestConfig): Promise<TResponse> {
    return await api.get<TResponse>(url, config).then((responseData) => responseData.data)
  },

  async post<TResponse, TBody = unknown>(
    url: string,
    body?: TBody,
    config?: AxiosRequestConfig,
  ): Promise<TResponse> {
    return await api.post<TResponse>(url, body, config).then((responseData) => responseData.data)
  },

  async delete<TResponse>(url: string, config?: AxiosRequestConfig): Promise<TResponse> {
    return await api.delete<TResponse>(url, config).then((responseData) => responseData.data)
  },
}
