import type { RequestOptions, JsonReturn } from "@/types";
import { useMainStore } from "@/stores";
import { getToken } from "@/utils/token";
import { useLocaleStore } from "@/stores/localeStore";

class RestClient {
  private baseUrl: string;

  constructor(baseUrl: string = "/api") {
    this.baseUrl = baseUrl;
  }

  async request<T>(endpoint: string, options: RequestOptions = {}): Promise<T> {
    try {
      useMainStore().loading = true;
      this.validateOptions(options);

      const url = this.buildUrl(endpoint, options.params);
      const headers = this.getHeaders(options.headers);
      const body = options.body ? JSON.stringify(options.body) : undefined;

      const response = await fetch(url, {
        method: options.method || "GET",
        headers,
        body,
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const ret = (await response.json()) as JsonReturn<T>;

      if (!ret.success) {
        throw new Error(`${ret.msg}`);
      }

      return ret.data as T;
    } catch (error) {
      console.error("Request failed:", error);
      throw error;
    } finally {
      useMainStore().loading = false;
    }
  }

  private validateOptions(options: RequestOptions): void {
    if (options.body && typeof options.body !== "object") {
      throw new Error("Request body must be an object");
    }

    if (options.params) {
      for (const key in options.params) {
        const value = options.params[key];
        if (value === undefined || value === null) {
          throw new Error(`Invalid parameter value for ${key}`);
        }
      }
    }
  }

  private buildUrl(
    endpoint: string,
    params?: Record<string, string | number | boolean>
  ): string {
    let url = `${this.baseUrl}${endpoint}`;

    if (params) {
      const queryString = Object.entries(params)
        .map(
          ([key, value]) =>
            `${encodeURIComponent(key)}=${encodeURIComponent(value)}`
        )
        .join("&");
      url += `?${queryString}`;
    }

    return url;
  }

  private getHeaders(customHeaders?: Record<string, string>): Headers {
    const headers = new Headers({
      "Content-Type": "application/json",
      Authorization: getToken(),
      "Accept-Language": useLocaleStore().getLanguage(),
      ...customHeaders,
    });
    return headers;
  }

  public get<T>(
    endpoint: string,
    params?: Record<string, string | number | boolean>
  ): Promise<T> {
    return this.request<T>(endpoint, { method: "GET", params });
  }

  public post<T>(endpoint: string, body: unknown): Promise<T> {
    return this.request<T>(endpoint, { method: "POST", body });
  }

  public put<T>(endpoint: string, body: unknown): Promise<T> {
    return this.request<T>(endpoint, { method: "PUT", body });
  }

  public delete<T>(endpoint: string, body: unknown): Promise<T> {
    return this.request<T>(endpoint, { method: "DELETE", body });
  }
}

export default RestClient;
