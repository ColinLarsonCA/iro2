import { GetCollabRequest, GetCollabResponse, ListCollabsRequest, ListCollabsResponse, SearchCollabsRequest, SearchCollabsResponse } from "../pb/collabcafe";
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
  async listCollabs(request: ListCollabsRequest): Promise<ListCollabsResponse> {
    return postRequest(
      this.service,
      "ListCollabs",
      ListCollabsRequest.toJSON(request),
    );
  }
}

