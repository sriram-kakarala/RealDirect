var uri = window.location.toString();
if (uri.indexOf("?") > 0) {
    clean_uri = uri.substring(0, uri.indexOf("?"));
    window.history.replaceState({}, document.title, clean_uri);
}

function loadmyassets() {
    console.log('myassets')
    var xhttp = new XMLHttpRequest();
    xhttp.onreadystatechange = function() {
        console.log("Ready State", this.readyState)
        if (this.readyState == 4) {
            console.log('myassets response')
            if(this.status == 200) {
                var data = JSON.parse(JSON.parse(this.responseText));
                console.log(data);
                //data = data.content;
                //data = [{"displayname":"First","docType":"asset","name":"first","owner":"sriram@gmail.com","price":12,"quantity":12}]

                var html = "<table border==\"1\"><tr>";
                for (key in data[0]) {
                    if(key == "displayname") {
                        key = "asset name";
                    }
                    if(key != "name" && key != "docType" && key != "owner") {
                        html += '<td class="tabheader" align="center">' + key + '</td>';
                    }
                }
                html += "</tr>";
                for (var i = 0; i < data.length; i++) {
                    html += '<tr  style="font-weight: bold padding: 5px" color: #000000;  align="center">';
                    for (key in data[i]) {
                        if(key != "name" && key != "docType" && key != "owner") {
                            html += '<td>' + data[i][key] + '</td>';
                        }
                    }
                    html += '</tr>';
                }
                html += "</table>";
                document.getElementById("assetsresponse").innerHTML = html;
            }
        }
    };

    
    xhttp.open("POST", "queryassetsByOwner");
    xhttp.send();
}

function loadallassets() {
    console.log('allassets')
    var xhttp = new XMLHttpRequest();
    xhttp.onreadystatechange = function() {
        console.log("Ready State", this.readyState)
        if (this.readyState == 4) {
            console.log('allassets response')
            if(this.status == 200) {
                var data = JSON.parse(JSON.parse(this.responseText));
                console.log(data);
                var html = "<table id=\"assettable\" border==\"1\"><tr>";
                for (key in data[0]) {
                    if(key == "displayname") {
                        key = "asset name";
                    }
                    if(key != "name" && key != "docType") {
                        html += '<td class="tabheader" align="center">' + key + '</td>';
                    }
                }
                html += '<td class="tabheader" align="center"> Trade </td>';
                html += "</tr>";

                var myself = getCookie("session").replace("%40", "@");

                for (var i = 0; i < data.length; i++) {
                    html += '<tr  style="font-weight: bold padding: 5px" color: #000000;  align="center">';
                    var owner = "";
                    var name = "";
                    for (key in data[i]) {
                        if(key == "name") {
                            name = data[i][key];
                        }
                        if(key != "name" && key != "docType") {
                            html += '<td>' + data[i][key] + '</td>';
                            if(key == "owner") {
                                owner = data[i][key];
                            }
                        }
                    }
                    
                    if(myself == owner) {
                        html += '<td> You Own It! </td>';
                    } else {
                        html += '<td id=' + name + ' class="purchase" onclick="purchase(event,\'' + name  + '\')"> Purchase </td>';
                    }
                    html += '</tr>';
                }


                html += "</table>";
                document.getElementById("allassetsresponse").innerHTML = html;
            }
        }
    };

    
    xhttp.open("POST", "queryAllassets");
    xhttp.send();
}

function purchase(evt, name) {
    console.log(name);
    console.log('createAsset')

    var xhttp = new XMLHttpRequest();
    xhttp.onreadystatechange = function() {
        console.log("Ready State", this.readyState)
        if (this.readyState == 4) {
            console.log('createAsset response')
            if(this.status == 200) {
                alert("Voot!! The Asset is Yours!!")
                loadallassets();
            } else {
                alert("Uh Oh!! Something went Wrong!!")
            }
        }
    };

    xhttp.open("POST", "transferasset", true);
    var params = JSON.stringify({ name: name});
    xhttp.setRequestHeader("Content-type", "application/json; charset=utf-8");
    xhttp.setRequestHeader("Content-length", params.length);
    xhttp.send(params);
}

function getCookie(name) {
    let cookie = {};
    document.cookie.split(';').forEach(function(el) {
      let [k,v] = el.split('=');
      cookie[k.trim()] = v;
    })
    return cookie[name];
  }

function createAsset() {
    console.log('createAsset')

    var assetName = document.getElementById("assetname").value;
    var assetQuantity = document.getElementById("assetquantity").value;
    var assetPrice = document.getElementById("assetprice").value;

    if(assetName == "") {
        alert("Asset Name cannot be Empty!!")
        return;
    }

    if(assetQuantity == "") {
        alert("Asset Quantity cannot be Empty!!")
        return;
    }

    if(assetPrice == "") {
        alert("Asset Price cannot be Empty!!")
        return;
    }

    var xhttp = new XMLHttpRequest();
    xhttp.onreadystatechange = function() {
        console.log("Ready State", this.readyState)
        if (this.readyState == 4) {
            console.log('createAsset response')
            if(this.status == 200) {
                document.getElementById("message").innerHTML = "Asset Created!!";
            } else {
                document.getElementById("message").innerHTML = "Asset Creation Failed!!";
            }
            document.getElementById("assetname").value = "";
            document.getElementById("assetquantity").value = "";
            document.getElementById("assetprice").value = "";
        }
    };

    
    xhttp.open("POST", "createAsset", true);
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

    if(assetType == "myAssets") {
        console.log("My Assets Clicked");
        loadmyassets();
    } else if(assetType = "all") {
        console.log("My Assets Clicked");
        loadallassets();
    }
  }
