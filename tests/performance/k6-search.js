import { check, sleep } from "k6";
import http from "k6/http";

export const options = {
  thresholds: {
    http_req_duration: ["p(99)<50"],
    http_req_failed: ["rate<0.001"],
  },
  vus: 50,
  duration: "5m",
};

const BASE = __ENV.BASE_URL || "http://localhost:80";

export default function () {
  const res = http.get(`${BASE}/api/search?q=laptop`, {
    headers: { "x-flag-search-v2": "100" },
  });

  check(res, {
    "status is 200": (r) => r.status === 200,
  });

  sleep(0.2);
}
