export const backendUrl = `http://localhost:8080`
class ScheduleService {
  list = async (accessToken) => {
    const response = await fetch(`${backendUrl}/v1alpha1/list_schedules?verbose=true`, {
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

  get = async (schedule, accessToken) => {
    const response = await fetch(
      `${backendUrl}/v1alpha1/get_schedule?schedule=${schedule}`,
      {
        method: 'GET',
        credentials: 'same-origin',
        headers: {
          Authorization: accessToken,
        },
      }
    );

    if (!response.ok) {
      console.error(response);
      const responseBody = await response.text();
      throw Error(responseBody || response.statusText);
    }

    return response.json();
  };

  delete = async (schedule, accessToken) => {
    const response = await fetch(
      `${backendUrl}/v1alpha1/del_schedule?schedule=${schedule}`,
      {
        method: 'GET',
        credentials: 'same-origin',
        headers: {
          Authorization: accessToken,
        },
      }
    );

    if (!response.ok) {
      const responseBody = await response.text();
      throw Error(responseBody || response.statusText);
    }

    return response;
  };

  add = async (schedule, accessToken) => {
    const response = await fetch(`${backendUrl}/v1alpha1/add_schedule`, {
      method: 'POST',
      credentials: 'same-origin',
      headers: {
        'Content-Type': 'application/json',
        Authorization: accessToken,
      },
      body: JSON.stringify(schedule),
    });

    if (!response.ok) {
      console.error(response);
      const responseBody = await response.text();
      throw Error(responseBody || response.statusText);
    }

    return response;
  };

  timezones = async (accessToken) => {
    const response = await fetch(`${backendUrl}/v1alpha1/time_zones`, {
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
}

export default ScheduleService;
