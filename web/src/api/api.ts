import { CollabCafeService } from "./collabcafe_api";
import { GreetingService } from "./greeting_api";

export const api = {
  Greeting: new GreetingService(),
  CollabCafe: new CollabCafeService(),
};
