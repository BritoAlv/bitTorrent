export async function get(url) {
    let response = await fetch(url)

    if (!response.ok) {
        throw new Error("There was an error while fetching");
    }

    return response.json();
}

export async function post(url, request) {
    let response = await fetch(url, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(request),
    });

    if (!response.ok) {
        throw new Error("There was an error while fetching");
    }

    return response.json();
}


// Timeout version of post
// export async function post(url, request, timeout = 2000) { // Default timeout: 10s
//     const controller = new AbortController();
//     const signal = controller.signal;

//     // Set a timeout to abort the request
//     const timeoutId = setTimeout(() => controller.abort(), timeout);

//     try {
//         let response = await fetch(url, {
//             method: "POST",
//             headers: {
//                 "Content-Type": "application/json",
//             },
//             body: JSON.stringify(request),
//             signal // Attach the abort signal
//         });

//         clearTimeout(timeoutId); // Clear timeout if request completes

//         if (!response.ok) {
//             throw new Error("There was an error while fetching");
//         }

//         return response.json(); // Return parsed JSON
//     } catch (error) {
//         if (error.name === "AbortError") {
//             throw new Error("Request timed out after " + timeout / 1000 + " seconds");
//         }
//         throw error; // Other fetch errors
//     }
// }
