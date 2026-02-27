/** Standard error response from the API. */
interface APIError {
	error: string;
	code: string;
}

/** API client error with structured details. */
export class APIClientError extends Error {
	readonly status: number;
	readonly code: string;

	constructor(status: number, body: APIError) {
		super(body.error);
		this.name = "APIClientError";
		this.status = status;
		this.code = body.code;
	}
}

/**
 * Type-safe fetch wrapper for the orchestrator API.
 * Uses Next.js rewrites in dev to proxy `/api` requests to the Go backend.
 */
export async function fetchAPI<T>(path: string, options?: RequestInit): Promise<T> {
	const url = `/api/v1${path}`;

	const res = await fetch(url, {
		...options,
		headers: {
			"Content-Type": "application/json",
			...options?.headers,
		},
	});

	if (!res.ok) {
		let body: APIError;
		try {
			body = (await res.json()) as APIError;
		} catch {
			body = { error: `HTTP ${res.status}: ${res.statusText}`, code: "UNKNOWN" };
		}
		throw new APIClientError(res.status, body);
	}

	return res.json() as Promise<T>;
}
