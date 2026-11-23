import { describe, it, expect, beforeEach, vi, afterEach } from "vitest";
import { createPinia, setActivePinia } from "pinia";
import RestClient from "./rest-client";
import * as tokenUtils from "./token";

describe("RestClient", () => {
  let client: RestClient;

  beforeEach(() => {
    setActivePinia(createPinia());
    client = new RestClient("/api");
    vi.restoreAllMocks();
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it("throws error when options parameter is invalid", async () => {
    await expect(
      client.request("/test", {
        params: { invalid: null as any },
      })
    ).rejects.toThrow("Invalid parameter value for invalid");

    await expect(
      client.request("/test", {
        body: "invalid-string" as any,
      })
    ).rejects.toThrow("Request body must be an object");
  });

  it("builds query parameters correctly", async () => {
    let capturedUrl = "";
    vi.stubGlobal(
      "fetch",
      vi.fn().mockImplementation((url: string) => {
        capturedUrl = url;
        return Promise.resolve({
          ok: true,
          json: () => Promise.resolve({ success: true, data: { status: "ok" } }),
        });
      })
    );

    const result = await client.get<{ status: string }>("/users", {
      page: 1,
      size: 20,
      keyword: "admin",
    });

    expect(capturedUrl).toBe("/api/users?page=1&size=20&keyword=admin");
    expect(result).toEqual({ status: "ok" });
  });

  it("attaches authorization header when token exists", async () => {
    vi.spyOn(tokenUtils, "getToken").mockReturnValue("bearer-token-123");

    let capturedHeaders: any = {};
    vi.stubGlobal(
      "fetch",
      vi.fn().mockImplementation((_url: string, init: any) => {
        capturedHeaders = init.headers;
        return Promise.resolve({
          ok: true,
          json: () => Promise.resolve({ success: true, data: true }),
        });
      })
    );

    await client.post("/action", { name: "test" });
    expect(capturedHeaders.get("Authorization")).toBe("bearer-token-123");
    expect(capturedHeaders.get("Content-Type")).toBe("application/json");
  });

  it("throws error when response HTTP status is not ok", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn().mockResolvedValue({
        ok: false,
        status: 500,
      })
    );

    await expect(client.get("/error-endpoint")).rejects.toThrow(
      "HTTP error! status: 500"
    );
  });

  it("throws error when API response success is false", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn().mockResolvedValue({
        ok: true,
        json: () =>
          Promise.resolve({
            success: false,
            msg: "用户未授权操作",
          }),
      })
    );

    await expect(client.get("/unauthorized")).rejects.toThrow("用户未授权操作");
  });

  it("supports put and delete shortcuts", async () => {
    let capturedMethod = "";
    vi.stubGlobal(
      "fetch",
      vi.fn().mockImplementation((_url: string, init: any) => {
        capturedMethod = init.method;
        return Promise.resolve({
          ok: true,
          json: () => Promise.resolve({ success: true, data: { updated: true } }),
        });
      })
    );

    await client.put("/resource/1", { name: "new" });
    expect(capturedMethod).toBe("PUT");

    await client.delete("/resource/1");
    expect(capturedMethod).toBe("DELETE");
  });
});
