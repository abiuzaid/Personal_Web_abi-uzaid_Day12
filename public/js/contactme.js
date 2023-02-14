function sendMessage(){
    let name = document.getElementById("name").value;
    let email = document.getElementById("email").value;
    let phoneNumber = document.getElementById("phone-number").value;
    let subject = document.getElementById("subject").value;
    let message = document.getElementById("message").value;
  
    if (name == "") {
      alert("You must fill your name!");
    } else if (email == "") {
      alert("You must fill your email!");
    } else if (phoneNumber == "") {
      alert("You must fill your phone number!");
    } else if (subject == "") {
      alert("You must fill the subject!");
    } else if (message == "") {
      alert("You must fill the message!");
    } else {
      const defaultEmail = "abiuzaid1427@gmail.com";
  
      let mailTo = document.createElement('a');
      mailTo.href = `mailto:${defaultEmail}?subject=${subject}&body=Greetings! My name is ${name}, ${message}, You can contact me at ${phoneNumber}.`;
      mailTo.target = "_blank";
      mailTo.click();
  
      let clientMessage = {
        name,
        email,
        phoneNumber,
        subject,
        message
      };
  
      console.log(clientMessage);
    }
  }