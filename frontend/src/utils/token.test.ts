import { describe, it, expect, beforeEach } from "vitest";
import { setToken, getToken, removeToken } from "./token";

describe("token utils", () => {
  beforeEach(() => {
    localStorage.clear();
  });

  it("returns empty string when token is not set", () => {
    expect(getToken()).toBe("");
  });

  it("stores and retrieves token correctly", () => {
    setToken("my-test-jwt-token-12345");
    expect(getToken()).toBe("my-test-jwt-token-12345");
  });

  it("removes token correctly", () => {
    setToken("token-to-be-deleted");
    expect(getToken()).toBe("token-to-be-deleted");
    removeToken();
    expect(getToken()).toBe("");
  });
});
