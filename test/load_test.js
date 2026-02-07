import http from "k6/http";
import { check, sleep } from "k6";

export let options = {
  stages: [
    { duration: "30s", target: 50 }, // Ramp-up to 50 users over 30 seconds
    { duration: "1m", target: 100 }, // Ramp-up to 100 users over 1 minute
    { duration: "1m", target: 200 }, // Ramp-up to 200 users over 1 minute
    { duration: "30s", target: 0 }, // Ramp-down to 0 users over 30 seconds
  ],
  thresholds: {
    http_req_duration: ["p(95)<500"],
  },
};

const BASE_URL = "http://localhost:8080/api"; // Replace with your API base URL

export default function () {
  let res = http.get(`${BASE_URL}/kategori`);

  if (res.status !== 200) {
    console.log(`Error! Status: ${res.status}, Body: ${res.body}`);
  }
  check(res, {
    "is status 200": (r) => r.status === 200,
    "response time < 500ms": (r) => r.timings.duration < 500,
  });
  sleep(1);
}
