import { backendUrl } from "./schedule";
class PolicyService {
  list = async (accessToken) => {
    const response = await fetch(`${backendUrl}/v1alpha1/list_policies?verbose=true`, {
      method: 'GET',
      credentials: 'same-origin',
      headers: {
        Authorization: accessToken,
      },
    });

    if (!response.ok) {
      console.error(response);
      const responseBody = await response.text();
      throw Error(responseBody || response.statusText);
    }

    return response.json();
  };

  get = async (policy, accessToken) => {
    const response = await fetch(`${backendUrl}/v1alpha1/get_policy?policy=${policy}`, {
      method: 'GET',
      credentials: 'same-origin',
      headers: {
        Authorization: accessToken,
      },
    });

    if (!response.ok) {
      console.error(response);
      const responseBody = await response.text();
      throw Error(responseBody || response.statusText);
    }

    return response.json();
  };

  delete = async (policy, accessToken) => {
    const response = await fetch(`${backendUrl}/v1alpha1/del_policy?policy=${policy}`, {
      method: 'GET',
      credentials: 'same-origin',
      headers: {
        Authorization: accessToken,
      },
    });

    if (!response.ok) {
      console.error(response);
      const responseBody = await response.text();
      throw Error(responseBody || response.statusText);
    }

    return response;
  };

  add = async (policy, accessToken) => {
    const response = await fetch(`${backendUrl}/v1alpha1/add_policy`, {
      method: 'POST',
      credentials: 'same-origin',
      headers: {
        'Content-Type': 'application/json',
        Authorization: accessToken,
      },
      body: JSON.stringify(policy),
    });

    if (!response.ok) {
      console.error(response);
      const responseBody = await response.text();
      throw Error(responseBody || response.statusText);
    }

    return response;
  };
}

export default PolicyService;
