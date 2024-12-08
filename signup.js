const submit = document.getElementById("submit")

submit.addEventListener("click", function() {
	const email = document.getElementById("email").value;
	const password = document.getElementById("password").value;
	console.log(email, password)

	fetch(`http://localhost:42069/signUp?email=${email}&password=${password}`, {
		method: "POST",
	})
		.then(function(response) {
			if (response.status == 201) { //Created
				document.getElementById("text").innerHTML = "Account created successfully!";
			}
			if (response.status == 401) { //Unautharized 
				document.getElementById("text").innerHTML = "An account with that email already exists!";
			}
		})
		.catch(error => console.error("Error:", error));
})

