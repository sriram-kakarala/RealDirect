var uri = window.location.toString();
if (uri.indexOf("?") > 0) {
    clean_uri = uri.substring(0, uri.indexOf("?"));
    window.history.replaceState({}, document.title, clean_uri);
}

function loadmyassets() {
    console.log('myassets')
    var xhttp = new XMLHttpRequest();
    xhttp.onreadystatechange = function() {
        if (this.readyState == 4 && this.status == 200) {
            console.log('myassets response')
        }
    };
    xhttp.open("POST", "signout");
    xhttp.send();
}

function loadallassets() {
    console.log('allassets')
    var xhttp = new XMLHttpRequest();
    xhttp.onreadystatechange = function() {
        if (this.readyState == 4 && this.status == 200) {
            console.log('allassets response')
            // document.getElementById("demo").innerHTML = this.responseText;
        }
    };
    xhttp.open("POST", "signout");
    xhttp.send();
}

function createAsset() {
    console.log('createAsset')
    var xhttp = new XMLHttpRequest();
    xhttp.onreadystatechange = function() {
        console.log("Ready State", this.readyState)
        if (this.readyState == 4 && this.status == 200) {
            console.log('createAsset response')
            // document.getElementById("demo").innerHTML = this.responseText;
        }
    };

    xhttp.open("POST", "createAsset", true);
    var assetName = document.getElementById("assetname").value;
    var assetQuantity = document.getElementById("assetquantity").value;
    var assetPrice = document.getElementById("assetprice").value;

    console.log('Asset Name ' + assetName);
    console.log('Asset Qty ' + assetQuantity);
    console.log('Asset Price ' + assetPrice);

    var params = JSON.stringify({ name: assetName, quantity: assetQuantity, price: assetPrice });
    xhttp.setRequestHeader("Content-type", "application/json; charset=utf-8");
    xhttp.setRequestHeader("Content-length", params.length);
    xhttp.send(params);
}

function loadAssets(evt, assetType) {
    // Declare all variables
    var i, tabcontent, tablinks;
  
    // Get all elements with class="tabcontent" and hide them
    tabcontent = document.getElementsByClassName("tabcontent");
    for (i = 0; i < tabcontent.length; i++) {
      tabcontent[i].style.display = "none";
    }
  
    // Get all elements with class="tablinks" and remove the class "active"
    tablinks = document.getElementsByClassName("tablinks");
    for (i = 0; i < tablinks.length; i++) {
      tablinks[i].className = tablinks[i].className.replace(" active", "");
    }
  
    // Show the current tab, and add an "active" class to the button that opened the tab
    document.getElementById(assetType).style.display = "block";
    evt.currentTarget.className += " active";
  }
