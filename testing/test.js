fetch("http://localhost:8080/transform", {
	method: "POST",
	headers: {
		"Content-Type": "application/json",
	},
	body: JSON.stringify({
		script_identifier: "sum-script-001",
		created_by: "vk",
		description: "Test:add price*quantity:(336.36)",
		name: "sum script",
	}),
}).then(async (response) => {
	const text = await response.text(); // don't assume it's JSON
	console.log("Raw Response:", text);

	try {
		const data = JSON.parse(text);
		console.log("Parsed JSON:", data);
	} catch (err) {
		console.error("JSON parse error:", err.message);
	}
}).catch((error) => console.error("Request failed:", error));
