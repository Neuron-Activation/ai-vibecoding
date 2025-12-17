import http from 'k6/http';
import { sleep, check } from 'k6';
import { Rate } from 'k6/metrics';

export let errorRate = new Rate('errors');

export let options = {
  vus: 50,               // количество виртуальных пользователей
  duration: '60s',       // длительность теста
  thresholds: {
    errors: ['rate<0.01'],           // <1% ошибок
    http_req_duration: ['p(95)<500'] // 95% запросов < 500ms
  }
};

const BASE = __ENV.BASE_URL || 'http://app:8080'; // в Docker host.docker.internal
const headers = { 'Content-Type': 'application/json' };

function createNotePayload(i) {
  return JSON.stringify({ title: `note ${i}`, content: 'lorem ipsum dolor sit amet ' + i });
}

export default function () {
  // weighted scenario: mostly GET /notes, some POST, occasional analytics
  let r = Math.random();
  if (r < 0.7) {
    // GET notes (list)
    let res = http.get(`${BASE}/notes`);
    errorRate.add(res.status >= 400);
    check(res, { 'get notes 200': (r) => r.status === 200 });
  } else if (r < 0.9) {
    // POST create
    let res = http.post(`${BASE}/notes`, createNotePayload(__ITER), { headers });
    errorRate.add(res.status >= 400);
    check(res, { 'create note 201or200': (r) => r.status === 200 || r.status === 201 });
  } else {
    // analytics
    let res = http.get(`${BASE}/analytics/summary`);
    errorRate.add(res.status >= 400);
    check(res, { 'analytics 200': (r) => r.status === 200 });
  }

  sleep(0.1); // небольшая пауза
}
