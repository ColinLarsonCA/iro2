import { GetCollabRequest, GetCollabResponse, SearchCollabsRequest, SearchCollabsResponse } from "../pb/collabcafe";
import { postRequest } from "./http";

export class CollabCafeService {
  service: string = "collabcafe.CollabCafeService";
  async getCollab(request: GetCollabRequest): Promise<GetCollabResponse> {
    return postRequest(
      this.service,
      "GetCollab",
      GetCollabRequest.toJSON(request),
    );
  }
  async searchCollabs(request: SearchCollabsRequest): Promise<SearchCollabsResponse> {
    return postRequest(
      this.service,
      "SearchCollabs",
      SearchCollabsRequest.toJSON(request),
    );
  }
}

