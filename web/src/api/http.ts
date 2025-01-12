export function postRequest(service: string, method: string, body: unknown) {
  const host = "http://localhost:8090";
  return fetch(`${host}/${service}/${method}`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(body),
  }).then((response) => {
    return response.json();
  });
}
