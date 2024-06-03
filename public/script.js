const map = L.map('map').setView([46.603354, 1.888334], 6); 
L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
  attribution: '© OpenStreetMap contributors'
}).addTo(map);

const randomLocationBtn = document.getElementById('random-location-btn');

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

  L.marker(randomLatLng).addTo(map)
    .bindPopup(`Vous avez découvert : [${lat.toFixed(5)}, ${lng.toFixed(5)}]`)
    .openPopup();

  map.setView(randomLatLng, 10);

  fetch('/update-coordinates', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ latitude: lat, longitude: lng }),
  })
  .then(response => response.json())
  .then(data => {
    console.log('Success:', data);
  })
  .catch((error) => {
    console.error('Error:', error);
  });
});
