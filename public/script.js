// Initialiser la carte
const map = L.map('map').setView([46.603354, 1.888334], 6); // Centre de la France

// Ajouter un layer à la carte
L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
  attribution: '© OpenStreetMap contributors'
}).addTo(map);

const randomLocationBtn = document.getElementById('random-location-btn');

// Coordonnées de la France métropolitaine
const franceBounds = {
  north: 51.089062,
  south: 42.328014,
  west: -5.142222,
  east: 8.231167
};

randomLocationBtn.addEventListener('click', () => {
  const lat = franceBounds.south + Math.random() * (franceBounds.north - franceBounds.south);
  const lng = franceBounds.west + Math.random() * (franceBounds.east - franceBounds.west);
  const randomLatLng = [lat, lng];

  // Ajouter un marqueur à la position aléatoire
  L.marker(randomLatLng).addTo(map)
    .bindPopup(`Vous avez découvert : [${lat.toFixed(5)}, ${lng.toFixed(5)}]`)
    .openPopup();

  // Centrer la carte sur la nouvelle position
  map.setView(randomLatLng, 10);
});


document.addEventListener("DOMContentLoaded", function () {
  const hotelInfoContainer = document.getElementById("hotel-info");

  fetch('/hotel')
    .then(response => response.json())
    .then(data => {
      const hotels = data.data.slice(0, 5); // Les cinq premiers hôtels dans la réponse JSON

      hotels.forEach(hotel => {
        // Créer des éléments HTML pour afficher les informations de l'hôtel
        const hotelContainer = document.createElement('div');
        hotelContainer.className = 'hotel';

        const hotelName = document.createElement('h2');
        hotelName.textContent = hotel.name;

        const hotelAddress = document.createElement('p');
        hotelAddress.textContent = `Country: ${hotel.address.countryCode}`;

        const hotelDistance = document.createElement('p');
        hotelDistance.textContent = `Distance: ${hotel.distance.value} ${hotel.distance.unit}`;

        // Ajouter les éléments au conteneur de l'hôtel
        hotelContainer.appendChild(hotelName);
        hotelContainer.appendChild(hotelAddress);
        hotelContainer.appendChild(hotelDistance);
        // Ajouter le conteneur de l'hôtel au conteneur principal
        hotelInfoContainer.appendChild(hotelContainer);
      });
    })
    .catch(error => {
      console.error('Error fetching hotel data:', error);
    });
});
