import { GetGreetingRequest, GetGreetingResponse } from "../pb/greeting";

class GreetingService {
  service: string = "greeting.GreetingService";
  async getGreeting(request: GetGreetingRequest): Promise<GetGreetingResponse> {
    return postRequest(
      this.service,
      "GetGreeting",
      GetGreetingRequest.toJSON(request),
    );
  }
}

export const api = {
  Greeting: new GreetingService(),
};

function postRequest(service: string, method: string, body: unknown) {
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
