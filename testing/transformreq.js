const fs = require("node:fs");

// Load JSON input
const data = fs.readFileSync("data.json", "utf-8");
const parsed = JSON.parse(data);

// Build request
const requestBody = {
	data: parsed, // this will be encoded properly as JSON.RawMessage
	script_identifier: "address-script-001",
	// optionally you can send `script` instead of script_identifier
};

fetch("http://localhost:8080/transform", {
	method: "POST",
	headers: {
		"Content-Type": "application/json",
	},
	body: JSON.stringify(requestBody),
})
	.then(async (res) => {
		const text = await res.text();
		try {
			console.log("✅ Response:", JSON.parse(text));
		} catch {
			console.log("❌ Raw response (not JSON):", text);
		}
	})
	.catch((err) => {
		console.error("❌ Fetch error:", err.message);
	});
