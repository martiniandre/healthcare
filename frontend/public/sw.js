const CACHE_NAME_KEY = "healthcare-shell-v2";
const NAVIGATION_FALLBACK_PATH = "/index.html";
const CACHE_BYPASS_PATH_PREFIXES = ["/api/", "/assets/", "/clinical.v1.", "/auth.v1."];
const ASSETS_TO_CACHE = [NAVIGATION_FALLBACK_PATH, "/favicon.svg", "/icons.svg", "/manifest.json"];

function isSameOriginGetRequest(request) {
  return request.method === "GET" && new URL(request.url).origin === self.location.origin;
}

function isCacheBypassedRequest(request) {
  const requestPath = new URL(request.url).pathname;
  return CACHE_BYPASS_PATH_PREFIXES.some((pathPrefix) => requestPath.startsWith(pathPrefix));
}

self.addEventListener("install", (event) => {
  event.waitUntil(
    caches.open(CACHE_NAME_KEY).then((cache) => {
      return cache.addAll(ASSETS_TO_CACHE);
    })
  );
});

self.addEventListener("activate", (event) => {
  event.waitUntil(
    caches.keys().then((cacheNames) => {
      return Promise.all(
        cacheNames.map((cacheName) => {
          if (cacheName !== CACHE_NAME_KEY) {
            return caches.delete(cacheName);
          }
        })
      );
    })
  );
});

self.addEventListener("fetch", (event) => {
  if (!isSameOriginGetRequest(event.request)) {
    return;
  }

  if (isCacheBypassedRequest(event.request)) {
    return;
  }

  const isNavigationRequest = event.request.mode === "navigate";

  event.respondWith(
    fetch(event.request)
      .then((networkResponse) => {
        if (isNavigationRequest && networkResponse.status === 200) {
          const responseToCache = networkResponse.clone();
          caches.open(CACHE_NAME_KEY).then((cache) => {
            cache.put(NAVIGATION_FALLBACK_PATH, responseToCache);
          });
        }
        return networkResponse;
      })
      .catch(() => {
        if (isNavigationRequest) {
          return caches.match(NAVIGATION_FALLBACK_PATH);
        }
        return Response.error();
      })
  );
});
