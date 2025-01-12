import { GetGreetingRequest, GetGreetingResponse } from "../pb/greeting";
import { postRequest } from "./http";

export class GreetingService {
  service: string = "greeting.GreetingService";
  async getGreeting(request: GetGreetingRequest): Promise<GetGreetingResponse> {
    return postRequest(
      this.service,
      "GetGreeting",
      GetGreetingRequest.toJSON(request),
    );
  }
}

