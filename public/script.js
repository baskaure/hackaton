const map = L.map('map').setView([46.603354, 1.888334], 6);
L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
  attribution: '© OpenStreetMap contributors'
}).addTo(map);

const randomLocationBtn = document.getElementById('random-location-btn');
let hotelMarker = null; 

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



// {
//   "type": "FeatureCollection",
//   "features": [
//     {
//       "type": "Feature",
//       "properties": {},
//       "geometry": {
//         "coordinates": [
//           [
//             [
//               2.5596690000205626,
//               51.09931136914983
//             ],
//             [
//               1.656558855932957,
//               50.91082422447218
//             ],
//             [
//               1.582314183391361,
//               50.194576243555076
//             ],
//             [
//               0.07110844621644219,
//               49.55547153393911
//             ],
//             [
//               0.4404656320381264,
//               49.47603652483795
//             ],
//             [
//               -1.2253465142589448,
//               49.384723873587575
//             ],
//             [
//               -1.3086814442448826,
//               49.704483906124324
//             ],
//             [
//               -1.930146394721362,
//               49.71776429358721
//             ],
//             [
//               -1.4562448874389986,
//               48.604339826903725
//             ],
//             [
//               -3.1785853695628816,
//               48.83938703333382
//             ],
//             [
//               -4.790532992670904,
//               48.50444268909851
//             ],
//             [
//               -4.770845187938761,
//               48.32088049152992
//             ],
//             [
//               -4.274283981637126,
//               48.295917504147326
//             ],
//             [
//               -4.617092426441701,
//               48.20093006174537
//             ],
//             [
//               -4.25806715040045,
//               48.124078035172175
//             ],
//             [
//               -4.702536608599871,
//               48.00129647660796
//             ],
//             [
//               -2.5179049071804513,
//               47.473220774498145
//             ],
//             [
//               -2.1306411650147083,
//               46.8190922670363
//             ],
//             [
//               -1.107300557093481,
//               46.11702141681366
//             ],
//             [
//               -1.2030391867158698,
//               44.50908697762128
//             ],
//             [
//               -1.6416633322764937,
//               43.36188688093836
//             ],
//             [
//               0.16965131109824938,
//               42.71835564166801
//             ],
//             [
//               0.8057849544499049,
//               42.907940803823465
//             ],
//             [
//               3.0392971135362643,
//               42.51959085499763
//             ],
//             [
//               3.046937304377991,
//               43.156937175932654
//             ],
//             [
//               4.063974597977307,
//               43.488211734862034
//             ],
//             [
//               6.402234465935635,
//               43.110780341404165
//             ],
//             [
//               7.373227174649713,
//               43.79061692414368
//             ],
//             [
//               7.592033048259225,
//               44.230220510017745
//             ],
//             [
//               6.902307944136567,
//               44.48221554674805
//             ],
//             [
//               6.934466901767678,
//               44.963991921594385
//             ],
//             [
//               6.567467901755265,
//               45.15617727329149
//             ],
//             [
//               7.110688774564551,
//               45.43446239255712
//             ],
//             [
//               6.9523095208091945,
//               45.88730522006858
//             ],
//             [
//               6.600035724964471,
//               46.34723006637046
//             ],
//             [
//               5.9998543200710515,
//               46.00289483133537
//             ],
//             [
//               6.0524633122788885,
//               46.845349095173844
//             ],
//             [
//               6.713370028850733,
//               47.217904503736435
//             ],
//             [
//               6.846091843460073,
//               47.64695837350456
//             ],
//             [
//               7.344373559490123,
//               47.68775842730403
//             ],
//             [
//               7.7578687745501895,
//               48.832215122102724
//             ],
//             [
//               4.485117256085061,
//               49.82068619792767
//             ],
//             [
//               2.97398358147467,
//               50.600194536315996
//             ],
//             [
//               2.5596690000205626,
//               51.09931136914983
//             ]
//           ]
//         ],
//         "type": "Polygon"
//       }
//     }
//   ]
// }