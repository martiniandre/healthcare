import { readFileSync } from "node:fs"
import { resolve } from "node:path"

interface ServiceWorkerFetchEvent {
  request: { method: string; mode: string; url: string }
  respondWith: (response: Promise<unknown>) => void
  respondedWith: Promise<unknown> | null
}

interface ServiceWorkerLifecycleEvent {
  waitUntil: (outcome: Promise<unknown>) => void
  settled: Promise<unknown> | null
}

interface LoadedServiceWorker {
  fetchListener: (event: ServiceWorkerFetchEvent) => void
  activateListener: (event: ServiceWorkerLifecycleEvent) => void
  installListener: (event: ServiceWorkerLifecycleEvent) => void
  cacheAddAll: ReturnType<typeof vi.fn>
  cachePut: ReturnType<typeof vi.fn>
  cachesDelete: ReturnType<typeof vi.fn>
  cachesMatch: ReturnType<typeof vi.fn>
  fetchMock: ReturnType<typeof vi.fn>
  networkErrorResponse: { kind: string }
}

const APP_ORIGIN = "https://frontend-six-virid-76.vercel.app"
const NAVIGATION_FALLBACK_PATH = "/index.html"
const CURRENT_SHELL_CACHE = "healthcare-shell-v2"
const POISONED_SHELL_CACHE = "healthcare-shell-v1"
const serviceWorkerSource = readFileSync(resolve(process.cwd(), "public/sw.js"), "utf-8")

const cachedShellResponse = { kind: "cached-index-html", headers: { "content-type": "text/html" } }
const networkErrorResponse = { kind: "network-error" }

function createNetworkResponse(kind: string, status: number) {
  return {
    kind,
    status,
    clone: () => ({ kind: `${kind}-clone`, status }),
  }
}

function loadServiceWorker(cachedCacheNames: string[] = []): LoadedServiceWorker {
  const listeners = new Map<string, (event: never) => void>()
  const serviceWorkerScope = {
    location: { origin: APP_ORIGIN },
    addEventListener: (eventType: string, listener: (event: never) => void) => {
      listeners.set(eventType, listener)
    },
  }
  const shellCache = { addAll: vi.fn(), put: vi.fn(), match: vi.fn() }
  const cachesStub = {
    open: vi.fn(async () => shellCache),
    keys: vi.fn(async () => cachedCacheNames),
    match: vi.fn(async () => cachedShellResponse),
    delete: vi.fn(async () => true),
  }
  const fetchMock = vi.fn()
  const responseStub = { error: () => networkErrorResponse }

  const evaluateServiceWorker = new Function("self", "caches", "fetch", "Response", serviceWorkerSource)
  evaluateServiceWorker(serviceWorkerScope, cachesStub, fetchMock, responseStub)

  const resolveListener = (eventType: string) => {
    const listener = listeners.get(eventType)
    if (!listener) {
      throw new Error(`service worker did not register a "${eventType}" listener`)
    }
    return listener as unknown as (event: never) => void
  }

  return {
    fetchListener: (event) => resolveListener("fetch")(event as never),
    activateListener: (event) => resolveListener("activate")(event as never),
    installListener: (event) => resolveListener("install")(event as never),
    cacheAddAll: shellCache.addAll,
    cachePut: shellCache.put,
    cachesDelete: cachesStub.delete,
    cachesMatch: cachesStub.match,
    fetchMock,
    networkErrorResponse,
  }
}

function dispatchFetch(worker: LoadedServiceWorker, path: string, mode = "no-cors", method = "GET") {
  const event: ServiceWorkerFetchEvent = {
    request: { method, mode, url: `${APP_ORIGIN}${path}` },
    respondWith: vi.fn(),
    respondedWith: null,
  }
  event.respondWith = vi.fn((response: Promise<unknown>) => {
    event.respondedWith = response
  })
  worker.fetchListener(event)
  return event
}

function dispatchCrossOriginFetch(worker: LoadedServiceWorker, url: string, mode = "no-cors") {
  const event: ServiceWorkerFetchEvent = {
    request: { method: "GET", mode, url },
    respondWith: vi.fn(),
    respondedWith: null,
  }
  event.respondWith = vi.fn((response: Promise<unknown>) => {
    event.respondedWith = response
  })
  worker.fetchListener(event)
  return event
}

function dispatchLifecycle(worker: LoadedServiceWorker, listener: "activate" | "install") {
  const event: ServiceWorkerLifecycleEvent = { waitUntil: vi.fn(), settled: null }
  event.waitUntil = vi.fn((outcome: Promise<unknown>) => {
    event.settled = outcome
  })
  if (listener === "activate") {
    worker.activateListener(event)
  } else {
    worker.installListener(event)
  }
  return event
}

describe("service worker shell strategy", () => {
  afterEach(() => {
    vi.clearAllMocks()
  })

  it("should not intercept hashed build assets so the browser always receives javascript", () => {
    const worker = loadServiceWorker()

    const event = dispatchFetch(worker, "/assets/index-DhJDc12i.js")

    expect(event.respondedWith).toBeNull()
    expect(worker.fetchMock).not.toHaveBeenCalled()
  })

  it("should not intercept api requests", () => {
    const worker = loadServiceWorker()

    const event = dispatchFetch(worker, "/api/v1/patients")

    expect(event.respondedWith).toBeNull()
  })

  it("should not intercept cross origin requests", () => {
    const worker = loadServiceWorker()

    const event = dispatchCrossOriginFetch(worker, "https://fonts.googleapis.com/css2?family=Poppins")

    expect(event.respondedWith).toBeNull()
  })

  it("should not intercept non get requests", () => {
    const worker = loadServiceWorker()

    const event = dispatchFetch(worker, "/api/v1/clinical.v1.Observation", "cors", "POST")

    expect(event.respondedWith).toBeNull()
  })

  it("should never resolve cached html for a module script when the network fails", async () => {
    const worker = loadServiceWorker()
    worker.fetchMock.mockRejectedValue(new Error("network down"))

    const event = dispatchFetch(worker, "/icons.svg")

    await expect(event.respondedWith).resolves.toBe(worker.networkErrorResponse)
    expect(worker.cachesMatch).not.toHaveBeenCalled()
  })

  it("should resolve a network error instead of html for a failed non navigation request", async () => {
    const worker = loadServiceWorker()
    worker.fetchMock.mockRejectedValue(new Error("network down"))

    const event = dispatchFetch(worker, "/manifest.json", "same-origin")

    await expect(event.respondedWith).resolves.toEqual({ kind: "network-error" })
  })

  it("should fall back to the cached shell only for a failed navigation", async () => {
    const worker = loadServiceWorker()
    worker.fetchMock.mockRejectedValue(new Error("network down"))

    const event = dispatchFetch(worker, "/patients", "navigate")

    await expect(event.respondedWith).resolves.toBe(cachedShellResponse)
    expect(worker.cachesMatch).toHaveBeenCalledWith(NAVIGATION_FALLBACK_PATH)
  })

  it("should pass a successful navigation response through untouched", async () => {
    const worker = loadServiceWorker()
    const networkResponse = createNetworkResponse("index-html", 200)
    worker.fetchMock.mockResolvedValue(networkResponse)

    const event = dispatchFetch(worker, "/patients", "navigate")

    await expect(event.respondedWith).resolves.toBe(networkResponse)
  })

  it("should cache a successful navigation under the shell path", async () => {
    const worker = loadServiceWorker()
    worker.fetchMock.mockResolvedValue(createNetworkResponse("index-html", 200))

    const event = dispatchFetch(worker, "/patients", "navigate")
    await event.respondedWith

    expect(worker.cachePut).toHaveBeenCalledWith(NAVIGATION_FALLBACK_PATH, expect.anything())
  })

  it("should not cache a navigation that resolved with a failure status", async () => {
    const worker = loadServiceWorker()
    worker.fetchMock.mockResolvedValue(createNetworkResponse("not-found", 404))

    const event = dispatchFetch(worker, "/patients", "navigate")
    await event.respondedWith

    expect(worker.cachePut).not.toHaveBeenCalled()
  })

  it("should pre cache the static shell assets on install", async () => {
    const worker = loadServiceWorker()

    const event = dispatchLifecycle(worker, "install")
    await event.settled

    expect(worker.cacheAddAll).toHaveBeenCalledWith([
      NAVIGATION_FALLBACK_PATH,
      "/favicon.svg",
      "/icons.svg",
      "/manifest.json",
    ])
  })

  it("should drop the stale shell cache that served html for module scripts", async () => {
    const worker = loadServiceWorker([POISONED_SHELL_CACHE, CURRENT_SHELL_CACHE])

    const event = dispatchLifecycle(worker, "activate")
    await event.settled

    expect(worker.cachesDelete).toHaveBeenCalledWith(POISONED_SHELL_CACHE)
    expect(worker.cachesDelete).not.toHaveBeenCalledWith(CURRENT_SHELL_CACHE)
  })
})
