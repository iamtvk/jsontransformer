const script = `{
	"name": FirstName & " " & Surname,
	"mobile": Phone[type = "mobile"].number
}`;

fetch("http://localhost:8080/create-script", {
	method: "POST",
	headers: {
		"Content-Type": "application/json",
	},
	body: JSON.stringify({
		script_identifier: "address-script-001",
		script: script,
		created_by: "vk",
		description: "Test script upload",
	}),
})
	.then((response) => response.json())
	.then((data) => console.log("Success:", data))
	.catch((error) => console.error("Error:", error));
