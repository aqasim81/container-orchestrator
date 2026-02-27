import { beforeEach, describe, expect, it, vi } from "vitest";
import { APIClientError, fetchAPI } from "./api";

describe("fetchAPI", () => {
	beforeEach(() => {
		vi.restoreAllMocks();
	});

	it("returns parsed JSON on success", async () => {
		const mockData = { status: "ok" };
		global.fetch = vi.fn().mockResolvedValue({
			ok: true,
			json: () => Promise.resolve(mockData),
		});

		const result = await fetchAPI<{ status: string }>("/healthz");

		expect(result).toEqual(mockData);
		expect(global.fetch).toHaveBeenCalledWith(
			"/api/v1/healthz",
			expect.objectContaining({
				headers: expect.objectContaining({
					"Content-Type": "application/json",
				}),
			}),
		);
	});

	it("throws APIClientError on failure", async () => {
		global.fetch = vi.fn().mockResolvedValue({
			ok: false,
			status: 404,
			json: () => Promise.resolve({ error: "not found", code: "NOT_FOUND" }),
		});

		await expect(fetchAPI("/missing")).rejects.toThrow(APIClientError);

		try {
			await fetchAPI("/missing");
		} catch (e) {
			const err = e as APIClientError;
			expect(err.status).toBe(404);
			expect(err.code).toBe("NOT_FOUND");
			expect(err.message).toBe("not found");
		}
	});

	it("handles non-JSON error responses gracefully", async () => {
		global.fetch = vi.fn().mockResolvedValue({
			ok: false,
			status: 502,
			statusText: "Bad Gateway",
			json: () => Promise.reject(new SyntaxError("Unexpected token")),
		});

		try {
			await fetchAPI("/broken");
		} catch (e) {
			const err = e as APIClientError;
			expect(err.status).toBe(502);
			expect(err.code).toBe("UNKNOWN");
			expect(err.message).toBe("HTTP 502: Bad Gateway");
		}
	});

	it("passes custom headers through", async () => {
		global.fetch = vi.fn().mockResolvedValue({
			ok: true,
			json: () => Promise.resolve({}),
		});

		await fetchAPI("/test", {
			headers: { "X-API-Key": "secret" },
		});

		expect(global.fetch).toHaveBeenCalledWith(
			"/api/v1/test",
			expect.objectContaining({
				headers: expect.objectContaining({
					"X-API-Key": "secret",
					"Content-Type": "application/json",
				}),
			}),
		);
	});
});
