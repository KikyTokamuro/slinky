let shorten = async () => {
    let link = document.getElementById("link").value
    let response = await fetch(`https://sliy.herokuapp.com/api/?url=${link}`)

    if (response.ok) { 
      let json = await response.json()
      let resultcard = document.getElementById("shortener-result")
      let result = document.getElementById("result")

      if (json["Err"] != undefined) {
        result.innerText = `${json["Link"]} ${json["Err"]}`
        resultcard.style.display = "block"
      } else {
        result.innerText = `${json["SLink"]}`
        resultcard.style.display = "block"
      }
    } else {
      alert("HTTP error: " + response.status)
    }
  }