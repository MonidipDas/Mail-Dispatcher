async function sendCampaign(){

let emails = document.getElementById("emails")
.value.split(",")

let subject = document.getElementById("subject").value

let body = document.getElementById("body").value

let res = await fetch("http://localhost:5000/send",{

method:"POST",

headers:{
"Content-Type":"application/json"
},

body:JSON.stringify({
emails:emails,
subject:subject,
body:body
})

})

let text = await res.text()

document.getElementById("status").innerText=text
}