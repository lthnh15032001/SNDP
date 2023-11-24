export const backendUrl = `http://localhost:9080`
class TunnelService {
    list = async (accessToken) => {
        const response = await fetch(`${backendUrl}/tunnel/`, {
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
}

export default TunnelService;
