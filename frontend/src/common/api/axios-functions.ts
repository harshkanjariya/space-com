import { apiClient } from "./axios-client.ts";

function defaultHeaders() {
  const headers: Record<string, string> = {};
  const token = localStorage.getItem("token");
  if (token) {
    headers["Authorization"] = "Bearer " + token;
  }
  if (process.env.NODE_ENV === "local") {
    headers["xx-kong-userid"] = "6523afb7dbce5a44eb0388eb";
  }
  return headers;
}

export async function get(
  path: string,
  query?: Record<string, string | number | boolean | (string | number | boolean)[] | undefined>,
  headers: Record<string, string> = defaultHeaders(),
) {
  const queryParts: any[] = [];
  if (query) {
    Object.keys(query).forEach((key) => {
      if (Array.isArray(query[key])) {
        query[key].forEach((value) => {
          queryParts.push(`${key}=${value}`);
        });
      } else if (typeof query[key] === "string" || typeof query[key] === "number" || typeof query[key] === "boolean") {
        queryParts.push(`${key}=${query[key]}`);
      }
    });
  }
  let url = path;
  if (queryParts.length) {
    url += "?";
    url += queryParts.join("&");
  }
  const response = await apiClient.get(url, {
    headers: {
      ...defaultHeaders(),
      ...headers,
    },
  });
  return response.data;
}

export async function post(path: string, data?: any, headers: Record<string, string> = {}) {
  const response = await apiClient.post(path, data, {
    headers: {
      ...defaultHeaders(),
      ...headers,
    },
  });
  return response.data;
}

export async function put(path: string, data?: any, headers: Record<string, string> = {}) {
  const response = await apiClient.put(path, data, {
    headers: {
      ...defaultHeaders(),
      ...headers,
    },
  });
  return response.data;
}

export async function deleteApi(path: string, headers: Record<string, string> = {}) {
  const response = await apiClient.delete(path, {
    headers: {
      ...defaultHeaders(),
      ...headers,
    },
  });
  return response.data;
}
