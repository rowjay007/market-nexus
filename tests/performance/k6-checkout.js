import { check, sleep } from "k6";
import http from "k6/http";

export const options = {
  thresholds: {
    http_req_duration: ["p(99)<200"],
    http_req_failed: ["rate<0.001"],
  },
  scenarios: {
    checkout_load: {
      executor: "ramping-vus",
      startVUs: 5,
      stages: [
        { duration: "1m", target: 25 },
        { duration: "3m", target: 50 },
        { duration: "1m", target: 0 },
      ],
    },
  },
};

const BASE = __ENV.BASE_URL || "http://localhost:80";

export default function () {
  const payload = JSON.stringify({
    id: `o-${__VU}-${__ITER}`,
    vendorId: "v-load",
    lines: [{ sku: "sku-load", quantity: 1 }],
  });

  const res = http.post(`${BASE}/api/orders`, payload, {
    headers: {
      "Content-Type": "application/json",
      "x-flag-ordering-v2": "100",
    },
  });

  check(res, {
    "status is 200/201": (r) => r.status === 200 || r.status === 201,
  });

  sleep(1);
}
