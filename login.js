const submit = document.getElementById("submit")

submit.addEventListener("click", function() {
	const email = document.getElementById("email").value;
	const password = document.getElementById("password").value;
	console.log(email, password)

	fetch(`http://localhost:42069/logIn?email=${email}&password=${password}`, {
		method: "POST",
	})
		.then(function(response) {
			if (response.status == 202) { //Accepted
				document.getElementById("text").innerHTML = "Logged in successfully";
			}
			if (response.status == 401) { //Unautharized 
				document.getElementById("text").innerHTML = "Wrong email or password!";
			}
		})
		.catch(error => console.error("Error:", error));
})

