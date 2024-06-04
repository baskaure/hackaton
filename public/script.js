const map = L.map('map').setView([46.603354, 1.888334], 6);
L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
  attribution: '© OpenStreetMap contributors'
}).addTo(map);

const randomLocationBtn = document.getElementById('random-location-btn');
let hotelMarker = null; 

const franceBounds = {
  north: 51.09931136914983,
  south: 42.328014,
  west: -5.142222,
  east: 8.231167
};

randomLocationBtn.addEventListener('click', () => {
  const lat = franceBounds.south + Math.random() * (franceBounds.north - franceBounds.south);
  const lng = franceBounds.west + Math.random() * (franceBounds.east - franceBounds.west);
  const randomLatLng = [lat, lng];

  if (hotelMarker) {
    map.removeLayer(hotelMarker);
  }

  hotelMarker = L.marker(randomLatLng).addTo(map)
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

  fetch('/hotel')
    .then(response => response.json())
    .then(data => {
      const hotels = data.data.slice(0, 5);
      updateHotelInfo(hotels); 
    })
    .catch(error => {
      console.error('Error fetching hotel data:', error);
    });
});

function updateHotelInfo(hotels) {
  const hotelInfoContainer = document.getElementById("hotel-info");
  hotelInfoContainer.innerHTML = "";

  if (hotels.length === 0) {
    return;
  }

  hotels.forEach(hotel => {
    const hotelContainer = document.createElement('div');
    hotelContainer.className = 'hotel';

    const hotelName = document.createElement('h2');
    hotelName.textContent = hotel.name;

    const hotelAddress = document.createElement('p');
    hotelAddress.textContent = `Country: ${hotel.address.countryCode}`;

    const hotelDistance = document.createElement('p');
    hotelDistance.textContent = `Distance: ${hotel.distance.value} ${hotel.distance.unit}`;

    hotelContainer.appendChild(hotelName);
    hotelContainer.appendChild(hotelAddress);
    hotelContainer.appendChild(hotelDistance);
    hotelInfoContainer.appendChild(hotelContainer);
  });
}

document.addEventListener("DOMContentLoaded", function () {
  fetch('/hotel')
    .then(response => response.json())
    .then(data => {
      const hotels = data.data.slice(0, 5);
      updateHotelInfo(hotels); 
    })
    .catch(error => {
      console.error('Error fetching hotel data:', error);
    });
});