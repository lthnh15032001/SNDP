export const backendUrl = `http://localhost:9080`
class UserService {
    list = async (accessToken) => {
        const response = await fetch(`${backendUrl}/user/`, {
            method: 'GET',
            credentials: 'omit',
            // mode: 'no-cors',
            headers: {
                Authorization: accessToken,
                // 'Access-Control-Allow-Origin': 'no-cors'
            },
        });

        if (!response.ok) {
            console.error(response);
            const responseBody = await response.text();
            throw Error(responseBody || response.statusText);
        }

        return response.json();
    };

    create = async (user, accessToken) => {
        const response = await fetch(`${backendUrl}/user/`, {
            method: 'POST',
            credentials: 'omit',
            // mode: 'no-cors',
            headers: {
                Authorization: accessToken,
                'Content-Type': 'application/json',
                // 'Access-Control-Allow-Origin': 'no-cors'
            },
            body: JSON.stringify(user),
        });

        if (!response.ok) {
            console.error(response);
            const responseBody = await response.text();
            throw Error(responseBody || response.statusText);
        }

        return response.json();
    };
    edit = async (user, id, accessToken) => {
        const response = await fetch(`${backendUrl}/user/${id}`, {
            method: 'PATCH',
            credentials: 'omit',
            // mode: 'no-cors',
            headers: {
                Authorization: accessToken,
                'Content-Type': 'application/json',
                // 'Access-Control-Allow-Origin': 'no-cors'
            },
            body: JSON.stringify(user),
        });

        if (!response.ok) {
            console.error(response);
            const responseBody = await response.text();
            throw Error(responseBody || response.statusText);
        }

        return response.json();
    };

    delete = async (id, accessToken) => {
        const response = await fetch(`${backendUrl}/user/?id=${id}`, {
            method: 'DELETE',

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

    getUser = async (id, accessToken) => {
        const response = await fetch(`${backendUrl}/user/${id}`, {
            method: 'GET',
            credentials: 'omit',
            // mode: 'no-cors',
            headers: {
                Authorization: accessToken,
                // 'Access-Control-Allow-Origin': 'no-cors'
            },
        });

        if (!response.ok) {
            console.error(response);
            const responseBody = await response.text();
            throw Error(responseBody || response.statusText);
        }

        return response.json();
    }
}

export default UserService;
