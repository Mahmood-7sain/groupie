// This function will be called when the Google Maps API is loaded
function initMap() {
    // Define the initial center of the map (you can adjust this)
    var mapCenter = { lat: 0, lng: 0 };

    // Create a new Google Map
    var map = new google.maps.Map(document.getElementById('map'), {
        zoom: 2, // Adjust the initial zoom level
        center: mapCenter
    });

    // Make an HTTP request to your Go server to get artist data
    fetch('/map?id=1') // Adjust the URL based on your server's configuration
        .then(response => response.json())
        .then(data => {
            // Loop through the received location coordinates and add markers to the map
            data.locCoords.forEach(coords => {
                var marker = new google.maps.Marker({
                    position: { lat: coords.lat, lng: coords.lng },
                    map: map,
                    title: data.name // You can customize the marker title as needed
                });
            });
        })
        .catch(error => console.error('Error fetching artist data:', error));
}